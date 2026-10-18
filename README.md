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
computes a GitVersion version, pushes its exact `v`-prefixed tag, and publishes
a GitHub Release with Conventional Commit notes rendered by a checksum-verified
`git-cliff`. It never creates or moves major/minor alias tags. A retry can finish
publishing when a previous run already pushed the tag.

Existing release descriptions are preserved on reruns. New notes emphasize
consumer changes, breaking changes, and dependency updates rather than routine
README/test/CI edits; the comparison link retains access to the full history.

This repository is public, so callers outside `vhco-pro`, including repositories
owned by `michielvha`, can use it. The job operates on the caller repository using
the caller's `GITHUB_TOKEN`; no personal access token or `secrets: inherit` is
needed. Both the caller and the reusable job must grant `contents: write`.
The caller's Actions policy must allow this repository and its public dependencies.

```yaml
name: Release
on:
  push:
    branches: [main]
permissions: {}
concurrency:
  group: release
  cancel-in-progress: false
jobs:
  release:
    permissions:
      contents: write
    uses: vhco-pro/.github/.github/workflows/reusable-release.yml@FULL_COMMIT_SHA
    with:
      workflow-ref: FULL_COMMIT_SHA
```

Use an immutable commit SHA in production, as with the other shared workflows.
  Replace both `FULL_COMMIT_SHA` placeholders with the same verified full commit
  SHA. The required `workflow-ref` checks out the helper scripts and config from
  that revision; the default checkout is the caller, not this repository. Helper
  source is moved to the runner's temporary directory before use. A placeholder
  or moving branch is not a valid `workflow-ref`.

  Release logic lives in [scripts/release/](scripts/release/) and note configuration
  in [release/cliff.toml](release/cliff.toml). Run `bash tests/test-release.sh` and
  shellcheck on those scripts to verify behavior without publishing. Repository
  instructions and the central conventions start in [AGENTS.md](AGENTS.md).

  Optional inputs are `runs-on` (a Linux runner, default `ubuntu-latest`) and
`dry-run` (boolean, default `false`). A dry run still computes the version and
renders prospective notes in the job summary, but creates no tag or release.
Workflow outputs are `version`, `tag`, `tagged`, and `release-url`.

### Releasing this repository

[.github/workflows/release.yml](.github/workflows/release.yml) calls this repo's
own reusable release workflow at the same commit. Script checks run first.
Relevant changes pushed to `main` publish an exact GitVersion tag and GitHub
Release; PRs preview the same flow without publication. Manual runs default to
dry-run. No separate patch counter or moving major/minor tagger is used.

## Tag-driven Go releases

[reusable-go-verify.yml](.github/workflows/reusable-go-verify.yml) runs formatting,
module tidiness, `go vet`, tests, builds, and blocking lint checks. It reads the
toolchain from the caller's `go.mod`. Optional inputs are `go-version-file` and
`golangci-lint-config`.

[reusable-go-release.yml](.github/workflows/reusable-go-release.yml) accepts an
exact `v`-prefixed semantic tag push, checks out that event's commit, runs the
shared verification workflow, and calls the existing GoReleaser action. It does
not run GitVersion or create tags. Existing action consumers and the notes-only
release workflow are unchanged.

```yaml
name: Go release
on:
  push:
    tags: ['v*']
permissions: {}
jobs:
  release:
    permissions:
      contents: write
    uses: vhco-pro/.github/.github/workflows/reusable-go-release.yml@FULL_COMMIT_SHA
    with:
      workflow-ref: FULL_COMMIT_SHA
    secrets:
      pgp-private-key: ${{ secrets.PGP_PRIVATE_KEY }}
```

Replace both SHA placeholders with the same published, verified workflow commit.
The tag glob is only an event filter; the helper rejects aliases and invalid
semantic versions. Caller source and helper source use separate immutable refs.
The caller supplies its GoReleaser and lint configuration, including artifact
signing. Publication creates a draft first and exposes it only after GoReleaser
succeeds. Prereleases never become the latest stable release.

A retry can rebuild only while no release exists for the tag. An existing draft
or public release stops the run instead of replacing its artifacts. Inspect
partial drafts and recover their original artifacts explicitly, or release a new
version. Never move the source tag. No cleanup runs automatically.

The caller owns tag creation. Use an event-emitting GitHub App identity for
automated tag pushes; the ordinary `GITHUB_TOKEN` does not trigger another tag-push
workflow. Keep PR version calculation read-only. Trusted feature-branch previews
can be requested manually and published through their prerelease tags.

Run `bash tests/test-go-release.sh`, the existing release tests, and shellcheck.
Shared release policy belongs in the [knowledge base](https://github.com/vhco-pro/docs),
not in copies of the policy in each consumer.

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
