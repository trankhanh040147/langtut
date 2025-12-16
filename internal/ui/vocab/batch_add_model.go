package vocab

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trankhanh040147/langtut/internal/ui"
	"github.com/trankhanh040147/langtut/internal/vocab"
)

type batchAddModel struct {
	ui.BaseModel
	words      []string
	currentIdx int
	addModel   *addModel
	library    *vocab.Library
	apiClient  MeaningInfoGenerator
	language   string
	done       bool
}

func NewBatchAddModel(words []string, lib *vocab.Library, apiClient MeaningInfoGenerator, language string) *batchAddModel {
	return &batchAddModel{
		BaseModel:  ui.BaseModel{},
		words:      words,
		currentIdx: 0,
		library:    lib,
		apiClient:  apiClient,
		language:   language,
		done:       false,
	}
}

func (m *batchAddModel) Init() tea.Cmd {
	if len(m.words) == 0 {
		m.done = true
		return tea.Quit
	}
	return m.startNextWord()
}

func (m *batchAddModel) startNextWord() tea.Cmd {
	if m.currentIdx >= len(m.words) {
		m.done = true
		return tea.Quit
	}

	term := m.words[m.currentIdx]
	m.addModel = NewAddModel(term, "", []string{}, m.library)
	m.addModel.SetAPIClient(m.apiClient)
	m.addModel.SetLanguage(m.language)
	return m.addModel.Init()
}

func (m *batchAddModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.BaseModel.SetSize(msg.Width, msg.Height)
		if m.addModel != nil {
			var cmd tea.Cmd
			_, cmd = m.addModel.Update(msg)
			return m, cmd
		}
		return m, nil

	case tea.KeyMsg:
		// Handle addModel's messages
		if m.addModel != nil {
			var cmd tea.Cmd
			var updated tea.Model
			updated, cmd = m.addModel.Update(msg)
			if addModel, ok := updated.(*addModel); ok {
				m.addModel = addModel
			}

			// Check if addModel is done
			if m.addModel.done {
				if m.addModel.saved {
					// Update library reference
					m.library = m.addModel.Library()
				}
				// Move to next word
				m.currentIdx++
				m.addModel = nil
				return m, m.startNextWord()
			}

			return m, cmd
		}
		return m, nil
	}

	// Forward other messages to addModel if it exists
	if m.addModel != nil {
		var cmd tea.Cmd
		_, cmd = m.addModel.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *batchAddModel) View() string {
	if m.addModel != nil {
		return m.addModel.View()
	}
	return ""
}

func (m *batchAddModel) Done() bool {
	return m.done
}

func (m *batchAddModel) Library() *vocab.Library {
	return m.library
}

