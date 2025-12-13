package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newWriteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "write",
		Short: "Writing practice commands",
		Long:  "Practice writing with fill-in phrases, sentence rewriting, and conversations.",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(os.Stderr, "Writing module coming soon. Use --help for available commands.\n")
			return cmd.Help()
		},
	}

	return cmd
}
