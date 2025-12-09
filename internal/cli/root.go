package cli

import (
	"os"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

var (
	// isInteractive indicates if we're in an interactive terminal
	isInteractive bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "langtut",
	Short: "A language learning CLI tool with AI tutoring",
	Long: `langtut is a language learning CLI tool that provides
interactive practice and AI-powered tutoring for vocabulary,
reading, and writing skills.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Check if stdout is a TTY
		isInteractive = isatty.IsTerminal(os.Stdout.Fd())
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

// IsInteractive returns whether we're in an interactive terminal
func IsInteractive() bool {
	return isInteractive
}

// InitRoot initializes the root command and adds subcommands
func InitRoot() {
	rootCmd.AddCommand(newConfigCmd())
	rootCmd.AddCommand(newVocabCmd())
	rootCmd.AddCommand(newReadCmd())
	rootCmd.AddCommand(newWriteCmd())
}
