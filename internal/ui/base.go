package ui

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/trankhanh040147/langtut/internal/constants"
)

// BaseModel provides base functionality for TUI models
type BaseModel struct {
	showHelp bool
	height   int
	width    int
}

// Init initializes the base model
func (m *BaseModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the base model
func (m *BaseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case constants.KeyHelp:
			m.showHelp = !m.showHelp
			return m, nil
		case constants.KeyCtrlC, constants.KeyQuit:
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the base model (to be overridden)
func (m *BaseModel) View() string {
	return ""
}

// ShouldShowTUI checks if TUI should be shown
func ShouldShowTUI() bool {
	// Check NO_COLOR env
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	// Check if stdout is a TTY
	stat, err := os.Stdout.Stat()
	if err != nil {
		return false
	}

	return (stat.Mode() & os.ModeCharDevice) != 0
}

// GetHelpText returns the help text
func GetHelpText() string {
	return constants.HelpText
}

// ShowHelp returns whether help should be shown
func (m *BaseModel) ShowHelp() bool {
	return m.showHelp
}

// SetShowHelp sets the help visibility
func (m *BaseModel) SetShowHelp(show bool) {
	m.showHelp = show
}

// Height returns the terminal height
func (m *BaseModel) Height() int {
	return m.height
}

// Width returns the terminal width
func (m *BaseModel) Width() int {
	return m.width
}
