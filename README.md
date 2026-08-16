# vhco-pro/.github

Org-wide community health files and reusable workflows for all
[vhco-pro](https://github.com/vhco-pro) repositories.

GitHub renders files in an org's `.github` repo as defaults for every
repo in the org that doesn't override them. The files here therefore
apply to all repos automatically — the Stackweaver satellites, the
Terraform providers, the container and release tooling, and the macOS
apps alike.

**Everything here must be project-agnostic.** Org defaults are inherited
all-or-nothing per repository: there is no way to serve a different
default to a subset of repos. Anything that only holds for one project
belongs in a clearly-scoped section (see the "Project-specific" part of
`SECURITY.md`) or in that repository's own override.

## What lives here

| Path | Purpose |
|------|---------|
| `profile/README.md`                       | Landing page for the `vhco-pro` org — indexes every public project |
| `SECURITY.md`                             | Org-wide reporting channel, scope, and disclosure process, plus a Stackweaver-specific annex |
| `CONTRIBUTING.md`                         | Contribution flow for both repo models: standard PR-taking repos, and the bot-synced Stackweaver mirrors |
| `CODE_OF_CONDUCT.md`                      | Contributor Covenant 2.1 |
| `SUPPORT.md`                              | How to get help |
| `TRADEMARK.md`                            | Trademark scope for the org, and the Stackweaver™ Trademark Policy |
| `.github/ISSUE_TEMPLATE/*.yml`            | Inherited issue forms |
| `.github/PULL_REQUEST_TEMPLATE.md`        | Inherited PR template |
| `.github/dependabot.yml`                  | Dependabot config for this repo itself |
| `.github/workflows/reusable-*.yml`        | Reusable workflows called by repos across the org |
| `workflow-templates/*`                    | Templates for bootstrapping new repos |

## The two repository models

Contribution guidance in this repo has to hold for both:

- **Standard repositories** — developed in place, accept pull requests:
  `distil`, `stackgraph`, `terraform-provider-garage`,
  `terraform-provider-stackweaver`, `stackweaver-operator`,
  `ssm-connect`, `dcv-session-agent`, `claude-companion`, `postbode`,
  `builders`, `swift-release-action`, `homebrew-tap`.
- **Stackweaver distribution mirrors** — one-way, bot-synced, no human
  PRs: `stackweaver-api`, `stackweaver-frontend`,
  `stackweaver-orchestrator`, `stackweaver-zitadel-init`,
  `stackweaver-runner`, `stackweaver-ansible-runner`,
  `stackweaver-helm`. These are exactly the repos with a `sync-*.yml`
  workflow in the Stackweaver monorepo.

When a new repository joins the org, check that the inherited templates
still read correctly for it, and update the lists above and in
`profile/README.md` and `CONTRIBUTING.md`.

## Source of truth

The Stackweaver-specific content here originates in the **private**
Stackweaver monorepo (`michielvha/stackweaver`), where the canonical
security plan and supporting documentation live:

- `docs/internal/security/osps-baseline-audit.md` — OSPS audit tracker
- `docs/internal/security/threat-model.md`
- `docs/internal/security/vuln-management-policy.md`
- `docs/internal/security/secrets-policy.md`
- `docs/internal/security/access-policy.md`
- `docs/internal/security/release-verification.md`
- `docs/internal/security/license-strategy.md`
- `docs/internal/security/core-auditor-access.md`
- `docs/architecture/repositories.md`

Those sections are kept in lockstep with the canonical sources. The
org-wide, project-agnostic content in this repo is authored here.
