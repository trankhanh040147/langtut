package vocab

import (
	"context"

	"github.com/trankhanh040147/langtut/internal/constants"
	"github.com/trankhanh040147/langtut/internal/ui"
	"github.com/trankhanh040147/langtut/internal/vocab"
)

// NewAddModelWithConfig creates a new addModel with API client and language configured
// This ensures consistent behavior between CLI and list mode add workflows
func NewAddModelWithConfig(term string, lib *vocab.Library, apiClient MeaningInfoGenerator, language string) *addModel {
	ctx, cancel := context.WithCancel(context.Background())
	return &addModel{
		BaseModel:          ui.BaseModel{},
		term:               term,
		definition:         "",
		examples:           []string{},
		currentField:       fieldTerm,
		editingField:       -1,
		library:            lib,
		apiClient:          apiClient,
		language:           language,
		isEditMode:         false,
		isAppendMode:       false,
		existingVocab:      nil,
		existingMeanings:   []vocab.Meaning{},
		originalMeaning:    nil,
		typeSuggestions:    []string{},
		selectedSuggestion: -1,
		showSuggestions:    false,
		ctx:                ctx,
		cancel:             cancel,
		keys:               constants.DefaultKeyMap(),
	}
}

