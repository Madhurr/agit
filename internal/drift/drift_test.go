package drift_test

import (
	"testing"

	"github.com/Madhurr/agit/internal/drift"
	"github.com/Madhurr/agit/internal/notes"
)

func TestDiffBothNil(t *testing.T) {
	result := drift.Diff("abc", nil, "def", nil)
	if result.Summary != "Neither commit has agent context" {
		t.Errorf("unexpected summary: %s", result.Summary)
	}
	if len(result.Changes) != 0 {
		t.Errorf("expected 0 changes, got %d", len(result.Changes))
	}
}

func TestDiffFromNil(t *testing.T) {
	to := &notes.CommitNote{
		Intent:     "add auth",
		Confidence: 0.85,
		Risks: []notes.Risk{
			{Area: "security", Description: "tokens expire", Severity: "high"},
		},
	}
	result := drift.Diff("abc", nil, "def", to)
	if !result.ToNote || result.FromNote {
		t.Error("note flags wrong")
	}
	if len(result.Changes) == 0 {
		t.Error("expected changes when going from nil to note")
	}
	// Should have intent added, confidence added, risk added
	fields := make(map[string]bool)
	for _, c := range result.Changes {
		fields[c.Field] = true
	}
	if !fields["intent"] {
		t.Error("missing intent change")
	}
	if !fields["confidence"] {
		t.Error("missing confidence change")
	}
	if !fields["risks"] {
		t.Error("missing risks change")
	}
}

func TestDiffToNil(t *testing.T) {
	from := &notes.CommitNote{Intent: "something"}
	result := drift.Diff("abc", from, "def", nil)
	if result.Summary != "Agent context removed (commit has no agit metadata)" {
		t.Errorf("unexpected summary: %s", result.Summary)
	}
}

func TestDiffNoChanges(t *testing.T) {
	note := &notes.CommitNote{
		Intent:     "same",
		Confidence: 0.9,
		Task:       "same task",
	}
	result := drift.Diff("abc", note, "def", note)
	if len(result.Changes) != 0 {
		t.Errorf("expected 0 changes for identical notes, got %d", len(result.Changes))
	}
}

func TestDiffIntentChanged(t *testing.T) {
	from := &notes.CommitNote{Intent: "old intent"}
	to := &notes.CommitNote{Intent: "new intent"}
	result := drift.Diff("a", from, "b", to)

	found := false
	for _, c := range result.Changes {
		if c.Field == "intent" && c.Kind == drift.Changed {
			found = true
			if c.Old != "old intent" || c.New != "new intent" {
				t.Errorf("intent diff: old=%q new=%q", c.Old, c.New)
			}
		}
	}
	if !found {
		t.Error("missing intent change")
	}
}

func TestDiffConfidenceChanged(t *testing.T) {
	from := &notes.CommitNote{Confidence: 0.5}
	to := &notes.CommitNote{Confidence: 0.9}
	result := drift.Diff("a", from, "b", to)

	found := false
	for _, c := range result.Changes {
		if c.Field == "confidence" {
			found = true
			if c.Old != "50%" || c.New != "90%" {
				t.Errorf("confidence: old=%q new=%q", c.Old, c.New)
			}
		}
	}
	if !found {
		t.Error("missing confidence change")
	}
}

func TestDiffRiskAdded(t *testing.T) {
	from := &notes.CommitNote{}
	to := &notes.CommitNote{
		Risks: []notes.Risk{{Area: "perf", Description: "slow", Severity: "medium"}},
	}
	result := drift.Diff("a", from, "b", to)

	found := false
	for _, c := range result.Changes {
		if c.Field == "risks" && c.Kind == drift.Added {
			found = true
		}
	}
	if !found {
		t.Error("missing risk added change")
	}
}

func TestDiffRiskResolved(t *testing.T) {
	from := &notes.CommitNote{
		Risks: []notes.Risk{{Area: "perf", Description: "slow", Severity: "high"}},
	}
	to := &notes.CommitNote{}
	result := drift.Diff("a", from, "b", to)

	found := false
	for _, c := range result.Changes {
		if c.Field == "risks" && c.Kind == drift.Resolved {
			found = true
		}
	}
	if !found {
		t.Error("missing risk resolved change")
	}
}

func TestDiffRiskEscalated(t *testing.T) {
	from := &notes.CommitNote{
		Risks: []notes.Risk{{Area: "perf", Description: "slow", Severity: "low"}},
	}
	to := &notes.CommitNote{
		Risks: []notes.Risk{{Area: "perf", Description: "very slow", Severity: "high"}},
	}
	result := drift.Diff("a", from, "b", to)

	found := false
	for _, c := range result.Changes {
		if c.Field == "risks" && c.Kind == drift.Changed {
			found = true
		}
	}
	if !found {
		t.Error("missing risk severity change")
	}

	if !drift.RiskEscalated(from, to) {
		t.Error("RiskEscalated should return true")
	}
}

func TestDiffUnknownResolved(t *testing.T) {
	from := &notes.CommitNote{Unknowns: []string{"token revocation"}}
	to := &notes.CommitNote{Unknowns: []string{}}
	result := drift.Diff("a", from, "b", to)

	found := false
	for _, c := range result.Changes {
		if c.Field == "unknowns" && c.Kind == drift.Resolved {
			found = true
			if c.Old != "token revocation" {
				t.Errorf("unexpected old value: %s", c.Old)
			}
		}
	}
	if !found {
		t.Error("missing unknown resolved change")
	}
}

func TestDiffAlternativesAdded(t *testing.T) {
	from := &notes.CommitNote{
		AlternativesConsidered: []notes.Alternative{
			{Approach: "JWT", RejectedReason: "too complex"},
		},
	}
	to := &notes.CommitNote{
		AlternativesConsidered: []notes.Alternative{
			{Approach: "JWT", RejectedReason: "too complex"},
			{Approach: "OAuth", RejectedReason: "overkill"},
		},
	}
	result := drift.Diff("a", from, "b", to)

	found := false
	for _, c := range result.Changes {
		if c.Field == "alternatives" && c.Kind == drift.Added && c.New == "OAuth" {
			found = true
		}
	}
	if !found {
		t.Error("missing alternative added change")
	}
}

func TestDiffTestResults(t *testing.T) {
	from := &notes.CommitNote{
		TestResults: &notes.TestResults{Passed: 40, Failed: 2, Skipped: 0},
	}
	to := &notes.CommitNote{
		TestResults: &notes.TestResults{Passed: 42, Failed: 0, Skipped: 0},
	}
	result := drift.Diff("a", from, "b", to)

	found := false
	for _, c := range result.Changes {
		if c.Field == "test_results" && c.Kind == drift.Changed {
			found = true
		}
	}
	if !found {
		t.Error("missing test results change")
	}
}

func TestConfidenceDelta(t *testing.T) {
	from := &notes.CommitNote{Confidence: 0.5}
	to := &notes.CommitNote{Confidence: 0.85}
	delta := drift.ConfidenceDelta(from, to)
	if delta < 0.34 || delta > 0.36 {
		t.Errorf("expected ~0.35, got %f", delta)
	}
}

func TestConfidenceDeltaNil(t *testing.T) {
	delta := drift.ConfidenceDelta(nil, nil)
	if delta != 0 {
		t.Errorf("expected 0 for nil, got %f", delta)
	}
}
