package vocab

import (
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/trankhanh040147/langtut/internal/constants"
	"github.com/trankhanh040147/langtut/internal/ui"
	"github.com/trankhanh040147/langtut/internal/vocab"
)

type listModel struct {
	ui.BaseModel
	library           *vocab.Library
	vocabs            []*vocab.Vocab
	filteredVocabs    []*vocab.Vocab
	selectedIdx       int
	searchQuery       string
	isSearching       bool
	showAddModal      bool
	showEditModal     bool
	showDeleteConfirm bool
	editVocab         *vocab.Vocab
	editMeaning       *vocab.Meaning
	addModel          *addModel
	err               error
	keys              constants.KeyMap
}

type vocabDeletedMsg struct {
	err error
}

func NewListModel(lib *vocab.Library) *listModel {
	vocabs := lib.GetAllVocabs()
	sort.Slice(vocabs, func(i, j int) bool {
		return strings.ToLower(vocabs[i].Term) < strings.ToLower(vocabs[j].Term)
	})

	return &listModel{
		BaseModel:      ui.BaseModel{},
		library:        lib,
		vocabs:         vocabs,
		filteredVocabs: vocabs,
		selectedIdx:    0,
		keys:           constants.DefaultKeyMap(),
	}
}

func (m *listModel) Init() tea.Cmd {
	return nil
}

func (m *listModel) applySearch() {
	if m.searchQuery == "" {
		m.filteredVocabs = m.vocabs
	} else {
		query := strings.ToLower(m.searchQuery)
		m.filteredVocabs = make([]*vocab.Vocab, 0, len(m.vocabs))
		for _, v := range m.vocabs {
			// Search in term
			if strings.Contains(strings.ToLower(v.Term), query) {
				m.filteredVocabs = append(m.filteredVocabs, v)
				continue
			}
			// Search in all meanings
			for _, meaning := range v.Meanings {
				if strings.Contains(strings.ToLower(meaning.Definition), query) ||
					strings.Contains(strings.ToLower(meaning.Context), query) {
					m.filteredVocabs = append(m.filteredVocabs, v)
					break
				}
			}
		}
	}
	// Adjust selectedIdx to valid range
	if len(m.filteredVocabs) == 0 {
		m.selectedIdx = -1
	} else {
		if m.selectedIdx >= len(m.filteredVocabs) {
			m.selectedIdx = len(m.filteredVocabs) - 1
		}
		if m.selectedIdx < 0 {
			m.selectedIdx = 0
		}
	}
}

func (m *listModel) reloadLibrary() {
	lib, err := vocab.Load()
	if err != nil {
		m.err = err
		return
	}

	m.library = lib
	vocabs := lib.GetAllVocabs()
	sort.Slice(vocabs, func(i, j int) bool {
		return strings.ToLower(vocabs[i].Term) < strings.ToLower(vocabs[j].Term)
	})
	m.vocabs = vocabs
	m.applySearch()
}

func (m *listModel) refreshVocabsFromLibrary() {
	// Re-fetch vocabs from the in-memory library object
	vocabs := m.library.GetAllVocabs()
	sort.Slice(vocabs, func(i, j int) bool {
		return strings.ToLower(vocabs[i].Term) < strings.ToLower(vocabs[j].Term)
	})
	m.vocabs = vocabs
	m.applySearch() // This already handles selection bounds
}

func (m *listModel) deleteVocabCmd() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		if err := vocab.Save(m.library); err != nil {
			return vocabDeletedMsg{err: err}
		}
		return vocabDeletedMsg{err: nil}
	})
}

func (m *listModel) hasValidSelection() bool {
	return m.selectedIdx >= 0 && m.selectedIdx < len(m.filteredVocabs)
}

func (m *listModel) safeWidth() int {
	w := m.Width()
	if w <= 0 {
		return 80 // Default minimum width
	}
	return w
}

// renderHorizontalRule renders a horizontal rule using lipgloss instead of strings.Repeat
func renderHorizontalRule(width int) string {
	if width <= 0 {
		return ""
	}
	return lipgloss.NewStyle().
		Width(width).
		Border(lipgloss.NormalBorder()).
		BorderTop(true).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false).
		Render("")
}
