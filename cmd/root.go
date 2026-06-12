package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ricebox",
	Short: "🍚 Ricebox - plug-and-play ricing for Hyprland",
	Long: `Ricebox lets you package, share, and switch between
complete Hyprland ricing setups with a single command.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}