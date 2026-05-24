# Stackweaver

> Multi-IaC orchestration platform — a self-hostable alternative to
> Terraform Cloud / Ansible AWX, supporting Terraform, OpenTofu and
> Ansible from a single control plane.

## Distribution

Stackweaver is developed in a private monorepo and published here as
seven independent satellite repositories. Pick the one you need:

| Component | Repo |
|-----------|------|
| Helm chart                | [`stackweaver-helm`](https://github.com/vhco-pro/stackweaver-helm)                       |
| Backend API               | [`stackweaver-api`](https://github.com/vhco-pro/stackweaver-api)                         |
| Orchestrator              | [`stackweaver-orchestrator`](https://github.com/vhco-pro/stackweaver-orchestrator)       |
| Terraform Runner          | [`stackweaver-runner`](https://github.com/vhco-pro/stackweaver-runner)                   |
| Ansible Runner            | [`stackweaver-ansible-runner`](https://github.com/vhco-pro/stackweaver-ansible-runner)   |
| Frontend                  | [`stackweaver-frontend`](https://github.com/vhco-pro/stackweaver-frontend)               |
| Zitadel bootstrap         | [`stackweaver-zitadel-init`](https://github.com/vhco-pro/stackweaver-zitadel-init)       |

## Quickstart

```bash
helm install stackweaver oci://ghcr.io/vhco-pro/charts/stackweaver --version <X.Y.Z>
```

See the helm satellite README for the full deployment guide.

## Security

- Security contacts and the private-vulnerability-reporting channel:
  [`SECURITY.md`](https://github.com/vhco-pro/.github/blob/main/SECURITY.md)
- How to verify a release was actually produced by us:
  [`profile/security/verifying-releases.md`](https://github.com/vhco-pro/.github/blob/main/profile/security/verifying-releases.md)
  *(coming with the first signed release)*

## Licence

- Original Stackweaver code (API, Orchestrator, Frontend, Helm chart,
  Zitadel init, shared `core/` module): **BSL 1.1** with Apache-2.0
  Change Date and a SaaS-exclusion Additional Use Grant.
- Ecosystem tooling (Terraform Runner, Ansible Runner, future
  Terraform provider): **Apache-2.0**; ships with a `NOTICE` file
  disclosing the BSL upstream linkage.

The Stackweaver project is maintained by **VH & Co BV**.

## Trademark

Stackweaver™ is a trademark of VH & Co. The Stackweaver name and word
mark identify the official project; the source-code licences above do
not grant any right to use the mark in product names, hosted services,
or company names. See the [Trademark Policy](https://github.com/vhco-pro/.github/blob/main/TRADEMARK.md)
for the full terms and the list of permitted and prohibited uses.
