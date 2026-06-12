package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yourusername/ricebox/rice"
)

var (
	dryRun bool
)

var applyCmd = &cobra.Command{
	Use:   "apply [rice-path]",
	Short: "Apply a rice to your system",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		applier, err := rice.NewApplier(args[0], dryRun)
		if err != nil {
			return err
		}
		return applier.Apply()
	},
}

func init() {
	applyCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be done without doing it")
	rootCmd.AddCommand(applyCmd)
}