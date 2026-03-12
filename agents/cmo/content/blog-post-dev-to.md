# Blog Post: dev.to / Hashnode / personal blog

**Title:** AI coding agents commit code. Their reasoning doesn't survive the session. Here's how to fix it.

**Tags:** ai, git, developer-tools, claude, agents

**Cover image suggestion:** Terminal screenshot showing `agit context show HEAD` output with the confidence/unknowns fields visible

---

## The problem nobody's talking about

Your AI coding agent just submitted a PR. 847 lines changed across 14 files. The commit messages say:

```
feat: refactor auth layer
fix: update API contracts
chore: update dependencies
```

You're reviewing the PR. You have questions:

- Why JWT instead of session cookies?
- The auth refactor touches the mobile API — is that intentional?
- That database query in `user_service.go` looks risky for large tables. Did the agent consider that?
- Is the token refresh flow actually implemented, or is that a TODO?

The agent knew the answers to all of these when it wrote the code. It doesn't anymore. The session ended. You're left with a diff and no context.

This isn't a hypothetical. It's the current state of every team using Claude Code, GitHub Copilot Workspace, Cursor, or Codex to write production code.

## What gets lost when a session ends

When an AI agent makes a commit, it knows:

1. **Why this approach** — not just what it did, but why it chose this over alternatives it considered and rejected
2. **Confidence level** — and why it's not 100% confident
3. **What might break** — downstream effects it identified but may not have fixed
4. **What it didn't know** — open questions, unresolved decisions, areas of uncertainty
5. **What it tried** — approaches it explored and discarded before landing here

None of this is in the commit. Commit messages are prose — they can carry some of this, but not in a form the *next* agent session can parse programmatically. And the next session starts from zero anyway.

## The solution: git notes

Git has had a notes feature since 2010. Notes let you attach arbitrary data to commits without modifying the commits themselves. They're stored in a separate ref (`refs/notes/*`) and travel with your repo via push/fetch.

Most developers have never used them because there wasn't a compelling reason to. AI agents give us one.

I built **agit** to make this easy. It wraps `git commit` and stores agent reasoning as a structured JSON note on each commit.

## How it works

Instead of:
```bash
git commit -m "feat: switch auth to JWT"
```

You (or your agent) runs:
```bash
agit commit \
  -m "feat: switch auth to JWT" \
  --intent "stateless sessions, no Redis dependency" \
  --confidence 0.82 \
  --tried "session cookies: needs shared store across pods" \
  --risk "high:token-refresh:not implemented yet" \
  --unknowns "revocation strategy undecided"
```

That's it. The commit is created normally. The metadata is stored as a git note.

## What you get back

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

```
$ agit log

a317512  feat: switch auth to JWT      8 Mar 2026  [82%] [risk:high]
                · dev
intent: stateless sessions, no Redis dependency
tried:  session cookies
```

```
$ agit diff HEAD~5..HEAD

Changed:
  Confidence: 68% → 91%

Resolved:
  ✓ Risk resolved: [high] token-refresh
  ✓ Unknown resolved: revocation strategy

Added:
  + Risk: [medium] race condition in worker pool
```

## The field that matters most: unknowns

The schema has the usual suspects — intent, confidence, risks, alternatives. But the field I keep coming back to is `unknowns`.

Most "what the AI did" tooling tells you what it did. `unknowns` tells you what it *didn't know*. That's different. An agent with 95% confidence and 0 unknowns is different from an agent with 95% confidence and 3 unknowns. The second one is a time bomb.

Making unknowns explicit is what surfaces the gaps that will hurt you later. "Revocation strategy undecided" is right there in the commit. Nobody can say they didn't know.

## Using it with your agent

**Claude Code** — add one line to `CLAUDE.md`:
```
Use agit commit instead of git commit. Read AGENTS.md for flags.
```

**Aider:**
```bash
aider --commit-cmd "agit commit -m"
```

**Cursor / Copilot / Continue / Devin** — if it has a terminal, it runs `agit commit`. Add `AGENTS.md` to your repo and most agents will read it automatically.

**PR comments** — drop `.github/workflows/agit-pr-context.yml` in your repo. Every PR automatically gets an Agent Context comment showing the reasoning for every commit in the branch. No app registration, no webhooks.

## Why not just put this in the commit message?

Three reasons:

1. **Structure.** The next agent session needs to parse this programmatically. `agit log --json` outputs machine-readable context. You can't parse that reliably from prose.

2. **Separation.** Commit messages are for humans reading `git log`. Agent reasoning metadata is for agents reading context. They serve different audiences with different needs.

3. **Volume.** The full context for a complex commit — alternatives considered, confidence rationale, test results, ripple effects — is too much for a commit message. It would swamp `git log`.

## Install

```bash
go install github.com/Madhurr/agit@v0.3.0
cd your-repo && agit init
```

Pre-built binaries for Linux, macOS, and Windows at [github.com/Madhurr/agit/releases](https://github.com/Madhurr/agit/releases).

## Why this matters now

AI coding agents are becoming a normal part of software development. The existing git tooling wasn't designed for a world where non-human agents make commits. `git log` gives you author, date, diff, and message. For human commits, that's usually enough context. For agent commits, it's almost never enough.

The reasoning gap is a genuine problem and it will get worse as agents handle more complex work. agit is my attempt at a minimal, zero-infrastructure fix — leverage what git already has, add structure to what agents already know, make it retrievable.

---

**GitHub:** [github.com/Madhurr/agit](https://github.com/Madhurr/agit) — MIT, contributions welcome.
