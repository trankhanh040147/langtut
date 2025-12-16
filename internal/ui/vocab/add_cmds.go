package vocab

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/trankhanh040147/langtut/internal/vocab"
)

func (m *addModel) generateMeaningInfo() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		if m.apiClient == nil {
			return meaningInfoGeneratedMsg{err: fmt.Errorf("API client not set")}
		}
		ctx := m.ctx
		if ctx == nil {
			ctx = context.Background()
		}
		meaning, err := m.apiClient.GenerateMeaningInfo(ctx, m.term, m.context, m.language)
		return meaningInfoGeneratedMsg{
			meaning: meaning,
			err:     err,
		}
	})
}

func (m *addModel) saveVocabCmd(libraryCopy *vocab.Library) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		if err := vocab.Save(libraryCopy); err != nil {
			return vocabSavedMsg{err: err}
		}
		return vocabSavedMsg{library: libraryCopy, err: nil}
	})
}
