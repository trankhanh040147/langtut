package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newVocabCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vocab",
		Short: "Vocabulary learning commands",
		Long:  "Practice vocabulary with guessing, typing, and phrase learning.",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(os.Stderr, "Vocabulary module coming soon. Use --help for available commands.\n")
			return cmd.Help()
		},
	}

	return cmd
}
