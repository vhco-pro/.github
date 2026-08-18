# Security policy

This policy applies to **every repository** in the
[`vhco-pro`](https://github.com/vhco-pro) GitHub organisation, which
hosts the open-source projects maintained by **VH & Co BV** —
[Stackweaver™](#project-specific-stackweaver) and its ecosystem, plus the
org's other tooling projects.

Project-specific terms that do not generalise — supported-version
windows, release-verification recipes, trust models — live in the
per-project sections at the end of this document.

## Reporting a vulnerability

**Preferred:** open a
[GitHub Private Vulnerability Report (PVR)](https://docs.github.com/en/code-security/security-advisories/working-with-repository-security-advisories/about-repository-security-advisories)
on the repository where you observed the issue. PVR keeps the report
invisible to the public until a fix is published, and it is enabled on
the `vhco-pro` repositories.

**Alternative:** email [`contact@vhco.pro`](mailto:contact@vhco.pro).
PVR remains the preferred path because it keeps the report invisible
until a fix ships. We do **not** publish a PGP key for this mailbox: the
org follows a Sigstore-only signing policy (no long-lived keys), and PVR
end-to-end-encrypts in transit and at rest via GitHub.

We aim to:

- acknowledge the report within **5 business days**;
- agree a remediation timeline within **10 business days**;
- coordinate a disclosure date with the reporter.

Please **do not** open public GitHub issues for vulnerabilities.

## Scope

In scope, for any repository in this org:

- The project's own source code, and the artifacts it publishes
  (container images, charts, binaries, packages, GitHub Actions).
- The project's CI and release pipeline, including the workflows that
  build and sign those artifacts and any GitHub App used to drive them.

Out of scope:

- Third-party services and dependencies a project integrates with — your
  cloud provider, your container registry, GitHub itself, and upstream
  software the project merely consumes.
- Your own deployment configuration: network policy, secret storage,
  access control, and host hardening.
- Findings that require an attacker to already hold administrative
  access to the system running the software.

## Verifying a release

Repositories that publish signed artifacts sign them with **cosign
keyless (Sigstore)** and, where the artifact is built by GitHub Actions,
attach a **SLSA build provenance attestation**. There are no long-lived
signing keys anywhere in this org.

The general verification shape is:

```bash
# Signature was produced by a workflow in this org
cosign verify \
  --certificate-identity-regexp '^https://github\.com/vhco-pro/<repo>' \
  --certificate-oidc-issuer 'https://token.actions.githubusercontent.com' \
  "<artifact>"

# Build provenance
gh attestation verify --owner vhco-pro --repo <repo> "<artifact>"
```

Each repository's README states which of these it publishes and the exact
commands for its artifacts. Not every project in this org publishes
release artifacts; the ones that do not are source-only.

## Coordinated disclosure

We publish a GitHub Security Advisory for every confirmed vulnerability
once a fix is available, request a CVE where the issue affects published
artifacts, and credit reporters who want to be credited.

---

## Project-specific: Stackweaver

> **Status (2026-05-23): the sections below marked _TBD_ are the
> audit-window stubs.** They are gated on decisions taken during the
> ongoing OSPS Baseline + OpenSSF Scorecard rollout. Tracking issue:
> [michielvha/stackweaver#224](https://github.com/michielvha/stackweaver/issues/224).

Covers the Stackweaver application (API, Orchestrator, OpenTofu Runner,
Ansible Runner, Frontend, Zitadel bootstrap, Helm chart), distributed via
the seven `vhco-pro/stackweaver-*` satellite repositories, and the
release pipeline behind them (the `stackweaver-release-bot` GitHub App,
the `sync-*.yml` workflows, the SLSA provenance attestations).

### Supported versions

_TBD — finalised at the first signed release. Working policy:_ the latest
minor receives security fixes; the previous minor receives security fixes
for **6 months** after the next minor ships.

### End-of-life

A version reaches end-of-life when it has been superseded by two newer
minor versions or after **18 months** from initial release, whichever
comes first. EOL versions receive no further security updates.

### Verifying a Stackweaver release

Every satellite release ships with a cosign keyless signature on every
published container image and Helm chart, a SLSA L3 build provenance
attestation binding the satellite commit SHA to the upstream monorepo
commit SHA and the Actions run that produced it, and an SBOM (SPDX +
CycloneDX) as a release asset.

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

A consumer-friendly long-form guide lives at
[`sw.vhco.pro/docs/security/verifying-releases`](https://sw.vhco.pro/docs/security/verifying-releases).

### Source-code review boundary

Stackweaver is developed in a private monorepo and published to the
satellites via a one-way bot-driven sync. Authoritative human review
happens **upstream** on the monorepo PR; the satellite commit is a
deterministic, cryptographically-attested re-publication produced by the
`stackweaver-release-bot` GitHub App. No human holds push rights to
satellite `main`.

The full threat model and the trust-model rationale behind this design
are published as part of the OSPS audit. The headline reason: in
2025-2026 the PR-time-automation supply-chain attack class
(tj-actions / Coinbase, March 2025) is materially worse than a
single-principal-push trust chain, so we chose the latter.

This boundary applies **only** to the seven satellites. Every other
repository in this org, Stackweaver-adjacent ones included
(`stackweaver-operator`, `terraform-provider-stackweaver`), is developed
in the open with ordinary pull-request review.

### Audit access to closed components

The closed `core/` Go module is not currently public. Security auditors
engaged by enterprise customers can receive read-only access under a
one-off NDA — see the project maintainers for the procedure.

---

Maintained by **VH & Co BV**. Owner: @michielvha.

Stackweaver™ is a trademark of VH & Co. See [`TRADEMARK.md`](TRADEMARK.md)
for the Trademark Policy.
