package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/madhurm/agit/internal/git"
	"github.com/madhurm/agit/internal/notes"
	"github.com/spf13/cobra"
)

type LogOutput struct {
	Commit git.LogEntry
	Note   *notes.CommitNote
}

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show semantic git log with agit notes",
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
			os.Exit(1)
		}

		count, _ := cmd.Flags().GetInt("count")
		jsonOutput, _ := cmd.Flags().GetBool("json")

		commits, err := git.Log(dir, count)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting log: %v\n", err)
			os.Exit(1)
		}

		var outputs []LogOutput
		for _, commit := range commits {
			note, _ := notes.Read(dir, commit.Hash)

			if jsonOutput {
				outputs = append(outputs, LogOutput{Commit: commit, Note: note})
				continue
			}

			printCommit(commit, note)
		}

		if jsonOutput {
			jsonData, err := json.MarshalIndent(outputs, "", "  ")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating JSON: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(string(jsonData))
		}
	},
}

func printCommit(commit git.LogEntry, note *notes.CommitNote) {
	yellow := color.New(color.FgYellow).SprintFunc()
	dim := color.New(color.Faint).SprintfFunc()
	timestamp := commit.Timestamp.Format("2 Jan 2006")

	if note == nil {
		fmt.Printf("%s  %s  %s\n                %s%s\n",
			yellow(commit.ShortHash),
			commit.Subject,
			dim(timestamp),
			dim("· "),
			commit.Author)
		return
	}

	confidenceColor := color.New()
	switch {
	case note.Confidence >= 0.8:
		confidenceColor = color.New(color.FgGreen)
	case note.Confidence >= 0.5:
		confidenceColor = color.New(color.FgYellow)
	default:
		confidenceColor = color.New(color.FgRed)
	}

	riskBadge := ""
	if len(note.Risks) > 0 {
		switch {
		case containsRisk(note.Risks, "high"):
			riskBadge = color.New(color.FgRed).Sprint(" [risk:high]")
		case containsRisk(note.Risks, "medium"):
			riskBadge = color.New(color.FgYellow).Sprint(" [risk:medium]")
		default:
			riskBadge = color.New(color.FgCyan).Sprint(" [risk:low]")
		}
	}

	fmt.Printf("%s  %s  %s%s%s\n                %s%s\n",
		yellow(commit.ShortHash),
		commit.Subject,
		dim(timestamp),
		confidenceColor.Sprintf(" [%.0f%%]", note.Confidence*100),
		riskBadge,
		dim("· "),
		commit.Author)

	if note.Intent != "" {
		green := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("intent: %s\n", green(note.Intent))
	}

	if len(note.AlternativesConsidered) > 0 {
		names := make([]string, len(note.AlternativesConsidered))
		for i, a := range note.AlternativesConsidered {
			names[i] = a.Approach
		}
		fmt.Printf("tried:  %s\n", dim(strings.Join(names, ", ")))
	}
}

func containsRisk(risks []notes.Risk, level string) bool {
	for _, r := range risks {
		if strings.ToLower(r.Severity) == level {
			return true
		}
	}
	return false
}

func init() {
	rootCmd.AddCommand(logCmd)
	logCmd.Flags().IntP("count", "n", 10, "number of commits to show")
	logCmd.Flags().Bool("json", false, "output as JSON array")
}
