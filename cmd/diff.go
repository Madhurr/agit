package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/madhurm/agit/internal/drift"
	"github.com/madhurm/agit/internal/git"
	"github.com/madhurm/agit/internal/notes"
)

var diffCmd = &cobra.Command{
	Use:   "diff [from] [to]",
	Short: "Semantic diff of agent context between two commits",
	Long: `Compare agent reasoning between two commits.

Shows what changed in intent, confidence, risks, unknowns, alternatives,
and other agent metadata — not the code diff, but the reasoning diff.

Examples:
  agit diff              # HEAD~1 vs HEAD
  agit diff abc1234      # abc1234 vs HEAD
  agit diff abc1234 def5678  # abc1234 vs def5678`,
	Args: cobra.MaximumNArgs(2),
	RunE: runDiff,
}

func init() {
	rootCmd.AddCommand(diffCmd)
	diffCmd.Flags().Bool("json", false, "output as JSON")
	diffCmd.Flags().Bool("files", false, "include file-level changes")
}

func runDiff(cmd *cobra.Command, args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working dir: %w", err)
	}

	jsonOutput, _ := cmd.Flags().GetBool("json")
	showFiles, _ := cmd.Flags().GetBool("files")

	// Resolve commit hashes
	var fromRef, toRef string

	switch len(args) {
	case 0:
		fromRef = "HEAD~1"
		toRef = "HEAD"
	case 1:
		// Support "abc..def" syntax (like git diff)
		if parts := strings.SplitN(args[0], "..", 2); len(parts) == 2 && parts[0] != "" && parts[1] != "" {
			fromRef = parts[0]
			toRef = parts[1]
		} else {
			fromRef = args[0]
			toRef = "HEAD"
		}
	case 2:
		fromRef = args[0]
		toRef = args[1]
	}

	fromHash, err := resolveRef(dir, fromRef)
	if err != nil {
		return fmt.Errorf("cannot resolve %q: %w", fromRef, err)
	}

	toHash, err := resolveRef(dir, toRef)
	if err != nil {
		return fmt.Errorf("cannot resolve %q: %w", toRef, err)
	}

	// Read notes
	fromNote, _ := notes.Read(dir, fromHash)
	toNote, _ := notes.Read(dir, toHash)

	// Compute semantic diff
	result := drift.Diff(fromHash, fromNote, toHash, toNote)

	// Optionally include file changes
	if showFiles {
		files, err := git.DiffFiles(dir, toHash)
		if err == nil {
			result.FilesAdded = files
		}
	}

	if jsonOutput {
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	printDiff(result, fromRef, toRef)
	return nil
}

func printDiff(result *drift.DiffResult, fromRef, toRef string) {
	header := color.New(color.Bold, color.FgYellow)
	dim := color.New(color.Faint)

	// Header
	header.Printf("agit diff %s..%s\n", shortRef(result.FromHash), shortRef(result.ToHash))
	fmt.Println()

	if !result.FromNote && !result.ToNote {
		dim.Println("Neither commit has agent context.")
		return
	}

	if !result.FromNote {
		color.Cyan("ℹ %s has no agit note — showing all context as new\n\n", fromRef)
	}

	if !result.ToNote {
		color.Yellow("⚠ %s has no agit note — agent context was dropped\n", toRef)
		return
	}

	if len(result.Changes) == 0 {
		color.Green("✓ No semantic changes in agent context\n")
		return
	}

	fmt.Printf("%s\n\n", result.Summary)

	// Group changes by kind for cleaner output
	var added, changed, removed, resolved []drift.FieldDiff
	for _, c := range result.Changes {
		switch c.Kind {
		case drift.Added:
			added = append(added, c)
		case drift.Changed:
			changed = append(changed, c)
		case drift.Removed:
			removed = append(removed, c)
		case drift.Resolved:
			resolved = append(resolved, c)
		}
	}

	if len(changed) > 0 {
		header.Println("Changed:")
		for _, c := range changed {
			printChange(c)
		}
		fmt.Println()
	}

	if len(added) > 0 {
		header.Println("Added:")
		for _, c := range added {
			printAdded(c)
		}
		fmt.Println()
	}

	if len(resolved) > 0 {
		header.Println("Resolved:")
		for _, c := range resolved {
			printResolved(c)
		}
		fmt.Println()
	}

	if len(removed) > 0 {
		header.Println("Removed:")
		for _, c := range removed {
			dim.Printf("  - %s\n", c.Detail)
		}
		fmt.Println()
	}

	if len(result.FilesAdded) > 0 {
		header.Println("Files changed:")
		for _, f := range result.FilesAdded {
			fmt.Printf("  %s\n", f)
		}
		fmt.Println()
	}
}

func printChange(c drift.FieldDiff) {
	switch c.Field {
	case "confidence":
		// Color based on direction
		if strings.Contains(c.Detail, "↑") {
			color.Green("  %s %s → %s", fieldLabel(c.Field), c.Old, c.New)
		} else {
			color.Red("  %s %s → %s", fieldLabel(c.Field), c.Old, c.New)
		}
		fmt.Println()
	case "risks":
		sev := severityColor(c.Severity)
		sev.Printf("  %s %s\n", fieldLabel(c.Field), c.Detail)
	case "test_results":
		fmt.Printf("  %s %s\n", fieldLabel(c.Field), c.Detail)
	default:
		fmt.Printf("  %s %q → %q\n", fieldLabel(c.Field), c.Old, c.New)
	}
}

func printAdded(c drift.FieldDiff) {
	switch c.Field {
	case "risks":
		sev := severityColor(c.Severity)
		sev.Printf("  + %s\n", c.Detail)
	default:
		color.Green("  + %s\n", c.Detail)
	}
}

func printResolved(c drift.FieldDiff) {
	color.Green("  ✓ %s\n", c.Detail)
}

func fieldLabel(field string) string {
	labels := map[string]string{
		"intent":                "Intent:",
		"confidence":            "Confidence:",
		"confidence_rationale":  "Rationale:",
		"task":                  "Task:",
		"agent":                 "Agent:",
		"risks":                 "Risk:",
		"unknowns":              "Unknown:",
		"alternatives":          "Alternative:",
		"ripple_effects":        "Ripple:",
		"key_decisions":         "Decision:",
		"test_results":          "Tests:",
	}
	if label, ok := labels[field]; ok {
		return label
	}
	return field + ":"
}

func severityColor(severity string) *color.Color {
	switch severity {
	case "high":
		return color.New(color.FgRed)
	case "medium":
		return color.New(color.FgYellow)
	default:
		return color.New(color.FgCyan)
	}
}

func resolveRef(dir, ref string) (string, error) {
	hash, err := git.RunGit(dir, "rev-parse", ref)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func shortRef(hash string) string {
	if len(hash) > 7 {
		return hash[:7]
	}
	return hash
}
