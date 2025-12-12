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

type addModel struct {
	ui.BaseModel
	word         string
	meaning      string
	examples     []string
	currentField int
	editingField int
	editBuffer   string
	library      *vocab.Library
	apiClient    interface {
		GenerateWordInfo(ctx context.Context, word, language string) (string, []string, error)
	}
	language     string
	isGenerating bool
	isEditMode   bool
	originalWord   *vocab.Word
	done         bool
	saved        bool
	err          error
}

type wordInfoGeneratedMsg struct {
	meaning  string
	examples []string
	err      error
}

func NewAddModel(word, meaning string, examples []string, lib *vocab.Library) *addModel {
	return &addModel{
		BaseModel:    ui.BaseModel{},
		word:         word,
		meaning:      meaning,
		examples:     examples,
		currentField: 0,
		editingField:  -1,
		library:      lib,
		language:     "English", // Default, can be from config
		isEditMode:   false, // Only set to true when editing existing word (originalWord != nil)
		originalWord: nil,
	}
}

func (m *addModel) SetAPIClient(client interface {
	GenerateWordInfo(ctx context.Context, word, language string) (string, []string, error)
}) {
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
		ctx := context.Background()
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
				maxField := 5 // Word(0), Meaning(1), Examples(2-4), Save(5)
				if m.currentField > maxField {
					m.currentField = 0
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
				maxField := 5 // Word(0), Meaning(1), Examples(2-4), Save(5)
				if m.currentField > maxField {
					m.currentField = 0
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
			return m, tea.Quit

		case constants.KeyEsc:
			m.done = true
			m.saved = false
			return m, nil

		case constants.KeyDown:
			m.currentField++
			maxField := 5 // Word(0), Meaning(1), Examples(2-4), Save(5)
			if m.currentField > maxField {
				m.currentField = 0
			}
			return m, nil

		case constants.KeyUp:
			m.currentField--
			if m.currentField < 0 {
				m.currentField = 5 // Word(0), Meaning(1), Examples(2-4), Save(5)
			}
			return m, nil

		case constants.KeyEnter:
			// Start editing current field or save
			if m.currentField == 5 {
				// Save button (always at index 5: Word(0), Meaning(1), Examples(2-4), Save(5))
				return m.saveWord()
			} else {
				// If word is valid and not currently editing, allow saving from any field
				// This allows immediate save without navigating to Save button
				if m.word != "" && m.editingField < 0 {
					return m.saveWord()
				}
				// Otherwise, enter edit mode for current field
				m.editingField = m.currentField
				m.loadCurrentFieldToBuffer()
			}
			return m, nil

		case "e":
			// Quick edit current field
			m.editingField = m.currentField
			m.loadCurrentFieldToBuffer()
			return m, nil

		case "tab":
			// Navigate to next field when not editing
			m.currentField++
			maxField := 5 // Word(0), Meaning(1), Examples(2-4), Save(5)
			if m.currentField > maxField {
				m.currentField = 0
			}
			return m, nil
		}
	}

	return m, nil
}

func (m *addModel) loadCurrentFieldToBuffer() {
	switch m.currentField {
	case 0:
		m.editBuffer = m.word
	case 1:
		m.editBuffer = m.meaning
	case 2:
		// First example
		if len(m.examples) > 0 {
			m.editBuffer = m.examples[0]
		} else {
			m.editBuffer = ""
		}
	default:
		// Additional examples
		idx := m.currentField - 2
		if idx < len(m.examples) {
			m.editBuffer = m.examples[idx]
		} else {
			m.editBuffer = ""
		}
	}
}

func (m *addModel) saveCurrentField() {
	switch m.editingField {
	case 0:
		m.word = m.editBuffer
	case 1:
		m.meaning = m.editBuffer
	case 2:
		// First example
		if len(m.examples) == 0 {
			m.examples = []string{m.editBuffer}
		} else {
			m.examples[0] = m.editBuffer
		}
	default:
		// Additional examples
		idx := m.editingField - 2
		for len(m.examples) <= idx {
			m.examples = append(m.examples, "")
		}
		m.examples[idx] = m.editBuffer
	}
}

func (m *addModel) saveWord() (tea.Model, tea.Cmd) {
	if m.word == "" {
		m.err = fmt.Errorf("word cannot be empty")
		return m, nil
	}

	// Create or update word
	var word *vocab.Word
	if m.isEditMode && m.originalWord != nil {
		word = m.originalWord
		word.Word = m.word
		word.Meaning = m.meaning
		word.Examples = m.examples
	} else {
		word = &vocab.Word{
			Word:      m.word,
			Meaning:   m.meaning,
			Language:  m.language,
			Examples:  m.examples,
			Tags:      []string{},
			CreatedAt: time.Now(),
		}
	}

	// If editing, delete old word first
	if m.isEditMode && m.originalWord != nil {
		m.library.DeleteWord(m.originalWord.Word)
	}

	m.library.AddWord(word)

	if err := vocab.Save(m.library); err != nil {
		m.err = err
		return m, nil
	}

	m.done = true
	m.saved = true
	return m, nil
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
	lines = append(lines, strings.Repeat("─", width-4))

	if m.isGenerating {
		lines = append(lines, "")
		lines = append(lines, "Generating meaning and examples...")
	} else {
		// Word field
		wordLabel := "Word:"
		if m.currentField == 0 {
			wordLabel = "▶ Word:"
		}
		if m.editingField == 0 {
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
		if m.currentField == 1 {
			meaningLabel = "▶ Meaning:"
		}
		if m.editingField == 1 {
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
		for i := 0; i < 3; i++ {
			fieldIdx := 2 + i
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
		if m.currentField == 5 {
			// Save button always at index 5: Word(0), Meaning(1), Examples(2-4), Save(5)
			saveLabel = "▶ [Save]"
		}
		lines = append(lines, saveLabel)
	}

	if m.err != nil {
		lines = append(lines, "")
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
		lines = append(lines, errorStyle.Render("Error: "+m.err.Error()))
	}

	lines = append(lines, "")
	helpText := "Enter: Edit field / Save | Esc: Cancel | Tab: Next field | j/k: Navigate"
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

	result := strings.Repeat("\n", topPadding)
	for _, line := range strings.Split(box, "\n") {
		result += strings.Repeat(" ", leftPadding) + line + "\n"
	}

	return result
}

