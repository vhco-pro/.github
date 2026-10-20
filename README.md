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
| `examples/*/*/`                           | Copyable caller workflows, one folder per arrangement |

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

### Adopt it

Copy one example folder into the repository's `.github/workflows/` and replace
`FULL_COMMIT_SHA` with a released commit of this repository. Use one arrangement
per repository; enabling both releases the same tag twice.

| Example | Use when |
| --- | --- |
| [examples/go-release/separate-workflows](examples/go-release/separate-workflows) | A tag push, including one pushed by hand, should start the release. Needs a GitHub App for the tag push |
| [examples/go-release/same-workflow](examples/go-release/same-workflow) | No GitHub App. `GITHUB_TOKEN` pushes the tag and a dependent job releases it |

Both run CI on branch pushes, tag `main` and open PRs labelled `release:preview`
after the checks pass, and calculate the version only on other branches. A
private repository that runs CI only on PRs changes the `on:` triggers; a
`pull_request` run never tags, because it has GitHub's merge ref instead of the
branch.

`tests/examples_test.go` fails when an example passes an input, secret or reads an
output that the shared workflows no longer declare.

### The shared workflows

[reusable-go-release.yml](.github/workflows/reusable-go-release.yml) checks out
the tag, calls the Go CI action, then calls the release action only if CI
passes. It runs for a tag push event, or for the `tag` input that a
same-workflow caller passes. Other optional inputs are `go-version-file` and
`golangci-lint-config`.

[reusable-version.yml](.github/workflows/reusable-version.yml) authenticates
when tagging, checks out the branch with full history, fails if the branch has
moved past the checked commit, and calls `michielvha/gitversion-tag-action`.
GitVersion and its configuration own versions and tag pushes. It inherits the
caller job's token permissions and exposes `tagged`, `tag` and `version`.
`dry-run` is the default and mints no token.

CI-only jobs call `michielvha/go-ci-action` directly; there is no shared CI-only
workflow.

In the same-workflow arrangement, retry a failed release with **Re-run failed
jobs**. **Re-run all jobs** finds the tag already pushed and releases nothing, and
a tag pushed by hand publishes nothing.

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
