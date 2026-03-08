# agit

[![Go 1.22](https://img.shields.io/badge/go-1.22-00ADD8?logo=go&logoColor=white)](https://golang.org) [![MIT](https://img.shields.io/badge/license-MIT-blue)](LICENSE) [![Latest Release](https://img.shields.io/github/v/release/Madhurr/agit)](https://github.com/Madhurr/agit/releases)

**Git commits tell you *what* changed. agit tells you *why*.**

When an AI agent makes a commit, it knows things: why it chose this approach over alternatives, how confident it was, what might break, what it left unresolved. That context disappears when the session ends.

agit stores it in the commit — using [git notes](https://git-scm.com/docs/git-notes), a native git feature that travels with push and fetch.

---

![terminal demo](docs/demo.png)

![PR context](docs/pr-preview.gif)

---

## Install

```bash
go install github.com/Madhurr/agit@v0.3.0
```

Pre-built binaries (Linux / macOS / Windows) at [releases](https://github.com/Madhurr/agit/releases).

```bash
cd your-repo && agit init
```

### Verify a binary release

Every release is built on GitHub Actions, signed with [cosign](https://docs.sigstore.dev/) keyless signing, and attested with SLSA provenance. To verify:

```bash
# 1. Checksum
sha256sum --check --ignore-missing checksums.txt

# 2. Cosign signature (proves GitHub Actions built it, not a dev machine)
cosign verify-blob checksums.txt \
  --signature checksums.txt.sig \
  --certificate checksums.txt.pem \
  --certificate-identity-regexp "https://github.com/Madhurr/agit/.github/workflows/release.yml" \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com"

# 3. SLSA provenance
gh attestation verify agit_linux_amd64.tar.gz --repo Madhurr/agit
```

See [SECURITY.md](./SECURITY.md) for full details.

---


## Usage

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

```
$ agit context show HEAD

Intent:      stateless sessions, no Redis dependency
Confidence:  82% — "core path works, refresh not covered"
Agent:       claude-code (claude-sonnet-4-6)

Alternatives rejected:
  • session cookies — needs shared store across pods

Risks:
  [high] token-refresh: not implemented yet

Unknowns:
  • revocation strategy undecided
```

---

## Commands

| Command | Description |
|---|---|
| `agit commit` | commit with reasoning attached |
| `agit log` | history with inline reasoning |
| `agit context show [hash]` | full reasoning for a commit |
| `agit diff [from]..[to]` | how reasoning evolved between commits |
| `agit init` | one-time repo setup |

### `agit diff`

```
$ agit diff HEAD~3..HEAD

Changed:
  Confidence: 68% → 91%

Resolved:
  ✓ Risk resolved: [high] token-refresh
  ✓ Unknown resolved: revocation strategy

Added:
  + Risk: [medium] race condition in worker pool
```

---

## GitHub PR comments

Drop `.github/workflows/agit-pr-context.yml` ([full file](.github/workflows/agit-pr-context.yml)) in any repo. Every PR gets an Agent Context comment automatically — no app registration, no webhooks.

---

## How it works

Notes are stored in `refs/notes/agit`, keyed by commit SHA. Each note is a JSON blob.

```bash
# push notes with your code
git push origin refs/notes/agit

# fetch on clone
git fetch origin refs/notes/agit:refs/notes/agit
```

No extra files in the working tree. Works offline. Plain JSON. The storage has been in git since 2010.

<details>
<summary>Schema</summary>

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

Works with any tool that can run shell commands. Add [AGENTS.md](AGENTS.md) to your repo to instruct agents to use `agit commit`.

```bash
export AGIT_AGENT_ID="claude-code"
export AGIT_MODEL="claude-sonnet-4-6"
```

**Claude Code** — add to your `CLAUDE.md` or system prompt:
```
Use agit commit instead of git commit. Read AGENTS.md for flags and guidance.
```

**Aider**
```bash
aider --commit-cmd "agit commit -m"
```

**Cursor / Copilot / Continue / Devin** — if it has a terminal, it can run `agit commit`. Point it at `AGENTS.md`.

**Any agent reading this repo** — run `agit log --json` before starting work to understand prior decisions.

---

## Contributing

```bash
git clone https://github.com/Madhurr/agit
cd agit && go build -o agit . && go test ./...
```

---

MIT © [Madhur](https://github.com/Madhurr)
