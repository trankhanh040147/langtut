package vocab

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trankhanh040147/langtut/internal/vocab"
)

func (m *addModel) loadCurrentFieldToBuffer() {
	switch m.currentField {
	case fieldTerm:
		m.editBuffer = m.term
	case fieldContext:
		m.editBuffer = m.context
	case fieldType:
		m.editBuffer = m.meaningType
	case fieldDefinition:
		m.editBuffer = m.definition
	case fieldExample1:
		// First example
		if len(m.examples) > 0 {
			m.editBuffer = m.examples[0]
		} else {
			m.editBuffer = ""
		}
	case fieldExample2:
		// Second example
		if len(m.examples) > 1 {
			m.editBuffer = m.examples[1]
		} else {
			m.editBuffer = ""
		}
	case fieldExample3:
		// Third example
		if len(m.examples) > 2 {
			m.editBuffer = m.examples[2]
		} else {
			m.editBuffer = ""
		}
	default:
		m.editBuffer = ""
	}
}

func (m *addModel) setExample(index int, value string) {
	// Ensure slice has enough capacity for the index
	for len(m.examples) <= index {
		m.examples = append(m.examples, "")
	}
	m.examples[index] = value
}

func (m *addModel) saveCurrentField() {
	switch m.editingField {
	case fieldTerm:
		m.term = m.editBuffer
	case fieldContext:
		m.context = m.editBuffer
	case fieldType:
		m.meaningType = m.editBuffer
	case fieldDefinition:
		m.definition = m.editBuffer
	case fieldExample1:
		m.setExample(0, m.editBuffer)
	case fieldExample2:
		m.setExample(1, m.editBuffer)
	case fieldExample3:
		m.setExample(2, m.editBuffer)
	default:
		// Ignore invalid field
	}
}

func (m *addModel) saveAndAdvance() tea.Cmd {
	wasTermField := m.editingField == fieldTerm
	wasContextField := m.editingField == fieldContext

	m.saveCurrentField()
	m.editingField = -1
	m.editBuffer = ""
	m.showSuggestions = false

	// Check for duplicate term after saving Term field
	if wasTermField {
		m.checkDuplicateTerm()
	}

	// Trigger generation after Context field is saved (or skipped)
	// This ensures user can optionally fill Context before generation
	var genCmd tea.Cmd
	if wasContextField {
		genCmd = m.triggerGenerationIfReady()
	}

	m.currentField++
	if m.currentField >= fieldCount {
		m.currentField = fieldTerm
	}

	return genCmd
}

func (m *addModel) checkDuplicateTerm() {
	if m.term == "" {
		return
	}
	if existingVocab, exists := m.library.GetVocab(m.term); exists {
		if !m.isEditMode {
			m.isAppendMode = true
			m.existingVocab = existingVocab
			m.existingMeanings = existingVocab.Meanings
		}
	} else {
		m.isAppendMode = false
		m.existingVocab = nil
		m.existingMeanings = []vocab.Meaning{}
	}
}

// triggerGenerationIfReady checks if Term is filled and triggers generation
// Generation happens after Term is saved, optionally including Context if provided
// Returns a tea.Cmd if generation should be triggered, nil otherwise
func (m *addModel) triggerGenerationIfReady() tea.Cmd {
	// Only generate for new meanings (not edit mode)
	if m.isEditMode {
		return nil
	}

	// Only generate if Term is filled and we haven't generated yet
	if m.term == "" || m.isGenerating {
		return nil
	}

	// Only generate if definition is empty (meaning we haven't generated yet)
	if m.definition != "" {
		return nil
	}

	// Only generate if API client is available
	if m.apiClient == nil {
		return nil
	}

	// Trigger generation
	m.isGenerating = true
	return m.generateMeaningInfo()
}

func (m *addModel) updateTypeSuggestions() {
	prefix := strings.ToLower(m.editBuffer)
	allTypes := vocab.GetPOSTypes()
	m.typeSuggestions = []string{}

	for _, t := range allTypes {
		if strings.HasPrefix(strings.ToLower(t), prefix) {
			m.typeSuggestions = append(m.typeSuggestions, t)
		}
	}

	// Limit to 5 suggestions
	if len(m.typeSuggestions) > 5 {
		m.typeSuggestions = m.typeSuggestions[:5]
	}

	m.showSuggestions = len(m.typeSuggestions) > 0 && prefix != ""
	if m.showSuggestions {
		if m.selectedSuggestion >= len(m.typeSuggestions) {
			m.selectedSuggestion = len(m.typeSuggestions) - 1
		}
		if m.selectedSuggestion < 0 && len(m.typeSuggestions) > 0 {
			m.selectedSuggestion = 0
		}
	} else {
		m.selectedSuggestion = -1
	}
}

func (m *addModel) handleAutocompleteNav(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Down) {
		if m.selectedSuggestion < len(m.typeSuggestions)-1 {
			m.selectedSuggestion++
		}
		return m, nil
	}

	if key.Matches(msg, m.keys.Up) {
		if m.selectedSuggestion > 0 {
			m.selectedSuggestion--
		}
		return m, nil
	}

	if key.Matches(msg, m.keys.Enter) {
		// Select suggestion
		if m.selectedSuggestion >= 0 && m.selectedSuggestion < len(m.typeSuggestions) {
			m.editBuffer = m.typeSuggestions[m.selectedSuggestion]
			m.showSuggestions = false
			m.selectedSuggestion = -1
		}
		// Save and advance
		return m, m.saveAndAdvance()
	}

	if key.Matches(msg, m.keys.Esc) {
		m.showSuggestions = false
		m.selectedSuggestion = -1
		m.editingField = -1
		m.editBuffer = ""
		return m, nil
	}

	if key.Matches(msg, m.keys.Tab) {
		// Select suggestion
		if m.selectedSuggestion >= 0 && m.selectedSuggestion < len(m.typeSuggestions) {
			m.editBuffer = m.typeSuggestions[m.selectedSuggestion]
			m.showSuggestions = false
			m.selectedSuggestion = -1
		}
		// Save and advance
		return m, m.saveAndAdvance()
	}

	return nil, nil
}

func (m *addModel) handleDefaultEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Enter) {
		// Save current field and advance to next
		// This may trigger generation if Term/Context was saved
		return m, m.saveAndAdvance()
	}

	if key.Matches(msg, m.keys.Esc) {
		// Cancel editing
		m.editingField = -1
		m.editBuffer = ""
		m.showSuggestions = false
		m.selectedSuggestion = -1
		return m, nil
	}

	if key.Matches(msg, m.keys.Backspace) {
		if len(m.editBuffer) > 0 {
			m.editBuffer = m.editBuffer[:len(m.editBuffer)-1]
			// Update autocomplete for type field
			if m.editingField == fieldType {
				m.updateTypeSuggestions()
			}
		}
		return m, nil
	}

	if key.Matches(msg, m.keys.Tab) {
		// Save and move to next
		// This may trigger generation if Term/Context was saved
		return m, m.saveAndAdvance()
	}

	// Handle character input
	if len(msg.Runes) > 0 {
		m.editBuffer += string(msg.Runes)
		// Update autocomplete for type field
		if m.editingField == fieldType {
			m.updateTypeSuggestions()
		}
	}
	return m, nil
}

func (m *addModel) handleFieldEditUpdate(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Editing a field
	if m.editingField == fieldType && m.showSuggestions {
		// Handle autocomplete navigation
		if model, cmd := m.handleAutocompleteNav(msg); model != nil || cmd != nil {
			return model, cmd
		}
	}

	// Handle default editing mode
	return m.handleDefaultEdit(msg)
}

func (m *addModel) saveVocab() (tea.Model, tea.Cmd) {
	if m.term == "" {
		m.err = fmt.Errorf("term cannot be empty")
		return m, nil
	}
	if m.definition == "" {
		m.err = fmt.Errorf("definition cannot be empty")
		return m, nil
	}
	if m.meaningType == "" {
		m.meaningType = vocab.TypeNoun // Default type
	}

	// Create a copy of the library to modify
	// This ensures we don't corrupt in-memory state if save fails
	libraryCopy := &vocab.Library{
		Vocabs:   make(map[string]*vocab.Vocab, len(m.library.Vocabs)),
		Metadata: m.library.Metadata,
	}
	for k, v := range m.library.Vocabs {
		libraryCopy.Vocabs[k] = v
	}

	// Create new meaning
	newMeaning := vocab.Meaning{
		Type:       m.meaningType,
		Context:    m.context,
		Definition: m.definition,
		Examples:   m.examples,
	}

	if m.isAppendMode && m.existingVocab != nil {
		// Append meaning to existing vocab
		vocabCopy := *m.existingVocab
		newMeaning.ID = vocabCopy.GetNextMeaningID()
		vocabCopy.Meanings = append(vocabCopy.Meanings, newMeaning)
		libraryCopy.AddVocab(&vocabCopy)
	} else if m.isEditMode && m.originalMeaning != nil && m.existingVocab != nil {
		// Update existing meaning
		vocabCopy := *m.existingVocab
		for i := range vocabCopy.Meanings {
			if vocabCopy.Meanings[i].ID == m.originalMeaning.ID {
				vocabCopy.Meanings[i] = newMeaning
				vocabCopy.Meanings[i].ID = m.originalMeaning.ID // Preserve ID
				break
			}
		}
		// If term changed, delete old and add new
		originalKey := vocab.NormalizeTerm(m.existingVocab.Term)
		newKey := vocab.NormalizeTerm(m.term)
		if originalKey != newKey {
			delete(libraryCopy.Vocabs, originalKey)
			vocabCopy.Term = m.term
		}
		libraryCopy.AddVocab(&vocabCopy)
	} else {
		// Create new vocab with single meaning
		newMeaning.ID = 1
		vocabToSave := &vocab.Vocab{
			ID:        vocab.NormalizeTerm(m.term),
			Term:      m.term,
			Meanings:  []vocab.Meaning{newMeaning},
			Language:  m.language,
			Tags:      []string{},
			CreatedAt: time.Now(),
		}
		libraryCopy.AddVocab(vocabToSave)
	}

	// Return command to save asynchronously
	return m, m.saveVocabCmd(libraryCopy)
}
