package cmd

import (
	"github.com/spf13/cobra"
)

// hostsCmd represents the hosts command
var hostsCmd = &cobra.Command{
	Use:   "hosts",
	Short: "Manage the host list",
	Long: `Manage the host list for pscan.

Add hosts with 'add' subcommand.
List hosts with 'list' subcommand.
Delete hosts with 'delete' subcommand.`,
}

func init() {
	rootCmd.AddCommand(hostsCmd)
}
