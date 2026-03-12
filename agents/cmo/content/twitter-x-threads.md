# Twitter/X Threads

## Thread 1: The Core Problem (best for virality)

**Tweet 1 (hook):**
When Claude Code commits your code, it knows things you'll never see:
- why it chose this approach over 3 others
- that it's only 68% confident
- that the auth refactor will break the mobile client
- that the revocation flow is completely unimplemented

Then the session ends. That knowledge is gone. 🧵

**Tweet 2:**
The commit message says "refactor: improve auth layer."

The agent knew:
• rejected JWT: "token refresh not implemented"
• rejected Redis sessions: "not scalable across pods"
• confidence 68%: "core path works, edge cases unclear"
• risk: "mobile client needs migration"
• unknown: "revocation strategy undecided"

**Tweet 3:**
This is the AI-native version of the "why?" gap in git history.

git blame: who changed it and when
git log: what changed
agit: why the agent did it, what it considered, how confident it was, what it left unresolved

**Tweet 4:**
I built agit to close this gap.

It wraps `git commit` and stores structured reasoning as a git note on each commit.

git notes are a native git feature from 2010. No new files. No .agit/ directory. Notes travel with push and fetch.

**Tweet 5:**
```
agit commit \
  -m "feat: switch auth to JWT" \
  --intent "stateless sessions, no Redis" \
  --confidence 0.82 \
  --tried "session cookies: shared store doesn't scale" \
  --risk "high:token-refresh:not implemented" \
  --unknowns "revocation strategy undecided"
```

**Tweet 6:**
Then later, `agit context show HEAD`:

```
Intent:     stateless sessions, no Redis dependency
Confidence: 82% — "core path works, refresh not covered"

Alternatives rejected:
  • session cookies — shared store doesn't scale

Risks:
  [high] token-refresh: not implemented yet

Unknowns:
  • revocation strategy undecided
```

**Tweet 7:**
The field I'm most proud of: `--unknowns`

Not "here's what the agent did." "Here's what it *didn't know*."

Three weeks later when someone asks why revocation is missing — that's your answer.

**Tweet 8:**
Works with Claude Code, Aider, Copilot, Cursor, Devin. Anything with a terminal.

```bash
go install github.com/Madhurr/agit@v0.3.0
cd your-repo && agit init
```

Claude Code: add "Use agit commit instead of git commit" to CLAUDE.md
Aider: `aider --commit-cmd "agit commit -m"`

**Tweet 9:**
`agit diff HEAD~5..HEAD` shows how reasoning evolved across a feature:

```
Confidence: 68% → 91%
✓ Risk resolved: [high] token-refresh
✓ Unknown resolved: revocation strategy
+ New risk: [medium] race condition in worker pool
```

Agent confidence as a timeline. Not just what changed — how the agent's understanding changed.

**Tweet 10:**
Open source, MIT. GitHub: https://github.com/Madhurr/agit

If you're using AI coding agents and reviewing their PRs with zero context about why they did what they did — this is for you.

---

## Thread 2: Technical angle (for engineer-heavy followers)

**Tweet 1:**
Hot take: the hardest problem with AI coding agents isn't the code quality. It's the reasoning gap.

The agent knows WHY. It forgets when the session ends. We've been shipping agent commits with no audit trail of decisions.

**Tweet 2:**
I've been thinking about where to store agent reasoning and the options are:
- commit message: unstructured, clutters git log
- .agit/ sidecar: pollutes working tree
- external DB: requires infra, doesn't travel with repo
- git notes: ✓

**Tweet 3:**
git notes have been in git since 2010. They attach to commits but live in a separate ref (refs/notes/agit). Zero working tree impact. Travel with push/fetch. Every existing git tool ignores them unless you ask.

Perfect for structured agent metadata.

**Tweet 4:**
Built agit around this. JSON schema with the fields that actually matter:

- intent (what the agent was trying to do)
- confidence + rationale
- alternatives_considered (what it tried and rejected)
- risks (severity-tagged)
- unknowns (what it didn't know)
- ripple_effects (what else will need changing)

**Tweet 5:**
The unknowns field is the one I keep coming back to.

An agent with 95% confidence that leaves 0 unknowns is different from an agent with 95% confidence that leaves 3 unknowns. The second one is a time bomb.

Making unknowns explicit is what surfaces that.

**Tweet 6:**
agit diff shows confidence drift across commits:

68% → 74% → 91%

You can see the agent get more confident as it resolves risks and unknowns. Or you can see it stay stuck at 70% across 8 commits — that's a signal something needs human attention.

**Tweet 7:**
GitHub: https://github.com/Madhurr/agit
`go install github.com/Madhurr/agit@v0.3.0`

MIT. Works with Claude Code, Aider, Cursor, Codex.

---

## Single tweet options (for standalone posts)

**Option A:**
AI coding agents know things at commit time that disappear when the session ends.

agit stores them as git notes: intent, confidence, alternatives rejected, risks, unknowns.

`agit context show HEAD` — WHY, not just WHAT.

github.com/Madhurr/agit

**Option B:**
You can see what Claude Code changed. You can't see why it changed it, what it considered, how confident it was, or what it left unresolved.

That's what agit fixes. Git notes. No repo pollution. Works with any agent.

github.com/Madhurr/agit

**Option C:**
The `--unknowns` flag in agit might be the most important field in the schema.

Not "here's what the AI agent did." "Here's what it didn't know."

Three weeks later when something breaks — that's where the answer is.
