package api

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"github.com/trankhanh040147/langtut/internal/constants"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
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

// Close closes the client
func (c *Client) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}
