# VH & Co

> Open-source infrastructure, IaC tooling, and developer tools from
> **VH & Co BV**.

## Stackweaver - multi-IaC orchestration platform

> A self-hostable alternative to Terraform Cloud / Ansible AWX,
> supporting OpenTofu and Ansible from a single control plane,
> API-compatible with Terraform Cloud/Enterprise tooling.

Stackweaver is developed in a private monorepo and published here as
eight independent satellite repositories. These are one-way distribution
mirrors: they take issues, not pull requests.

| Component | Repo | OpenSSF Scorecard |
|-----------|------|-------------------|
| Helm chart                | [`stackweaver-helm`](https://github.com/vhco-pro/stackweaver-helm) | [![Scorecard](https://api.scorecard.dev/projects/github.com/vhco-pro/stackweaver-helm/badge)](https://scorecard.dev/viewer/?uri=github.com/vhco-pro/stackweaver-helm) |
| Backend API               | [`stackweaver-api`](https://github.com/vhco-pro/stackweaver-api) | [![Scorecard](https://api.scorecard.dev/projects/github.com/vhco-pro/stackweaver-api/badge)](https://scorecard.dev/viewer/?uri=github.com/vhco-pro/stackweaver-api) |
| Orchestrator              | [`stackweaver-orchestrator`](https://github.com/vhco-pro/stackweaver-orchestrator) | [![Scorecard](https://api.scorecard.dev/projects/github.com/vhco-pro/stackweaver-orchestrator/badge)](https://scorecard.dev/viewer/?uri=github.com/vhco-pro/stackweaver-orchestrator) |
| OpenTofu Runner           | [`stackweaver-opentofu-runner`](https://github.com/vhco-pro/stackweaver-opentofu-runner) | [![Scorecard](https://api.scorecard.dev/projects/github.com/vhco-pro/stackweaver-opentofu-runner/badge)](https://scorecard.dev/viewer/?uri=github.com/vhco-pro/stackweaver-opentofu-runner) |
| Ansible Runner            | [`stackweaver-ansible-runner`](https://github.com/vhco-pro/stackweaver-ansible-runner) | [![Scorecard](https://api.scorecard.dev/projects/github.com/vhco-pro/stackweaver-ansible-runner/badge)](https://scorecard.dev/viewer/?uri=github.com/vhco-pro/stackweaver-ansible-runner) |
| Frontend                  | [`stackweaver-frontend`](https://github.com/vhco-pro/stackweaver-frontend) | [![Scorecard](https://api.scorecard.dev/projects/github.com/vhco-pro/stackweaver-frontend/badge)](https://scorecard.dev/viewer/?uri=github.com/vhco-pro/stackweaver-frontend) |
| Zitadel bootstrap         | [`stackweaver-zitadel-init`](https://github.com/vhco-pro/stackweaver-zitadel-init) | [![Scorecard](https://api.scorecard.dev/projects/github.com/vhco-pro/stackweaver-zitadel-init/badge)](https://scorecard.dev/viewer/?uri=github.com/vhco-pro/stackweaver-zitadel-init) |
| Secret bootstrap          | [`stackweaver-secrets-init`](https://github.com/vhco-pro/stackweaver-secrets-init) | [![Scorecard](https://api.scorecard.dev/projects/github.com/vhco-pro/stackweaver-secrets-init/badge)](https://scorecard.dev/viewer/?uri=github.com/vhco-pro/stackweaver-secrets-init) |

```bash
helm install stackweaver oci://ghcr.io/vhco-pro/charts/stackweaver --version <X.Y.Z>
```

See the Helm satellite README for the full deployment guide, or
[`sw.vhco.pro/docs`](https://sw.vhco.pro/docs) for the user documentation.

**Ecosystem** - developed in the open, pull requests welcome:

| Repo | What it is |
|------|------------|
| [`terraform-provider-stackweaver`](https://github.com/vhco-pro/terraform-provider-stackweaver) | Terraform provider for Stackweaver, derived from `terraform-provider-tfe` and kept in sync with it |
| [`stackweaver-operator`](https://github.com/vhco-pro/stackweaver-operator) | Kubernetes operator for managing Stackweaver deployments |

## Infrastructure & IaC tooling

| Repo | What it is |
|------|------------|
| [`stackgraph`](https://github.com/vhco-pro/stackgraph) | Infrastructure diagram generator for OpenTofu/Terraform - production-quality architecture diagrams from state files, plan JSON, and HCL source |
| [`terraform-provider-garage`](https://github.com/vhco-pro/terraform-provider-garage) | Terraform provider for Garage, the S3-compatible object store (Admin API v2) |
| [`builders`](https://github.com/vhco-pro/builders) | PDS-backed golden-image builders for the homelab |

## Containers & release engineering

| Repo | What it is |
|------|------------|
| [`distil`](https://github.com/vhco-pro/distil) | Zero-CVE container images. Fully automated, built with apko + Wolfi |
| [`swift-release-action`](https://github.com/vhco-pro/swift-release-action) | Reusable macOS Swift app release pipeline - reusable workflow plus a secret-free composite build/sign/package action |

## macOS & AWS tools

| Repo | What it is |
|------|------------|
| [`ssm-connect`](https://github.com/vhco-pro/ssm-connect) | Config-driven macOS menu-bar app that auto-establishes AWS SSM port-forward tunnels to EC2 workstations (SSO auth, bundled session-manager-plugin) |
| [`dcv-session-agent`](https://github.com/vhco-pro/dcv-session-agent) | On-box agent for multi-user Amazon DCV on a single self-managed EC2 host - per-user virtual sessions and AWS-SSO-identity token auth, no broker, no passwords |
| [`claude-companion`](https://github.com/vhco-pro/claude-companion) | macOS menu-bar companion for Claude Code - tool-call auto-approval against a shared blacklist, session monitoring, usage and cost tracking |
| [`postbode`](https://github.com/vhco-pro/postbode) | Gmail to ClearFacts/QPS purchase-invoice agent, running as a macOS launchd daemon |
| [`homebrew-tap`](https://github.com/vhco-pro/homebrew-tap) | Homebrew tap for the org's macOS tools |

## Contributing

Most repositories here accept pull requests in the normal way. The seven
`stackweaver-*` distribution mirrors are the exception - they are
bot-synced and take issues instead. See
[`CONTRIBUTING.md`](https://github.com/vhco-pro/.github/blob/main/CONTRIBUTING.md)
for which is which, and
[`SUPPORT.md`](https://github.com/vhco-pro/.github/blob/main/SUPPORT.md)
for where to ask questions.

## Security

- Reporting channel, scope, and disclosure process:
  [`SECURITY.md`](https://github.com/vhco-pro/.github/blob/main/SECURITY.md)
- Release artifacts are signed with cosign keyless (Sigstore) and carry
  SLSA build provenance. There are no long-lived signing keys in this
  org.

## Licence

Licences vary per repository - check the `LICENSE` file in each. In
summary:

- **Stackweaver core** (API, Orchestrator, Frontend, Helm chart, Zitadel
  init, shared `core/` module): **BSL 1.1**, with an Apache-2.0 Change
  Date and a SaaS-exclusion Additional Use Grant.
- **Stackweaver ecosystem tooling** (the runners): **Apache-2.0**,
  shipping a `NOTICE` disclosing the BSL upstream linkage.
  `terraform-provider-stackweaver` is **MPL-2.0**, inherited from its
  upstream.
- **Everything else**: **Apache-2.0**, except `builders` (**GPL-3.0**).

## Trademark

Stackweaver™ is a trademark of VH & Co. The Stackweaver name and word
mark identify the official project; the source-code licences above do not
grant any right to use the mark in product names, hosted services, or
company names. See the
[Trademark Policy](https://github.com/vhco-pro/.github/blob/main/TRADEMARK.md)
for the full terms.

---

Maintained by **VH & Co BV** · [contact@vhco.pro](mailto:contact@vhco.pro)
