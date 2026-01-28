package dialog

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/trankhanh040147/langtut/internal/ui/common"
)

const TopicInputID = "topic_input"

type TopicInput struct {
	com   *common.Common
	help  help.Model
	input textinput.Model

	keyMap struct {
		Submit key.Binding
		Close  key.Binding
	}
}

var _ Dialog = (*TopicInput)(nil)

func NewTopicInput(com *common.Common) *TopicInput {
	t := &TopicInput{com: com}

	help := help.New()
	help.Styles = com.Styles.DialogHelpStyles()
	t.help = help

	t.input = textinput.New()
	t.input.SetVirtualCursor(false)
	t.input.Placeholder = "Enter topic (leave empty for random)"
	t.input.SetStyles(com.Styles.TextInput)
	t.input.Focus()

	t.keyMap.Submit = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "submit"),
	)
	t.keyMap.Close = CloseKey

	return t
}

func (t *TopicInput) ID() string {
	return TopicInputID
}

func (t *TopicInput) HandleMsg(msg tea.Msg) Action {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, t.keyMap.Close):
			return ActionClose{}
		case key.Matches(msg, t.keyMap.Submit):
			topic := t.input.Value()
			return ActionTopicInput{Topic: topic}
		default:
			var cmd tea.Cmd
			t.input, cmd = t.input.Update(msg)
			return ActionCmd{cmd}
		}
	}
	return nil
}

func (t *TopicInput) Cursor() *tea.Cursor {
	return InputCursor(t.com.Styles, t.input.Cursor())
}

func (t *TopicInput) Draw(scr uv.Screen, area uv.Rectangle) *tea.Cursor {
	s := t.com.Styles
	width := max(0, min(defaultDialogMaxWidth, area.Dx()))
	innerWidth := width - s.Dialog.View.GetHorizontalFrameSize()

	t.input.SetWidth(max(0, innerWidth-s.Dialog.InputPrompt.GetHorizontalFrameSize()-1))
	t.help.SetWidth(innerWidth)

	rc := NewRenderContext(s, width)
	rc.Title = "Writing Topic"

	prompt := s.Dialog.PrimaryText.Render("What topic would you like to practice?")
	rc.AddPart(prompt)

	inputView := s.Dialog.InputPrompt.Render(t.input.View())
	rc.AddPart(inputView)

	rc.Help = t.help.View(t)

	view := rc.Render()
	cur := t.Cursor()
	DrawCenterCursor(scr, area, view, cur)
	return cur
}

func (t *TopicInput) ShortHelp() []key.Binding {
	return []key.Binding{
		t.keyMap.Submit,
		t.keyMap.Close,
	}
}

func (t *TopicInput) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			t.keyMap.Submit,
			t.keyMap.Close,
		},
	}
}
