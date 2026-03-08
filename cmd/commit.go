package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/Madhurr/agit/internal/git"
	"github.com/Madhurr/agit/internal/notes"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Create a commit with agent metadata as git notes",
	RunE:  runCommit,
}

func init() {
	rootCmd.AddCommand(commitCmd)
	
	// Required flag
	commitCmd.Flags().StringP("message", "m", "", "commit message (required)")
	_ = commitCmd.MarkFlagRequired("message")

	// Optional flags with defaults
	commitCmd.Flags().String("intent", "", "what the agent was trying to accomplish")
	commitCmd.Flags().Float64("confidence", 0, "confidence 0.0-1.0")
	commitCmd.Flags().String("confidence-rationale", "", "why this confidence level")
	commitCmd.Flags().StringArray("tried", []string{}, "repeatable: \"approach: rejection reason\"")
	commitCmd.Flags().StringArray("risk", []string{}, "repeatable: \"severity: area: description\"")
	commitCmd.Flags().String("task", "", "the original human prompt/task")
	commitCmd.Flags().StringArray("unknowns", []string{}, "things agent was unsure about (repeatable)")
	commitCmd.Flags().StringArray("ripple", []string{}, "files noticed but not changed (repeatable)")
	
	// Default environment-based flags
	agentID := os.Getenv("AGIT_AGENT_ID")
	if agentID == "" {
		agentID = "unknown"
	}
	agentModel := os.Getenv("AGIT_MODEL")
	if agentModel == "" {
		agentModel = "unknown"
	}
	sessionID := os.Getenv("AGIT_SESSION_ID")

	commitCmd.Flags().String("agent-id", agentID, "agent identifier")
	commitCmd.Flags().String("agent-model", agentModel, "model name")
	commitCmd.Flags().String("session-id", sessionID, "session identifier")
	
	// JSON note input
	commitCmd.Flags().String("json-note", "", "path to JSON file OR \"-\" to read from stdin")
}

func runCommit(cmd *cobra.Command, args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Get flags
	message, _ := cmd.Flags().GetString("message")
	intent, _ := cmd.Flags().GetString("intent")
	confidence, _ := cmd.Flags().GetFloat64("confidence")
	confidenceRationale, _ := cmd.Flags().GetString("confidence-rationale")
	tried, _ := cmd.Flags().GetStringArray("tried")
	risk, _ := cmd.Flags().GetStringArray("risk")
	task, _ := cmd.Flags().GetString("task")
	unknowns, _ := cmd.Flags().GetStringArray("unknowns")
	ripple, _ := cmd.Flags().GetStringArray("ripple")
	agentID, _ := cmd.Flags().GetString("agent-id")
	agentModel, _ := cmd.Flags().GetString("agent-model")
	sessionID, _ := cmd.Flags().GetString("session-id")
	jsonNotePath, _ := cmd.Flags().GetString("json-note")

	var note notes.CommitNote

	// Handle JSON input
	if jsonNotePath != "" {
		var data []byte
		if jsonNotePath == "-" {
			data, err = io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("failed to read from stdin: %w", err)
			}
		} else {
			data, err = os.ReadFile(jsonNotePath)
			if err != nil {
				return fmt.Errorf("failed to read JSON file: %w", err)
			}
		}

		if err := json.Unmarshal(data, &note); err != nil {
			return fmt.Errorf("failed to parse JSON note: %w", err)
		}
	} else {
		// Build note from flags
		note.Intent = intent
		note.Confidence = confidence
		note.ConfidenceRationale = confidenceRationale
		note.Task = task
		note.Agent = notes.AgentInfo{
			ID:        agentID,
			Model:     agentModel,
			SessionID: sessionID,
		}

		// Parse --tried flags: "approach: rejection reason"
		for _, t := range tried {
			parts := strings.SplitN(t, ":", 2)
			if len(parts) == 0 {
				continue
			}
			approach := strings.TrimSpace(parts[0])
			rejection := ""
			if len(parts) > 1 {
				rejection = strings.TrimSpace(parts[1])
			}
			note.AlternativesConsidered = append(note.AlternativesConsidered,
				notes.Alternative{Approach: approach, RejectedReason: rejection})
		}

		// Parse --risk flags: "severity: area: description"
		for _, r := range risk {
			parts := strings.SplitN(r, ":", 3)
			severity := strings.TrimSpace(parts[0])
			area, description := "", ""
			if len(parts) > 1 {
				area = strings.TrimSpace(parts[1])
			}
			if len(parts) > 2 {
				description = strings.TrimSpace(parts[2])
			} else if len(parts) == 1 {
				description = severity
				severity = ""
			}
			note.Risks = append(note.Risks, notes.Risk{
				Severity:    severity,
				Area:        area,
				Description: description,
			})
		}

		note.Unknowns = unknowns
		note.RippleEffects = ripple
	}

	// Stage all files and commit
	if err := git.StageAll(dir); err != nil {
		return fmt.Errorf("failed to stage files: %w", err)
	}
	
	commitHash, err := git.CommitWithMessage(dir, message)
	if err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	// Set note metadata
	note.SchemaVersion = "1.0"
	note.CommitHash = commitHash

	// Write the note
	if err := notes.Write(dir, commitHash, &note); err != nil {
		return fmt.Errorf("failed to write git note: %w", err)
	}

	// Print success message with colors
	shortHash := commitHash[:7]
	color.Green("✓ Committed: " + shortHash)
	
	if intent != "" {
		color.Cyan("  intent: " + intent)
	}
	
	if confidence > 0 {
		confStr := strconv.FormatFloat(confidence, 'f', -1, 64)
		color.Cyan("  confidence: " + confStr)
	}

	fmt.Println(shortHash)
	return nil
}
