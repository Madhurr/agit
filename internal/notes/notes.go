// Package notes reads and writes agit metadata to git notes.
package notes

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Madhurr/agit/internal/git"
)

// Alternative represents an approach the agent considered but rejected.
type Alternative struct {
	Approach       string `json:"approach"`
	RejectedReason string `json:"rejected_reason"`
}

// KeyDecision represents a significant architectural or implementation choice.
type KeyDecision struct {
	Decision  string `json:"decision"`
	Rationale string `json:"rationale"`
}

// Risk represents a potential issue the agent identified.
type Risk struct {
	Area        string `json:"area"`
	Description string `json:"description"`
	Severity    string `json:"severity"` // "low", "medium", "high"
}

// TestResults captures test run outcomes at commit time.
type TestResults struct {
	Passed  int    `json:"passed"`
	Failed  int    `json:"failed"`
	Skipped int    `json:"skipped"`
	Command string `json:"command"`
}

// AgentInfo identifies the agent that made the commit.
type AgentInfo struct {
	ID        string `json:"id"`
	Model     string `json:"model"`
	SessionID string `json:"session_id"`
}

// CommitNote is the full agent context stored as a git note on each commit.
type CommitNote struct {
	SchemaVersion          string        `json:"schema_version"`
	CommitHash             string        `json:"commit_hash"`
	Agent                  AgentInfo     `json:"agent"`
	Task                   string        `json:"task"`
	Intent                 string        `json:"intent"`
	Confidence             float64       `json:"confidence"`
	ConfidenceRationale    string        `json:"confidence_rationale"`
	AlternativesConsidered []Alternative `json:"alternatives_considered"`
	KeyDecisions           []KeyDecision `json:"key_decisions"`
	Risks                  []Risk        `json:"risks"`
	ContextConsulted       []string      `json:"context_consulted"`
	TestResults            *TestResults  `json:"test_results,omitempty"`
	RippleEffects          []string      `json:"ripple_effects"`
	Unknowns               []string      `json:"unknowns"`
}

// Write serializes note to JSON and stores it on commitHash.
// Uses: git notes --ref=agit add -f -m <json> <commitHash>
func Write(dir, commitHash string, note *CommitNote) error {
	data, err := json.Marshal(note)
	if err != nil {
		return fmt.Errorf("marshal note: %w", err)
	}
	_, err = git.RunGit(dir, "notes", "--ref=agit", "add", "-f", "-m", string(data), commitHash)
	return err
}

// Read retrieves and deserializes the agit note for commitHash.
// Returns (nil, nil) if no note exists — not an error.
func Read(dir, commitHash string) (*CommitNote, error) {
	out, err := git.RunGit(dir, "notes", "--ref=agit", "show", commitHash)
	if err != nil {
		if strings.Contains(err.Error(), "No note found") ||
			strings.Contains(err.Error(), "no note found") {
			return nil, nil
		}
		return nil, err
	}
	var note CommitNote
	if err := json.Unmarshal([]byte(out), &note); err != nil {
		return nil, fmt.Errorf("unmarshal note: %w", err)
	}
	return &note, nil
}

// Exists returns true if an agit note exists for commitHash.
func Exists(dir, commitHash string) bool {
	note, err := Read(dir, commitHash)
	return err == nil && note != nil
}

// Delete removes the agit note for commitHash.
func Delete(dir, commitHash string) error {
	_, err := git.RunGit(dir, "notes", "--ref=agit", "remove", commitHash)
	return err
}
