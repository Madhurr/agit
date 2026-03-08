# Security Policy

## Reporting a Vulnerability

Open a [GitHub Security Advisory](https://github.com/Madhurr/agit/security/advisories/new). Do not file a public issue for security vulnerabilities.

We'll respond within 48 hours and aim to ship a fix within 7 days for critical issues.

## Verifying a Release

Every agit release is built entirely on GitHub Actions — no local builds, no manual uploads. You can verify this:

### 1. Verify the checksum

```bash
# Download binary + checksums
curl -LO https://github.com/Madhurr/agit/releases/latest/download/agit_linux_amd64.tar.gz
curl -LO https://github.com/Madhurr/agit/releases/latest/download/checksums.txt

# Verify
sha256sum --check --ignore-missing checksums.txt
```

### 2. Verify the cosign signature (keyless)

```bash
# Install cosign: https://docs.sigstore.dev/cosign/system_config/installation/
cosign verify-blob checksums.txt \
  --signature checksums.txt.sig \
  --certificate checksums.txt.pem \
  --certificate-identity-regexp "https://github.com/Madhurr/agit/.github/workflows/release.yml" \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com"
```

A successful verification proves:
- The checksums were signed by GitHub Actions (not a compromised developer machine)
- The signing identity matches this exact repository and workflow
- The signature was created at release time via GitHub's OIDC token

### 3. Verify SLSA provenance

```bash
gh attestation verify agit_linux_amd64.tar.gz \
  --repo Madhurr/agit
```

This confirms the binary was produced by GitHub Actions from the source commit in this repository (SLSA Build Level 2).

## Supply Chain

- **No CGO**: binaries are fully static, no C library dependencies
- **Reproducible builds**: `mod_timestamp` is pinned to commit time — the same source commit always produces the same binary hash
- **Vulnerability scanning**: `govulncheck` runs on every commit and every release, blocking release if known CVEs are found in dependencies
- **SBOM**: a CycloneDX SBOM is attached to every release listing all Go module dependencies with exact versions
