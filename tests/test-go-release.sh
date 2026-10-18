#!/usr/bin/env bash
set -euo pipefail

SOURCE_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
TEMP=$(mktemp -d)
trap 'rm -r "$TEMP"' EXIT
export GITHUB_EVENT_NAME=push GITHUB_REF_TYPE=tag
export GITHUB_REPOSITORY=example/project
export GITHUB_EVENT_PATH="$TEMP/event.json"
export GITHUB_OUTPUT="$TEMP/output" GITHUB_STEP_SUMMARY="$TEMP/summary"
export GH_LOG="$TEMP/gh.log" GH_PAYLOAD='[[]]' GH_FAIL=false
printf '{"deleted":false}\n' > "$GITHUB_EVENT_PATH"
touch "$GITHUB_OUTPUT" "$GITHUB_STEP_SUMMARY" "$GH_LOG"

gh() {
  printf '%s\n' "$*" >> "$GH_LOG"
  [[ "$GH_FAIL" == false ]] || return 1
  case "$1" in
    api) printf '%s\n' "$GH_PAYLOAD" ;;
    release) [[ "$2" == edit ]] ;;
    *) return 1 ;;
  esac
}
export -f gh

git init -q -b main "$TEMP/repo"
cd "$TEMP/repo"
git config user.name 'Go Release Test'
git config user.email test@example.com
git commit -q --allow-empty -m 'feat: release fixture'
export GITHUB_SHA
GITHUB_SHA=$(git rev-parse HEAD)

run() {
  bash "$SOURCE_ROOT/scripts/release/go-release.sh" "$1"
}
reject() {
  : > "$GH_LOG"
  if run "$1" > "$TEMP/error" 2>&1; then
    printf 'Expected rejection for %s (%s).\n' "$GITHUB_REF_NAME" "$1" >&2
    exit 1
  fi
  if grep -q '^release edit ' "$GH_LOG"; then
    printf 'Rejected request mutated a release.\n' >&2
    exit 1
  fi
}

# Publication builds the exact commit identified by a validated stable or prerelease semantic tag.
for TAG in v1.2.3 v0.0.0 v1.2.3-feature.demo.1 v1.2.3-rc.0 v1.2.3+build.01 v1.2.3+build-with-hyphen; do
  git tag "$TAG"
  export GITHUB_REF_NAME="$TAG"
  run preflight
done
for TAG in v1 v1.2 1.2.3 v01.2.3 v1.02.3 v1.2.03 v1.2.3-01 v1.2.3-rc.01 v1.2.3- v1.2.3+ v1.2.3-rc..1; do
  export GITHUB_REF_NAME="$TAG"
  reject preflight
done
export GITHUB_REF_NAME=v1.2.3
export GITHUB_EVENT_NAME=pull_request
reject preflight
export GITHUB_EVENT_NAME=workflow_dispatch
reject preflight
export GITHUB_EVENT_NAME=push GITHUB_REF_TYPE=branch
reject preflight
export GITHUB_REF_TYPE=tag
printf '{"deleted":true}\n' > "$GITHUB_EVENT_PATH"
reject preflight
printf '{"deleted":false}\n' > "$GITHUB_EVENT_PATH"
git commit -q --allow-empty -m 'fix: different source'
reject preflight
GITHUB_SHA=$(git rev-parse HEAD)
reject preflight
git tag v1.2.4
export GITHUB_REF_NAME=v1.2.4
run preflight

# Existing releases and partial release drafts are not overwritten by a retry.
export GH_PAYLOAD='[[{"tag_name":"v1.2.4","draft":false}],[]]'
reject preflight
export GH_PAYLOAD='[[{"tag_name":"v1.2.4","draft":true}],[]]'
reject preflight
export GH_FAIL=true
reject preflight
export GH_FAIL=false GH_PAYLOAD='not-json'
reject preflight
export GH_PAYLOAD='{"unexpected":"shape"}'
reject preflight
export GH_PAYLOAD='[[]]'
run preflight
reject finalize

# Successful publication exposes signed artifacts as a stable release or a prerelease matching the source tag.
for TAG in v1.2.4 v2.0.0-feature.demo.1 v2.0.0+build-with-hyphen; do
  if [[ "$TAG" != v1.2.4 ]]; then git tag "$TAG"; fi
  export GITHUB_REF_NAME="$TAG"
  export GH_PAYLOAD="[[{\"tag_name\":\"$TAG\",\"draft\":true,\"assets\":[{\"name\":\"checksums.txt.sig\"}],\"html_url\":\"https://github.com/example/project/releases/tag/$TAG\"}]]"
  : > "$GH_LOG"
  run finalize
  grep -q "^release edit $TAG .*--draft=false" "$GH_LOG"
  if [[ "$TAG" == v2.0.0-feature.demo.1 ]]; then
    grep -q -- '--prerelease=true --latest=false' "$GH_LOG"
  else
    grep -q -- '--prerelease=false --latest=true' "$GH_LOG"
  fi
done
export GITHUB_REF_NAME=v1.2.4
export GH_PAYLOAD='[[{"tag_name":"v1.2.4","draft":false,"assets":[{}],"html_url":"https://example.test/release"}]]'
reject finalize
export GH_PAYLOAD='[[{"tag_name":"v1.2.4","draft":true,"assets":[],"html_url":"https://example.test/release"}]]'
reject finalize
export GH_PAYLOAD='[[{"tag_name":"v1.2.4","draft":true,"assets":[{}],"html_url":"https://example.test/release"}]]'
export GH_FAIL=true
reject finalize
printf 'Go release helper checks passed: exact tag/source, prereleases, retry protection, draft finalization, and API failures.\n'
