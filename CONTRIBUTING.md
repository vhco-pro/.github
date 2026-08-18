# Contributing to vhco-pro projects

Thanks for your interest in contributing! This document applies to every
repository in the [`vhco-pro`](https://github.com/vhco-pro) organisation.

Repositories here fall into two categories, and which one you are looking
at determines how contributions reach it. Check this first, because it is
the one genuinely unusual thing about this org.

## Two kinds of repository

### 1. Standard repositories — pull requests welcome

Most repositories in this org are developed in the open, in place, and
accept pull requests in the normal way: fork, branch, PR. This includes
`distil`, `stackgraph`, `terraform-provider-garage`,
`terraform-provider-stackweaver`, `stackweaver-operator`, `ssm-connect`,
`dcv-session-agent`, `claude-companion`, `postbode`, `builders`,
`swift-release-action`, and `homebrew-tap`.

The workflow is:

1. Open an issue describing the bug or the feature, unless the change is
   trivial. For anything substantial, agreeing on the approach before
   code is written saves everyone time.
2. Fork the repository and create a branch.
3. Make the change, with tests, and make sure the repository's CI passes.
4. Open a pull request using the PR template.

### 2. Stackweaver distribution mirrors — no pull requests

Stackweaver itself is developed in a **private monorepo**
(`michielvha/stackweaver`) and published to **seven public satellite
repositories** in this org: `stackweaver-api`, `stackweaver-frontend`,
`stackweaver-orchestrator`, `stackweaver-zitadel-init`,
`stackweaver-opentofu-runner`, `stackweaver-ansible-runner`, and
`stackweaver-helm`.

These satellites are one-way distribution mirrors. Their content is a
deterministic re-publication of a reviewed monorepo commit, signed by the
`stackweaver-release-bot` GitHub App. This means:

- **You cannot send a PR to a satellite repo.** Any push to a satellite
  is overwritten by the next monorepo sync. The sync bot is the only
  writer; no human holds push rights to satellite `main`.
- **All Stackweaver source changes happen in the monorepo.** Maintainers
  triage satellite-side issues and port accepted changes upstream.

What you *can* do on a satellite: file issues, report vulnerabilities via
Private Vulnerability Reporting, open Discussions, and suggest
documentation improvements. Those are all valuable and all get routed
upstream.

Note that `stackweaver-operator` and `terraform-provider-stackweaver` are
part of the Stackweaver ecosystem but are **not** mirrors — they are
developed in place and take pull requests like any standard repository.

## Code of Conduct

All repositories in this org follow the
[Contributor Covenant 2.1](./CODE_OF_CONDUCT.md). By participating you
agree to its terms.

## Quality expectations

These apply to pull requests on standard repositories, and to changes
ported upstream from satellite issues:

- **Tests.** Any non-trivial behaviour change includes new or updated
  automated tests, in whatever framework the repository already uses.
- **Lint clean.** Every repository runs linters in CI (`golangci-lint`
  and `govulncheck` for Go, ESLint for JS/TS, SwiftLint for Swift, and so
  on). PRs must be green before merge.
- **Conventional commits.** `feat(scope): subject`,
  `fix(scope): subject`, `refactor(scope): subject`, and friends.
  Several repositories derive their release version from commit
  messages, so this is load-bearing, not cosmetic.
- **No copyrighted material** you do not have the right to contribute.
- **No secrets.** Secret scanning is enforced.

## Licensing

Licences differ per repository — check the `LICENSE` file in the one you
are contributing to. By opening a pull request you agree to license your
contribution under that repository's licence.

For Stackweaver specifically: contributions to the BSL-licensed
components (api, orchestrator, frontend, helm chart, zitadel-init, and
the shared `core/` module) are accepted under BSL 1.1 with the
Stackweaver SaaS-exclusion Additional Use Grant. Contributions to the
Apache-licensed components (the runners) are accepted under Apache-2.0.

## Trademarks

Contributing code does not grant any right to use the project's names or
marks. See [`TRADEMARK.md`](./TRADEMARK.md).

## Questions?

See [`SUPPORT.md`](./SUPPORT.md), or ping the maintainer (@michielvha).
