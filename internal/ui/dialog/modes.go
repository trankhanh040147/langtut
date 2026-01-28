package dialog

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/sahilm/fuzzy"
	"github.com/trankhanh040147/langtut/internal/ui/common"
	"github.com/trankhanh040147/langtut/internal/ui/list"
	"github.com/trankhanh040147/langtut/internal/ui/styles"
)

const ModesID = "modes"

type SessionMode string

const (
	SessionModeCoder        SessionMode = "coder"
	SessionModeWritingTutor SessionMode = "writing_tutor"
)

type Modes struct {
	com  *common.Common
	help help.Model
	list *list.FilterableList

	keyMap struct {
		Select   key.Binding
		Next     key.Binding
		Previous key.Binding
		UpDown   key.Binding
		Close    key.Binding
	}
}

var _ Dialog = (*Modes)(nil)

func NewModes(com *common.Common) *Modes {
	m := &Modes{com: com}

	help := help.New()
	help.Styles = com.Styles.DialogHelpStyles()
	m.help = help

	items := []list.FilterableItem{
		newModeItem(com, SessionModeWritingTutor, "Writing Tutor", "IELTS writing practice with feedback"),
		newModeItem(com, SessionModeCoder, "Coder", "AI coding assistant (default mode)"),
	}

	m.list = list.NewFilterableList(items...)
	m.list.Focus()
	m.list.SetSelected(0)

	m.keyMap.Select = key.NewBinding(
		key.WithKeys("enter", "ctrl+y"),
		key.WithHelp("enter", "choose"),
	)
	m.keyMap.Next = key.NewBinding(
		key.WithKeys("down", "ctrl+n"),
		key.WithHelp("↓", "next"),
	)
	m.keyMap.Previous = key.NewBinding(
		key.WithKeys("up", "ctrl+p"),
		key.WithHelp("↑", "previous"),
	)
	m.keyMap.UpDown = key.NewBinding(
		key.WithKeys("up", "down"),
		key.WithHelp("↑↓", "choose"),
	)
	m.keyMap.Close = CloseKey

	return m
}

func (m *Modes) ID() string {
	return ModesID
}

func (m *Modes) HandleMsg(msg tea.Msg) Action {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keyMap.Close):
			return ActionClose{}
		case key.Matches(msg, m.keyMap.Previous):
			m.list.Focus()
			if m.list.IsSelectedFirst() {
				m.list.SelectLast()
				m.list.ScrollToBottom()
				break
			}
			m.list.SelectPrev()
			m.list.ScrollToSelected()
		case key.Matches(msg, m.keyMap.Next):
			m.list.Focus()
			if m.list.IsSelectedLast() {
				m.list.SelectFirst()
				m.list.ScrollToTop()
				break
			}
			m.list.SelectNext()
			m.list.ScrollToSelected()
		case key.Matches(msg, m.keyMap.Select):
			if item := m.list.SelectedItem(); item != nil {
				modeItem := item.(*modeItem)
				return ActionSelectMode{Mode: modeItem.mode}
			}
		}
	}
	return nil
}

func (m *Modes) Cursor() *tea.Cursor {
	return nil
}

func (m *Modes) Draw(scr uv.Screen, area uv.Rectangle) *tea.Cursor {
	t := m.com.Styles
	width := max(0, min(defaultDialogMaxWidth, area.Dx()))
	height := max(0, min(defaultDialogHeight, area.Dy()))
	innerWidth := width - t.Dialog.View.GetHorizontalFrameSize() - 2
	heightOffset := t.Dialog.Title.GetVerticalFrameSize() + titleContentHeight +
		t.Dialog.HelpView.GetVerticalFrameSize() +
		t.Dialog.View.GetVerticalFrameSize()

	m.list.SetSize(innerWidth, height-heightOffset)
	m.help.SetWidth(innerWidth)

	rc := NewRenderContext(t, width)
	rc.Title = "Select Mode"

	listView := t.Dialog.List.Height(m.list.Height()).Render(m.list.Render())
	rc.AddPart(listView)
	rc.Help = m.help.View(m)

	view := rc.Render()
	DrawCenterCursor(scr, area, view, nil)
	return nil
}

func (m *Modes) ShortHelp() []key.Binding {
	return []key.Binding{
		m.keyMap.UpDown,
		m.keyMap.Select,
		m.keyMap.Close,
	}
}

func (m *Modes) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			m.keyMap.UpDown,
			m.keyMap.Select,
			m.keyMap.Close,
		},
	}
}

type modeItem struct {
	mode        SessionMode
	name        string
	description string
	t           *styles.Styles
	m           fuzzy.Match
	focused     bool
}

func newModeItem(com *common.Common, mode SessionMode, name, description string) *modeItem {
	return &modeItem{
		mode:        mode,
		name:        name,
		description: description,
		t:           com.Styles,
	}
}

func (i *modeItem) ID() string {
	return string(i.mode)
}

func (i *modeItem) Filter() string {
	return i.name
}

func (i *modeItem) SetMatch(m fuzzy.Match) {
	i.m = m
}

func (i *modeItem) Focus() {
	i.focused = true
}

func (i *modeItem) Blur() {
	i.focused = false
}

func (i *modeItem) Render(width int) string {
	style := i.t.Dialog.NormalItem
	if i.focused {
		style = i.t.Dialog.SelectedItem
	}
	title := i.name
	if i.description != "" {
		title += " - " + i.description
	}
	return style.Width(width).Render(title)
}
