<div align="center">

<img src="https://raw.githubusercontent.com/vhco-pro/.github/main/profile/assets/vhco.png" alt="VH & Co" width="150" />

# VH & Co

Open-source infrastructure and IaC tooling. Makers of [Stackweaver](https://sw.vhco.pro).

[![Website](https://img.shields.io/badge/website-vhco.pro-06b6d4)](https://vhco.pro)
[![Docs](https://img.shields.io/badge/docs-sw.vhco.pro-3b82f6)](https://sw.vhco.pro/docs)
[![Contributing](https://img.shields.io/badge/contributing-guide-6366f1)](https://github.com/vhco-pro/.github/blob/main/CONTRIBUTING.md)
[![Security](https://img.shields.io/badge/security-policy-a855f7)](https://github.com/vhco-pro/.github/blob/main/SECURITY.md)

</div>

## Stackweaver

<img src="https://raw.githubusercontent.com/vhco-pro/.github/main/profile/assets/stackweaver.png" alt="Stackweaver" width="110" align="right" />

A self-hostable alternative to Terraform Cloud and Ansible AWX. It runs OpenTofu and Ansible from one control plane and is API-compatible with Terraform Cloud/Enterprise tooling.

```bash
helm install stackweaver oci://ghcr.io/vhco-pro/charts/stackweaver --version <X.Y.Z>
```

[Documentation](https://sw.vhco.pro/docs) · [Deployment guide](https://github.com/vhco-pro/stackweaver-helm#readme) · [Verifying releases](https://sw.vhco.pro/docs/security/verifying-releases)

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

The `stackweaver-*` repositories above are release mirrors published from the Stackweaver source tree. They take issues, not pull requests.

**Ecosystem** - developed in the open, pull requests welcome:

| Repo | What it is |
|------|------------|
| [`terraform-provider-stackweaver`](https://github.com/vhco-pro/terraform-provider-stackweaver) | Terraform provider for Stackweaver, derived from `terraform-provider-tfe` and kept in sync with it |
| [`stackweaver-operator`](https://github.com/vhco-pro/stackweaver-operator) | Kubernetes operator for managing Stackweaver deployments |
| [`stackweaver-registry`](https://github.com/vhco-pro/stackweaver-registry) | Multi-format artifact registry with upstream caching, SSO and RBAC |

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

## Contributing and security

Pull requests are welcome everywhere except the Stackweaver release mirrors. Start with [CONTRIBUTING](https://github.com/vhco-pro/.github/blob/main/CONTRIBUTING.md), ask questions through [SUPPORT](https://github.com/vhco-pro/.github/blob/main/SUPPORT.md), and report vulnerabilities through [SECURITY](https://github.com/vhco-pro/.github/blob/main/SECURITY.md). Releases are signed with cosign keyless (Sigstore) and carry SLSA build provenance.

<sub>Licences vary per repository, see each `LICENSE`. Stackweaver™ is a trademark of VH & Co ([policy](https://github.com/vhco-pro/.github/blob/main/TRADEMARK.md)). [contact@vhco.pro](mailto:contact@vhco.pro)</sub>
