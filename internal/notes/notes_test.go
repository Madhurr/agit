package notes_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/madhurm/agit/internal/notes"
)

// initTestRepo creates a temporary git repo with one commit and returns its path.
func initTestRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=Test",
			"GIT_AUTHOR_EMAIL=test@test.com",
			"GIT_COMMITTER_NAME=Test",
			"GIT_COMMITTER_EMAIL=test@test.com",
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("command %v failed: %v\n%s", args, err, out)
		}
	}

	run("git", "init", "-b", "main")
	run("git", "config", "user.email", "test@test.com")
	run("git", "config", "user.name", "Test")

	// Create a file and commit it
	if err := os.WriteFile(filepath.Join(dir, "hello.txt"), []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	run("git", "add", ".")
	run("git", "commit", "-m", "initial commit")

	return dir
}

func headHash(t *testing.T, dir string) string {
	t.Helper()
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git rev-parse HEAD: %v", err)
	}
	hash := string(out)
	if len(hash) > 0 && hash[len(hash)-1] == '\n' {
		hash = hash[:len(hash)-1]
	}
	return hash
}

func TestWriteAndRead(t *testing.T) {
	dir := initTestRepo(t)
	hash := headHash(t, dir)

	note := &notes.CommitNote{
		SchemaVersion:       "1.0",
		CommitHash:          hash,
		Task:                "test task",
		Intent:              "test intent",
		Confidence:          0.85,
		ConfidenceRationale: "thoroughly tested",
		Agent: notes.AgentInfo{
			ID:        "test-agent",
			Model:     "test-model",
			SessionID: "test-session",
		},
		AlternativesConsidered: []notes.Alternative{
			{Approach: "approach-a", RejectedReason: "too slow"},
		},
		Risks: []notes.Risk{
			{Area: "performance", Description: "may be slow", Severity: "low"},
		},
		RippleEffects: []string{"file-a.go"},
		Unknowns:      []string{"edge case X"},
	}

	if err := notes.Write(dir, hash, note); err != nil {
		t.Fatalf("Write: %v", err)
	}

	got, err := notes.Read(dir, hash)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if got == nil {
		t.Fatal("Read returned nil note, want non-nil")
	}

	if got.Intent != note.Intent {
		t.Errorf("Intent: got %q, want %q", got.Intent, note.Intent)
	}
	if got.Confidence != note.Confidence {
		t.Errorf("Confidence: got %v, want %v", got.Confidence, note.Confidence)
	}
	if got.Agent.ID != note.Agent.ID {
		t.Errorf("Agent.ID: got %q, want %q", got.Agent.ID, note.Agent.ID)
	}
	if got.Agent.Model != note.Agent.Model {
		t.Errorf("Agent.Model: got %q, want %q", got.Agent.Model, note.Agent.Model)
	}
	if len(got.AlternativesConsidered) != 1 {
		t.Fatalf("AlternativesConsidered: got %d items, want 1", len(got.AlternativesConsidered))
	}
	if got.AlternativesConsidered[0].Approach != "approach-a" {
		t.Errorf("Alternative.Approach: got %q, want %q", got.AlternativesConsidered[0].Approach, "approach-a")
	}
	if len(got.Risks) != 1 || got.Risks[0].Severity != "low" {
		t.Errorf("Risks: unexpected value %+v", got.Risks)
	}
	if len(got.RippleEffects) != 1 || got.RippleEffects[0] != "file-a.go" {
		t.Errorf("RippleEffects: got %v, want [file-a.go]", got.RippleEffects)
	}
	if len(got.Unknowns) != 1 || got.Unknowns[0] != "edge case X" {
		t.Errorf("Unknowns: got %v, want [edge case X]", got.Unknowns)
	}
}

func TestReadNilForMissingNote(t *testing.T) {
	dir := initTestRepo(t)
	hash := headHash(t, dir)

	got, err := notes.Read(dir, hash)
	if err != nil {
		t.Fatalf("Read on missing note: expected nil error, got %v", err)
	}
	if got != nil {
		t.Fatalf("Read on missing note: expected nil, got %+v", got)
	}
}

func TestExists(t *testing.T) {
	dir := initTestRepo(t)
	hash := headHash(t, dir)

	if notes.Exists(dir, hash) {
		t.Fatal("Exists: expected false before write")
	}

	note := &notes.CommitNote{SchemaVersion: "1.0", CommitHash: hash}
	if err := notes.Write(dir, hash, note); err != nil {
		t.Fatalf("Write: %v", err)
	}

	if !notes.Exists(dir, hash) {
		t.Fatal("Exists: expected true after write")
	}
}

func TestOverwrite(t *testing.T) {
	dir := initTestRepo(t)
	hash := headHash(t, dir)

	note1 := &notes.CommitNote{Intent: "first intent"}
	if err := notes.Write(dir, hash, note1); err != nil {
		t.Fatalf("first Write: %v", err)
	}

	note2 := &notes.CommitNote{Intent: "second intent"}
	if err := notes.Write(dir, hash, note2); err != nil {
		t.Fatalf("second Write: %v", err)
	}

	got, err := notes.Read(dir, hash)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if got.Intent != "second intent" {
		t.Errorf("Intent: got %q, want %q", got.Intent, "second intent")
	}
}

func TestTestResults(t *testing.T) {
	dir := initTestRepo(t)
	hash := headHash(t, dir)

	note := &notes.CommitNote{
		TestResults: &notes.TestResults{
			Passed:  42,
			Failed:  0,
			Skipped: 3,
			Command: "go test ./...",
		},
	}
	if err := notes.Write(dir, hash, note); err != nil {
		t.Fatalf("Write: %v", err)
	}

	got, err := notes.Read(dir, hash)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if got.TestResults == nil {
		t.Fatal("TestResults: got nil")
	}
	if got.TestResults.Passed != 42 {
		t.Errorf("TestResults.Passed: got %d, want 42", got.TestResults.Passed)
	}
	if got.TestResults.Command != "go test ./..." {
		t.Errorf("TestResults.Command: got %q", got.TestResults.Command)
	}
}
