# Show HN: agit — git commits that explain why the AI did it

**Title:** Show HN: agit – store AI agent reasoning in git commits (using git notes)

**URL:** https://github.com/Madhurr/agit

---

## Post body (plain text for HN submission)

When Claude Code, Copilot, or Cursor makes a commit, it knows things: why it chose this approach over alternatives, how confident it was, what could break downstream, what it left unresolved. That context disappears when the session ends. The next agent — or the next human — opens the file and sees code with no reasoning attached.

agit fixes that. It wraps `git commit` and stores agent reasoning as a git note on each commit:

```
$ agit commit \
  -m "feat: switch auth to JWT" \
  --intent "stateless sessions, no Redis dependency" \
  --confidence 0.82 \
  --tried "session cookies: needs shared store across pods" \
  --risk "high:token-refresh:not implemented yet" \
  --unknowns "revocation strategy undecided"
```

Then later:

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

**Why git notes?** They're a native git feature that's been there since 2010. No new files in your working tree, no `.agit/` directory to gitignore. Notes travel with `git push` and `git fetch`. Works offline. Plain JSON. Every existing git tool ignores them unless you explicitly ask.

**The field I'm most proud of:** `unknowns`. Not "here's what the agent did" — "here's what the agent *didn't know*." This is the thing that actually kills you three weeks later when someone asks why the revocation strategy is missing.

**How to use with Claude Code:** Add one line to CLAUDE.md: `Use agit commit instead of git commit.`

**With Aider:** `aider --commit-cmd "agit commit -m"`

**With Cursor/Copilot/Devin:** If it has a terminal, it runs `agit commit`. Point it at the included AGENTS.md.

Install: `go install github.com/Madhurr/agit@v0.3.0`

---

## Talking points for comments

If someone asks "why not just put this in the commit message?":
> Commit messages are for humans reading git log. This is structured data for *agents* reading context before starting work. `agit log --json` gives the next session everything it needs programmatically — model, confidence, unknowns, risks. You can't parse that from prose.

If someone asks "isn't this what commit messages are for?":
> A commit message says "switched auth to JWT." agit says: switched to JWT (confidence 82%), rejected session cookies because shared store doesn't scale, token refresh is unimplemented (high risk), revocation strategy is completely open. These are different things.

If someone asks "doesn't this bloat the repo?":
> No. git notes live in refs/notes/agit, completely separate from commit objects. `git log` doesn't show them. Clones don't get them unless you explicitly fetch. Add 1 line to .git/config to always fetch them if you want.

If someone asks "how does this compare to conventional commits / semantic commits?":
> Conventional commits are a formatting convention for commit *messages* to enable changelog generation. agit is structured metadata for agent *reasoning* — different axis entirely. They're complementary.

---

## Timing

Post on a Tuesday or Wednesday, 8–10am ET. HN front page is most competitive Mon morning and Fri afternoon.

**Target thread:** Show HN (not Ask HN, not just a link submission — use Show HN to get top placement in the Show HN section)
