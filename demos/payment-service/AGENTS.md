# AGENTS.md — Payment Service Demo

This demo shows agit used by **Cursor** (claude-3.5-sonnet) building a Stripe payment integration.

## Commits made by the agent

### `f24fd88` feat: add Stripe payment intent creation and confirmation
- **Agent:** cursor-agent / claude-3.5-sonnet
- **Confidence:** 78% — 3DS redirect untested with real cards, webhook signature verification missing
- **Rejected:** Direct card charging (PCI scope), Charges API legacy (deprecated, no 3DS)
- **Risk [high]:** webhook endpoint not validating Stripe-Signature — any POST processed
- **Risk [medium]:** no idempotency keys — retries could double-charge
- **Unknowns:** Setup Intents vs Stripe Billing for recurring, refund flow not designed

### `b8e4b14` feat: add webhook handler for payment events
- **Agent:** cursor-agent / claude-3.5-sonnet
- **Confidence:** 65% — Stripe-Signature validation TODO, no DB retry logic
- **Rejected:** Polling (30s delay), Synchronous confirmation (incompatible with 3DS)
- **Risk [high]:** stripe-signature not validated, idempotency not checked
- **Unknowns:** GDPR implications of raw webhook storage, Stripe retry behavior on DB failure

## How to use agit here

```bash
go install github.com/madhurm/agit@latest
agit init
agit commit -m "feat: add refund endpoint" \
  --intent "Handle partial and full refunds via Stripe Refunds API" \
  --confidence 0.75 \
  --risk "medium:idempotency:refund retries could double-refund without idempotency key"
```
