package vocab

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/trankhanh040147/langtut/internal/ui"
)

func (m *addModel) View() string {
	if m.ShowHelp() {
		return ui.RenderHelp(m.Width(), m.Height())
	}

	width := 70

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
	ruleWidth := width - 4
	if ruleWidth < 0 {
		ruleWidth = 0
	}
	lines = append(lines, renderHorizontalRule(ruleWidth))

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
		lines = append(lines, renderHorizontalRule(ruleWidth))
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
				contextValue = "(optional - tag or sentence)"
			}
			lines = append(lines, contextLabel+" "+contextValue)
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

	// Center the box using lipgloss.Place
	return lipgloss.Place(m.Width(), m.Height(), lipgloss.Center, lipgloss.Center, box)
}
