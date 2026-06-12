package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yourusername/ricebox/sandbox"
)

var sandboxCmd = &cobra.Command{
	Use:   "sandbox",
	Short: "Manage Docker sandbox environment",
}

var sandboxUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Start the sandbox",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := sandbox.BuildImage(); err != nil {
			return err
		}
		_, err := sandbox.Start()
		return err
	},
}

var sandboxDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop the sandbox",
	RunE: func(cmd *cobra.Command, args []string) error {
		return sandbox.Stop()
	},
}

func init() {
	sandboxCmd.AddCommand(sandboxUpCmd)
	sandboxCmd.AddCommand(sandboxDownCmd)
	rootCmd.AddCommand(sandboxCmd)
}