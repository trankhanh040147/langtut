package vocab

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/trankhanh040147/langtut/internal/ui"
)

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
	lines = append(lines, renderHorizontalRule(width))

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
	lines = append(lines, renderHorizontalRule(width))

	// Meanings
	for idx, meaning := range v.Meanings {
		if idx > 0 {
			lines = append(lines, "")
			lines = append(lines, renderHorizontalRule(width))
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

	// Center the box using lipgloss.Place
	return lipgloss.Place(m.Width(), m.Height(), lipgloss.Center, lipgloss.Center, box)
}

