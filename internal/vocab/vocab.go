package vocab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/trankhanh040147/langtut/internal/constants"
)

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
	if err := sonic.Unmarshal(data, lib); err != nil {
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

	// Marshal to JSON using sonic
	jsonBytes, err := sonic.Marshal(lib)
	if err != nil {
		return fmt.Errorf("failed to marshal vocab library: %w", err)
	}

	// Format with indentation
	var buf bytes.Buffer
	if err := json.Indent(&buf, jsonBytes, "", "  "); err != nil {
		return fmt.Errorf("failed to indent JSON: %w", err)
	}
	data := buf.Bytes()

	// Atomic write: write to temporary file first, then rename
	tmpPath := vocabPath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0600); err != nil {
		// Clean up temp file if it was partially created
		if removeErr := os.Remove(tmpPath); removeErr != nil {
			fmt.Fprintf(os.Stderr, "level=warn msg=\"failed to remove temp file\" file=\"%s\" err=\"%v\"\n", tmpPath, removeErr)
		}
		return fmt.Errorf("failed to write vocab file: %w", err)
	}

	// Atomically rename temp file to final destination
	if err := os.Rename(tmpPath, vocabPath); err != nil {
		// Clean up temp file on rename failure
		if removeErr := os.Remove(tmpPath); removeErr != nil {
			fmt.Fprintf(os.Stderr, "level=warn msg=\"failed to remove temp file\" file=\"%s\" err=\"%v\"\n", tmpPath, removeErr)
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
