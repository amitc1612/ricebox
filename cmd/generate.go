package cmd

import (
	"fmt"

	"github.com/amitc1612/ricebox/rice"
	"github.com/spf13/cobra"
)

var (
	genDescription string
)

var generateCmd = &cobra.Command{
	Use:   "generate [name]",
	Short: "Generate a rice from your current configs",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		output := fmt.Sprintf("./%s", name)
		return rice.Generate(output, name, genDescription)
	},
}

func init() {
	generateCmd.Flags().StringVarP(&genDescription, "description", "d", "", "Description of the rice")
	rootCmd.AddCommand(generateCmd)
}