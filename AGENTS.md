# Repository instructions

Read the personal conventions at
https://github.com/vhco-pro/docs/blob/93d832328a547aefc9b793bb444d6e9171aa2451/AGENTS.md
and follow the relevant routes. For shared workflows, read `conventions/coding.md`,
`conventions/delivery.md`, and `knowledge/shared-delivery.md` at that revision.
For docs changes, also read `conventions/documentation.md`. Report inaccessible
sources or conflicts rather than guessing.

## Shared release checks

Run `go test -count=1 ./...` and
`actionlint .github/workflows/reusable-go-release.yml .github/workflows/release.yml`.
The Go workflow contracts use `go.yaml.in/yaml/v3`; no JavaScript test runtime
is required.

The reusable release workflow checks out helpers from the caller-supplied
`workflow-ref`, which must be the same full SHA as its workflow pin. Helpers
and configuration belong in this repo, never the caller's root. Substantial
logic lives in scripts. Preserve dry runs, exact-only tags, caller-scoped tokens,
and retry behavior. Existing unrelated workflows are not part of this refactor.
Do not merge or release without explicit authorization.