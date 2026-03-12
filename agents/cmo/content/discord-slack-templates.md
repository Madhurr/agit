# Discord / Slack Community Messages

## Aider Discord (#general or #tools)

Hey all — built something that might be useful for Aider users.

When Aider commits, it knows why it made each decision, what alternatives it considered, how confident it was, and what's unresolved. That all disappears when the session ends.

Built **agit** to capture it as structured git notes. Works natively with Aider:

```bash
aider --commit-cmd "agit commit -m"
```

Then `agit context show HEAD` shows the full reasoning. `agit log` shows history with inline confidence and risks. `agit diff HEAD~5..HEAD` shows how reasoning evolved across commits.

Notes are stored in `refs/notes/agit` — no working tree pollution, travel with push/fetch.

GitHub: https://github.com/Madhurr/agit

---

## Continue.dev Discord (#show-and-tell or #tools)

If you're using Continue with local models to write code — the context you're building up in the session is gone when you close it. The commits don't carry any of it.

Built a wrapper called **agit** that stores agent reasoning as git notes. Continue can use it wherever it runs `git commit` (terminal mode, slash commands, etc.).

Key fields: intent, confidence, alternatives considered, risks, unknowns, ripple_effects.

`agit log --json` lets the next session read the full reasoning chain before starting.

https://github.com/Madhurr/agit

---

## Claude.ai / Anthropic Community

If you're using Claude Code on anything more than toy projects, you've probably hit the session-reset problem: Claude makes commits with rich reasoning, session ends, next session starts from scratch.

Built **agit** to preserve that. Add one line to `CLAUDE.md`:
```
Use agit commit instead of git commit.
```

That's it. Every commit Claude Code makes will have the full reasoning attached as a git note — intent, confidence, alternatives, risks, unknowns.

`agit context show HEAD` shows it. `agit log` shows the history. `agit log --json` is machine-readable for the next session.

GitHub: https://github.com/Madhurr/agit

---

## Generic developer Slack / Discord (adapted to context)

**Short version (for busy channels):**
Built a tool for teams using AI coding agents: **agit** stores agent reasoning (intent, confidence, alternatives, unknowns) as git notes on each commit. No repo pollution — git notes are native, travel with push/fetch.

`agit context show HEAD` shows why the agent did what it did. `agit diff` shows how confidence evolved across commits.

`go install github.com/Madhurr/agit@v0.3.0`

GitHub: https://github.com/Madhurr/agit

---

## Rules for Discord/Slack engagement

1. **Read the room first.** Drop into channels where people are already discussing AI code review pain, not just any channel.
2. **Lead with the problem.** "When Aider commits, all context is gone when the session ends" — not "I built a tool."
3. **Don't spam.** One post per community. Follow up only if people engage.
4. **Respond to questions.** If people ask "why not commit messages?" or "doesn't this bloat the repo?" — have good answers ready (see HN talking points).
5. **No cross-posting the same message.** Each community gets a version tailored to their tool.
