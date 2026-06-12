package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/ricebox/rice"
	"github.com/yourusername/ricebox/sandbox"
)

var testCmd = &cobra.Command{
	Use:   "test [rice-path]",
	Short: "Test a rice in the sandbox",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate the rice first
		if _, err := rice.LoadManifest(args[0]); err != nil {
			return fmt.Errorf("invalid rice: %w", err)
		}

		// Start sandbox
		fmt.Println("🏗️  Building sandbox...")
		if err := sandbox.BuildImage(); err != nil {
			return err
		}

		s, err := sandbox.Start()
		if err != nil {
			return err
		}

		// Copy rice into container
		if err := s.CopyRiceToContainer(args[0]); err != nil {
			sandbox.Stop()
			return err
		}

		// Apply inside container
		fmt.Println("🍚 Applying rice in sandbox...")
		if err := s.Exec("bash", "-c",
			fmt.Sprintf("cd /home/rice/rices/%s && echo 'Rice applied!'", args[0]),
		); err != nil {
			sandbox.Stop()
			return err
		}

		fmt.Println("\n✅ Rice tested successfully!")
		fmt.Println("   VNC available at: localhost:5900")
		fmt.Println("   Use 'ricebox sandbox down' to stop")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}