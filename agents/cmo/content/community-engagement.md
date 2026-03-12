# Community Engagement Strategy

## Target Communities

### Tier 1 (highest ROI — post immediately once approved)

| Community | Channel | Why |
|---|---|---|
| Hacker News | Show HN | Developer tool discovery channel #1. AI tools get traction. |
| r/programming | Reddit | 6M developers. "I built" posts with demos do well. |
| r/LocalLLaMA | Reddit | Power users of local AI agents — exactly our users. |
| @simonw / @karpathy crowd | Twitter/X | AI coding tool amplifiers. Retweets = thousands of devs. |
| Aider Discord | Discord | Aider users are already using AI to commit code. Ideal fit. |
| Continue.dev Discord | Discord | Same profile as Aider — local AI coding tools. |

### Tier 2 (after initial traction)

| Community | Channel | Why |
|---|---|---|
| r/ExperiencedDevs | Reddit | Senior engineers dealing with AI agent code review pain. |
| Dev.to | Blog | Long-form technical posts get indexed and drive steady traffic. |
| GitHub trending | GitHub | Get on the trending page — one viral HN post can do it. |
| Claude.ai subreddit | Reddit | Claude Code users = primary audience. |
| Cursor community | Discord/Forum | Cursor users commit AI code constantly. |

### Tier 3 (medium-term)

| Community | Channel | Why |
|---|---|---|
| Developer newsletters | Email | TLDR, Pointer, Changelog — one mention = 100K+ devs. |
| YouTube | Video | Demo video: 3 min showing the problem + solution. Evergreen content. |
| GitHub Copilot communities | Various | Enterprise teams hitting this exact problem. |

---

## Engagement Tactics

### "I built this" posts
Lead with the problem, not the solution. The hook is: "AI agents lose all reasoning when the session ends." That's the pain. Then reveal the solution.

### GitHub README is the landing page
Most traffic from HN/Reddit will go to the GitHub repo. The README is already strong. Make sure the demo screenshot is compelling and visible above the fold.

### Engage in existing threads
Search for threads on Reddit/HN where people complain about:
- "AI agent code is a black box"
- "Can't tell why Copilot made this change"
- "AI PR reviews are useless without context"
- "How do you audit AI-generated commits?"

Drop into those threads with a concrete solution. Don't spam — only where genuinely relevant.

### AGENTS.md as a distribution vector
Every time a developer adds `AGENTS.md` to their repo, it tells future AI agents to use `agit commit`. This creates a network effect: one install spreads to everyone who clones that repo and uses AI agents on it.

### The PR workflow angle
The GitHub Actions workflow (`.github/workflows/agit-pr-context.yml`) is a zero-friction distribution path. No account, no registration, just drop a YAML file. Target repos where Claude Code / Copilot is already used.

---

## Key Messages by Audience

### For individual developers using Claude Code / Cursor
"You're losing everything the agent knew when it committed. One line in CLAUDE.md fixes it."

### For engineering managers reviewing AI PRs
"Your team's AI agents commit code with zero reasoning context. agit makes every AI commit auditable."

### For open source maintainers accepting AI contributions
"Contributors using AI agents can now attach their agent's full reasoning to every PR. Reviewers can see confidence levels, alternatives considered, and open unknowns."

### For security/compliance teams
"Every AI-generated commit now has an audit trail: what the agent intended, how confident it was, what risks it flagged, what it didn't know."

---

## Viral Loops

1. **AGENTS.md propagation** — devs add agit to their repo → `AGENTS.md` tells new AI agents to use agit → those agents spread usage to new repos

2. **PR comment visibility** — `agit-pr-context.yml` puts agent reasoning in every PR comment, visible to all reviewers — creates awareness organically

3. **`agit log --json` in agent prompts** — when agents read `agit log --json` at session start to understand prior decisions, they become evangelists (they tell users the context is there and useful)

4. **GitHub trending** — one HN front page post typically generates enough stars to hit GitHub trending for the day/week → compounds discovery

---

## Content Calendar (post-approval)

**Week 1:**
- Day 1: HN Show HN post
- Day 2: r/programming post
- Day 3: r/LocalLLaMA post
- Day 4: Twitter/X thread (thread 1)
- Day 5: Drop into Aider Discord + Continue.dev Discord

**Week 2:**
- Dev.to blog post (the long-form piece)
- Twitter/X thread (thread 2 — technical angle)
- r/ExperiencedDevs post
- Follow up on HN thread comments

**Week 3:**
- Reach out to developer newsletters (TLDR, Pointer, Changelog)
- Create a 3-minute demo video
- GitHub README: add a "Why?" section if not already there

---

## Metrics to Track

- GitHub stars (primary)
- `go install` download counts (via GitHub releases)
- HN ranking / comment engagement
- Reddit upvotes / cross-posts
- Inbound GitHub issues and PRs (best signal of engaged users)
- AGENTS.md usage across public repos (search GitHub)
