# agit — Marketing Drafts

## Hacker News — Show HN

**Title:** Show HN: agit – Store AI agent reasoning in git notes alongside commits

**Text:**

I've been using AI coding agents (Claude Code, Cursor) daily for months. The biggest problem isn't the code quality — it's that the reasoning disappears.

An agent makes a commit: "feat: add JWT auth middleware." But *why* JWT over sessions? What else was considered? What risks were flagged? That context is gone the moment the session ends.

agit is a `git commit` wrapper that stores structured reasoning in git notes — intent, confidence, alternatives tried, risks, unknowns. It's a standard git ref (`refs/notes/agit`), so it travels with push/fetch. No extra files, no database.

```
agit commit -m "feat: add JWT auth" \
  --intent "Stateless auth to avoid session storage across pods" \
  --confidence 0.82 \
  --tried "session-based: needs shared Redis — rejected" \
  --risk "high:token-expiry:refresh not implemented"
```

Then `agit log` shows commits with inline reasoning, `agit diff` shows how reasoning evolved between commits, and a GitHub Actions workflow auto-posts agent context on PRs.

Written in Go. Works with any agent that can run shell commands. ~1200 lines of code.

https://github.com/Madhurr/agit

---

## Reddit — r/programming

**Title:** I built a git wrapper that stores AI agent reasoning alongside commits using git notes

**Post:**

Every AI coding agent commit I've seen looks like this:

```
abc1234  feat: add user authentication
def5678  fix: resolve payment edge case
```

No reasoning. No context. Two weeks later when something breaks, you have to read an entire session transcript to figure out what the agent was thinking.

I built `agit` — it wraps `git commit` and attaches structured metadata (intent, confidence score, alternatives tried, risks, unknowns) as git notes. It's a native git feature from 2010 that almost nobody uses.

The cool part: `agit diff` shows semantic diffs between commits — not what code changed, but how the *reasoning* changed. Confidence went up? Risk resolved? Unknown became known?

It also ships a GitHub Actions workflow that auto-posts agent context on PRs, so reviewers see the full reasoning without leaving GitHub.

Go binary, MIT licensed, works with Claude Code, Cursor, Aider, or anything that has a terminal.

https://github.com/Madhurr/agit

Would love feedback from anyone else dealing with the "why did the AI do this?" problem.

---

## Reddit — r/ClaudeAI

**Title:** Built a tool to preserve Claude Code's reasoning in git — so next session knows what was decided

**Post:**

One thing that bugs me about Claude Code: it makes great decisions, explains its reasoning... and then the session ends and all that context is gone.

Next session opens the same repo and has zero idea why JWT was chosen over sessions, or what risks the previous session flagged.

I built `agit` — a `git commit` wrapper that stores agent reasoning in git notes:

- Intent: what the agent was trying to do
- Confidence: how sure it was (0.0-1.0)
- Alternatives: what was tried and rejected
- Risks: what could go wrong
- Unknowns: what the agent wasn't sure about

Set `AGIT_AGENT_ID=claude-code` and `AGIT_MODEL=claude-sonnet-4-6` in your env, and Claude's metadata auto-fills.

Next Claude Code session can run `agit log --json` or `agit context show HEAD --json` to read the full history of reasoning — not just what files changed.

https://github.com/Madhurr/agit

---

## Twitter/X Thread

**Tweet 1:**
Your AI agent made 50 commits today.

Can you explain any of them?

I built agit — stores agent reasoning (intent, confidence, risks, alternatives tried) in git notes alongside every commit.

No extra files. No database. Just git.

github.com/Madhurr/agit

**Tweet 2:**
The problem:

```
abc1234 feat: add JWT auth
```

vs with agit:

```
abc1234 feat: add JWT auth [82%] [risk:high]
  intent: Stateless auth to avoid session storage
  tried: session-based (rejected), OAuth-only (rejected)
  risk: token refresh not implemented
```

**Tweet 3:**
agit diff shows how reasoning *evolved* between commits:

```
Confidence: 50% → 95%
Risk resolved: [high] token-expiry
Unknown resolved: revocation strategy
New decision: switched to RS256
```

Not what code changed. What the *thinking* changed.

**Tweet 4:**
Ships a GitHub Actions workflow — every PR gets an automatic Agent Context comment.

Zero config. Just add the workflow file.

No app registration. No webhooks. No hosting.

**Tweet 5:**
Works with any AI coding agent:
- Claude Code
- Cursor
- Aider
- Copilot
- Anything with a terminal

Go binary, MIT licensed, built on git notes (stable since 2010).

github.com/Madhurr/agit

---

## Dev.to Blog Post Title Ideas

1. "The Missing Layer in AI-Generated Code: Why Your Agent's Commits Need Reasoning"
2. "I Built a Git Wrapper That Makes AI Coding Agents Accountable"
3. "git notes: The 15-Year-Old Feature That Solves AI Coding's Biggest Problem"
