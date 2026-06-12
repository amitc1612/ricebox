package cmd

import (
	"github.com/amitc1612/ricebox/rice"
	"github.com/spf13/cobra"
)

var dryRun bool

var applyCmd = &cobra.Command{
	Use:   "apply [rice-path]",
	Short: "Apply a rice to your system",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return rice.Apply(args[0], dryRun)
	},
}

func init() {
	applyCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without applying them")
	rootCmd.AddCommand(applyCmd)
}