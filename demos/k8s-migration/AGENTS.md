# AGENTS.md — K8s Migration Demo

This demo shows agit used by **Devin** migrating a monolith to Kubernetes microservices.

## Commits made by the agent

### `361a850` feat: extract auth service from monolith as standalone K8s deployment
- **Agent:** devin / devin-1
- **Confidence:** 74% — K8s manifests untested against real cluster, Secret rotation not documented
- **Rejected:** Sidecar pattern (still coupled), shared library (RFC-47 rules it out), Lambda (cold start breaks P99 < 50ms SLA)
- **Risk [high]:** JWT_SECRET rotation requires rolling restart — zero-downtime rotation not implemented
- **Risk [high]:** dual-write period with monolith creates token consistency risk
- **Ripple effects:** monolith auth middleware needs feature flag, API gateway needs /auth/* routing, all services need JWT validation against auth-service public key
- **Unknowns:** token refresh strategy, active session behavior during cutover

### `529f016` feat: add Istio VirtualService for canary traffic splitting on auth-service
- **Agent:** devin / devin-1
- **Confidence:** 69% — config syntactically valid but untested, outlier thresholds untuned
- **Rejected:** Argo Rollouts (operational complexity), blue/green (DNS TTL delays), feature flags (12-factor violation)
- **Risk [high]:** cluster on Istio 1.16 but config uses 1.18 API — may fail silently on apply
- **Risk [medium]:** x-canary header unauthenticated — anyone can self-select into v2
- **Unknowns:** mTLS enforcement between services not verified, no PodDisruptionBudget

## How to use agit here

```bash
go install github.com/madhurm/agit@latest
agit init
agit commit -m "fix: add PodDisruptionBudget for auth-service" \
  --intent "Ensure rolling updates never take all 3 replicas down simultaneously" \
  --confidence 0.92 \
  --risk "low:min-available:minAvailable=2 means 1 pod always up during updates"
```
