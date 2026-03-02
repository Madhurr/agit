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

var contextJsonOutput bool

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Manage agit contexts",
	Long:  `Commands for managing and viewing agit commit contexts.`,
}

var showCmd = &cobra.Command{
	Use:   "show [hash]",
	Short: "Show agit context details for a commit",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runContextShow,
}

func init() {
	rootCmd.AddCommand(contextCmd)
	contextCmd.AddCommand(showCmd)
	showCmd.Flags().BoolVar(&contextJsonOutput, "json", false, "output raw JSON")
}

func runContextShow(cmd *cobra.Command, args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working dir: %w", err)
	}

	var commitHash string
	if len(args) == 0 {
		commitHash, err = git.HeadHash(dir)
		if err != nil {
			return fmt.Errorf("get HEAD: %w", err)
		}
	} else {
		commitHash = args[0]
	}

	note, err := notes.Read(dir, commitHash)
	if err != nil {
		return err
	}
	if note == nil {
		fmt.Printf("No agit note found for %s. Commit was not made with agit.\n", commitHash)
		return nil
	}

	if contextJsonOutput {
		data, err := json.MarshalIndent(note, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal note: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	shortHash := commitHash
	if len(commitHash) > 7 {
		shortHash = commitHash[:7]
	}

	sectionHeader("Commit:")
	fmt.Printf("  %s — %s\n", shortHash, commitHash)

	if note.Task != "" {
		sectionHeader("Task:")
		fmt.Printf("  %s\n", note.Task)
	}

	if note.Intent != "" {
		sectionHeader("Intent:")
		fmt.Printf("  %s\n", note.Intent)
	}

	if note.Confidence > 0 {
		pct := fmt.Sprintf("%.0f%%", note.Confidence*100)
		var c *color.Color
		switch {
		case note.Confidence >= 0.8:
			c = color.New(color.FgGreen)
		case note.Confidence >= 0.5:
			c = color.New(color.FgYellow)
		default:
			c = color.New(color.FgRed)
		}
		sectionHeader("Confidence:")
		if note.ConfidenceRationale != "" {
			fmt.Printf("  %s — %q\n", c.Sprint(pct), note.ConfidenceRationale)
		} else {
			fmt.Printf("  %s\n", c.Sprint(pct))
		}
	}

	if note.Agent.ID != "" || note.Agent.Model != "" {
		sectionHeader("Agent:")
		var parts []string
		if note.Agent.ID != "" {
			parts = append(parts, note.Agent.ID)
		}
		if note.Agent.Model != "" {
			parts = append(parts, fmt.Sprintf("(%s)", note.Agent.Model))
		}
		if note.Agent.SessionID != "" {
			parts = append(parts, fmt.Sprintf("session:%s", note.Agent.SessionID))
		}
		fmt.Printf("  %s\n", strings.Join(parts, " "))
	}

	if len(note.AlternativesConsidered) > 0 {
		sectionHeader("Alternatives rejected:")
		for _, alt := range note.AlternativesConsidered {
			fmt.Printf("  • %s — %s\n", alt.Approach, alt.RejectedReason)
		}
	}

	if len(note.KeyDecisions) > 0 {
		sectionHeader("Key decisions:")
		for _, dec := range note.KeyDecisions {
			fmt.Printf("  • %s: %s\n", dec.Decision, dec.Rationale)
		}
	}

	if len(note.Risks) > 0 {
		sectionHeader("Risks:")
		for _, r := range note.Risks {
			var sc *color.Color
			switch r.Severity {
			case "high":
				sc = color.New(color.FgRed)
			case "medium":
				sc = color.New(color.FgYellow)
			default:
				sc = color.New(color.FgCyan)
			}
			fmt.Printf("  [%s] %s: %s\n", sc.Sprint(r.Severity), r.Area, r.Description)
		}
	}

	if len(note.ContextConsulted) > 0 {
		sectionHeader("Context consulted:")
		for _, item := range note.ContextConsulted {
			fmt.Printf("  • %s\n", item)
		}
	}

	if len(note.RippleEffects) > 0 {
		sectionHeader("Ripple effects (check these):")
		for _, item := range note.RippleEffects {
			fmt.Printf("  • %s\n", item)
		}
	}

	if len(note.Unknowns) > 0 {
		sectionHeader("Unknowns:")
		for _, item := range note.Unknowns {
			fmt.Printf("  • %s\n", item)
		}
	}

	if note.TestResults != nil {
		sectionHeader("Test results:")
		fmt.Printf("  %d passed, %d failed, %d skipped\n",
			note.TestResults.Passed, note.TestResults.Failed, note.TestResults.Skipped)
	}

	return nil
}

func sectionHeader(title string) {
	color.New(color.Bold, color.FgYellow).Println(title)
}
