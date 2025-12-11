package constants

const (
	// KeyBindings
	KeyHelp   = "?"
	KeyDown   = "j"
	KeyUp     = "k"
	KeyTop    = "g"
	KeyBottom = "G"
	KeySearch = "/"
	KeyQuit   = "q"
	KeyEnter  = "enter"
	KeyEsc    = "esc"
	KeyCtrlC  = "ctrl+c"
	// Vocab-specific keys
	KeyEdit   = "e"
	KeyDelete = "d"
	KeyAdd    = "a"
)

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
`
