#!/usr/bin/env bash
set -euo pipefail

SOURCE_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
TEMP=$(mktemp -d)
trap 'rm -rf "$TEMP"' EXIT
export RUNNER_TEMP="$TEMP" RUNNER_ARCH=X64
export GITHUB_OUTPUT="$TEMP/outputs" GITHUB_PATH="$TEMP/path" GITHUB_STEP_SUMMARY="$TEMP/summary"
export GITHUB_SERVER_URL=https://github.com GITHUB_REPOSITORY=example/project
export CONFIG_FILE="$SOURCE_ROOT/release/cliff.toml"
bash "$SOURCE_ROOT/scripts/release/install-git-cliff.sh"
export PATH="$TEMP/release-tools:$PATH"
git init -q -b main "$TEMP/repo"
cd "$TEMP/repo"
git config user.name 'Release Test'
git config user.email test@example.com
git commit -q --allow-empty -m 'feat: first capability'

render() {
  : > "$GITHUB_OUTPUT"
  VERSION="$1" DRY_RUN="$2" bash "$SOURCE_ROOT/scripts/release/select-range.sh"
  export TAG RANGE PREVIOUS_TAG
  TAG=$(sed -n 's/^tag=//p' "$GITHUB_OUTPUT")
  RANGE=$(sed -n 's/^range=//p' "$GITHUB_OUTPUT")
  PREVIOUS_TAG=$(sed -n 's/^previous-tag=//p' "$GITHUB_OUTPUT")
  bash "$SOURCE_ROOT/scripts/release/render-notes.sh"
}
render 1.0.0 true
grep -q 'First capability' "$TEMP/release-notes.md"
test -z "$(git tag)"
git tag v1.0.0
git commit -q --allow-empty -m 'fix: repair something'
git tag v1.0.1-preview.1
git commit -q --allow-empty -m 'feat(api)!: incompatible API'
render 2.0.0 true
test "$PREVIOUS_TAG" = v1.0.0
grep -q 'Breaking changes' "$TEMP/release-notes.md"
grep -q 'Repair something' "$TEMP/release-notes.md"
grep -q 'compare/v1.0.0...v2.0.0' "$TEMP/release-notes.md"
git tag v2.0.0
git commit -q --allow-empty -m 'fix: must not appear in the old release'
render 2.0.0 false
if grep -q 'Must not appear' "$TEMP/release-notes.md"; then exit 1; fi

mkdir "$TEMP/bin"
cp "$SOURCE_ROOT/tests/mock-gh.sh" "$TEMP/bin/gh"
chmod 755 "$TEMP/bin/gh"
export PATH="$TEMP/bin:$PATH" GH_LOG="$TEMP/gh-log" GH_REPO=example/project NOTES="$TEMP/release-notes.md"
export RELEASE_CASE=existing
: > "$GH_LOG"
bash "$SOURCE_ROOT/scripts/release/publish-release.sh"
grep -q '^release edit v2.0.0 ' "$GH_LOG"
if grep -q '^release create ' "$GH_LOG"; then exit 1; fi
export RELEASE_CASE=absent TAG=v2.1.0-preview.1
: > "$GH_LOG"
bash "$SOURCE_ROOT/scripts/release/publish-release.sh"
grep -q '^release create v2.1.0-preview.1 .*--prerelease --latest=false' "$GH_LOG"
export RELEASE_CASE=failure
if bash "$SOURCE_ROOT/scripts/release/publish-release.sh" >/dev/null 2>&1; then exit 1; fi
printf 'Release scripts pass: first release, dry-run, reachable ranges, retry bounds, paginated existing releases, prereleases, and API failure.\n'