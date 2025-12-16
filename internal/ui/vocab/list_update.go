package vocab

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trankhanh040147/langtut/internal/api"
	"github.com/trankhanh040147/langtut/internal/config"
)

func (m *listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.showAddModal && m.addModel != nil {
		var cmd tea.Cmd
		var updated tea.Model
		updated, cmd = m.addModel.Update(msg)
		if addModel, ok := updated.(*addModel); ok {
			m.addModel = addModel
		}
		if m.addModel.done {
			if m.addModel.saved {
				// Vocab was saved, refresh from in-memory library
				m.refreshVocabsFromLibrary()
			}
			m.showAddModal = false
			m.addModel = nil
			// Don't propagate tea.Quit from addModel - just close the modal
			return m, nil
		}
		return m, cmd
	}

	if m.showEditModal && m.addModel != nil {
		var cmd tea.Cmd
		var updated tea.Model
		updated, cmd = m.addModel.Update(msg)
		if addModel, ok := updated.(*addModel); ok {
			m.addModel = addModel
		}
		if m.addModel.done {
			if m.addModel.saved {
				// Vocab was saved, refresh from in-memory library
				m.refreshVocabsFromLibrary()
				// Maintain selection (applySearch already handles bounds)
			}
			m.showEditModal = false
			m.addModel = nil
			m.editVocab = nil
			m.editMeaning = nil
			// Don't propagate tea.Quit from addModel - just close the modal
			return m, nil
		}
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.BaseModel.SetSize(msg.Width, msg.Height)
		return m, nil

	case vocabDeletedMsg:
		m.showDeleteConfirm = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			// Reload
			m.reloadLibrary()
			// Selection bounds already handled by applySearch()
		}
		return m, nil

	case tea.KeyMsg:
		if m.showDeleteConfirm {
			if key.Matches(msg, m.keys.Confirm) {
				// Delete vocab
				if m.hasValidSelection() {
					v := m.filteredVocabs[m.selectedIdx]
					m.library.DeleteVocab(v.Term)
					return m, m.deleteVocabCmd()
				}
				m.showDeleteConfirm = false
			} else if key.Matches(msg, m.keys.Cancel) || key.Matches(msg, m.keys.Esc) {
				m.showDeleteConfirm = false
			}
			return m, nil
		}

		if m.isSearching {
			if key.Matches(msg, m.keys.Enter) {
				m.isSearching = false
				m.searchQuery = ""
				m.applySearch()
				return m, nil
			}

			if key.Matches(msg, m.keys.Esc) {
				m.isSearching = false
				m.searchQuery = ""
				m.applySearch()
				return m, nil
			}

			if key.Matches(msg, m.keys.Backspace) {
				if len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
					m.applySearch()
				}
				return m, nil
			}

			if len(msg.Runes) > 0 {
				m.searchQuery += string(msg.Runes)
				m.applySearch()
			}
			return m, nil
		}

		if key.Matches(msg, m.keys.Help) {
			m.SetShowHelp(!m.ShowHelp())
			return m, nil
		}

		if key.Matches(msg, m.keys.Esc) {
			// Close help overlay if shown
			if m.ShowHelp() {
				m.SetShowHelp(false)
				return m, nil
			}
		}

		if key.Matches(msg, m.keys.CtrlC) || key.Matches(msg, m.keys.Quit) {
			return m, tea.Quit
		}

		if key.Matches(msg, m.keys.Down) {
			if m.selectedIdx < len(m.filteredVocabs)-1 {
				m.selectedIdx++
			}
			return m, nil
		}

		if key.Matches(msg, m.keys.Up) {
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
			return m, nil
		}

		if key.Matches(msg, m.keys.Top) {
			m.selectedIdx = 0
			return m, nil
		}

		if key.Matches(msg, m.keys.Bottom) {
			if len(m.filteredVocabs) > 0 {
				m.selectedIdx = len(m.filteredVocabs) - 1
			}
			return m, nil
		}

		if key.Matches(msg, m.keys.Search) {
			m.isSearching = true
			m.searchQuery = ""
			return m, nil
		}

		if key.Matches(msg, m.keys.Add) {
			// Load config and create API client (same as CLI mode)
			cfg, err := config.Load()
			if err != nil {
				m.err = fmt.Errorf("failed to load config: %w", err)
				return m, nil
			}

			if cfg.APIKey == "" {
				m.err = fmt.Errorf("API key not set. Please configure it first")
				return m, nil
			}

			apiClient, err := api.NewClient(cfg.APIKey)
			if err != nil {
				m.err = fmt.Errorf("failed to create API client: %w", err)
				return m, nil
			}

			language := cfg.TargetLanguage
			if language == "" {
				language = "English"
			}

			// Use shared initialization function
			m.showAddModal = true
			m.addModel = NewAddModelWithConfig("", m.library, apiClient, language)
			return m, m.addModel.Init()
		}

		if key.Matches(msg, m.keys.Edit) {
			if m.hasValidSelection() {
				// Load config for API client (may be needed for regeneration)
				cfg, err := config.Load()
				var apiClient MeaningInfoGenerator
				if err == nil && cfg.APIKey != "" {
					if client, err := api.NewClient(cfg.APIKey); err == nil {
						apiClient = client
					}
				}

				language := "English"
				if cfg != nil && cfg.TargetLanguage != "" {
					language = cfg.TargetLanguage
				}

				m.showEditModal = true
				m.editVocab = m.filteredVocabs[m.selectedIdx]
				// Edit first meaning by default (can enhance later to select meaning)
				if len(m.editVocab.Meanings) > 0 {
					m.editMeaning = &m.editVocab.Meanings[0]
					m.addModel = NewAddModelWithConfig(
						m.editVocab.Term,
						m.library,
						apiClient,
						language,
					)
					m.addModel.definition = m.editMeaning.Definition
					m.addModel.examples = m.editMeaning.Examples
					m.addModel.meaningType = m.editMeaning.Type
					m.addModel.context = m.editMeaning.Context
					m.addModel.originalMeaning = m.editMeaning
					m.addModel.existingVocab = m.editVocab
					m.addModel.isEditMode = true
					return m, m.addModel.Init()
				}
			}
			return m, nil
		}

		if key.Matches(msg, m.keys.Delete) {
			if m.hasValidSelection() {
				m.showDeleteConfirm = true
			}
			return m, nil
		}

		if key.Matches(msg, m.keys.Enter) {
			// View full details (already shown in right pane)
			return m, nil
		}
	}

	return m, nil
}
