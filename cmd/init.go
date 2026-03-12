package cmd

import (
	"fmt"
	"os"

	"github.com/Madhurr/agit/internal/git"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize agit in the current git repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting working directory: %w", err)
		}

		if !git.IsRepo(dir) {
			return fmt.Errorf("%s", color.RedString("✗ not a git repository"))
		}

		// Configure remote origin fetch for agit notes (best-effort: origin may not exist yet)
		git.RunGit(dir, "config", "--add", "remote.origin.fetch", "+refs/notes/agit:refs/notes/agit") //nolint:errcheck

		// Configure notes rewrite ref to carry agit notes during rebase/cherry-pick
		if _, err := git.RunGit(dir, "config", "notes.rewriteRef", "refs/notes/agit"); err != nil {
			return fmt.Errorf("configuring notes.rewriteRef: %w", err)
		}

		// Print success message and instructions
		color.Green("✓ agit initialized in this repository")
		fmt.Println("")
		fmt.Println("To push notes to remote:")
		color.Cyan("  git push origin refs/notes/agit")
		fmt.Println("")
		fmt.Println("To fetch notes from remote:")
		color.Cyan("  git fetch origin refs/notes/agit:refs/notes/agit")
		fmt.Println("")
		fmt.Println("Agents should now use 'agit commit' instead of 'git commit'.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
