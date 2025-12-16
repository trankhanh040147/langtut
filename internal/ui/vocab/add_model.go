package vocab

import (
	"context"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trankhanh040147/langtut/internal/constants"
	"github.com/trankhanh040147/langtut/internal/ui"
	"github.com/trankhanh040147/langtut/internal/vocab"
)

const (
	fieldTerm = iota
	fieldType
	fieldContext
	fieldDefinition
	fieldExample1
	fieldExample2
	fieldExample3
	fieldSave
	fieldCount // Total number of fields
)

// MeaningInfoGenerator defines the interface for generating meaning information
type MeaningInfoGenerator interface {
	GenerateMeaningInfo(ctx context.Context, term, context, language string) (vocab.Meaning, error)
}

type addModel struct {
	ui.BaseModel
	term             string
	meaningType      string
	context          string
	definition       string
	examples         []string
	currentField     int
	editingField     int
	editBuffer       string
	library          *vocab.Library
	apiClient        MeaningInfoGenerator
	language         string
	isGenerating     bool
	isEditMode       bool
	isAppendMode     bool
	existingVocab    *vocab.Vocab
	existingMeanings []vocab.Meaning
	originalMeaning  *vocab.Meaning
	// Autocomplete state
	typeSuggestions    []string
	selectedSuggestion int
	showSuggestions    bool
	done               bool
	saved              bool
	err                error
	ctx                context.Context
	cancel             context.CancelFunc
	keys               constants.KeyMap
}

type meaningInfoGeneratedMsg struct {
	meaning vocab.Meaning
	err     error
}

type vocabSavedMsg struct {
	library *vocab.Library
	err     error
}

func NewAddModel(term, meaning string, examples []string, lib *vocab.Library) *addModel {
	ctx, cancel := context.WithCancel(context.Background())
	return &addModel{
		BaseModel:          ui.BaseModel{},
		term:               term,
		definition:         meaning,
		examples:           examples,
		currentField:       fieldTerm,
		editingField:       -1,
		library:            lib,
		language:           "English", // Default, can be from config
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

func (m *addModel) SetAPIClient(client MeaningInfoGenerator) {
	m.apiClient = client
}

func (m *addModel) SetLanguage(lang string) {
	m.language = lang
}

func (m *addModel) SetTerm(term string) {
	m.term = term
}

func (m *addModel) SetDefinition(definition string) {
	m.definition = definition
}

func (m *addModel) SetExamples(examples []string) {
	m.examples = examples
}

func (m *addModel) Saved() bool {
	return m.saved
}

func (m *addModel) Library() *vocab.Library {
	return m.library
}

func (m *addModel) Init() tea.Cmd {
	// Check for duplicate term
	if m.term != "" && !m.isEditMode {
		if existingVocab, exists := m.library.GetVocab(m.term); exists {
			m.isAppendMode = true
			m.existingVocab = existingVocab
			m.existingMeanings = existingVocab.Meanings
		}
	}

	// Auto-generate for new meanings
	if m.term != "" && m.definition == "" && m.apiClient != nil && !m.isEditMode {
		m.isGenerating = true
		return m.generateMeaningInfo()
	}

	return nil
}

func (m *addModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.BaseModel.SetSize(msg.Width, msg.Height)
		return m, nil

	case meaningInfoGeneratedMsg:
		m.isGenerating = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.meaningType = msg.meaning.Type
			m.definition = msg.meaning.Definition
			m.context = msg.meaning.Context
			m.examples = msg.meaning.Examples
		}
		return m, nil

	case vocabSavedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		// On success, replace the model's library with the successfully saved version
		m.library = msg.library
		m.done = true
		m.saved = true
		if m.cancel != nil {
			m.cancel()
		}
		return m, tea.Quit

	case tea.KeyMsg:
		if m.editingField >= 0 {
			return m.handleFieldEditUpdate(msg)
		}

		if key.Matches(msg, m.keys.Help) {
			m.SetShowHelp(!m.ShowHelp())
			return m, nil
		}

		if key.Matches(msg, m.keys.CtrlC) || key.Matches(msg, m.keys.Quit) {
			m.done = true
			m.saved = false
			if m.cancel != nil {
				m.cancel()
			}
			return m, tea.Quit
		}

		if key.Matches(msg, m.keys.Esc) {
			// Close help overlay if shown
			if m.ShowHelp() {
				m.SetShowHelp(false)
				return m, nil
			}
			// Otherwise close modal
			m.done = true
			m.saved = false
			if m.cancel != nil {
				m.cancel()
			}
			return m, nil
		}

		if key.Matches(msg, m.keys.Down) {
			m.currentField++
			if m.currentField >= fieldCount {
				m.currentField = fieldTerm
			}
			return m, nil
		}

		if key.Matches(msg, m.keys.Up) {
			m.currentField--
			if m.currentField < 0 {
				m.currentField = fieldSave
			}
			return m, nil
		}

		if key.Matches(msg, m.keys.Enter) {
			// Start editing current field or save current field and advance
			if m.currentField == fieldSave {
				// Save button
				return m.saveVocab()
			} else {
				// Enter edit mode for current field
				m.editingField = m.currentField
				m.loadCurrentFieldToBuffer()
				// Initialize autocomplete for type field
				if m.currentField == fieldType {
					m.updateTypeSuggestions()
				}
			}
			return m, nil
		}

		if key.Matches(msg, m.keys.CtrlS) {
			// Dedicated save key - save vocab immediately from any field
			return m.saveVocab()
		}

		if key.Matches(msg, m.keys.Edit) {
			// Quick edit current field
			m.editingField = m.currentField
			m.loadCurrentFieldToBuffer()
			if m.currentField == fieldType {
				m.updateTypeSuggestions()
			}
			return m, nil
		}

		if key.Matches(msg, m.keys.Tab) {
			// Navigate to next field when not editing
			m.currentField++
			if m.currentField >= fieldCount {
				m.currentField = fieldTerm
			}
			return m, nil
		}
	}

	return m, nil
}
