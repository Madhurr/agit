# AGENTS.md — ML Pipeline Demo

This demo shows agit used by **GitHub Copilot Workspace** (gpt-4o) refactoring a recommendation engine data pipeline.

## Commits made by the agent

### `0d4a2be` refactor: rewrite feature engineering with recency-weighted embeddings
- **Agent:** copilot-workspace / gpt-4o
- **Confidence:** 72% — offline NDCG@10 improved 8.3% but backtesting window only 30 days
- **Rejected:** Transformer sequential model (400ms P99), Attention (12GB RAM/instance), mean pooling (stale recommendations)
- **Risk [high]:** cold-start users (<5 interactions) get near-zero embeddings — no fallback
- **Risk [medium]:** decay rate 0.1 hardcoded, needs tuning per vertical
- **Unknowns:** optimal decay rate per content type, popularity debiasing placement (retrieval vs ranking)

### `ce41bf1` feat: add pipeline data validation and train/serve skew detection
- **Agent:** copilot-workspace / gpt-4o
- **Confidence:** 81% — threshold=0.5 is a guess, needs 2 weeks calibration against production data
- **Rejected:** Great Expectations (800MB dep), assert statements (silently swallowed in prod)
- **Risk [medium]:** drift tolerance may cause false positives in seasonal traffic spikes
- **Unknowns:** hard-fail vs log-and-continue on validation failure — unclear from product requirements

## How to use agit here

```bash
go install github.com/madhurm/agit@latest
agit init
agit commit -m "feat: add cold-start fallback to popularity model" \
  --intent "Users with <5 interactions fall back to global popularity ranking" \
  --confidence 0.80 \
  --tried "content-based filtering: rejected — no item metadata available at query time" \
  --unknowns "optimal threshold for switching from cold-start to personalized model"
```
