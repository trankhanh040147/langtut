package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/trankhanh040147/langtut/internal/constants"
)

var (
	helpStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1, 2).
		Width(60)
)

// RenderHelp renders the help overlay
func RenderHelp(width, height int) string {
	helpText := constants.HelpText

	// Wrap help text to fit width
	lines := strings.Split(helpText, "\n")
	wrapped := make([]string, 0, len(lines))
	for _, line := range lines {
		if len(line) > width-10 {
			// Simple word wrap
			words := strings.Fields(line)
			currentLine := ""
			for _, word := range words {
				if len(currentLine)+len(word)+1 > width-10 {
					if currentLine != "" {
						wrapped = append(wrapped, currentLine)
					}
					currentLine = word
				} else {
					if currentLine == "" {
						currentLine = word
					} else {
						currentLine += " " + word
					}
				}
			}
			if currentLine != "" {
				wrapped = append(wrapped, currentLine)
			}
		} else {
			wrapped = append(wrapped, line)
		}
	}

	helpText = strings.Join(wrapped, "\n")
	rendered := helpStyle.Render(helpText)

	// Center the help box using style width (60) instead of calculated width
	renderedLines := strings.Split(rendered, "\n")
	boxHeight := len(renderedLines)
	boxWidth := 60 // Use style width directly for consistent centering

	// Calculate padding
	topPadding := (height - boxHeight) / 2
	if topPadding < 0 {
		topPadding = 0
	}
	leftPadding := (width - boxWidth) / 2
	if leftPadding < 0 {
		leftPadding = 0
	}

	// Add vertical padding
	result := strings.Repeat("\n", topPadding)
	for _, line := range renderedLines {
		result += fmt.Sprintf("%s%s\n", strings.Repeat(" ", leftPadding), line)
	}

	return result
}

// GetKeyBindingDescription returns a formatted key binding description
func GetKeyBindingDescription(key, description string) string {
	keyStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	return fmt.Sprintf("%s  %s", keyStyle.Render(key), description)
}
