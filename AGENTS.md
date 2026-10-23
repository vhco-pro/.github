# Repository instructions

Read the personal conventions at
https://github.com/vhco-pro/docs/blob/93d832328a547aefc9b793bb444d6e9171aa2451/AGENTS.md
and follow the relevant routes. For shared workflows, read `conventions/coding.md`,
`conventions/delivery.md`, and `knowledge/shared-delivery.md` at that revision.
For docs changes, also read `conventions/documentation.md`. Report inaccessible
sources or conflicts rather than guessing.

## Shared release checks

Run `go tool actionlint -ignore SC2012 .github/workflows/*.yml`. The SC2012
note in `reusable-codeql-swift.yml` predates the check.

Do not merge or release without explicit authorization.
