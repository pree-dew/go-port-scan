package cmd

import (
	"fmt"

	"pscan/scan"

	"github.com/spf13/cobra"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Run a port scan on the hosts list",
	Long:  `Provide a list of port that you want to scan on list of host.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile, err := cmd.Flags().GetString("hosts-file")
		if err != nil {
			return err
		}

		ports, err := cmd.Flags().GetIntSlice("ports")
		if err != nil {
			return err
		}

		return scanAction(hostsFile, ports)
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().IntSliceP("ports", "p", []int{22, 80, 443}, "Ports to scan")
}

func scanAction(hostsFile string, ports []int) error {
	hl := &scan.HostsList{}

	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	results := scan.Run(hl, ports)

	for _, r := range results {
		if r.NotFound {
			fmt.Printf("Host %s not found\n", r.Host)
			continue
		}

		fmt.Printf("Host %s\n", r.Host)
		for _, p := range r.PortStates {
			fmt.Printf("Port %d is %s\n", p.Port, p.Open)
		}
	}

	return nil
}
