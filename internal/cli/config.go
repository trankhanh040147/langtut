package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/trankhanh040147/langtut/internal/config"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
		Long:  "Manage langtut configuration settings.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			fmt.Fprintf(os.Stdout, "Configuration:\n")
			fmt.Fprintf(os.Stdout, "  API Key: %s\n", maskAPIKey(cfg.APIKey))
			fmt.Fprintf(os.Stdout, "  Target Language: %s\n", cfg.TargetLanguage)
			if len(cfg.Presets) > 0 {
				fmt.Fprintf(os.Stdout, "  Presets: %v\n", cfg.Presets)
			}

			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available configuration keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			fmt.Fprintf(os.Stdout, "Available configuration keys:\n\n")

			// api_key
			fmt.Fprintf(os.Stdout, "  api_key          Gemini API key\n")
			fmt.Fprintf(os.Stdout, "                    Current value: %s\n\n", maskAPIKey(cfg.APIKey))

			// target_language
			targetLang := cfg.TargetLanguage
			if targetLang == "" {
				targetLang = "(not set)"
			}
			fmt.Fprintf(os.Stdout, "  target_language  Target language for learning\n")
			fmt.Fprintf(os.Stdout, "                    Current value: %s\n", targetLang)

			return nil
		},
	}

	showCmd := &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			fmt.Fprintf(os.Stdout, "API Key: %s\n", maskAPIKey(cfg.APIKey))
			fmt.Fprintf(os.Stdout, "Target Language: %s\n", cfg.TargetLanguage)
			return nil
		},
	}

	setCmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set configuration value",
		Long: `Set a configuration value.

Valid keys:
  api_key          Gemini API key
  target_language  Target language for learning

Examples:
  langtut config set api_key YOUR_API_KEY
  langtut config set target_language Spanish`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			key := args[0]
			value := args[1]

			switch key {
			case "api_key":
				cfg.APIKey = value
			case "target_language":
				cfg.TargetLanguage = value
			default:
				return fmt.Errorf("unknown config key: %s\n\nValid keys: api_key, target_language", key)
			}

			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Fprintf(os.Stderr, "Configuration updated\n")
			return nil
		},
	}

	cmd.AddCommand(listCmd)
	cmd.AddCommand(showCmd)
	cmd.AddCommand(setCmd)

	return cmd
}

func maskAPIKey(key string) string {
	if key == "" {
		return "(not set)"
	}
	if len(key) <= 8 {
		return "***"
	}
	return key[:4] + "..." + key[len(key)-4:]
}
