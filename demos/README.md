# agit Demos

Real-world scenarios showing agit used by different AI coding tools.

Each demo is a realistic project with commits made by an AI agent using `agit commit`. Browse the `AGENTS.md` in each folder to see the full reasoning chain that would normally be lost.

---

## [payment-service/](payment-service/) — Cursor + claude-3.5-sonnet
Stripe payment integration with PaymentIntents, 3DS support, webhook handling.

Key insight: agent flagged its own security gap (missing Stripe-Signature validation) at 65% confidence — that risk is preserved in the commit, visible to the next reviewer or agent session.

---

## [ml-pipeline/](ml-pipeline/) — GitHub Copilot Workspace + gpt-4o
Recommendation engine data pipeline refactor: recency-weighted embeddings + train/serve skew detection.

Key insight: agent tried 3 approaches before settling on exponential decay weighting — Transformer (too slow), Attention (too much RAM), mean pooling (produces stale recommendations). All preserved.

---

## [k8s-migration/](k8s-migration/) — Devin
Monolith → Kubernetes migration. Auth service extraction + Istio canary deployment.

Key insight: agent knew the Istio config might silently fail (version mismatch risk) but shipped anyway at 69% confidence. That uncertainty is now tracked, not buried.

---

## Running the demos yourself

```bash
# Install agit
go install github.com/madhurm/agit@latest

# Try any of these repos
cd demos/payment-service
git init && agit init

# Make commits with full reasoning
agit commit -m "fix: add Stripe-Signature validation" \
  --intent "Close the security gap flagged in the webhook commit" \
  --confidence 0.95 \
  --task "Fix critical: webhook accepts any POST without signature check"

# See the semantic history
agit log
agit context show HEAD
```
