package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/trankhanh040147/langtut/internal/api"
	"github.com/trankhanh040147/langtut/internal/config"
	"github.com/trankhanh040147/langtut/internal/ui"
	vocabui "github.com/trankhanh040147/langtut/internal/ui/vocab"
	"github.com/trankhanh040147/langtut/internal/vocab"
)

func newVocabCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vocab",
		Short: "Vocabulary learning commands",
		Long:  "Practice vocabulary with guessing, typing, and phrase learning.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Default to list view
			return runVocabList()
		},
	}

	// Add subcommand
	addCmd := &cobra.Command{
		Use:   "add <word> [word2] [word3]...",
		Short: "Add word(s) to vocabulary library",
		Long:  "Add one or more words to your vocabulary library. AI will generate meaning and examples.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVocabAdd(args)
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List vocabulary words",
		Long:  "Open interactive TUI to browse, search, edit, and delete vocabulary words.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVocabList()
		},
	}

	cmd.AddCommand(addCmd)
	cmd.AddCommand(listCmd)

	return cmd
}

func runVocabList() error {
	if !ui.ShouldShowTUI() {
		return fmt.Errorf("TUI mode requires an interactive terminal")
	}

	lib, err := vocab.Load()
	if err != nil {
		return fmt.Errorf("failed to load vocabulary library: %w", err)
	}

	model := vocabui.NewListModel(lib)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run TUI: %w", err)
	}

	return nil
}

func runVocabAdd(words []string) error {
	if !ui.ShouldShowTUI() {
		return fmt.Errorf("TUI mode requires an interactive terminal")
	}

	// Load config for API key
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.APIKey == "" {
		if err := config.PromptAPIKey(cfg); err != nil {
			return fmt.Errorf("API key required: %w", err)
		}
		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("failed to reload config: %w", err)
		}
	}

	// Create API client
	apiClient, err := api.NewClient(cfg.APIKey)
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}
	defer apiClient.Close()

	// Load library
	lib, err := vocab.Load()
	if err != nil {
		return fmt.Errorf("failed to load vocabulary library: %w", err)
	}

	// Get language from config
	language := cfg.TargetLanguage
	if language == "" {
		language = "English"
	}

	// Pass all words directly to batch model
	// TUI will handle duplicate detection and append mode
	batchModel := vocabui.NewBatchAddModel(words, lib, apiClient, language)
	p := tea.NewProgram(batchModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run TUI: %w", err)
	}

	// Update library reference from batch model
	lib = batchModel.Library()

	return nil
}
