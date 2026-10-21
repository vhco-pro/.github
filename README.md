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

## Shared releases

[.github/workflows/reusable-release.yml](.github/workflows/reusable-release.yml)
is deprecated. Use
[vhco-pro/gh-release-workflow](https://github.com/vhco-pro/gh-release-workflow)
and copy its example. This repository releases itself with it. Refs that already
pin `reusable-release.yml` keep working; it is removed in the next major release.

### Releasing this repository

[.github/workflows/release.yml](.github/workflows/release.yml) runs the script
and Go tests, then calls `gh-release-workflow`. Relevant changes pushed to `main`
publish an exact GitVersion tag and GitHub Release; PRs preview the same flow
without publication. Manual runs default to dry-run. Repository instructions and
the central conventions start in [AGENTS.md](AGENTS.md).

## Go delivery pilot

The pilot composes intent-focused actions: `michielvha/go-ci-action` for blocking
Go checks, `michielvha/goreleaser-action` for artifact publication, and
`michielvha/gitversion-tag-action` called directly for versioning and tagging.
Shared workflows own checkout, authentication, and ordering, not copies of those
actions' commands or calls into their internal scripts.

The release-only next major is isolated from existing GoReleaser action
consumers. Shared release guidance belongs in the
[knowledge base](https://github.com/vhco-pro/docs).

[reusable-go-release.yml](.github/workflows/reusable-go-release.yml) checks out
the tag event's commit, calls the Go CI action, then calls the release action
only if CI passes. It contains no shell commands or helper checkouts.

```yaml
name: Go release
on:
  push:
    tags: ['v*.*.*']
permissions: {}
jobs:
  release:
    permissions:
      contents: write
    uses: vhco-pro/.github/.github/workflows/reusable-go-release.yml@FULL_COMMIT_SHA
    secrets:
      pgp-private-key: ${{ secrets.PGP_PRIVATE_KEY }}
```

Use a verified commit in place of the placeholder. Optional inputs are
`go-version-file` and `golangci-lint-config`. There is no `workflow-ref` input:
the Go workflow has no scripts to retrieve. The notes-only workflow's existing
contract is unchanged.

Tag creation remains separate. Call `michielvha/gitversion-tag-action` directly;
do not invoke its internal scripts. Use an App token for automated tag pushes
that must start the tag-event publisher. PR version previews use `dry-run: true`
and receive no publication credentials.

GoReleaser owns release creation, upload, and publication. Keep asset and draft
replacement disabled in the caller's configuration and inspect partial
publication before retrying. No cleanup is enabled by this workflow.

The Go workflow contracts run with `go test -count=1 ./...` and use
`go.yaml.in/yaml/v3`. They check action ordering, source-commit identity,
signing-secret scope, immutable refs, and release serialization without freezing
unrelated input additions or step labels. The notes-only release script tests
remain Bash.

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
  `stackweaver-secrets-init`, `stackweaver-opentofu-runner`,
  `stackweaver-ansible-runner`, `stackweaver-helm`. These are exactly the
  repos with a `sync-*.yml` workflow in the Stackweaver monorepo.

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
