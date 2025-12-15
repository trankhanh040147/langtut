package vocab

import "time"

// POS type constants
const (
	TypeVerb        = "verb"
	TypeNoun        = "noun"
	TypeAdjective   = "adjective"
	TypeAdverb      = "adverb"
	TypeIdiom       = "idiom"
	TypePhrasalVerb = "phrasal_verb"
	TypeCollocation = "collocation"
)

// Meaning represents a single meaning/definition of a vocabulary term
type Meaning struct {
	ID         int      `json:"id"`
	Type       string   `json:"type"`    // POS type (verb, noun, etc.)
	Context    string   `json:"context"` // Tag or sentence context
	Definition string   `json:"definition"`
	Examples   []string `json:"examples"`
}

// Vocab represents a vocabulary term with multiple meanings
type Vocab struct {
	ID        string    `json:"id"`
	Term      string    `json:"term"`
	Meanings  []Meaning `json:"meanings"`
	Language  string    `json:"language"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
}

// GetPOSTypes returns a slice of valid POS types for autocomplete
func GetPOSTypes() []string {
	return []string{
		TypeVerb,
		TypeNoun,
		TypeAdjective,
		TypeAdverb,
		TypeIdiom,
		TypePhrasalVerb,
		TypeCollocation,
	}
}
