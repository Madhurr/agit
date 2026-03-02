# agit — Agent-Native Git

> **The missing commit layer for AI coding agents.**
> agit preserves the *why* behind every AI commit — so the next agent session (or human) knows exactly what was decided, what was tried, and what risks were left open.

[![Go](https://img.shields.io/badge/go-1.22-blue)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

---

## The Problem

Every day, AI coding agents make millions of commits. Every one of those commits looks like this:

```
abc1234  feat: add user authentication
def5678  fix: resolve payment edge case
ghi9012  refactor: extract service layer
```

No reasoning. No alternatives considered. No risks flagged. No unknowns acknowledged.

**The context that matters — the *why* — is gone the moment the agent session ends.**

When you open a new session and ask "why did the agent choose JWT over sessions?", the answer is gone. When a PR reviewer looks at an agent commit, they see output with no reasoning. When a bug surfaces two weeks later, tracing it back to an agent decision requires reading an entire session transcript.

This is the gap. agit closes it.

---

## The Solution

agit is a `git commit` wrapper that stores full AI agent reasoning in git notes — a native git feature that travels with the repo.

```bash
# Instead of: git commit -m "feat: add JWT auth"
agit commit \
  -m "feat: add JWT auth middleware" \
  --intent "Stateless auth to avoid session storage across pods" \
  --confidence 0.82 \
  --confidence-rationale "Core logic tested, token refresh path incomplete" \
  --tried "session-based: needs shared Redis — rejected" \
  --risk "high:token-expiry:refresh not implemented" \
  --unknowns "token revocation strategy not decided" \
  --ripple "all protected route handlers need middleware applied"
```

Now `agit log` shows this:

```
abc1234  feat: add JWT auth middleware  12 Mar 2026  [82%] [risk:high]
                · madhur
intent: Stateless auth to avoid session storage across pods
tried:  session-based, OAuth-only
```

And `agit context show abc1234` shows the full picture:

```
Task:
  Add stateless auth to protect API endpoints

Intent:
  Stateless auth via JWT to avoid session storage across pods

Confidence:
  82% — "Core logic tested, token refresh path incomplete"

Alternatives rejected:
  • session-based — needs shared Redis store
  • OAuth-only — internal service-to-service calls also need auth

Risks:
  [high] token-expiry: short-lived tokens cause UX issues if refresh not implemented
  [medium] key-rotation: no rotation strategy yet

Unknowns:
  • token revocation strategy not decided
  • rate limiting on /auth/token not planned

Ripple effects (check these):
  • all protected route handlers need middleware applied
  • middleware/auth.go needs updating for new token shape
```

All of this is stored in `refs/notes/agit` — a standard git ref that travels with `git push/fetch`. No extra files. No databases. Just git.

---

## GitHub App — Agent Context in Every PR

The [agit GitHub App](https://github.com/marketplace/agit) reads `refs/notes/agit` and posts an **Agent Context** comment on every PR, giving reviewers full reasoning without leaving GitHub.

```
┌─────────────────────────────────────────────────────────────────┐
│ 🤖 Agent Context  ·  2/3 commits have agent metadata            │
│ ──────────────────────────────────────────────────────────────  │
│ abc1234  feat: add JWT auth middleware                           │
│                                                                  │
│ Confidence  🟢 82%  "refresh path not yet covered"              │
│ Agent       claude-sonnet-4-6  (session: s_abc123)              │
│ Intent      Stateless auth to avoid session storage             │
│                                                                  │
│ Alternatives rejected                                           │
│   ✗ session-based — needs shared Redis                          │
│   ✗ OAuth-only — internal calls also need auth                  │
│                                                                  │
│ ⚠ Risks                                                         │
│   🔴 [high] token-expiry — refresh not implemented              │
│   🟡 [medium] key-rotation — no strategy yet                    │
│                                                                  │
│ ❓ Unknowns                                                      │
│   · token revocation strategy not decided                       │
│   · rate limiting on /auth/token not planned                    │
│                                                                  │
│ 👀 Check these (ripple effects)                                  │
│   · all protected route handlers need middleware                 │
│   · middleware/auth.go — token shape changed                    │
└─────────────────────────────────────────────────────────────────┘
```

Works for every AI coding tool. Works on every GitHub repo. Requires only that agents use `agit commit`.

---

## Install

### Go install (fastest)
```bash
go install github.com/madhurm/agit@latest
```

### Homebrew
```bash
brew install madhurm/tap/agit
```

### Download binary
Download from [GitHub Releases](https://github.com/madhurm/agit/releases) for your platform (Linux/macOS/Windows, amd64/arm64).

---

## Quick Start

```bash
# 1. Initialize in your repo (configures git to push/fetch notes)
cd your-repo
agit init

# 2. Make your first context-aware commit
agit commit -m "feat: initial implementation" \
  --intent "What you were trying to do" \
  --confidence 0.9

# 3. See the semantic log
agit log

# 4. Inspect any commit
agit context show HEAD
agit context show HEAD --json   # machine-readable
```

---

## For AI Agents

Add an `AGENTS.md` file to your repo (see [AGENTS.md](AGENTS.md) in this repo for the template). It tells any AI coding agent how to use agit in your project.

### Claude Code
```bash
# Set in your shell or .env
export CLAUDE_AGENT_ID="claude-code"
export CLAUDE_MODEL="claude-sonnet-4-6"
export CLAUDE_SESSION_ID="$(uuidgen)"

# Now agit commit auto-fills agent metadata
agit commit -m "feat: ..." --intent "..." --confidence 0.85
```

### Aider
```bash
aider --commit-cmd "agit commit" ...
```

### GitHub Copilot Workspace
Configure `agit commit` as the commit command in workspace settings.

### Cursor / Continue.dev / Any other tool
If your agent can run shell commands, it can use `agit commit`. See [AGENTS.md](AGENTS.md) for the full flag reference.

---

## Commands

### `agit commit`
```
Usage: agit commit -m <message> [flags]

Flags:
  -m, --message string               Commit message (required)
      --task string                  Original human task/prompt
      --intent string                Agent's goal for this commit
      --confidence float             Confidence 0.0–1.0
      --confidence-rationale string  Why this confidence level
      --tried string                 Repeatable: "approach: rejected-reason"
      --risk string                  Repeatable: "severity:area:description"
      --unknowns string              Repeatable: things agent was unsure about
      --ripple string                Repeatable: files not changed but affected
      --agent-id string              Agent identifier
      --agent-model string           Model name
      --session-id string            Agent session ID
      --json-note string             Path to JSON metadata file (or - for stdin)
```

### `agit log`
```
Usage: agit log [-n <count>] [--json]

Shows recent commits with agit metadata inline.
--json outputs machine-readable JSON array (for agents reading history).
```

### `agit context show`
```
Usage: agit context show [hash] [--json]

Shows full agent context for a commit (defaults to HEAD).
--json outputs the raw CommitNote JSON.
```

### `agit init`
```
Usage: agit init

Configures git in the current repo to:
  - Fetch refs/notes/agit from origin
  - Carry agit notes through rebase/cherry-pick
```

---

## How Storage Works

agit stores metadata using **git notes** — a native git feature designed for exactly this: attaching metadata to commits without changing commit hashes.

```
refs/notes/agit
  └── <commit-sha>  → JSON CommitNote
```

- **No extra files** in your repo
- **Travels with the repo**: push/pull like any ref
  ```bash
  git push origin refs/notes/agit
  git fetch origin refs/notes/agit:refs/notes/agit
  ```
- **Survives rebase**: `agit init` configures `notes.rewriteRef` so notes follow rebased commits
- **Works offline**: everything is local git

---

## AGENTS.md — The Open Standard

agit introduces **AGENTS.md** — a convention for AI-native repositories. Drop it in your repo root to tell any AI coding agent how to commit with context in your project.

This repo itself has an [AGENTS.md](AGENTS.md). Copy the template to your own projects.

---

## Schema

The CommitNote JSON stored per commit:

```json
{
  "schema_version": "1.0",
  "commit_hash": "abc1234...",
  "agent": { "id": "claude-code", "model": "claude-sonnet-4-6", "session_id": "..." },
  "task": "Original human prompt",
  "intent": "Agent's goal",
  "confidence": 0.85,
  "confidence_rationale": "Why this confidence level",
  "alternatives_considered": [{ "approach": "...", "rejected_reason": "..." }],
  "key_decisions": [{ "decision": "...", "rationale": "..." }],
  "risks": [{ "area": "...", "description": "...", "severity": "low|medium|high" }],
  "context_consulted": ["file.go:42", "README.md"],
  "test_results": { "passed": 42, "failed": 0, "skipped": 3, "command": "go test ./..." },
  "ripple_effects": ["middleware/auth.go — may need updating"],
  "unknowns": ["token revocation strategy not decided"]
}
```

---

## Project Status

- [x] `agit commit` — full metadata capture
- [x] `agit log` — semantic history
- [x] `agit context show` — full context display
- [x] `agit init` — repo initialization
- [x] GitHub App — agent context in PR comments
- [x] GoReleaser — multi-platform binary releases
- [ ] `agit diff` — semantic diff between any two commits
- [ ] VS Code extension — inline context in editor
- [ ] GitLab App

---

## Contributing

```bash
git clone https://github.com/madhurm/agit
cd agit
go build -o agit .

# Run tests
go test ./...

# Make a contribution (using agit itself, naturally)
agit commit -m "feat: ..." --intent "..." --confidence 0.9
```

See [AGENTS.md](AGENTS.md) for how agents should contribute to this project.

---

## License

MIT — see [LICENSE](LICENSE)

---

*agit is built to last. git notes have been in git since 2010. The data format is plain JSON. No vendor lock-in, no cloud dependency, no API keys required. Just git.*
