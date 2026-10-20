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

## Go delivery pilot

The pilot composes intent-focused actions: `michielvha/go-ci-action` for blocking
Go checks, `michielvha/goreleaser-action` for artifact publication, and
`michielvha/gitversion-tag-action` called directly for versioning and tagging.
Shared workflows own checkout, authentication, and ordering, not copies of those
actions' commands or calls into their internal scripts.

Legacy GoReleaser consumers retain their original immutable refs. Shared
release guidance belongs in the
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

Use a verified commit in place of the placeholder. Optional inputs are `tag`
(same-workflow arrangement, below),
`go-version-file` and `golangci-lint-config`. There is no `workflow-ref` input:
the Go workflow has no scripts to retrieve. The notes-only workflow's existing
contract is unchanged.

### Caller checks and versioning

CI-only jobs call `michielvha/go-ci-action` directly; there is no shared CI-only
workflow. Callers own their triggers and repository-specific checks.

[reusable-version.yml](.github/workflows/reusable-version.yml) is a thin wrapper:
App authentication when tagging, full-history checkout of the branch, a guard
that the branch has not moved past the checked commit, and one call to
`michielvha/gitversion-tag-action`. GitVersion and its config own versions and
tag pushes.

Call it from a branch `push` (or `workflow_dispatch`) run, after the caller's
checks through `needs:`, so GitVersion sees the real branch. A `pull_request`
run has GitHub's merge ref instead. Pass `dry-run: false` only when the caller
decides the branch should tag; dry-run is the default and mints no token.
The wrapper inherits the caller job's token permissions and exposes `tagged`,
`tag` and `version` outputs.

Choose one of two arrangements per repository:

**Separate workflows.** The caller grants `contents: read` and passes App
credentials. The App's tag push starts the caller's tag-triggered release
workflow, which calls `reusable-go-release.yml`. A tag pushed by hand releases
the same way.

```yaml
  version:
    needs: release-gate          # main, or an open release:preview PR at this commit
    permissions:
      contents: read
    uses: vhco-pro/.github/.github/workflows/reusable-version.yml@FULL_COMMIT_SHA
    with:
      dry-run: ${{ needs.release-gate.outputs.tag != 'true' }}
      release-app-client-id: ${{ vars.RELEASE_APP_CLIENT_ID }}
    secrets:
      release-app-private-key: ${{ secrets.RELEASE_APP_PRIVATE_KEY }}
```

**Same workflow.** The caller grants `contents: write` and passes no App. A tag
pushed with `GITHUB_TOKEN` starts no other workflow, so a dependent job releases
the tag this run created:

```yaml
  version:
    needs: release-gate
    permissions:
      contents: write
    uses: vhco-pro/.github/.github/workflows/reusable-version.yml@FULL_COMMIT_SHA
    with:
      dry-run: ${{ needs.release-gate.outputs.tag != 'true' }}

  release:
    needs: version
    if: ${{ needs.version.outputs.tagged == 'true' }}
    permissions:
      contents: write
    uses: vhco-pro/.github/.github/workflows/reusable-go-release.yml@FULL_COMMIT_SHA
    with:
      tag: ${{ needs.version.outputs.tag }}
    secrets:
      pgp-private-key: ${{ secrets.PGP_PRIVATE_KEY }}
```

To retry a failed release in the same-workflow arrangement, use **Re-run failed
jobs**: it reuses this attempt's `tagged` and `tag`. **Re-run all jobs** finds the
tag already pushed, reports `tagged: false` and skips the release. A tag pushed
by hand publishes nothing in this arrangement. GoReleaser never replaces an
existing release's assets in either arrangement.

Do not enable both arrangements for the same repository; one tag would be
released twice.

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
