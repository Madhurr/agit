package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "agit",
	Short: "Agent-native Git middleware",
	Long:  `agit adds semantic context to Git for AI agent workflows.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
