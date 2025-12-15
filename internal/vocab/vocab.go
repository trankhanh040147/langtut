package vocab

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/trankhanh040147/langtut/internal/constants"
)

// Word is a legacy type kept for backward compatibility during transition
// Deprecated: Use Vocab instead
type Word struct {
	Word      string    `json:"word"`
	Meaning   string    `json:"meaning"`
	Language  string    `json:"language"`
	Examples  []string  `json:"examples"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
}

// Metadata contains library metadata
type Metadata struct {
	Version     string    `json:"version"`
	LastUpdated time.Time `json:"last_updated"`
}

// Library represents the vocabulary library
type Library struct {
	Vocabs   map[string]*Vocab `json:"vocabs"`
	Metadata Metadata          `json:"metadata"`
}

// Load loads the vocabulary library from JSON file
func Load() (*Library, error) {
	vocabPath := GetVocabPath()
	vocabDir := constants.GetConfigDir()

	// Ensure config directory exists
	if err := os.MkdirAll(vocabDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	lib := &Library{
		Vocabs: make(map[string]*Vocab),
		Metadata: Metadata{
			Version:     "1.0",
			LastUpdated: time.Now(),
		},
	}

	// If vocab file doesn't exist, return empty library
	if _, err := os.Stat(vocabPath); os.IsNotExist(err) {
		return lib, nil
	}

	// Read vocab file
	data, err := os.ReadFile(vocabPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read vocab file: %w", err)
	}

	// Parse JSON
	if err := json.Unmarshal(data, lib); err != nil {
		// Backup corrupted file
		backupPath := vocabPath + ".corrupted"
		if writeErr := os.WriteFile(backupPath, data, 0600); writeErr == nil {
			// If backup succeeds, return empty library
			return lib, nil
		}
		return nil, fmt.Errorf("failed to parse vocab file: %w", err)
	}

	// Ensure vocabs map is initialized
	if lib.Vocabs == nil {
		lib.Vocabs = make(map[string]*Vocab)
	}

	return lib, nil
}

// Save saves the vocabulary library to JSON file using atomic write.
// It writes to a temporary file first, then atomically renames it to the final destination.
// This ensures data integrity even if the process is interrupted during the write operation.
func Save(lib *Library) error {
	vocabPath := GetVocabPath()
	vocabDir := constants.GetConfigDir()

	// Ensure config directory exists
	if err := os.MkdirAll(vocabDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Update metadata
	lib.Metadata.LastUpdated = time.Now()

	// Marshal to JSON
	data, err := json.MarshalIndent(lib, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal vocab library: %w", err)
	}

	// Atomic write: write to temporary file first, then rename
	tmpPath := vocabPath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0600); err != nil {
		// Clean up temp file if it was partially created
		if removeErr := os.Remove(tmpPath); removeErr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to remove temp file %s: %v\n", tmpPath, removeErr)
		}
		return fmt.Errorf("failed to write vocab file: %w", err)
	}

	// Atomically rename temp file to final destination
	if err := os.Rename(tmpPath, vocabPath); err != nil {
		// Clean up temp file on rename failure
		if removeErr := os.Remove(tmpPath); removeErr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to remove temp file %s: %v\n", tmpPath, removeErr)
		}
		return fmt.Errorf("failed to rename vocab file: %w", err)
	}

	return nil
}

// NormalizeTerm normalizes a term to a case-insensitive key
func NormalizeTerm(term string) string {
	return strings.ToLower(strings.TrimSpace(term))
}

// NormalizeWord is kept for backward compatibility
// Deprecated: Use NormalizeTerm instead
func NormalizeWord(word string) string {
	return NormalizeTerm(word)
}

// GetVocabPath returns the full vocab file path
func GetVocabPath() string {
	return fmt.Sprintf("%s/vocab.json", constants.GetConfigDir())
}

// AddVocab adds a vocab to the library
func (lib *Library) AddVocab(v *Vocab) {
	key := NormalizeTerm(v.Term)
	lib.Vocabs[key] = v
}

// GetVocab retrieves a vocab from the library
func (lib *Library) GetVocab(term string) (*Vocab, bool) {
	key := NormalizeTerm(term)
	v, ok := lib.Vocabs[key]
	return v, ok
}

// DeleteVocab removes a vocab from the library
func (lib *Library) DeleteVocab(term string) bool {
	key := NormalizeTerm(term)
	if _, exists := lib.Vocabs[key]; exists {
		delete(lib.Vocabs, key)
		return true
	}
	return false
}

// GetAllVocabs returns all vocabs as a slice
func (lib *Library) GetAllVocabs() []*Vocab {
	vocabs := make([]*Vocab, 0, len(lib.Vocabs))
	for _, v := range lib.Vocabs {
		vocabs = append(vocabs, v)
	}
	return vocabs
}

// GetNextMeaningID returns the next available meaning ID for a vocab
func (v *Vocab) GetNextMeaningID() int {
	maxID := 0
	for _, m := range v.Meanings {
		if m.ID > maxID {
			maxID = m.ID
		}
	}
	return maxID + 1
}

// Backward compatibility aliases
// Deprecated: Use AddVocab instead
func (lib *Library) AddWord(word *Word) {
	// Convert Word to Vocab for backward compatibility
	v := &Vocab{
		ID:        NormalizeTerm(word.Word),
		Term:      word.Word,
		Language:  word.Language,
		Tags:      word.Tags,
		CreatedAt: word.CreatedAt,
	}
	if word.Meaning != "" {
		v.Meanings = []Meaning{
			{
				ID:         1,
				Type:       TypeNoun, // Default type
				Definition: word.Meaning,
				Examples:   word.Examples,
			},
		}
	}
	lib.AddVocab(v)
}

// Deprecated: Use GetVocab instead
func (lib *Library) GetWord(word string) (*Word, bool) {
	v, ok := lib.GetVocab(word)
	if !ok {
		return nil, false
	}
	// Convert Vocab to Word (use first meaning)
	w := &Word{
		Word:      v.Term,
		Language:  v.Language,
		Tags:      v.Tags,
		CreatedAt: v.CreatedAt,
	}
	if len(v.Meanings) > 0 {
		w.Meaning = v.Meanings[0].Definition
		w.Examples = v.Meanings[0].Examples
	}
	return w, true
}

// Deprecated: Use DeleteVocab instead
func (lib *Library) DeleteWord(word string) bool {
	return lib.DeleteVocab(word)
}

// Deprecated: Use GetAllVocabs instead
func (lib *Library) GetAllWords() []*Word {
	vocabs := lib.GetAllVocabs()
	words := make([]*Word, 0, len(vocabs))
	for _, v := range vocabs {
		w := &Word{
			Word:      v.Term,
			Language:  v.Language,
			Tags:      v.Tags,
			CreatedAt: v.CreatedAt,
		}
		if len(v.Meanings) > 0 {
			w.Meaning = v.Meanings[0].Definition
			w.Examples = v.Meanings[0].Examples
		}
		words = append(words, w)
	}
	return words
}
