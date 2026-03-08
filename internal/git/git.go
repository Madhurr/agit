// Package git provides a thin wrapper around the git CLI for agit operations.
package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// RunGit runs a git command in dir, returns stdout output trimmed.
// Returns error containing stderr if exit code != 0.
func RunGit(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%w: %s", err, strings.TrimSpace(stderr.String()))
	}

	return strings.TrimSpace(stdout.String()), nil
}

// RepoRoot returns absolute path to the git repo root containing dir.
// Calls: git rev-parse --show-toplevel
func RepoRoot(dir string) (string, error) {
	output, err := RunGit(dir, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	return output, nil
}

// HeadHash returns full 40-char SHA of HEAD.
// Calls: git rev-parse HEAD
func HeadHash(dir string) (string, error) {
	output, err := RunGit(dir, "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	return output, nil
}

// IsRepo returns true if dir is inside a git repository.
func IsRepo(dir string) bool {
	_, err := RepoRoot(dir)
	return err == nil
}

// CommitWithMessage runs: git add -u && git commit -m message
// Returns full SHA of the new commit (calls HeadHash after commit).
func CommitWithMessage(dir, message string) (string, error) {
	if _, err := RunGit(dir, "add", "-u"); err != nil {
		return "", err
	}
	if _, err := RunGit(dir, "commit", "-m", message); err != nil {
		return "", err
	}
	return HeadHash(dir)
}

// StageAll runs: git add -A
func StageAll(dir string) error {
	_, err := RunGit(dir, "add", "-A")
	return err
}

// LogEntry represents a single commit in git history.
type LogEntry struct {
	Hash      string
	ShortHash string
	Subject   string
	Author    string
	Timestamp time.Time
}

// Log returns last n commits.
// Parses: git log --format="%H|%h|%s|%an|%aI" -n <n>
// Split each line on "|" into exactly 5 fields.
// Parse Timestamp with time.Parse(time.RFC3339, field[4])
func Log(dir string, n int) ([]LogEntry, error) {
	output, err := RunGit(dir, "log", "--format=%H|%h|%s|%an|%aI", "-n", fmt.Sprintf("%d", n))
	if err != nil {
		return nil, err
	}

	lines := strings.Split(output, "\n")
	var entries []LogEntry

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.SplitN(line, "|", 5)
		if len(fields) != 5 {
			return nil, fmt.Errorf("invalid log entry format: %s", line)
		}

		timestamp, err := time.Parse(time.RFC3339, fields[4])
		if err != nil {
			return nil, fmt.Errorf("failed to parse timestamp: %w", err)
		}

		entries = append(entries, LogEntry{
			Hash:      fields[0],
			ShortHash: fields[1],
			Subject:   fields[2],
			Author:    fields[3],
			Timestamp: timestamp,
		})
	}

	return entries, nil
}

// DiffFiles returns file paths changed in commitHash vs its parent.
// Calls: git diff-tree --no-commit-id -r --name-only <commitHash>
// Returns slice of trimmed file path strings.
func DiffFiles(dir, commitHash string) ([]string, error) {
	output, err := RunGit(dir, "diff-tree", "--no-commit-id", "-r", "--name-only", commitHash)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(output, "\n")
	var files []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		files = append(files, line)
	}

	return files, nil
}
