package constants

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines key bindings for UI components
type KeyMap struct {
	Tab       key.Binding
	Backspace key.Binding
	Edit      key.Binding
	Confirm   key.Binding
	Cancel    key.Binding
	Down      key.Binding
	Up        key.Binding
	Enter     key.Binding
	Esc       key.Binding
	CtrlS     key.Binding
	Help      key.Binding
	CtrlC     key.Binding
	Quit      key.Binding
	Top       key.Binding
	Bottom    key.Binding
	Search    key.Binding
	Add       key.Binding
	Delete    key.Binding
}

// DefaultKeyMap returns the default keymap
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Tab: key.NewBinding(
			key.WithKeys("tab"),
		),
		Backspace: key.NewBinding(
			key.WithKeys("backspace"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
		),
		Confirm: key.NewBinding(
			key.WithKeys("y", "Y"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("n", "N"),
		),
		Down: key.NewBinding(
			key.WithKeys("j"),
		),
		Up: key.NewBinding(
			key.WithKeys("k"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
		),
		Esc: key.NewBinding(
			key.WithKeys("esc"),
		),
		CtrlS: key.NewBinding(
			key.WithKeys("ctrl+s"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
		),
		CtrlC: key.NewBinding(
			key.WithKeys("ctrl+c"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q"),
		),
		Top: key.NewBinding(
			key.WithKeys("g"),
		),
		Bottom: key.NewBinding(
			key.WithKeys("G"),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
		),
		Add: key.NewBinding(
			key.WithKeys("a"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
		),
	}
}

// HelpText contains the help overlay text
const HelpText = `Keyboard Shortcuts:
  ?          Show/hide this help
  j          Move down
  k          Move up
  g          Go to top
  G          Go to bottom
  /          Search
  q          Quit
  Esc        Cancel/Back
  Ctrl+C     Force quit

Vocab List:
  e          Edit selected word
  d          Delete selected word
  a          Add new word
  Enter      View full details

Add/Edit Word Modal:
  Enter      Edit field / Save field and next
  Ctrl+S     Save word immediately
  Tab        Save field and move to next
  Esc        Cancel editing / Close modal
  j/k        Navigate between fields
`
