# agit

Git commits tell you *what* changed. agit tells you *why*.

![terminal](docs/demo.png)

![PR comment](docs/pr-preview.gif)

---

When an AI agent makes a commit, it knows a lot: why it chose this approach over others, how confident it was, what might break, what it left unresolved. That context disappears the moment the session ends.

agit stores it in the commit — using [git notes](https://git-scm.com/docs/git-notes), a native git feature that travels with push and fetch.

```bash
agit commit \
  -m "feat: switch auth to JWT" \
  --intent "stateless sessions, no Redis dependency" \
  --confidence 0.82 \
  --tried "session cookies: needs shared store across pods" \
  --risk "high:token-refresh:not implemented yet" \
  --unknowns "revocation strategy undecided"
```

```
$ agit log

a317512  feat: switch auth to JWT  8 Mar 2026  [82%] [risk:high]
                · dev
intent: stateless sessions, no Redis dependency
tried:  session cookies
```

[![Go 1.22](https://img.shields.io/badge/go-1.22-blue)](https://golang.org) [![MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

---

## Install

```bash
go install github.com/Madhurr/agit@v0.3.0
```

Pre-built binaries at [releases](https://github.com/Madhurr/agit/releases).

```bash
cd your-repo && agit init
```

---

## Commands

### `agit commit`

```
agit commit -m <message> [flags]

--intent string              why this change
--confidence float           0.0 to 1.0
--confidence-rationale string
--tried string               "approach: reason" — repeatable
--risk string                "severity:area:description" — repeatable
--unknowns string            repeatable
--ripple string              affected but not modified — repeatable
--agent-id / --agent-model / --session-id
--json-note string           path to full JSON instead of flags
```

Agent metadata is auto-filled from environment variables (`AGIT_AGENT_ID`, `AGIT_MODEL`, `AGIT_SESSION_ID`).

### `agit log`

Commit history with reasoning inline. `--json` for machine-readable output.

### `agit context show [hash]`

Full reasoning for a commit. Defaults to HEAD. `--json` for raw JSON.

### `agit diff [from] [to]`

How reasoning changed between commits.

```
$ agit diff HEAD~3..HEAD

Changed:
  Confidence: 68% → 91%

Added:
  + Risk: [medium] race condition in worker pool

Resolved:
  ✓ Risk resolved: [high] token-refresh
  ✓ Unknown resolved: revocation strategy
```

Supports `agit diff`, `agit diff <hash>`, `agit diff <from>..<to>`.

### `agit init`

One-time repo setup: configures git to fetch `refs/notes/agit` and carry notes through rebase.

---

## PR comments

Add `.github/workflows/agit-pr-context.yml` ([see it here](.github/workflows/agit-pr-context.yml)) to any repo. Every PR gets an Agent Context comment — full reasoning, no setup required.

---

## How it works

Notes live under `refs/notes/agit`, keyed by commit SHA.

```
refs/notes/agit
  └── <commit-sha>  →  JSON
```

```bash
# push notes with your code
git push origin refs/notes/agit

# fetch them when cloning
git fetch origin refs/notes/agit:refs/notes/agit
```

No extra files in the working tree. Works offline. Plain JSON. The storage mechanism has been in git since 2010.

<details>
<summary>Full schema</summary>

```json
{
  "schema_version": "1.0",
  "commit_hash": "...",
  "agent": { "id": "claude-code", "model": "claude-sonnet-4-6", "session_id": "..." },
  "task": "what the human asked for",
  "intent": "what the agent was trying to do",
  "confidence": 0.85,
  "confidence_rationale": "...",
  "alternatives_considered": [{ "approach": "...", "rejected_reason": "..." }],
  "key_decisions": [{ "decision": "...", "rationale": "..." }],
  "risks": [{ "area": "...", "description": "...", "severity": "low|medium|high" }],
  "test_results": { "passed": 42, "failed": 0, "skipped": 3, "command": "go test ./..." },
  "ripple_effects": ["..."],
  "unknowns": ["..."]
}
```

</details>

---

## Agent setup

Works with any tool that runs shell commands. Drop [AGENTS.md](AGENTS.md) in your repo — it tells agents to use `agit commit` and documents the flags.

```bash
# Claude Code, Cursor, Copilot, Aider, Devin, anything
export AGIT_AGENT_ID="claude-code"
export AGIT_MODEL="claude-sonnet-4-6"
```

---

## Contributing

```bash
git clone https://github.com/Madhurr/agit
cd agit && go build -o agit . && go test ./...
```

---

MIT
