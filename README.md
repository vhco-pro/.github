# vhco-pro/.github

Org-default community-health files and reusable workflows for the
[Stackweaver](https://github.com/vhco-pro) satellite repositories.

GitHub renders files in an org's `.github` repo as defaults for every
repo in the org that doesn't override them. The files here therefore
apply to all `vhco-pro/stackweaver-*` repos automatically.

## What lives here

| Path | Purpose |
|------|---------|
| `profile/README.md`                       | Landing page for the `vhco-pro` org |
| `SECURITY.md`                             | Security contacts, PVR pointer, supported versions, verification |
| `CONTRIBUTING.md`                         | Contribution flow: monorepo source-of-truth + satellite distribution |
| `CODE_OF_CONDUCT.md`                      | Contributor Covenant 2.1 |
| `SUPPORT.md`                              | How to get help |
| `.github/ISSUE_TEMPLATE/*.yml`            | Inherited issue forms |
| `.github/PULL_REQUEST_TEMPLATE.md`        | Inherited PR template |
| `.github/dependabot.yml`                  | Dependabot config for this repo itself |
| `.github/workflows/reusable-*.yml`        | Reusable workflows called by each satellite |
| `workflow-templates/*`                    | Templates for bootstrapping new repos |

## Source of truth

This repo's content originates in the **private** Stackweaver monorepo
(`michielvha/stackweaver`). The canonical security plan and supporting
documentation live there:

- `docs/internal/security/osps-baseline-audit.md` — OSPS audit tracker
- `docs/internal/security/threat-model.md`
- `docs/internal/security/vuln-management-policy.md`
- `docs/internal/security/secrets-policy.md`
- `docs/internal/security/access-policy.md`
- `docs/internal/security/release-verification.md`
- `docs/internal/security/license-strategy.md`
- `docs/internal/security/core-auditor-access.md`
- `docs/architecture/repositories.md`

Files in this repo are kept in lockstep with the canonical sources.
