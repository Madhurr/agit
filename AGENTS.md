# AGENTS.md

This repo uses [agit](https://github.com/Madhurr/agit). Use `agit commit` instead of `git commit`.

## Setup

```bash
go install github.com/Madhurr/agit@v0.3.0
agit init
```

## How to commit

```bash
agit commit \
  -m "feat: add rate limiting" \
  --intent "..." \
  --confidence 0.85 \
  --tried "approach: reason" \
  --risk "severity:area:description" \
  --unknowns "..." \
  --ripple "..."
```

## What to write in each field

**`--intent`** — *why* you made this change, not what it does. The commit message says what. Intent says why.
- ❌ `"implement rate limiting middleware"`
- ✅ `"protect the API from abuse without adding a Redis dependency"`

**`--confidence`** — how sure you are the implementation is correct. Be honest.
- `0.95+` only if you've tested it and there are no known gaps
- `0.70–0.85` for most changes: logic is sound but edge cases exist
- `< 0.70` if you're making assumptions or skipping something important
- **Do not default to 0.9 every time.** Calibrate it.

**`--confidence-rationale`** — one sentence on why you gave that score. What did you verify? What didn't you test?
- `"core path tested, token refresh not covered"`
- `"logic is straightforward but not benchmarked under load"`

**`--tried`** — other approaches you considered and rejected. Format: `"approach: why rejected"`.
- `"session cookies: needs shared store, won't work across pods"`
- `"Redis rate limiter: adds infra dependency we don't have"`
- Only include things you actually considered. Don't fabricate alternatives.

**`--risk`** — what could go wrong. Format: `"severity:area:description"`. Severity: `low`, `medium`, `high`.
- `"high:data-loss:rollback deletes rows, no soft delete"`
- `"medium:performance:N+1 query on user lookup, not optimized yet"`
- `"low:ux:error message is generic, could be more specific"`
- Flag real risks. Don't write `"low:none:no risks"`.

**`--unknowns`** — things you weren't sure about and didn't resolve.
- `"token revocation strategy not decided"`
- `"unclear if this needs to be idempotent"`
- `"haven't checked behavior when user has no permissions"`
- This is not a weakness. It's useful information.

**`--ripple`** — files or systems you didn't touch but that may be affected.
- `"middleware/auth.go — expects old token shape"`
- `"mobile clients — refresh flow will break if token TTL changes"`

## Environment variables

```bash
export AGIT_AGENT_ID="your-agent-name"
export AGIT_MODEL="model-name"
export AGIT_SESSION_ID="session-id"   # optional
```

## Reading the history

```bash
agit log                    # recent commits with reasoning inline
agit context show HEAD      # full reasoning for a commit
agit diff HEAD~3..HEAD      # how reasoning evolved
agit log --json             # machine-readable, for agents reading history
```

Before starting work on this repo, run `agit log` to understand past decisions.
