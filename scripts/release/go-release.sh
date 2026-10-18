#!/usr/bin/env bash
set -euo pipefail

fail() {
  printf '::error::%s\n' "$*" >&2
  exit 1
}

[[ "${GITHUB_EVENT_NAME:-}" == push && "${GITHUB_REF_TYPE:-}" == tag ]] ||
  fail 'Go publication requires an exact tag push event.'
jq -e '.deleted == false' "$GITHUB_EVENT_PATH" >/dev/null ||
  fail 'A deleted tag cannot authorize publication.'

TAG=${GITHUB_REF_NAME:?}
NUMBER='(0|[1-9][0-9]*)'
IDENTIFIER='(0|[1-9][0-9]*|[0-9]*[A-Za-z-][0-9A-Za-z-]*)'
SEMVER="^v${NUMBER}\\.${NUMBER}\\.${NUMBER}(-${IDENTIFIER}(\\.${IDENTIFIER})*)?(\\+[0-9A-Za-z-]+(\\.[0-9A-Za-z-]+)*)?$"
[[ "$TAG" =~ $SEMVER ]] || fail "Not an exact semantic version tag: $TAG"

TAG_SHA=$(git rev-parse --verify "refs/tags/$TAG^{commit}")
HEAD_SHA=$(git rev-parse HEAD)
[[ "$TAG_SHA" == "$GITHUB_SHA" && "$HEAD_SHA" == "$GITHUB_SHA" ]] ||
  fail "Tag $TAG does not identify the event's checked-out commit."

RELEASES=$(gh api "repos/$GITHUB_REPOSITORY/releases" --paginate --slurp)
RELEASE=$(printf '%s' "$RELEASES" | jq -c --arg tag "$TAG" '
  if type != "array" or any(.[]; type != "array") then
    error("Expected paginated release arrays")
  else
    [.[][] | select(.tag_name == $tag)]
    | if length > 1 then error("Duplicate releases for the exact tag")
      elif length == 0 then null
      else .[0]
      end
  end')

case "${1:-}" in
  preflight)
    [[ "$RELEASE" == null ]] ||
      fail "Release $TAG already exists, possibly as a partial draft. Preserve its artifacts; inspect recovery or use a new version."
    printf 'Validated new release %s at %s.\n' "$TAG" "$TAG_SHA"
    ;;
  finalize)
    [[ "$RELEASE" != null ]] || fail "GoReleaser created no draft for $TAG."
    printf '%s' "$RELEASE" | jq -e '
      .draft == true and (.assets | type == "array" and length > 0)
      and (.html_url | type == "string" and length > 0)' >/dev/null ||
      fail "Release $TAG is not an unpublished draft with uploaded artifacts."
    VERSION=${TAG#v}
    VERSION=${VERSION%%+*}
    PRERELEASE=false
    LATEST=true
    if [[ "$VERSION" == *-* ]]; then
      PRERELEASE=true
      LATEST=false
    fi
    gh release edit "$TAG" --repo "$GITHUB_REPOSITORY" \
      --draft=false "--prerelease=$PRERELEASE" "--latest=$LATEST"
    URL=$(printf '%s' "$RELEASE" | jq -er .html_url)
    printf 'release-url=%s\n' "$URL" >> "$GITHUB_OUTPUT"
    printf 'Published [%s](%s) from %s.\n' "$TAG" "$URL" "$TAG_SHA" >> "$GITHUB_STEP_SUMMARY"
    ;;
  *)
    fail 'Expected preflight or finalize.'
    ;;
esac
