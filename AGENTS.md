# Repository instructions

Read the personal conventions at
https://github.com/vhco-pro/docs/blob/5f1133d65a8c605099543ca25688e9b180a068b5/AGENTS.md
and follow the relevant routes. For shared workflows, read `conventions/coding.md`,
`conventions/delivery.md`, and `knowledge/shared-delivery.md` at that revision.
For docs changes, also read `conventions/documentation.md`. Report inaccessible
sources or conflicts rather than guessing.

## Shared release checks

Run `shellcheck scripts/release/*.sh tests/*.sh`, `bash tests/test-release.sh`,
and `actionlint .github/workflows/reusable-release.yml .github/workflows/release-preview.yml`.
The tests download and verify the pinned git-cliff and use disposable Git repos
and a mocked GitHub CLI; they do not publish a release.

The reusable release workflow checks out helpers from the caller-supplied
`workflow-ref`, which must be the same full SHA as its workflow pin. Helpers
and configuration belong in this repo, never the caller's root. Substantial
logic lives in scripts. Preserve dry runs, exact-only tags, caller-scoped tokens,
and retry behavior. Existing unrelated workflows are not part of this refactor.
Do not merge or release without explicit authorization.