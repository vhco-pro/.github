# Security policy

This policy applies to every repository in the
[`vhco-pro`](https://github.com/vhco-pro) GitHub organisation —
i.e. all `stackweaver-*` satellite repositories that together comprise
the [Stackweaver](https://github.com/vhco-pro) distribution.

> **Status (2026-05-23): this document is the audit-window stub.**
> Some sections marked _TBD_ are gated on decisions taken during the
> ongoing OSPS Baseline + OpenSSF Scorecard rollout. Tracking issue:
> [michielvha/stackweaver#224](https://github.com/michielvha/stackweaver/issues/224).

## Reporting a vulnerability

**Preferred:** open a [GitHub Private Vulnerability Report (PVR)](https://docs.github.com/en/code-security/security-advisories/working-with-repository-security-advisories/about-repository-security-advisories) on the satellite where you observed the issue. PVR
keeps the report invisible to the public until a fix is published. Each
`vhco-pro/stackweaver-*` repository has PVR enabled.

**Alternative:** email `security@stackweaver.io` _(TBD — the mailbox and
PGP key are being provisioned; use PVR until this section is updated)_.

We aim to:

- acknowledge the report within **5 business days**;
- agree a remediation timeline within **10 business days**;
- coordinate a disclosure date with the reporter.

Please **do not** open public GitHub issues for vulnerabilities.

## Scope

In scope:

- The Stackweaver application (API, Orchestrator, Terraform Runner,
  Ansible Runner, Frontend, Zitadel bootstrap, Helm chart).
- The CI / release pipeline (the `stackweaver-release-bot` GitHub App,
  the `sync-*.yml` workflows, the SLSA provenance attestations).

Out of scope:

- Third-party services Stackweaver integrates with (Zitadel itself,
  GitHub, your container registry, your cloud provider).
- Your own deployment configuration (network policies, secret stores).

## Supported versions

_TBD — finalised at the first signed release. Working policy:_ the
latest minor receives security fixes; the previous minor receives
security fixes for **6 months** after the next minor ships.

## End-of-life

A version reaches end-of-life when it has been superseded by two newer
minor versions or after **18 months** from initial release, whichever
comes first. EOL versions receive no further security updates.

## Verifying a release

Every Stackweaver satellite release ships with:

1. A **cosign keyless signature** on every published container image
   and Helm chart, signed via the project's Sigstore identity (no
   long-lived signing keys).
2. A **SLSA L3 build provenance attestation** binding the satellite
   commit SHA to the upstream monorepo commit SHA and the GitHub
   Actions workflow run that produced it.
3. An **SBOM** (SPDX + CycloneDX) published as a release asset.

To verify:

```bash
IMAGE=ghcr.io/vhco-pro/stackweaver-api:<tag>

# 1. Sigstore signature is by the project release pipeline
cosign verify \
  --certificate-identity-regexp '^https://github\.com/vhco-pro/stackweaver-.+' \
  --certificate-oidc-issuer 'https://token.actions.githubusercontent.com' \
  "$IMAGE"

# 2. SLSA build provenance attestation
gh attestation verify --owner vhco-pro --repo stackweaver-api "$IMAGE"
```

A consumer-friendly long-form guide will appear at
`profile/security/verifying-releases.md` with the first signed release.

## Source-code review boundary

Stackweaver is developed in a private monorepo and published to these
satellites via a one-way bot-driven sync. Authoritative human review
happens **upstream** on the monorepo PR; the satellite commit is a
deterministic, cryptographically-attested re-publication produced by
the `stackweaver-release-bot` GitHub App. No human holds push rights to
satellite `main`.

The full threat model and the trust-model rationale behind this design
are published as part of the OSPS audit. The headline reason: in
2025-2026 the PR-time-automation supply-chain attack class
(tj-actions / Coinbase, March 2025) is materially worse than a
single-principal-push trust chain, so we chose the latter.

## Audit access to closed components

The closed `core/` Go module is not currently public. Security auditors
engaged by enterprise customers can receive read-only access under a
one-off NDA — see the project maintainers for the procedure.

---

Maintained by **VH & Co BV**. Owner: @michielvha.
