# CLAUDE.md — agit

This file tells Claude Code how to work in the agit repository.

## Commits

**Always use `agit commit` instead of `git commit`.** agit dog-foods itself — this repo's entire history is built with agit.

```bash
agit commit \
  -m "feat: short imperative description" \
  --task "$ORIGINAL_HUMAN_TASK" \
  --intent "what you were trying to accomplish" \
  --confidence 0.85 \
  --confidence-rationale "why this confidence level" \
  --tried "approach1: why rejected" \
  --risk "severity:area:description" \
  --ripple "files not changed but may be affected" \
  --unknowns "things you were uncertain about"
```

See [AGENTS.md](AGENTS.md) for full flag reference. AGENTS.md is the authoritative guide — CLAUDE.md just points here.

## Build

```bash
cd /home/madhur/agit
go build -o agit .
```

## Test

```bash
go test ./...
```

## Run

```bash
./agit --help
./agit log
./agit context show HEAD
```

## Environment

Set these so `agit commit` auto-fills agent metadata:

```bash
export CLAUDE_AGENT_ID="claude-code"
export CLAUDE_MODEL="claude-sonnet-4-6"  # or whatever model is running
export CLAUDE_SESSION_ID="<current session id>"
```

## Architecture

- `cmd/` — Cobra CLI commands (commit, context, log, init)
- `internal/git/` — Git CLI wrapper (RunGit, HeadHash, Log, StageAll, etc.)
- `internal/notes/` — CommitNote schema, git notes read/write
- `github-app/` — TypeScript Probot app for GitHub PR comments
- `scripts/gen.py` — Code generation via local Ollama/GLM model
- `prompts/` — Prompts used for GLM code generation

## Key Patterns

- `git.RunGit(dir, args...)` returns `(string, error)` — always `_, err := git.RunGit(...)` not `err :=`
- `notes.Read()` returns `(*CommitNote, nil)` if no note — check for nil before accessing fields
- `AgentInfo` is a struct (not pointer) — check `note.Agent.ID != ""` not `note.Agent != nil`
- `commit.Timestamp` is `time.Time` — use `.Format(...)` directly

## GLM Codegen Workflow

For boilerplate implementation (not architecture):

```bash
# Write a prompt to prompts/XX_name.txt first
python3 scripts/gen.py prompts/XX_name.txt output/path.go
# Then review, fix cross-package refs, build
```

GLM generates ~70% correct code. Common fixes needed: field names, package imports, function signatures. Always verify against actual struct/function definitions before accepting.
