package vocab

import (
	"fmt"
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
}

type wordGeneratedMsg struct {
	word    *vocab.Word
	success bool
	err     error
}

type wordSavedMsg struct {
	success bool
	err     error
}

type wordDeletedMsg struct {
	success bool
	err     error
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
	}
}

func (m *listModel) Init() tea.Cmd {
	return nil
}

func (m *listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.showAddModal && m.addModel != nil {
		var cmd tea.Cmd
		var updated tea.Model
		updated, cmd = m.addModel.Update(msg)
		if addModel, ok := updated.(*addModel); ok {
			m.addModel = addModel
		}
		if m.addModel.done {
			if m.addModel.saved {
				// Vocab was saved, refresh from in-memory library
				m.refreshVocabsFromLibrary()
			}
			m.showAddModal = false
			m.addModel = nil
			// Don't propagate tea.Quit from addModel - just close the modal
			return m, nil
		}
		return m, cmd
	}

	if m.showEditModal && m.addModel != nil {
		var cmd tea.Cmd
		var updated tea.Model
		updated, cmd = m.addModel.Update(msg)
		if addModel, ok := updated.(*addModel); ok {
			m.addModel = addModel
		}
		if m.addModel.done {
			if m.addModel.saved {
				// Vocab was saved, refresh from in-memory library
				m.refreshVocabsFromLibrary()
				// Maintain selection (applySearch already handles bounds)
			}
			m.showEditModal = false
			m.addModel = nil
			m.editVocab = nil
			m.editMeaning = nil
			// Don't propagate tea.Quit from addModel - just close the modal
			return m, nil
		}
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.BaseModel.SetSize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		if m.showDeleteConfirm {
			switch msg.String() {
			case "y", "Y":
				// Delete vocab
				if m.hasValidSelection() {
					v := m.filteredVocabs[m.selectedIdx]
					m.library.DeleteVocab(v.Term)
					if err := vocab.Save(m.library); err != nil {
						m.err = err
					} else {
						// Reload
						m.reloadLibrary()
						// Selection bounds already handled by applySearch()
					}
				}
				m.showDeleteConfirm = false
			case "n", "N", constants.KeyEsc:
				m.showDeleteConfirm = false
			}
			return m, nil
		}

		if m.isSearching {
			switch msg.String() {
			case constants.KeyEnter:
				m.isSearching = false
				m.searchQuery = ""
				m.applySearch()
			case constants.KeyEsc:
				m.isSearching = false
				m.searchQuery = ""
				m.applySearch()
			case "backspace":
				if len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
					m.applySearch()
				}
			default:
				if len(msg.Runes) > 0 {
					m.searchQuery += string(msg.Runes)
					m.applySearch()
				}
			}
			return m, nil
		}

		switch msg.String() {
		case constants.KeyHelp:
			m.SetShowHelp(!m.ShowHelp())
			return m, nil

		case constants.KeyEsc:
			// Close help overlay if shown
			if m.ShowHelp() {
				m.SetShowHelp(false)
				return m, nil
			}

		case constants.KeyCtrlC, constants.KeyQuit:
			return m, tea.Quit

		case constants.KeyDown:
			if m.selectedIdx < len(m.filteredVocabs)-1 {
				m.selectedIdx++
			}
			return m, nil

		case constants.KeyUp:
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
			return m, nil

		case constants.KeyTop:
			m.selectedIdx = 0
			return m, nil

		case constants.KeyBottom:
			if len(m.filteredVocabs) > 0 {
				m.selectedIdx = len(m.filteredVocabs) - 1
			}
			return m, nil

		case constants.KeySearch:
			m.isSearching = true
			m.searchQuery = ""
			return m, nil

		case constants.KeyAdd:
			m.showAddModal = true
			m.addModel = NewAddModel("", "", []string{}, m.library)
			return m, m.addModel.Init()

		case constants.KeyEdit:
			if m.hasValidSelection() {
				m.showEditModal = true
				m.editVocab = m.filteredVocabs[m.selectedIdx]
				// Edit first meaning by default (can enhance later to select meaning)
				if len(m.editVocab.Meanings) > 0 {
					m.editMeaning = &m.editVocab.Meanings[0]
					m.addModel = NewAddModel(
						m.editVocab.Term,
						m.editMeaning.Definition,
						m.editMeaning.Examples,
						m.library,
					)
					m.addModel.meaningType = m.editMeaning.Type
					m.addModel.context = m.editMeaning.Context
					m.addModel.originalMeaning = m.editMeaning
					m.addModel.existingVocab = m.editVocab
					m.addModel.isEditMode = true
					return m, m.addModel.Init()
				}
			}
			return m, nil

		case constants.KeyDelete:
			if m.hasValidSelection() {
				m.showDeleteConfirm = true
			}
			return m, nil

		case constants.KeyEnter:
			// View full details (already shown in right pane)
			return m, nil
		}
	}

	return m, nil
}

func (m *listModel) applySearch() {
	if m.searchQuery == "" {
		m.filteredVocabs = m.vocabs
	} else {
		query := strings.ToLower(m.searchQuery)
		m.filteredVocabs = []*vocab.Vocab{}
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

func (m *listModel) View() string {
	if m.showAddModal && m.addModel != nil {
		return m.addModel.View()
	}

	if m.showEditModal && m.addModel != nil {
		return m.addModel.View()
	}

	if m.showDeleteConfirm {
		return m.renderDeleteConfirm()
	}

	if m.ShowHelp() {
		return ui.RenderHelp(m.Width(), m.Height())
	}

	// Split pane layout
	width := m.safeWidth()
	leftWidth := width * 40 / 100
	if leftWidth < 1 {
		leftWidth = 1
	}
	rightWidth := width - leftWidth - 1
	if rightWidth < 1 {
		rightWidth = 1
	}

	leftPane := m.renderVocabList(leftWidth, m.Height())
	rightPane := m.renderVocabDetails(rightWidth, m.Height())

	// Combine panes
	content := lipgloss.JoinHorizontal(lipgloss.Left, leftPane, "│", rightPane)

	if m.err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")).
			Padding(0, 1)
		content += "\n\n" + errorStyle.Render("Error: "+m.err.Error())
	}

	return content
}

func (m *listModel) renderVocabList(width, height int) string {
	var lines []string

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(0, 1)
	header := headerStyle.Render("Terms")
	if m.isSearching {
		header += fmt.Sprintf(" (Search: %s)", m.searchQuery)
	} else {
		header += fmt.Sprintf(" (%d)", len(m.filteredVocabs))
	}
	lines = append(lines, header)
	repeatCount := width
	if repeatCount < 0 {
		repeatCount = 0
	}
	lines = append(lines, strings.Repeat("─", repeatCount))

	// Vocab list
	listHeight := height - 3
	startIdx := 0
	if m.selectedIdx >= listHeight {
		startIdx = m.selectedIdx - listHeight + 1
	}

	for i := startIdx; i < len(m.filteredVocabs) && i < startIdx+listHeight; i++ {
		v := m.filteredVocabs[i]
		termText := v.Term
		// Add meaning count badge
		if len(v.Meanings) > 1 {
			countBadge := lipgloss.NewStyle().
				Foreground(lipgloss.Color("8")).
				Render(fmt.Sprintf(" (%d)", len(v.Meanings)))
			termText += countBadge
		}

		// Highlight search matches
		if m.searchQuery != "" {
			query := strings.ToLower(m.searchQuery)
			termLower := strings.ToLower(v.Term)
			if strings.Contains(termLower, query) {
				idx := strings.Index(termLower, query)
				before := v.Term[:idx]
				match := v.Term[idx : idx+len(m.searchQuery)]
				after := v.Term[idx+len(m.searchQuery):]
				highlightStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("3"))
				termText = before + highlightStyle.Render(match) + after
				if len(v.Meanings) > 1 {
					countBadge := lipgloss.NewStyle().
						Foreground(lipgloss.Color("8")).
						Render(fmt.Sprintf(" (%d)", len(v.Meanings)))
					termText += countBadge
				}
			}
		}

		if i == m.selectedIdx {
			selectedStyle := lipgloss.NewStyle().
				Background(lipgloss.Color("4")).
				Foreground(lipgloss.Color("15")).
				Padding(0, 1).
				Width(width - 2)
			termText = selectedStyle.Render("▶ " + termText)
		} else {
			normalStyle := lipgloss.NewStyle().
				Padding(0, 1).
				Width(width - 2)
			termText = normalStyle.Render("  " + termText)
		}

		lines = append(lines, termText)
	}

	return lipgloss.NewStyle().Width(width).Height(height).Render(strings.Join(lines, "\n"))
}

func (m *listModel) renderVocabDetails(width, height int) string {
	var lines []string

	if m.selectedIdx < 0 || m.selectedIdx >= len(m.filteredVocabs) {
		lines = append(lines, "No term selected")
		return lipgloss.NewStyle().Width(width).Height(height).Render(strings.Join(lines, "\n"))
	}

	v := m.filteredVocabs[m.selectedIdx]

	// Term
	termStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(0, 1)
	lines = append(lines, termStyle.Render(v.Term))
	repeatCount := width
	if repeatCount < 0 {
		repeatCount = 0
	}
	lines = append(lines, strings.Repeat("─", repeatCount))

	// Meanings
	for idx, meaning := range v.Meanings {
		if idx > 0 {
			lines = append(lines, "")
			lines = append(lines, strings.Repeat("─", repeatCount))
		}

		// Meaning header with type
		meaningHeader := fmt.Sprintf("[%s] %s", meaning.Type, meaning.Definition)
		if meaning.Context != "" {
			contextBadge := lipgloss.NewStyle().
				Foreground(lipgloss.Color("6")).
				Background(lipgloss.Color("0")).
				Padding(0, 1).
				Render(meaning.Context)
			meaningHeader += " " + contextBadge
		}
		meaningHeaderStyle := lipgloss.NewStyle().Bold(true).Padding(0, 1)
		lines = append(lines, meaningHeaderStyle.Render(meaningHeader))

		// Examples for this meaning
		if len(meaning.Examples) > 0 {
			lines = append(lines, "")
			examplesLabel := lipgloss.NewStyle().Bold(true).Render("Examples:")
			lines = append(lines, examplesLabel)
			for i, ex := range meaning.Examples {
				exStyle := lipgloss.NewStyle().Padding(0, 2).Foreground(lipgloss.Color("8")).Width(width - 4)
				lines = append(lines, exStyle.Render(fmt.Sprintf("%d. %s", i+1, ex)))
			}
		}
	}

	// Tags
	if len(v.Tags) > 0 {
		lines = append(lines, "")
		tagsLabel := lipgloss.NewStyle().Bold(true).Render("Tags:")
		lines = append(lines, tagsLabel)
		tagStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("6")).
			Background(lipgloss.Color("0")).
			Padding(0, 1).
			Margin(0, 1)
		tags := []string{}
		for _, tag := range v.Tags {
			tags = append(tags, tagStyle.Render(tag))
		}
		lines = append(lines, strings.Join(tags, " "))
	}

	// Created at
	lines = append(lines, "")
	dateStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Padding(0, 1)
	lines = append(lines, dateStyle.Render("Added: "+v.CreatedAt.Format("2006-01-02")))

	return lipgloss.NewStyle().Width(width).Height(height).Render(strings.Join(lines, "\n"))
}

func (m *listModel) renderDeleteConfirm() string {
	width := 50
	height := 5

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("1")).
		Padding(1, 2).
		Width(width)

	var content string
	if m.hasValidSelection() {
		content = "Delete term: " + m.filteredVocabs[m.selectedIdx].Term + "?\n\n"
	} else {
		content = "Delete term: (invalid selection)?\n\n"
	}
	content += "Press 'y' to confirm, 'n' to cancel"

	box := boxStyle.Render(content)

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
