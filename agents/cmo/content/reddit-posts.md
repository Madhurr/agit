# Reddit Posts

## r/programming

**Title:** I built agit: AI coding agents commit code but lose all reasoning when the session ends. This stores it.

**Body:**

The problem: When Claude Code, Copilot, or Cursor writes a commit, it knows:
- Why it chose this approach over alternatives
- How confident it was
- What could break downstream
- What it left unresolved

That all disappears. The next session (or the next human) gets code with no context. "Why was auth switched to JWT?" ¯\\_(ツ)_/¯

The solution I built: `agit` wraps `git commit` and stores structured reasoning as a git note on each commit. Notes are a native git feature from 2010 — no new files, no repo pollution, travels with push/fetch.

```bash
agit commit \
  -m "feat: switch auth to JWT" \
  --intent "stateless sessions, no Redis dependency" \
  --confidence 0.82 \
  --tried "session cookies: needs shared store across pods" \
  --risk "high:token-refresh:not implemented yet" \
  --unknowns "revocation strategy undecided"
```

Then `agit context show HEAD` or `agit log` shows the full reasoning. `agit diff HEAD~5..HEAD` shows how confidence changed across a feature branch.

The field I find most useful: `--unknowns`. Not "here's what the agent did" — "here's what it *didn't know*." That's what bites you weeks later.

GitHub: https://github.com/Madhurr/agit
Install: `go install github.com/Madhurr/agit@v0.3.0`

Works with Claude Code, Aider, Cursor, Copilot, Devin — anything that can run a shell command.

---

## r/LocalLLaMA

**Title:** agit: store local LLM coding agent reasoning in git commits (git notes, no repo pollution)

**Body:**

If you run local coding agents (Continue, Aider with local models, Cursor with local LLMs) you've probably hit this: the agent makes changes, commits, session ends. Next session has no idea what the previous one was thinking.

Built `agit` to fix this. It wraps `git commit` and stores the agent's reasoning as a git note — structured JSON attached to the commit, stored in `refs/notes/agit`. No new files in the working tree. Notes travel with push/fetch so the whole team gets context.

Key fields:
- `intent` — what the agent was actually trying to do
- `confidence` + `confidence_rationale` — how sure it was and why
- `alternatives_considered` — what it tried and rejected
- `risks` — what might break
- `unknowns` — what it didn't know (the most important field)
- `ripple_effects` — downstream changes expected

Works with any agent that can run shell commands. Add one line to your `.aider.conf.yml`:
```
commit-cmd: agit commit -m
```

Or in Continue's config, override the commit command.

`agit log --json` lets the next session read the full reasoning chain before starting work.

GitHub: https://github.com/Madhurr/agit

---

## r/ExperiencedDevs

**Title:** The context problem with AI coding agents: they know things when they commit that you'll never see again

**Body:**

Thought experiment: you review a PR where an AI agent touched 12 files, refactored auth, and changed 3 API contracts. The commit message says "refactor: improve auth layer."

The agent, when it made that commit, knew:
- It considered 3 other approaches and rejected them
- It's 71% confident — there's a specific edge case it wasn't sure about
- One of the changed contracts will require a migration in the mobile app
- The token refresh flow is completely unimplemented

None of that is in the commit. It can't be — commit messages aren't structured data, and the session ended.

This is the AI-native version of the "why?" gap in git history. Tools like `git blame` help you find *who* and *when*. `agit` captures the *why* while the agent still knows it.

I built it using git notes (native git feature, been there since 2010) so there's nothing new to store or sync. Notes travel with push/fetch.

`agit context show <hash>` shows the full reasoning. `agit diff from..to` shows how agent confidence evolved across commits.

GitHub: https://github.com/Madhurr/agit
Install: `go install github.com/Madhurr/agit@v0.3.0`

Interested in the tradeoffs people see here. Is structured agent context in the repo useful, or does it add noise?

---

## r/MachineLearning (more applied/engineering)

**Title:** [Project] Preserving AI agent reasoning across coding sessions using git notes

**Body:**

Short version: AI coding agents lose all reasoning when a session ends. This project stores it in git commits.

Longer: The gap I'm trying to close is between what an AI agent *knows* at commit time and what's recoverable afterward. Right now nothing is recoverable — the model is gone, the session is gone, the context window is gone. What's left is a diff.

`agit` wraps `git commit` and attaches a JSON note to each commit:

```json
{
  "intent": "stateless sessions, no Redis dependency",
  "confidence": 0.82,
  "confidence_rationale": "core path works, refresh not covered",
  "alternatives_considered": [
    { "approach": "session cookies", "rejected_reason": "needs shared store across pods" }
  ],
  "risks": [{ "area": "token-refresh", "description": "not implemented", "severity": "high" }],
  "unknowns": ["revocation strategy undecided"],
  "ripple_effects": ["mobile clients will need token refresh handling"]
}
```

Storage: git notes, `refs/notes/agit`. No repo pollution. Travels with push/fetch. Has been in git since 2010.

The interesting design question was where to put this. Options I considered:
- Commit message: too unstructured for programmatic consumption, clutters `git log`
- `.agit/` sidecar directory: pollutes working tree, conflicts with existing tooling
- External database: requires infrastructure, doesn't travel with repo
- git notes: zero working tree impact, travels with the repo, native git

GitHub: https://github.com/Madhurr/agit
