package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	buildVersion = "dev"
	buildCommit  = "none"
	buildDate    = "unknown"
)

// SetVersion is called from main with values injected by GoReleaser ldflags.
func SetVersion(v, c, d string) {
	buildVersion = v
	buildCommit = c
	buildDate = d
}

var rootCmd = &cobra.Command{
	Use:   "agit",
	Short: "Agent-native Git middleware",
	Long:  `agit adds semantic context to Git for AI agent workflows.`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print agit version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("agit %s (commit: %s, built: %s)\n", buildVersion, buildCommit, buildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
