package ui

import (
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trankhanh040147/langtut/internal/constants"
)

// BaseModel provides base functionality for TUI models
type BaseModel struct {
	showHelp bool
	height   int
	width    int
	keys     constants.KeyMap
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
		// Initialize keys on first use (BaseModel keys are lazily initialized)
		// Check if keys are uninitialized by testing if Help binding matches nothing
		if len(m.keys.Help.Keys()) == 0 {
			m.keys = constants.DefaultKeyMap()
		}
		if key.Matches(msg, m.keys.Help) {
			m.showHelp = !m.showHelp
			return m, nil
		}
		if key.Matches(msg, m.keys.CtrlC) || key.Matches(msg, m.keys.Quit) {
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

// SetSize sets the terminal size
func (m *BaseModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}
