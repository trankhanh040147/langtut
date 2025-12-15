package vocab

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

// WordInfoGenerator is kept for backward compatibility
// Deprecated: Use MeaningInfoGenerator instead
type WordInfoGenerator interface {
	GenerateWordInfo(ctx context.Context, word, language string) (string, []string, error)
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
	wordInfoClient   WordInfoGenerator // For backward compatibility
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
}

type meaningInfoGeneratedMsg struct {
	meaning vocab.Meaning
	err     error
}

type wordInfoGeneratedMsg struct {
	meaning  string
	examples []string
	err      error
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
	}
}

func (m *addModel) SetAPIClient(client MeaningInfoGenerator) {
	m.apiClient = client
}

// SetWordInfoClient is for backward compatibility
func (m *addModel) SetWordInfoClient(client WordInfoGenerator) {
	m.wordInfoClient = client
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

	// Backward compatibility: use old API if new one not available
	if m.term != "" && m.definition == "" && m.wordInfoClient != nil && m.apiClient == nil && !m.isEditMode {
		m.isGenerating = true
		return m.generateWordInfo()
	}

	return nil
}

func (m *addModel) generateMeaningInfo() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		if m.apiClient == nil {
			return meaningInfoGeneratedMsg{err: fmt.Errorf("API client not set")}
		}
		ctx := m.ctx
		if ctx == nil {
			ctx = context.Background()
		}
		meaning, err := m.apiClient.GenerateMeaningInfo(ctx, m.term, m.context, m.language)
		return meaningInfoGeneratedMsg{
			meaning: meaning,
			err:     err,
		}
	})
}

func (m *addModel) generateWordInfo() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		if m.wordInfoClient == nil {
			return wordInfoGeneratedMsg{err: fmt.Errorf("API client not set")}
		}
		ctx := m.ctx
		if ctx == nil {
			ctx = context.Background()
		}
		meaning, examples, err := m.wordInfoClient.GenerateWordInfo(ctx, m.term, m.language)
		return wordInfoGeneratedMsg{
			meaning:  meaning,
			examples: examples,
			err:      err,
		}
	})
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

	case wordInfoGeneratedMsg:
		m.isGenerating = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.definition = msg.meaning
			m.examples = msg.examples
			// Default type for backward compatibility
			if m.meaningType == "" {
				m.meaningType = vocab.TypeNoun
			}
		}
		return m, nil

	case tea.KeyMsg:
		if m.editingField >= 0 {
			// Editing a field
			if m.editingField == fieldType && m.showSuggestions {
				// Handle autocomplete navigation
				switch msg.String() {
				case constants.KeyDown:
					if m.selectedSuggestion < len(m.typeSuggestions)-1 {
						m.selectedSuggestion++
					}
					return m, nil
				case constants.KeyUp:
					if m.selectedSuggestion > 0 {
						m.selectedSuggestion--
					}
					return m, nil
				case constants.KeyEnter, "tab":
					// Select suggestion
					if m.selectedSuggestion >= 0 && m.selectedSuggestion < len(m.typeSuggestions) {
						m.editBuffer = m.typeSuggestions[m.selectedSuggestion]
						m.showSuggestions = false
						m.selectedSuggestion = -1
					}
					// Save and advance
					m.saveCurrentField()
					m.editingField = -1
					m.editBuffer = ""
					m.currentField++
					if m.currentField >= fieldCount {
						m.currentField = fieldTerm
					}
					return m, nil
				case constants.KeyEsc:
					m.showSuggestions = false
					m.selectedSuggestion = -1
					m.editingField = -1
					m.editBuffer = ""
					return m, nil
				}
			}

			switch msg.String() {
			case constants.KeyEnter:
				// Check for duplicate term if term field was saved
				wasTermField := m.editingField == fieldTerm
				// Save current field and advance to next
				m.saveCurrentField()
				m.editingField = -1
				m.editBuffer = ""
				m.showSuggestions = false
				// Check for duplicate after saving term
				if wasTermField {
					m.checkDuplicateTerm()
				}
				// Advance to next field
				m.currentField++
				if m.currentField >= fieldCount {
					m.currentField = fieldTerm
				}
			case constants.KeyEsc:
				// Cancel editing
				m.editingField = -1
				m.editBuffer = ""
				m.showSuggestions = false
				m.selectedSuggestion = -1
			case "backspace":
				if len(m.editBuffer) > 0 {
					m.editBuffer = m.editBuffer[:len(m.editBuffer)-1]
					// Update autocomplete for type field
					if m.editingField == fieldType {
						m.updateTypeSuggestions()
					}
				}
			case "tab":
				// Save and move to next
				m.saveCurrentField()
				m.editingField = -1
				m.editBuffer = ""
				m.showSuggestions = false
				m.currentField++
				if m.currentField >= fieldCount {
					m.currentField = fieldTerm
				}
			default:
				if len(msg.Runes) > 0 {
					m.editBuffer += string(msg.Runes)
					// Update autocomplete for type field
					if m.editingField == fieldType {
						m.updateTypeSuggestions()
					}
				}
			}
			return m, nil
		}

		switch msg.String() {
		case constants.KeyHelp:
			m.SetShowHelp(!m.ShowHelp())
			return m, nil

		case constants.KeyCtrlC, constants.KeyQuit:
			m.done = true
			m.saved = false
			if m.cancel != nil {
				m.cancel()
			}
			return m, tea.Quit

		case constants.KeyEsc:
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

		case constants.KeyDown:
			m.currentField++
			if m.currentField >= fieldCount {
				m.currentField = fieldTerm
			}
			return m, nil

		case constants.KeyUp:
			m.currentField--
			if m.currentField < 0 {
				m.currentField = fieldSave
			}
			return m, nil

		case constants.KeyEnter:
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

		case constants.KeyCtrlS:
			// Dedicated save key - save vocab immediately from any field
			return m.saveVocab()

		case "e":
			// Quick edit current field
			m.editingField = m.currentField
			m.loadCurrentFieldToBuffer()
			if m.currentField == fieldType {
				m.updateTypeSuggestions()
			}
			return m, nil

		case "tab":
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

func (m *addModel) loadCurrentFieldToBuffer() {
	switch m.currentField {
	case fieldTerm:
		m.editBuffer = m.term
	case fieldType:
		m.editBuffer = m.meaningType
	case fieldContext:
		m.editBuffer = m.context
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

func (m *addModel) saveCurrentField() {
	switch m.editingField {
	case fieldTerm:
		m.term = m.editBuffer
	case fieldType:
		m.meaningType = m.editBuffer
	case fieldContext:
		m.context = m.editBuffer
	case fieldDefinition:
		m.definition = m.editBuffer
	case fieldExample1:
		// First example
		if len(m.examples) == 0 {
			m.examples = []string{m.editBuffer}
		} else {
			m.examples[0] = m.editBuffer
		}
	case fieldExample2:
		// Second example
		for len(m.examples) <= 1 {
			m.examples = append(m.examples, "")
		}
		m.examples[1] = m.editBuffer
	case fieldExample3:
		// Third example
		for len(m.examples) <= 2 {
			m.examples = append(m.examples, "")
		}
		m.examples[2] = m.editBuffer
	default:
		// Ignore invalid field
	}
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

	// Attempt to save the copy
	if err := vocab.Save(libraryCopy); err != nil {
		m.err = err
		// Original m.library remains untouched and consistent
		return m, nil
	}

	// On success, replace the model's library with the successfully saved version
	m.library = libraryCopy
	m.done = true
	m.saved = true
	if m.cancel != nil {
		m.cancel()
	}
	return m, tea.Quit
}

func (m *addModel) View() string {
	if m.ShowHelp() {
		return ui.RenderHelp(m.Width(), m.Height())
	}

	width := 70
	height := 20

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1, 2).
		Width(width)

	var lines []string

	title := "Add Meaning"
	if m.isEditMode {
		title = "Edit Meaning"
	} else if m.isAppendMode {
		title = "Append Meaning"
	}
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	lines = append(lines, titleStyle.Render(title))
	repeatCount := width - 4
	if repeatCount < 0 {
		repeatCount = 0
	}
	lines = append(lines, strings.Repeat("─", repeatCount))

	// Show existing meanings if in append mode
	if m.isAppendMode && len(m.existingMeanings) > 0 {
		lines = append(lines, "")
		existingLabel := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6")).Render("Existing Meanings:")
		lines = append(lines, existingLabel)
		for _, meaning := range m.existingMeanings {
			meaningText := fmt.Sprintf("[%s] %s", meaning.Type, meaning.Definition)
			if meaning.Context != "" {
				contextBadge := lipgloss.NewStyle().
					Foreground(lipgloss.Color("8")).
					Render(fmt.Sprintf(" (%s)", meaning.Context))
				meaningText += contextBadge
			}
			lines = append(lines, "  "+meaningText)
		}
		lines = append(lines, "")
		lines = append(lines, strings.Repeat("─", repeatCount))
	}

	if m.isGenerating {
		lines = append(lines, "")
		lines = append(lines, "Generating meaning and examples...")
	} else {
		// Term field
		termLabel := "Term:"
		if m.currentField == fieldTerm {
			termLabel = "▶ Term:"
		}
		termStyle := lipgloss.NewStyle()
		if m.isAppendMode {
			termStyle = termStyle.Foreground(lipgloss.Color("8"))
		}
		if m.editingField == fieldTerm && !m.isAppendMode {
			lines = append(lines, termLabel+" "+m.editBuffer+"█")
		} else {
			termValue := m.term
			if termValue == "" {
				termValue = "(enter term)"
			}
			if m.isAppendMode {
				termValue += " 🔒"
			}
			lines = append(lines, termStyle.Render(termLabel+" "+termValue))
		}

		// Type field
		typeLabel := "Type:"
		if m.currentField == fieldType {
			typeLabel = "▶ Type:"
		}
		if m.editingField == fieldType {
			typeLine := typeLabel + " " + m.editBuffer + "█"
			lines = append(lines, typeLine)
			// Show autocomplete suggestions
			if m.showSuggestions && len(m.typeSuggestions) > 0 {
				for i, suggestion := range m.typeSuggestions {
					prefix := "  "
					if i == m.selectedSuggestion {
						prefix = "▶ "
					}
					suggestionStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
					lines = append(lines, suggestionStyle.Render(prefix+suggestion))
				}
			}
		} else {
			typeValue := m.meaningType
			if typeValue == "" {
				typeValue = "(enter type)"
			}
			lines = append(lines, typeLabel+" "+typeValue)
		}

		// Context field
		contextLabel := "Context:"
		if m.currentField == fieldContext {
			contextLabel = "▶ Context:"
		}
		if m.editingField == fieldContext {
			lines = append(lines, contextLabel+" "+m.editBuffer+"█")
		} else {
			contextValue := m.context
			if contextValue == "" {
				contextValue = "(enter context - tag or sentence)"
			}
			lines = append(lines, contextLabel+" "+contextValue)
		}

		// Definition field
		definitionLabel := "Definition:"
		if m.currentField == fieldDefinition {
			definitionLabel = "▶ Definition:"
		}
		if m.editingField == fieldDefinition {
			lines = append(lines, definitionLabel)
			lines = append(lines, "  "+m.editBuffer+"█")
		} else {
			definitionValue := m.definition
			if definitionValue == "" {
				definitionValue = "(enter definition)"
			}
			lines = append(lines, definitionLabel+" "+definitionValue)
		}

		// Examples
		lines = append(lines, "")
		lines = append(lines, "Examples:")
		exampleFields := []int{fieldExample1, fieldExample2, fieldExample3}
		for i, fieldIdx := range exampleFields {
			exLabel := fmt.Sprintf("  %d.", i+1)
			if m.currentField == fieldIdx {
				exLabel = "▶ " + exLabel
			}
			if m.editingField == fieldIdx {
				lines = append(lines, exLabel+" "+m.editBuffer+"█")
			} else {
				exValue := ""
				if i < len(m.examples) {
					exValue = m.examples[i]
				}
				if exValue == "" {
					exValue = "(enter example)"
				}
				lines = append(lines, exLabel+" "+exValue)
			}
		}

		// Save button
		lines = append(lines, "")
		saveLabel := "[Save]"
		if m.currentField == fieldSave {
			saveLabel = "▶ [Save]"
		}
		lines = append(lines, saveLabel)
	}

	if m.saved {
		lines = append(lines, "")
		successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
		lines = append(lines, successStyle.Render("✓ Saved!"))
	}

	if m.err != nil {
		lines = append(lines, "")
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
		lines = append(lines, errorStyle.Render("Error: "+m.err.Error()))
	}

	lines = append(lines, "")
	helpText := "Enter: Edit field | Ctrl+S: Save | Esc: Cancel | Tab: Next field | j/k: Navigate"
	lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render(helpText))

	box := boxStyle.Render(strings.Join(lines, "\n"))

	// Center the box
	topPadding := (m.Height() - height) / 2
	if topPadding < 0 {
		topPadding = 0
	}
	leftPadding := (m.Width() - width) / 2
	if leftPadding < 0 {
		leftPadding = 0
	}

	topRepeat := topPadding
	if topRepeat < 0 {
		topRepeat = 0
	}
	leftRepeat := leftPadding
	if leftRepeat < 0 {
		leftRepeat = 0
	}
	result := strings.Repeat("\n", topRepeat)
	for _, line := range strings.Split(box, "\n") {
		result += strings.Repeat(" ", leftRepeat) + line + "\n"
	}

	return result
}
