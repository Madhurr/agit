# agit

Git commits tell you *what* changed. agit tells you *why*.

`agit commit` wraps `git commit` and attaches structured reasoning — intent, confidence, risks, alternatives tried, unknowns — as [git notes](https://git-scm.com/docs/git-notes). No extra files. No database. Just git.

```bash
agit commit \
  -m "feat: add JWT auth middleware" \
  --intent "Stateless auth to avoid session storage across pods" \
  --confidence 0.82 \
  --tried "session-based: needs shared Redis — rejected" \
  --risk "high:token-expiry:refresh not implemented" \
  --unknowns "token revocation strategy not decided"
```

```
$ agit log

abc1234  feat: add JWT auth middleware  12 Mar 2026  [82%] [risk:high]
                · madhur
intent: Stateless auth to avoid session storage across pods
tried:  session-based, OAuth-only
```

[![Go 1.22](https://img.shields.io/badge/go-1.22-blue)](https://golang.org) [![MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

---

## Why

AI coding agents write millions of commits daily. Every one looks like this:

```
abc1234  feat: add user authentication
def5678  fix: resolve payment edge case
```

The reasoning is gone the moment the session ends. Next session, next developer, next agent — nobody knows *why* JWT was chosen over sessions, what alternatives were considered, or what risks were flagged.

agit fixes this by storing agent reasoning alongside every commit, using a git feature that's been stable since 2010.

---

## Install

```bash
# Go
go install github.com/madhurm/agit@latest

# Binary (Linux/macOS/Windows)
# → github.com/Madhurr/agit/releases
```

## Setup

```bash
cd your-repo
agit init    # configures notes fetch/rebase
```

That's it. Now use `agit commit` instead of `git commit`.

---

## Commands

**`agit commit`** — commit with reasoning attached

```bash
agit commit -m "message" \
  --intent "what you were trying to do" \
  --confidence 0.85 \
  --confidence-rationale "tests pass, edge cases uncovered" \
  --tried "approach: why rejected" \        # repeatable
  --risk "severity:area:description" \      # repeatable
  --unknowns "things you're unsure about" \ # repeatable
  --ripple "files affected but not changed"  # repeatable
```

Agent metadata (`--agent-id`, `--agent-model`, `--session-id`) auto-fills from environment variables. Pass `--json-note <file>` to pipe in structured JSON instead of flags.

**`agit log`** — git log with reasoning inline

```
$ agit log -n 5

d73d771  feat: implement agit diff  8 Mar 2026  [95%] [risk:low]
                · madhur
intent: Semantic diff between commits for reasoning evolution
tried:  text diff of JSON, stored diffs

2f5c0ee  fix: disable brew tap  7 Mar 2026  [99%]
                · madhur
intent: Unblock binary releases
```

`agit log --json` for machine-readable output.

**`agit context show`** — full reasoning for a commit

```
$ agit context show d73d771

Task:     Implement agit diff for comparing agent reasoning
Intent:   Semantic diff between any two commits
Confidence: 95% — "14 tests pass, CLI verified end-to-end"
Agent:    anton (claude-opus-4-6)

Alternatives rejected:
  • Text diff of JSON — loses semantic meaning
  • Stored diffs — unnecessary complexity

Risks:
  [low] merge-commits: multiple parents may confuse HEAD~1

Unknowns:
  • Should GitHub PR comments also show diffs?

Ripple effects:
  • formatter.ts — could add diff view
```

`agit context show --json` for the raw CommitNote.

**`agit diff`** — how reasoning evolved between commits

```
$ agit diff 2f5c0ee..d73d771

12 change(s) in agent context

Changed:
  Confidence: 99% → 95%
  Agent: claude-code → anton

Added:
  + 2 alternatives rejected
  + 2 new risks identified

Resolved:
  ✓ Risk resolved: [medium] github-api
  ✓ Unknown resolved: API rate limits
```

Supports `agit diff`, `agit diff <hash>`, `agit diff <from>..<to>`, `agit diff <from> <to>`. Add `--json` or `--files`.

**`agit init`** — one-time repo setup

Configures `git fetch` to pull `refs/notes/agit` and `notes.rewriteRef` to carry notes through rebase.

---

## PR Context on GitHub

Drop this workflow in `.github/workflows/agit-pr-context.yml` and every PR gets an automatic comment with full agent reasoning:

```yaml
name: agit PR Context
on:
  pull_request:
    types: [opened, synchronize]
permissions:
  contents: read
  pull-requests: write
```

[Full workflow →](.github/workflows/agit-pr-context.yml)

Result on your PR:

> ### 🤖 Agent Context
> **1/2 commits have agent context**
>
> | Confidence | 🟢 95% — "All tests pass, CLI verified" |
> | Agent | anton (claude-opus-4-6) |
>
> **Intent:** Semantic diff between commits
>
> **Risks:** 🟢 [low] merge-commits — multiple parents
>
> **❓ Unknowns:** PR formatter diff view

No app registration. No webhook hosting. Just a workflow file.

A self-hosted [Probot app](github-app/) is also included for orgs that need it at scale.

---

## How it works

agit stores metadata in **git notes** under `refs/notes/agit`. Each note is a JSON blob keyed by commit SHA.

```
refs/notes/agit
  └── <commit-sha> → JSON
```

- No files added to your repo
- Push/fetch like any ref: `git push origin refs/notes/agit`
- Survives rebase (with `notes.rewriteRef` configured by `agit init`)
- Works offline — everything is local git
- Standard feature since git 1.6.6 (2010)

### Schema

```json
{
  "schema_version": "1.0",
  "commit_hash": "abc1234...",
  "agent": { "id": "claude-code", "model": "claude-sonnet-4-6", "session_id": "..." },
  "task": "Original task",
  "intent": "What the agent was trying to do",
  "confidence": 0.85,
  "confidence_rationale": "Why this confidence",
  "alternatives_considered": [{ "approach": "...", "rejected_reason": "..." }],
  "key_decisions": [{ "decision": "...", "rationale": "..." }],
  "risks": [{ "area": "...", "description": "...", "severity": "low|medium|high" }],
  "test_results": { "passed": 42, "failed": 0, "skipped": 3, "command": "go test ./..." },
  "ripple_effects": ["files affected but not changed"],
  "unknowns": ["things the agent wasn't sure about"]
}
```

---

## Agent integration

agit works with any tool that can run shell commands.

**Claude Code** — set env vars, use `agit commit`:
```bash
export AGIT_AGENT_ID="claude-code"
export AGIT_MODEL="claude-sonnet-4-6"
```

**Aider** — `aider --commit-cmd "agit commit"`

**Cursor / Copilot / Continue.dev / Devin** — if it has a terminal, it can use agit.

**Any repo** — add [AGENTS.md](AGENTS.md) to tell agents to use `agit commit`. It's a convention, like `.editorconfig` but for AI.

---

## Project status

Done:
- `agit commit` — metadata capture via flags or JSON
- `agit log` — semantic history with inline notes
- `agit context show` — full reasoning display
- `agit diff` — semantic diff between commits
- `agit init` — repo setup
- GitHub Actions workflow — PR comments
- Probot GitHub App — self-hosted alternative
- GoReleaser — Linux/macOS/Windows binaries

Next:
- `agit blame` — reasoning per line
- VS Code extension
- GitLab CI integration

---

## Contributing

```bash
git clone https://github.com/Madhurr/agit
cd agit
go build -o agit .
go test ./...

# Contribute using agit (naturally)
agit commit -m "feat: ..." --intent "..." --confidence 0.9
```

---

## License

MIT

---

Built on git notes (stable since 2010). Plain JSON. No vendor lock-in. No cloud dependency. No API keys.
