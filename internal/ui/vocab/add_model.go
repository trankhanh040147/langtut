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
	fieldWord = iota
	fieldMeaning
	fieldExample1
	fieldExample2
	fieldExample3
	fieldSave
	fieldCount // Total number of fields
)

// WordInfoGenerator defines the interface for generating word information
type WordInfoGenerator interface {
	GenerateWordInfo(ctx context.Context, word, language string) (string, []string, error)
}

type addModel struct {
	ui.BaseModel
	word         string
	meaning      string
	examples     []string
	currentField int
	editingField int
	editBuffer   string
	library      *vocab.Library
	apiClient    WordInfoGenerator
	language     string
	isGenerating bool
	isEditMode   bool
	originalWord *vocab.Word
	done         bool
	saved        bool
	err          error
	ctx          context.Context
	cancel       context.CancelFunc
}

type wordInfoGeneratedMsg struct {
	meaning  string
	examples []string
	err      error
}

func NewAddModel(word, meaning string, examples []string, lib *vocab.Library) *addModel {
	ctx, cancel := context.WithCancel(context.Background())
	return &addModel{
		BaseModel:    ui.BaseModel{},
		word:         word,
		meaning:      meaning,
		examples:     examples,
		currentField: fieldWord,
		editingField: -1,
		library:      lib,
		language:     "English", // Default, can be from config
		isEditMode:   false,     // Only set to true when editing existing word (originalWord != nil)
		originalWord: nil,
		ctx:          ctx,
		cancel:       cancel,
	}
}

func (m *addModel) SetAPIClient(client WordInfoGenerator) {
	m.apiClient = client
}

func (m *addModel) SetLanguage(lang string) {
	m.language = lang
}

func (m *addModel) SetWord(word string) {
	m.word = word
}

func (m *addModel) SetMeaning(meaning string) {
	m.meaning = meaning
}

func (m *addModel) SetExamples(examples []string) {
	m.examples = examples
}

func (m *addModel) Saved() bool {
	return m.saved
}

func (m *addModel) Init() tea.Cmd {
	if m.word != "" && m.meaning == "" && m.apiClient != nil && !m.isEditMode {
		// Auto-generate for new words
		m.isGenerating = true
		return m.generateWordInfo()
	}
	return nil
}

func (m *addModel) generateWordInfo() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		if m.apiClient == nil {
			return wordInfoGeneratedMsg{err: fmt.Errorf("API client not set")}
		}
		ctx := m.ctx
		if ctx == nil {
			ctx = context.Background()
		}
		meaning, examples, err := m.apiClient.GenerateWordInfo(ctx, m.word, m.language)
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

	case wordInfoGeneratedMsg:
		m.isGenerating = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.meaning = msg.meaning
			m.examples = msg.examples
		}
		return m, nil

	case tea.KeyMsg:
		if m.editingField >= 0 {
			// Editing a field
			switch msg.String() {
			case constants.KeyEnter:
				// Save current field and advance to next
				m.saveCurrentField()
				m.editingField = -1
				m.editBuffer = ""
				// Advance to next field
				m.currentField++
				if m.currentField >= fieldCount {
					m.currentField = fieldWord
				}
			case constants.KeyEsc:
				// Cancel editing
				m.editingField = -1
				m.editBuffer = ""
			case "backspace":
				if len(m.editBuffer) > 0 {
					m.editBuffer = m.editBuffer[:len(m.editBuffer)-1]
				}
			case "tab":
				// Save and move to next
				m.saveCurrentField()
				m.editingField = -1
				m.editBuffer = ""
				m.currentField++
				if m.currentField >= fieldCount {
					m.currentField = fieldWord
				}
			default:
				if len(msg.Runes) > 0 {
					m.editBuffer += string(msg.Runes)
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
				m.currentField = fieldWord
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
				return m.saveWord()
			} else {
				// Enter edit mode for current field
				m.editingField = m.currentField
				m.loadCurrentFieldToBuffer()
			}
			return m, nil

		case constants.KeyCtrlS:
			// Dedicated save key - save word immediately from any field
			return m.saveWord()

		case "e":
			// Quick edit current field
			m.editingField = m.currentField
			m.loadCurrentFieldToBuffer()
			return m, nil

		case "tab":
			// Navigate to next field when not editing
			m.currentField++
			if m.currentField >= fieldCount {
				m.currentField = fieldWord
			}
			return m, nil
		}
	}

	return m, nil
}

func (m *addModel) loadCurrentFieldToBuffer() {
	switch m.currentField {
	case fieldWord:
		m.editBuffer = m.word
	case fieldMeaning:
		m.editBuffer = m.meaning
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
	case fieldWord:
		m.word = m.editBuffer
	case fieldMeaning:
		m.meaning = m.editBuffer
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

func (m *addModel) saveWord() (tea.Model, tea.Cmd) {
	if m.word == "" {
		m.err = fmt.Errorf("word cannot be empty")
		return m, nil
	}

	// Create a copy of the library to modify
	// This ensures we don't corrupt in-memory state if save fails
	libraryCopy := &vocab.Library{
		Words:    make(map[string]*vocab.Word, len(m.library.Words)),
		Metadata: m.library.Metadata,
	}
	for k, v := range m.library.Words {
		libraryCopy.Words[k] = v
	}

	// If editing, delete old word from copy if key changed
	if m.isEditMode && m.originalWord != nil {
		originalKey := vocab.NormalizeWord(m.originalWord.Word)
		newKey := vocab.NormalizeWord(m.word)
		if originalKey != newKey {
			delete(libraryCopy.Words, originalKey)
		}
	}

	// Create or update word object
	var wordToSave *vocab.Word
	if m.isEditMode && m.originalWord != nil {
		// Update existing word
		wordToSave = m.originalWord
		wordToSave.Word = m.word
		wordToSave.Meaning = m.meaning
		wordToSave.Examples = m.examples
	} else {
		// Create new word
		wordToSave = &vocab.Word{
			Word:      m.word,
			Meaning:   m.meaning,
			Language:  m.language,
			Examples:  m.examples,
			Tags:      []string{},
			CreatedAt: time.Now(),
		}
	}

	libraryCopy.AddWord(wordToSave)

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

	title := "Add Word"
	if m.isEditMode {
		title = "Edit Word"
	}
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	lines = append(lines, titleStyle.Render(title))
	repeatCount := width - 4
	if repeatCount < 0 {
		repeatCount = 0
	}
	lines = append(lines, strings.Repeat("─", repeatCount))

	if m.isGenerating {
		lines = append(lines, "")
		lines = append(lines, "Generating meaning and examples...")
	} else {
		// Word field
		wordLabel := "Word:"
		if m.currentField == fieldWord {
			wordLabel = "▶ Word:"
		}
		if m.editingField == fieldWord {
			lines = append(lines, wordLabel+" "+m.editBuffer+"█")
		} else {
			wordValue := m.word
			if wordValue == "" {
				wordValue = "(enter word)"
			}
			lines = append(lines, wordLabel+" "+wordValue)
		}

		// Meaning field
		meaningLabel := "Meaning:"
		if m.currentField == fieldMeaning {
			meaningLabel = "▶ Meaning:"
		}
		if m.editingField == fieldMeaning {
			lines = append(lines, meaningLabel)
			lines = append(lines, "  "+m.editBuffer+"█")
		} else {
			meaningValue := m.meaning
			if meaningValue == "" {
				meaningValue = "(enter meaning)"
			}
			lines = append(lines, meaningLabel+" "+meaningValue)
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
