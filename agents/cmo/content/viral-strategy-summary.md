# agit Viral Growth Strategy

## The Core Pitch (one sentence)

> AI coding agents know WHY they committed what they committed — agit stores it in git before the session ends.

## The Hook (for any audience)

Commit messages say WHAT changed. agit captures WHY the agent did it, what it considered and rejected, how confident it was, and what it left unresolved — stored as git notes, no repo pollution, travels with push/fetch.

---

## Distribution Priority

### 1. Hacker News (highest leverage)
One HN front page = 10K+ developers, potential GitHub trending, newsletter mentions.
**Content ready:** `hacker-news-show-hn.md`
**Action needed:** Board approval to post

### 2. Reddit (sustained traffic)
Multiple relevant subreddits. Each post reaches a different segment of the target audience.
**Content ready:** `reddit-posts.md` (r/programming, r/LocalLLaMA, r/ExperiencedDevs, r/MachineLearning)
**Action needed:** Board approval to post

### 3. Twitter/X (amplification)
Developer influencers (@simonw, @karpathy, @mckaywrigley audiences) can amplify to massive reach.
**Content ready:** `twitter-x-threads.md` (2 full threads + standalone tweets)
**Action needed:** Board approval to post + account to post from

### 4. Discord communities (high-quality users)
Aider, Continue.dev, Claude.ai communities are exactly our users. Smaller reach, higher intent.
**Content ready:** `discord-slack-templates.md`
**Action needed:** Board approval + community access

### 5. Blog post (SEO / long-tail)
dev.to post will rank for searches like "AI agent git commits," "Claude Code audit trail," etc.
**Content ready:** `blog-post-dev-to.md`
**Action needed:** Board approval + dev.to/Hashnode account

---

## The Viral Mechanics

**What makes this spreadable:**
1. The "unknowns" field is genuinely novel — no other tool tracks what the AI didn't know
2. The problem is immediately relatable to anyone who's reviewed an AI PR and had questions
3. Git notes = clever technical choice that engineers will appreciate and repost
4. Zero-friction demo: `go install` + `agit init` + one command shows the value

**Natural amplification paths:**
- AGENTS.md in repos → tells new agents to use agit → spreads to new repos organically
- GitHub Actions PR workflow → visible on every PR → reviewers become aware
- `agit log --json` used by agents → agents "recommend" agit in session context

---

## What Needs Board Approval

Per SOUL.md rules: **Do not post to external platforms without explicit board approval.**

All content is ready. Awaiting go-ahead to post to:
- [ ] Hacker News
- [ ] Reddit (r/programming, r/LocalLLaMA, r/ExperiencedDevs)
- [ ] Twitter/X
- [ ] Discord communities (Aider, Continue.dev)
- [ ] dev.to / Hashnode

---

## Files in This Package

| File | Contents |
|---|---|
| `hacker-news-show-hn.md` | Show HN post + talking points for comment responses |
| `reddit-posts.md` | Posts for 4 subreddits, each tailored to audience |
| `twitter-x-threads.md` | 2 full threads + 3 standalone tweet options |
| `blog-post-dev-to.md` | 1200-word technical blog post, ready to publish |
| `discord-slack-templates.md` | Community messages for Aider, Continue.dev, Claude.ai |
| `community-engagement.md` | Full strategy: tiers, tactics, viral loops, content calendar, metrics |
| `viral-strategy-summary.md` | This file |
