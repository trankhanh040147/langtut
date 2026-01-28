package agent

import (
	"context"
	_ "embed"
	"strings"

	"github.com/trankhanh040147/langtut/internal/agent/prompt"
	"github.com/trankhanh040147/langtut/internal/config"
)

//go:embed templates/coder.md.tpl
var coderPromptTmpl []byte

//go:embed templates/task.md.tpl
var taskPromptTmpl []byte

//go:embed templates/initialize.md.tpl
var initializePromptTmpl []byte

//go:embed templates/writing-tutor.tpl
var writingTutorPromptTmpl []byte

func coderPrompt(opts ...prompt.Option) (*prompt.Prompt, error) {
	systemPrompt, err := prompt.NewPrompt("coder", string(coderPromptTmpl), opts...)
	if err != nil {
		return nil, err
	}
	return systemPrompt, nil
}

func taskPrompt(opts ...prompt.Option) (*prompt.Prompt, error) {
	systemPrompt, err := prompt.NewPrompt("task", string(taskPromptTmpl), opts...)
	if err != nil {
		return nil, err
	}
	return systemPrompt, nil
}

func InitializePrompt(cfg config.Config) (string, error) {
	systemPrompt, err := prompt.NewPrompt("initialize", string(initializePromptTmpl))
	if err != nil {
		return "", err
	}
	return systemPrompt.Build(context.Background(), "", "", cfg)
}

func WritingTutorPrompt(topic string) (string, error) {
	tmpl := strings.ReplaceAll(string(writingTutorPromptTmpl), "{{.Topic}}", topic)
	return tmpl, nil
}

// CoderPrompt returns a new coder prompt instance.
func CoderPrompt(opts ...prompt.Option) (*prompt.Prompt, error) {
	return coderPrompt(opts...)
}
