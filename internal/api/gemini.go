package api

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/google/generative-ai-go/genai"
	"github.com/trankhanh040147/langtut/internal/constants"
	"github.com/trankhanh040147/langtut/internal/vocab"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var (
	// jsonBlockRegex matches markdown code blocks (```json ... ``` or ``` ... ```)
	jsonBlockRegex = regexp.MustCompile(`(?s)^\s*` + "```" + `(?:json)?\s*(.*?)\s*` + "```" + `\s*$`)
	// exampleListRegex matches numbered list items or bullet points
	exampleListRegex = regexp.MustCompile(`^\d+\.\s*(.+)$|^[-*]\s*(.+)$`)
)

// Client wraps the Gemini API client
type Client struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

// NewClient creates a new Gemini API client
func NewClient(apiKey string) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	model := client.GenerativeModel(constants.APIModel)

	return &Client{
		client: client,
		model:  model,
	}, nil
}

// StreamChat streams a chat completion
func (c *Client) StreamChat(ctx context.Context, prompt string) (<-chan string, <-chan error) {
	ch := make(chan string, 10)
	errCh := make(chan error, 1)

	go func() {
		defer close(ch)
		defer close(errCh)

		// Create context with timeout
		timeoutCtx, cancel := context.WithTimeout(ctx, constants.APITimeout)
		defer cancel()

		iter := c.model.GenerateContentStream(timeoutCtx, genai.Text(prompt))
		for {
			resp, err := iter.Next()
			if err != nil {
				// iterator.Done signals normal end-of-stream, not an error
				if err == iterator.Done {
					break
				}
				errCh <- fmt.Errorf("stream error: %w", err)
				return
			}

			if resp == nil {
				break
			}

			// Extract text from response
			for _, cand := range resp.Candidates {
				if cand.Content != nil {
					for _, part := range cand.Content.Parts {
						if text, ok := part.(genai.Text); ok {
							select {
							case ch <- string(text):
							case <-timeoutCtx.Done():
								errCh <- timeoutCtx.Err()
								return
							}
						}
					}
				}
			}
		}
	}()

	return ch, errCh
}

// Chat sends a non-streaming chat completion
func (c *Client) Chat(ctx context.Context, prompt string) (string, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, constants.APITimeout)
	defer cancel()

	resp, err := c.model.GenerateContent(timeoutCtx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("API error: %w", err)
	}

	var result string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				if text, ok := part.(genai.Text); ok {
					result += string(text)
				}
			}
		}
	}

	return result, nil
}

// GenerateMeaningInfo generates structured meaning information for a term
func (c *Client) GenerateMeaningInfo(ctx context.Context, term, context, language string) (vocab.Meaning, error) {
	var meaning vocab.Meaning

	prompt := fmt.Sprintf(`Generate meaning information for the term "%s" in %s.`, term, language)
	if context != "" {
		prompt += fmt.Sprintf(` Context: %s.`, context)
	}
	prompt += fmt.Sprintf(`

Return the response as a valid JSON object with this exact structure:
{
  "type": "verb|noun|adjective|adverb|idiom|phrasal_verb|collocation",
  "definition": "clear and concise definition",
  "context": "%s",
  "examples": ["example sentence 1", "example sentence 2", "example sentence 3"]
}

The type must be one of: verb, noun, adjective, adverb, idiom, phrasal_verb, collocation.
The definition should be clear and concise.
Examples should be practical and demonstrate common usage.
Return ONLY the JSON object, no markdown formatting or additional text.`, context)

	response, err := c.Chat(ctx, prompt)
	if err != nil {
		return meaning, fmt.Errorf("failed to generate meaning info: %w", err)
	}

	// Try to parse JSON response
	response = strings.TrimSpace(response)
	// Remove markdown code blocks if present using regex
	// Matches: ```json ... ``` or ``` ... ``` with optional whitespace
	if matches := jsonBlockRegex.FindStringSubmatch(response); len(matches) > 1 {
		response = strings.TrimSpace(matches[1])
	}

	var result struct {
		Type       string   `json:"type"`
		Definition string   `json:"definition"`
		Context    string   `json:"context"`
		Examples   []string `json:"examples"`
	}

	if err := sonic.Unmarshal([]byte(response), &result); err != nil {
		// Log error before falling back to legacy parser
		fmt.Fprintf(os.Stderr, "level=warn msg=\"failed to parse API JSON response\" err=\"%v\"\n", err)
		// Fallback to old parsing logic if JSON parse fails
		definition, examples := parseWordInfoResponse(response)
		meaning.Type = vocab.TypeNoun // Default type
		meaning.Definition = definition
		meaning.Context = context
		meaning.Examples = examples
		return meaning, nil
	}

	// Validate type against constants
	validTypes := vocab.GetPOSTypes()
	typeValid := false
	for _, validType := range validTypes {
		if result.Type == validType {
			typeValid = true
			break
		}
	}
	if !typeValid {
		result.Type = vocab.TypeNoun // Default to noun if invalid
	}

	meaning.Type = result.Type
	meaning.Definition = result.Definition
	meaning.Context = result.Context
	if len(result.Examples) == 0 {
		meaning.Examples = []string{"No examples generated"}
	} else {
		meaning.Examples = result.Examples
	}

	return meaning, nil
}

// GenerateWordInfo generates meaning and examples for a word
// Deprecated: Use GenerateMeaningInfo instead
func (c *Client) GenerateWordInfo(ctx context.Context, word, language string) (string, []string, error) {
	prompt := fmt.Sprintf(`Generate meaning and 3 example sentences for the word "%s" in %s.

Return the response in the following markdown format:
## Meaning
[meaning here]

## Examples
1. [first example sentence]
2. [second example sentence]
3. [third example sentence]

Make the meaning clear and concise. Examples should be practical and demonstrate common usage.`, word, language)

	response, err := c.Chat(ctx, prompt)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate word info: %w", err)
	}

	// Parse markdown response
	meaning, examples := parseWordInfoResponse(response)
	return meaning, examples, nil
}

// parseWordInfoResponse parses the markdown response to extract meaning and examples
func parseWordInfoResponse(response string) (string, []string) {
	lines := strings.Split(response, "\n")
	var meaning strings.Builder
	var examples []string
	var inMeaningSection bool
	var inExamplesSection bool

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check for section headers
		if strings.HasPrefix(line, "## Meaning") {
			inMeaningSection = true
			inExamplesSection = false
			continue
		}
		if strings.HasPrefix(line, "## Examples") {
			inMeaningSection = false
			inExamplesSection = true
			continue
		}

		// Collect meaning
		if inMeaningSection && !inExamplesSection {
			if meaning.Len() > 0 {
				meaning.WriteString(" ")
			}
			meaning.WriteString(line)
		}

		// Collect examples
		if inExamplesSection {
			// Match numbered list items: "1. example" or "- example" or "* example"
			matches := exampleListRegex.FindStringSubmatch(line)
			if len(matches) > 0 {
				example := matches[1]
				if example == "" && len(matches) > 2 {
					example = matches[2]
				}
				if example != "" {
					examples = append(examples, strings.TrimSpace(example))
				}
			} else if len(examples) > 0 {
				// Continuation of previous example
				examples[len(examples)-1] += " " + line
			}
		}
	}

	// If we didn't find sections, try to extract from plain text
	if meaning.Len() == 0 && len(examples) == 0 {
		// Try to find meaning in first paragraph
		for idx, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") {
				if meaning.Len() > 0 {
					meaning.WriteString(" ")
				}
				meaning.WriteString(line)
				// Look for examples after meaning
				for j := idx + 1; j < len(lines) && j < idx+5; j++ {
					exLine := strings.TrimSpace(lines[j])
					if exLine != "" && !strings.HasPrefix(exLine, "#") {
						matches := exampleListRegex.FindStringSubmatch(exLine)
						if len(matches) > 0 {
							example := matches[1]
							if example == "" && len(matches) > 2 {
								example = matches[2]
							}
							if example != "" {
								examples = append(examples, strings.TrimSpace(example))
							}
						}
					}
				}
				break
			}
		}
	}

	meaningStr := strings.TrimSpace(meaning.String())
	if meaningStr == "" {
		meaningStr = "No meaning generated"
	}

	// Ensure we have at least one example
	if len(examples) == 0 {
		examples = []string{"No examples generated"}
	}

	return meaningStr, examples
}

// Close closes the client
func (c *Client) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}
