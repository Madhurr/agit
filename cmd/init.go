package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/madhurm/agit/internal/git"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize agit in the current git repository",
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting working directory: %v\n", err)
			return
		}

		if !git.IsRepo(dir) {
			fmt.Println(color.RedString("✗ not a git repository"))
			return
		}

		// Configure remote origin fetch for agit notes (ignore error if origin doesn't exist)
		git.RunGit(dir, "config", "--add", "remote.origin.fetch", "+refs/notes/agit:refs/notes/agit")

		// Configure notes rewrite ref to carry agit notes during rebase/cherry-pick
		if _, err := git.RunGit(dir, "config", "notes.rewriteRef", "refs/notes/agit"); err != nil {
			fmt.Printf("Error configuring notes.rewriteRef: %v\n", err)
			return
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
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
