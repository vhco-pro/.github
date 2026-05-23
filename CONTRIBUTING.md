# Contributing to Stackweaver

Thanks for your interest in contributing! This document covers how
contributions flow into Stackweaver, including the multi-repo topology
that makes contributing here slightly unusual.

## How the repositories fit together

Stackweaver is developed in a **private monorepo**
(`michielvha/stackweaver`) and published to **seven public satellite
repositories** under the `vhco-pro` org (this org). Satellites are
one-way distribution mirrors — their content is a deterministic
re-publication of a reviewed monorepo commit, signed by the
`stackweaver-release-bot` GitHub App.

This means:

- **You cannot send a PR directly to a satellite repo.** Any push to a
  satellite is overwritten by the next monorepo sync. Satellites do not
  accept human pushes; the sync bot is the only writer.
- **All source changes happen in the monorepo.** Maintainers triage
  satellite-side issues and port accepted changes upstream.

## What you can do here

- **File issues** on the satellite repo most closely related to your
  problem. The repo's issue templates will route you correctly.
- **Report security vulnerabilities** via GitHub Private Vulnerability
  Reporting on the relevant satellite (see [`SECURITY.md`](./SECURITY.md)).
- **Open Discussions** on the satellite repo for design questions,
  feature ideas, or "how do I…" questions.
- **Suggest documentation improvements** — these are valuable. File an
  issue on the satellite that owns the docs you care about; if the
  change is small the maintainer will land it directly upstream.

## What you cannot do here

- Send a PR that modifies code in a satellite repo. The next monorepo
  sync would clobber it.
- Push to satellite `main` — branch ownership is bot-only by design.

## Code of Conduct

This project follows the [Contributor Covenant 2.1](./CODE_OF_CONDUCT.md).
By participating you agree to its terms.

## Quality expectations for upstream changes

Once a contributed change reaches the monorepo, maintainers expect:

- **Tests.** Any non-trivial behaviour change includes new or updated
  automated tests (Go unit tests, Vitest, or Playwright as appropriate).
- **Lint clean.** The monorepo runs `golangci-lint`, ESLint, and
  `govulncheck` in CI; PRs must be green before merge.
- **Conventional commits.** `feat(scope): subject`, `fix(scope): subject`,
  `refactor(scope): subject`, etc.
- **No copyrighted material** you don't have the right to contribute.
- **No secrets.** Secret scanning is enforced on monorepo PRs.

## Licence flow-through

Contributions to BSL-licensed parts of Stackweaver (api, orchestrator,
frontend, helm chart, zitadel-init, `core/`) are accepted under BSL 1.1
with the Stackweaver SaaS-exclusion Additional Use Grant.
Contributions to Apache-licensed parts (the runners) are accepted under
Apache-2.0. See `SECURITY.md` and the monorepo's licence-strategy doc.

## Questions?

Open a Discussion on the satellite repo nearest to your question, or
ping the maintainer (@michielvha).
