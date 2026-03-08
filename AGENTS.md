# AGENTS.md

This repo uses [agit](https://github.com/Madhurr/agit) for commits. Use `agit commit` instead of `git commit`.

## Setup

```bash
go install github.com/madhurm/agit@latest
agit init
```

## Commit

```bash
agit commit \
  -m "feat: add rate limiting" \
  --intent "protect API from abuse without Redis dependency" \
  --confidence 0.85 \
  --tried "Redis-backed counters: adds infra dependency" \
  --risk "medium:memory:in-process counter resets on restart" \
  --unknowns "rate limit per user vs per IP not decided"
```

Flags: `--intent`, `--confidence`, `--confidence-rationale`, `--tried`, `--risk`, `--unknowns`, `--ripple`, `--task`, `--agent-id`, `--agent-model`

All repeatable flags (`--tried`, `--risk`, `--unknowns`, `--ripple`) can be used multiple times.

## Read

```bash
agit log                        # recent commits with inline reasoning
agit context show HEAD          # full reasoning for a commit
agit diff HEAD~3..HEAD          # how reasoning evolved
agit log --json                 # machine-readable
```

## Environment

```bash
export AGIT_AGENT_ID="claude-code"
export AGIT_MODEL="claude-sonnet-4-6"
export AGIT_SESSION_ID="..."    # optional, for tracing
```
