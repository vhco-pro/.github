# vhco-pro/.github

Org-wide community health files and reusable workflows for all
[vhco-pro](https://github.com/vhco-pro) repositories.

GitHub uses these files for repositories owned by `vhco-pro`. A valid local
issue form or chooser configuration replaces the entire inherited chooser.
GitHub resolves security, support, and contribution files separately.
Product-specific policies belong with the product.

## What lives here

| Path | Purpose |
|------|---------|
| `profile/README.md`                       | Organization profile and project index |
| `SECURITY.md`                             | Default private reporting instructions |
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

Release workflows live in their own repositories; copy their examples:

| Repository | Use |
| --- | --- |
| [vhco-pro/gh-release-workflow](https://github.com/vhco-pro/gh-release-workflow) | Tag and GitHub Release, for repositories whose tag is the release |
| [vhco-pro/version-workflow](https://github.com/vhco-pro/version-workflow) | GitVersion tag after the caller's checks, with labelled PR previews |
| [vhco-pro/go-release-workflow](https://github.com/vhco-pro/go-release-workflow) | Go checks and signed GoReleaser release for an exact tag |

They replace this repository's former `reusable-release.yml` and
`reusable-go-release.yml`.

### Releasing this repository

[.github/workflows/release.yml](.github/workflows/release.yml) lints the
workflows with actionlint, then calls `gh-release-workflow`. Relevant changes
pushed to `main` publish an exact GitVersion tag and GitHub Release; PRs preview
the same flow without publication. Manual runs default to dry-run. Repository
instructions and the central conventions start in [AGENTS.md](AGENTS.md).

## Contributions

See the [contributing guide](CONTRIBUTING.md) for where to propose changes and
the [Stackweaver release sync architecture](https://stackweaver.sh/docs/security/sync-architecture)
for how its public distributions are maintained.
