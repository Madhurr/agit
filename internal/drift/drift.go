// Package drift computes semantic diffs between two agit CommitNotes.
package drift

import (
	"fmt"
	"math"
	"strings"

	"github.com/madhurm/agit/internal/notes"
)

// ChangeKind classifies a diff entry.
type ChangeKind string

const (
	Added    ChangeKind = "added"
	Removed  ChangeKind = "removed"
	Changed  ChangeKind = "changed"
	Resolved ChangeKind = "resolved" // unknown became known, risk mitigated, etc.
)

// FieldDiff represents a single changed field between two notes.
type FieldDiff struct {
	Field    string     `json:"field"`
	Kind     ChangeKind `json:"kind"`
	Old      string     `json:"old,omitempty"`
	New      string     `json:"new,omitempty"`
	Detail   string     `json:"detail,omitempty"` // human-readable summary
	Severity string     `json:"severity,omitempty"` // for risk changes
}

// DiffResult holds the full semantic diff between two commits.
type DiffResult struct {
	FromHash    string      `json:"from_hash"`
	ToHash      string      `json:"to_hash"`
	FromNote    bool        `json:"from_has_note"`
	ToNote      bool        `json:"to_has_note"`
	Changes     []FieldDiff `json:"changes"`
	FilesAdded  []string    `json:"files_added,omitempty"`
	FilesRemoved []string   `json:"files_removed,omitempty"`
	Summary     string      `json:"summary"`
}

// Diff computes the semantic difference between two CommitNotes.
// Either note can be nil (commit without agit metadata).
func Diff(fromHash string, from *notes.CommitNote, toHash string, to *notes.CommitNote) *DiffResult {
	result := &DiffResult{
		FromHash: fromHash,
		ToHash:   toHash,
		FromNote: from != nil,
		ToNote:   to != nil,
	}

	if from == nil && to == nil {
		result.Summary = "Neither commit has agent context"
		return result
	}

	if from == nil {
		result.Summary = "Agent context added (no prior context)"
		if to.Intent != "" {
			result.Changes = append(result.Changes, FieldDiff{
				Field: "intent", Kind: Added, New: to.Intent,
				Detail: fmt.Sprintf("Intent set: %s", to.Intent),
			})
		}
		addNoteAsNew(result, to)
		return result
	}

	if to == nil {
		result.Summary = "Agent context removed (commit has no agit metadata)"
		return result
	}

	// Both notes exist — compare field by field
	diffIntent(result, from, to)
	diffConfidence(result, from, to)
	diffTask(result, from, to)
	diffAgent(result, from, to)
	diffAlternatives(result, from, to)
	diffRisks(result, from, to)
	diffUnknowns(result, from, to)
	diffRipple(result, from, to)
	diffDecisions(result, from, to)
	diffTests(result, from, to)

	if len(result.Changes) == 0 {
		result.Summary = "No semantic changes in agent context"
	} else {
		result.Summary = fmt.Sprintf("%d change(s) in agent context", len(result.Changes))
	}

	return result
}

func diffIntent(r *DiffResult, from, to *notes.CommitNote) {
	if from.Intent != to.Intent {
		if from.Intent == "" {
			r.Changes = append(r.Changes, FieldDiff{
				Field: "intent", Kind: Added, New: to.Intent,
				Detail: fmt.Sprintf("Intent set: %s", to.Intent),
			})
		} else if to.Intent == "" {
			r.Changes = append(r.Changes, FieldDiff{
				Field: "intent", Kind: Removed, Old: from.Intent,
				Detail: "Intent removed",
			})
		} else {
			r.Changes = append(r.Changes, FieldDiff{
				Field: "intent", Kind: Changed, Old: from.Intent, New: to.Intent,
				Detail: fmt.Sprintf("Intent changed: %q → %q", from.Intent, to.Intent),
			})
		}
	}
}

func diffConfidence(r *DiffResult, from, to *notes.CommitNote) {
	if from.Confidence != to.Confidence {
		delta := to.Confidence - from.Confidence
		direction := "↑"
		if delta < 0 {
			direction = "↓"
		}
		r.Changes = append(r.Changes, FieldDiff{
			Field:  "confidence",
			Kind:   Changed,
			Old:    fmt.Sprintf("%.0f%%", from.Confidence*100),
			New:    fmt.Sprintf("%.0f%%", to.Confidence*100),
			Detail: fmt.Sprintf("Confidence %s %.0f%% → %.0f%% (%+.0f%%)", direction, from.Confidence*100, to.Confidence*100, delta*100),
		})
	}
	if from.ConfidenceRationale != to.ConfidenceRationale && to.ConfidenceRationale != "" {
		if from.ConfidenceRationale == "" {
			r.Changes = append(r.Changes, FieldDiff{
				Field: "confidence_rationale", Kind: Added, New: to.ConfidenceRationale,
				Detail: fmt.Sprintf("Rationale added: %s", to.ConfidenceRationale),
			})
		} else {
			r.Changes = append(r.Changes, FieldDiff{
				Field: "confidence_rationale", Kind: Changed,
				Old: from.ConfidenceRationale, New: to.ConfidenceRationale,
				Detail: "Confidence rationale updated",
			})
		}
	}
}

func diffTask(r *DiffResult, from, to *notes.CommitNote) {
	if from.Task != to.Task {
		kind := Changed
		if from.Task == "" {
			kind = Added
		}
		r.Changes = append(r.Changes, FieldDiff{
			Field: "task", Kind: kind, Old: from.Task, New: to.Task,
			Detail: fmt.Sprintf("Task: %q → %q", from.Task, to.Task),
		})
	}
}

func diffAgent(r *DiffResult, from, to *notes.CommitNote) {
	if from.Agent.ID != to.Agent.ID || from.Agent.Model != to.Agent.Model {
		oldAgent := formatAgent(from.Agent)
		newAgent := formatAgent(to.Agent)
		if oldAgent != newAgent {
			r.Changes = append(r.Changes, FieldDiff{
				Field: "agent", Kind: Changed, Old: oldAgent, New: newAgent,
				Detail: fmt.Sprintf("Agent changed: %s → %s", oldAgent, newAgent),
			})
		}
	}
}

func diffAlternatives(r *DiffResult, from, to *notes.CommitNote) {
	fromSet := make(map[string]string)
	for _, a := range from.AlternativesConsidered {
		fromSet[a.Approach] = a.RejectedReason
	}
	toSet := make(map[string]string)
	for _, a := range to.AlternativesConsidered {
		toSet[a.Approach] = a.RejectedReason
	}

	for approach, reason := range toSet {
		if _, exists := fromSet[approach]; !exists {
			r.Changes = append(r.Changes, FieldDiff{
				Field: "alternatives", Kind: Added, New: approach,
				Detail: fmt.Sprintf("New alternative rejected: %s — %s", approach, reason),
			})
		}
	}

	for approach := range fromSet {
		if _, exists := toSet[approach]; !exists {
			r.Changes = append(r.Changes, FieldDiff{
				Field: "alternatives", Kind: Removed, Old: approach,
				Detail: fmt.Sprintf("Alternative no longer listed: %s", approach),
			})
		}
	}
}

func diffRisks(r *DiffResult, from, to *notes.CommitNote) {
	fromRisks := make(map[string]notes.Risk)
	for _, risk := range from.Risks {
		fromRisks[risk.Area] = risk
	}
	toRisks := make(map[string]notes.Risk)
	for _, risk := range to.Risks {
		toRisks[risk.Area] = risk
	}

	for area, toRisk := range toRisks {
		fromRisk, existed := fromRisks[area]
		if !existed {
			r.Changes = append(r.Changes, FieldDiff{
				Field: "risks", Kind: Added, New: formatRisk(toRisk),
				Severity: toRisk.Severity,
				Detail:   fmt.Sprintf("New risk: [%s] %s — %s", toRisk.Severity, area, toRisk.Description),
			})
		} else if fromRisk.Severity != toRisk.Severity {
			r.Changes = append(r.Changes, FieldDiff{
				Field:    "risks",
				Kind:     Changed,
				Old:      formatRisk(fromRisk),
				New:      formatRisk(toRisk),
				Severity: toRisk.Severity,
				Detail:   fmt.Sprintf("Risk %s: severity %s → %s", area, fromRisk.Severity, toRisk.Severity),
			})
		}
	}

	for area, fromRisk := range fromRisks {
		if _, exists := toRisks[area]; !exists {
			r.Changes = append(r.Changes, FieldDiff{
				Field:    "risks",
				Kind:     Resolved,
				Old:      formatRisk(fromRisk),
				Severity: fromRisk.Severity,
				Detail:   fmt.Sprintf("Risk resolved: [%s] %s", fromRisk.Severity, area),
			})
		}
	}
}

func diffUnknowns(r *DiffResult, from, to *notes.CommitNote) {
	fromSet := toStringSet(from.Unknowns)
	toSet := toStringSet(to.Unknowns)

	for u := range toSet {
		if !fromSet[u] {
			r.Changes = append(r.Changes, FieldDiff{
				Field: "unknowns", Kind: Added, New: u,
				Detail: fmt.Sprintf("New unknown: %s", u),
			})
		}
	}

	for u := range fromSet {
		if !toSet[u] {
			r.Changes = append(r.Changes, FieldDiff{
				Field: "unknowns", Kind: Resolved, Old: u,
				Detail: fmt.Sprintf("Unknown resolved: %s", u),
			})
		}
	}
}

func diffRipple(r *DiffResult, from, to *notes.CommitNote) {
	fromSet := toStringSet(from.RippleEffects)
	toSet := toStringSet(to.RippleEffects)

	for item := range toSet {
		if !fromSet[item] {
			r.Changes = append(r.Changes, FieldDiff{
				Field: "ripple_effects", Kind: Added, New: item,
				Detail: fmt.Sprintf("New ripple effect: %s", item),
			})
		}
	}

	for item := range fromSet {
		if !toSet[item] {
			r.Changes = append(r.Changes, FieldDiff{
				Field: "ripple_effects", Kind: Resolved, Old: item,
				Detail: fmt.Sprintf("Ripple effect addressed: %s", item),
			})
		}
	}
}

func diffDecisions(r *DiffResult, from, to *notes.CommitNote) {
	fromSet := make(map[string]string)
	for _, d := range from.KeyDecisions {
		fromSet[d.Decision] = d.Rationale
	}

	for _, d := range to.KeyDecisions {
		if _, exists := fromSet[d.Decision]; !exists {
			r.Changes = append(r.Changes, FieldDiff{
				Field: "key_decisions", Kind: Added, New: d.Decision,
				Detail: fmt.Sprintf("New decision: %s — %s", d.Decision, d.Rationale),
			})
		}
	}
}

func diffTests(r *DiffResult, from, to *notes.CommitNote) {
	if from.TestResults == nil && to.TestResults == nil {
		return
	}
	if from.TestResults == nil && to.TestResults != nil {
		r.Changes = append(r.Changes, FieldDiff{
			Field: "test_results", Kind: Added,
			New:    fmt.Sprintf("%d passed, %d failed, %d skipped", to.TestResults.Passed, to.TestResults.Failed, to.TestResults.Skipped),
			Detail: "Test results added",
		})
		return
	}
	if from.TestResults != nil && to.TestResults != nil {
		if from.TestResults.Failed != to.TestResults.Failed || from.TestResults.Passed != to.TestResults.Passed {
			failDelta := to.TestResults.Failed - from.TestResults.Failed
			passDelta := to.TestResults.Passed - from.TestResults.Passed
			var parts []string
			if passDelta != 0 {
				parts = append(parts, fmt.Sprintf("passed %+d", passDelta))
			}
			if failDelta != 0 {
				parts = append(parts, fmt.Sprintf("failed %+d", failDelta))
			}
			r.Changes = append(r.Changes, FieldDiff{
				Field:  "test_results",
				Kind:   Changed,
				Old:    fmt.Sprintf("%d passed, %d failed", from.TestResults.Passed, from.TestResults.Failed),
				New:    fmt.Sprintf("%d passed, %d failed", to.TestResults.Passed, to.TestResults.Failed),
				Detail: fmt.Sprintf("Tests: %s", strings.Join(parts, ", ")),
			})
		}
	}
}

// addNoteAsNew adds all fields from a note as "added" changes.
func addNoteAsNew(r *DiffResult, note *notes.CommitNote) {
	if note.Confidence > 0 {
		r.Changes = append(r.Changes, FieldDiff{
			Field: "confidence", Kind: Added,
			New:    fmt.Sprintf("%.0f%%", note.Confidence*100),
			Detail: fmt.Sprintf("Confidence: %.0f%%", note.Confidence*100),
		})
	}
	for _, risk := range note.Risks {
		r.Changes = append(r.Changes, FieldDiff{
			Field: "risks", Kind: Added, New: formatRisk(risk),
			Severity: risk.Severity,
			Detail:   fmt.Sprintf("Risk: [%s] %s — %s", risk.Severity, risk.Area, risk.Description),
		})
	}
	for _, u := range note.Unknowns {
		r.Changes = append(r.Changes, FieldDiff{
			Field: "unknowns", Kind: Added, New: u,
			Detail: fmt.Sprintf("Unknown: %s", u),
		})
	}
}

// Helper functions

func formatAgent(a notes.AgentInfo) string {
	parts := []string{}
	if a.ID != "" && a.ID != "unknown" {
		parts = append(parts, a.ID)
	}
	if a.Model != "" && a.Model != "unknown" {
		parts = append(parts, fmt.Sprintf("(%s)", a.Model))
	}
	if len(parts) == 0 {
		return "unknown"
	}
	return strings.Join(parts, " ")
}

func formatRisk(r notes.Risk) string {
	return fmt.Sprintf("[%s] %s: %s", r.Severity, r.Area, r.Description)
}

func toStringSet(items []string) map[string]bool {
	set := make(map[string]bool)
	for _, item := range items {
		set[item] = true
	}
	return set
}

// ConfidenceDelta returns the absolute change in confidence.
func ConfidenceDelta(from, to *notes.CommitNote) float64 {
	if from == nil || to == nil {
		return 0
	}
	return math.Abs(to.Confidence - from.Confidence)
}

// RiskEscalated returns true if any risk severity increased.
func RiskEscalated(from, to *notes.CommitNote) bool {
	if from == nil || to == nil {
		return false
	}
	fromRisks := make(map[string]string)
	for _, r := range from.Risks {
		fromRisks[r.Area] = r.Severity
	}
	severityRank := map[string]int{"low": 1, "medium": 2, "high": 3}
	for _, r := range to.Risks {
		if oldSev, exists := fromRisks[r.Area]; exists {
			if severityRank[r.Severity] > severityRank[oldSev] {
				return true
			}
		}
	}
	return false
}
