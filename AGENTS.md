# AGENTS.md

> **AGENTS.md** is a convention for AI-native repositories. Place this file in your repo root to tell any AI coding agent how to make context-preserving commits. Inspired by CLAUDE.md but designed to be tool-agnostic.

## What Is This?

This repo uses **[agit](https://github.com/madhurm/agit)** — a git wrapper that preserves AI agent reasoning in git notes. Instead of a blank commit, every agent commit stores *why* the change was made, what was tried, what risks were identified, and what the agent wasn't sure about.

When another agent session opens this repo, it can read `agit log` and `agit context show` to understand the full history of decisions — not just what changed.

## Setup (One Time)

```bash
# Install agit
go install github.com/madhurm/agit@latest

# Initialize in this repo
agit init
```

## How to Commit

**Instead of `git commit`, always use `agit commit`.**

### Minimal commit (required flags only)
```bash
agit commit -m "feat: add user authentication"
```

### Full commit (recommended for meaningful changes)
```bash
agit commit \
  -m "feat: add JWT authentication middleware" \
  --task "Add stateless auth to protect API endpoints" \
  --intent "Implement JWT to avoid session storage across pods" \
  --confidence 0.82 \
  --confidence-rationale "Core logic tested, token refresh path not yet covered" \
  --tried "session-based: rejected — needs shared Redis store" \
  --tried "OAuth-only: rejected — internal service-to-service calls also need auth" \
  --risk "high:token-expiry:short-lived tokens cause UX issues if refresh not implemented" \
  --risk "medium:key-rotation:no key rotation strategy yet" \
  --ripple "middleware/auth.go — needs updating for new token shape" \
  --ripple "all protected route handlers — need middleware applied" \
  --unknowns "token revocation strategy not decided" \
  --unknowns "rate limiting on /auth/token not planned"
```

### Pass context as JSON (for agents that build metadata programmatically)
```bash
agit commit -m "refactor: extract payment service" --json-note context.json
```

Where `context.json` follows the [CommitNote schema](https://github.com/madhurm/agit#schema).

## Flag Reference

| Flag | Type | Description |
|------|------|-------------|
| `-m` | string | **Required.** Commit message |
| `--task` | string | The original human task/prompt that triggered this work |
| `--intent` | string | What the agent was trying to accomplish |
| `--confidence` | float | Agent's confidence: 0.0 (guessing) to 1.0 (certain) |
| `--confidence-rationale` | string | Why this confidence level |
| `--tried` | repeatable | `"approach: rejected-reason"` — alternatives considered and why rejected |
| `--risk` | repeatable | `"severity:area:description"` — severity: low/medium/high |
| `--ripple` | repeatable | Files/areas not changed but that may be affected |
| `--unknowns` | repeatable | Things the agent wasn't sure about |
| `--agent-id` | string | Agent identifier (auto-detected from env if set) |
| `--agent-model` | string | Model name (auto-detected from env if set) |
| `--session-id` | string | Agent session ID for tracing |
| `--json-note` | string | Path to JSON file with full CommitNote (or `-` for stdin) |

## Reading Context

```bash
# See recent commits with agent metadata
agit log

# See full context for a specific commit
agit context show HEAD
agit context show abc1234

# Machine-readable output (for agents reading history)
agit log --json
agit context show HEAD --json
```

## Environment Variables

Set these so agit automatically fills agent metadata without flags:

| Variable | Used by |
|----------|---------|
| `CLAUDE_AGENT_ID` | Claude Code |
| `CLAUDE_MODEL` | Claude Code |
| `CLAUDE_SESSION_ID` | Claude Code |
| `COPILOT_AGENT_ID` | GitHub Copilot |
| `CURSOR_AGENT_ID` | Cursor |
| `AIDER_MODEL` | Aider |
| `DEVIN_AGENT_ID` | Devin |

agit reads `--agent-id` / `--agent-model` / `--session-id` from these env vars automatically if the flags aren't passed explicitly.

## Supported AI Tools

agit works with any AI coding agent that can run shell commands:

- **Claude Code** — use `--agent-id $CLAUDE_AGENT_ID --agent-model $CLAUDE_MODEL`
- **GitHub Copilot / Copilot Workspace** — set `COPILOT_AGENT_ID` in environment
- **Cursor** — configure in `.cursor/tools.json` to use `agit commit`
- **Aider** — use `--commit-cmd "agit commit"` flag
- **Continue.dev** — add to custom actions
- **Devin** — set as default commit tool in repo settings
- **Any other agent** — if it can run `agit commit`, it works

## Schema

The full JSON schema stored in each git note:

```json
{
  "schema_version": "1.0",
  "commit_hash": "<sha>",
  "agent": {
    "id": "claude-code",
    "model": "claude-sonnet-4-6",
    "session_id": "abc123"
  },
  "task": "Original human prompt that triggered this session",
  "intent": "What the agent was trying to accomplish",
  "confidence": 0.85,
  "confidence_rationale": "Why this confidence level",
  "alternatives_considered": [
    { "approach": "approach name", "rejected_reason": "why rejected" }
  ],
  "key_decisions": [
    { "decision": "what was decided", "rationale": "why" }
  ],
  "risks": [
    { "area": "component", "description": "risk description", "severity": "low|medium|high" }
  ],
  "context_consulted": ["file.go:42", "README.md", "existing auth middleware"],
  "test_results": { "passed": 42, "failed": 0, "skipped": 3, "command": "go test ./..." },
  "ripple_effects": ["files not changed but may be affected"],
  "unknowns": ["things the agent wasn't sure about"]
}
```

Notes are stored in `refs/notes/agit` — a standard git ref, no extra files in your repo.

## Why This Matters

Without AGENTS.md / agit:
- Agent commits look like any human commit
- Future agents (and humans) have no idea *why* decisions were made
- Code review of agent PRs is guesswork
- Debugging agent mistakes requires reading the full session transcript

With AGENTS.md / agit:
- Every commit carries full reasoning context
- New agents can read `agit log` to understand project history
- Code reviewers see intent, risks, and unknowns in the PR comment (via [agit GitHub App](https://github.com/marketplace/agit))
- Agent debugging is dramatically faster

---

*AGENTS.md is an open convention. Copy this template to your own repos. The spec lives at [github.com/madhurm/agit](https://github.com/madhurm/agit).*
