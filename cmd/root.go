package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ricebox",
	Short: "🍚 ricebox — plug-and-play ricing for Hyprland",
	Long: `ricebox lets you package, share, and switch between
complete Hyprland ricing setups with a single command.

Generate a rice from your current setup:
  ricebox generate my-rice

Apply a rice to your system:
  ricebox apply ./my-rice

Apply with dry-run to preview changes:
  ricebox apply ./my-rice --dry-run`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}