package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/ricebox/rice"
)

var generateCmd = &cobra.Command{
	Use:   "generate [name]",
	Short: "Generate a rice from current configs",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		output := fmt.Sprintf("./%s", name)

		gen := rice.NewGenerator(output, name)
		if err := gen.Generate(); err != nil {
			return err
		}

		fmt.Printf("\n✅ Rice generated: %s/\n", output)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}