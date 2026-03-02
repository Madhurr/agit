package git_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/madhurm/agit/internal/git"
)

// initTestRepo creates a temp dir, initializes a git repo with one commit, returns the dir.
func initTestRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	run := func(args ...string) string {
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
		result := string(out)
		if len(result) > 0 && result[len(result)-1] == '\n' {
			result = result[:len(result)-1]
		}
		return result
	}

	run("git", "init", "-b", "main")
	run("git", "config", "user.email", "test@test.com")
	run("git", "config", "user.name", "Test")

	if err := os.WriteFile(filepath.Join(dir, "hello.txt"), []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	run("git", "add", ".")
	run("git", "commit", "-m", "first commit")

	return dir
}

func TestIsRepo(t *testing.T) {
	dir := initTestRepo(t)

	if !git.IsRepo(dir) {
		t.Errorf("IsRepo(%q): expected true for initialized git repo", dir)
	}

	notRepo := t.TempDir()
	if git.IsRepo(notRepo) {
		t.Errorf("IsRepo(%q): expected false for non-git dir", notRepo)
	}
}

func TestHeadHash(t *testing.T) {
	dir := initTestRepo(t)

	hash, err := git.HeadHash(dir)
	if err != nil {
		t.Fatalf("HeadHash: %v", err)
	}
	if len(hash) != 40 {
		t.Errorf("HeadHash: got %q (len %d), want 40-char SHA", hash, len(hash))
	}
}

func TestLog(t *testing.T) {
	dir := initTestRepo(t)

	// Add a second commit
	cmd := exec.Command("git", "config", "user.email", "test@test.com")
	cmd.Dir = dir
	cmd.Run()
	if err := os.WriteFile(filepath.Join(dir, "world.txt"), []byte("world"), 0644); err != nil {
		t.Fatal(err)
	}
	c := exec.Command("git", "add", ".")
	c.Dir = dir
	c.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=Test",
		"GIT_AUTHOR_EMAIL=test@test.com",
		"GIT_COMMITTER_NAME=Test",
		"GIT_COMMITTER_EMAIL=test@test.com",
	)
	c.Run()
	c2 := exec.Command("git", "commit", "-m", "second commit")
	c2.Dir = dir
	c2.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=Test",
		"GIT_AUTHOR_EMAIL=test@test.com",
		"GIT_COMMITTER_NAME=Test",
		"GIT_COMMITTER_EMAIL=test@test.com",
	)
	c2.Run()

	entries, err := git.Log(dir, 10)
	if err != nil {
		t.Fatalf("Log: %v", err)
	}
	if len(entries) < 2 {
		t.Fatalf("Log: got %d entries, want at least 2", len(entries))
	}

	// Most recent commit first
	if entries[0].Subject != "second commit" {
		t.Errorf("Log[0].Subject: got %q, want %q", entries[0].Subject, "second commit")
	}
	if entries[1].Subject != "first commit" {
		t.Errorf("Log[1].Subject: got %q, want %q", entries[1].Subject, "first commit")
	}
	if len(entries[0].Hash) != 40 {
		t.Errorf("Log[0].Hash: got %q (len %d), want 40-char SHA", entries[0].Hash, len(entries[0].Hash))
	}
	if len(entries[0].ShortHash) != 7 {
		t.Errorf("Log[0].ShortHash: got %q (len %d), want 7 chars", entries[0].ShortHash, len(entries[0].ShortHash))
	}
	if entries[0].Timestamp.IsZero() {
		t.Error("Log[0].Timestamp: got zero time, want non-zero")
	}
}

func TestLogCount(t *testing.T) {
	dir := initTestRepo(t)

	entries, err := git.Log(dir, 1)
	if err != nil {
		t.Fatalf("Log: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("Log(dir, 1): got %d entries, want 1", len(entries))
	}
}

func TestRepoRoot(t *testing.T) {
	dir := initTestRepo(t)

	root, err := git.RepoRoot(dir)
	if err != nil {
		t.Fatalf("RepoRoot: %v", err)
	}
	if root == "" {
		t.Error("RepoRoot: got empty string")
	}
	// root should contain the temp dir path
	if len(root) == 0 {
		t.Error("RepoRoot: empty")
	}
}

func TestStageAll(t *testing.T) {
	dir := initTestRepo(t)

	// Create a new file
	if err := os.WriteFile(filepath.Join(dir, "new.txt"), []byte("new"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := git.StageAll(dir); err != nil {
		t.Fatalf("StageAll: %v", err)
	}

	// Verify the file is staged
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git diff --cached: %v", err)
	}
	if string(out) == "" {
		t.Error("StageAll: no files staged")
	}
}
