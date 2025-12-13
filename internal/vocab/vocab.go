package vocab

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/trankhanh040147/langtut/internal/constants"
)

// Word represents a vocabulary word
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
	Words    map[string]*Word `json:"words"`
	Metadata Metadata         `json:"metadata"`
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
		Words: make(map[string]*Word),
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

	// Ensure words map is initialized
	if lib.Words == nil {
		lib.Words = make(map[string]*Word)
	}

	return lib, nil
}

// Save saves the vocabulary library to JSON file
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

	// Write to file
	if err := os.WriteFile(vocabPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write vocab file: %w", err)
	}

	return nil
}

// NormalizeWord normalizes a word to a case-insensitive key
func NormalizeWord(word string) string {
	return strings.ToLower(strings.TrimSpace(word))
}

// GetVocabPath returns the full vocab file path
func GetVocabPath() string {
	return fmt.Sprintf("%s/vocab.json", constants.GetConfigDir())
}

// AddWord adds a word to the library
func (lib *Library) AddWord(word *Word) {
	key := NormalizeWord(word.Word)
	lib.Words[key] = word
}

// GetWord retrieves a word from the library
func (lib *Library) GetWord(word string) (*Word, bool) {
	key := NormalizeWord(word)
	w, ok := lib.Words[key]
	return w, ok
}

// DeleteWord removes a word from the library
func (lib *Library) DeleteWord(word string) bool {
	key := NormalizeWord(word)
	if _, exists := lib.Words[key]; exists {
		delete(lib.Words, key)
		return true
	}
	return false
}

// GetAllWords returns all words as a slice, sorted by word
func (lib *Library) GetAllWords() []*Word {
	words := make([]*Word, 0, len(lib.Words))
	for _, word := range lib.Words {
		words = append(words, word)
	}
	return words
}



