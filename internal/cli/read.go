package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newReadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Reading practice commands",
		Long:  "Read blogs, articles, and watch videos with AI assistance.",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(os.Stderr, "Reading module coming soon. Use --help for available commands.\n")
			return cmd.Help()
		},
	}

	return cmd
}
