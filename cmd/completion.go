package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate completion script for your shell",
	Long: `To load your completions run

For bash:
$ source <(pscan completion)

For zsh:
$ source <(pscan completion)

To load completions automatically on login add
"source <(pscan completion)" to your ~/.bashrc or ~/.zshrc.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return completeAction(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}

func completeAction(out *os.File) error {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "bash"
	}

	switch shell {
	case "/bin/bash":
		return rootCmd.GenBashCompletion(out)
	case "/bin/zsh":
		return rootCmd.GenZshCompletion(out)
	default:
		return rootCmd.GenBashCompletion(out)
	}
}
