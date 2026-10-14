#!/usr/bin/env bash
set -euo pipefail

FLAGS=()
if [[ "$TAG" == *-* ]]; then
  FLAGS+=(--prerelease --latest=false)
fi
RELEASES=$(gh api "repos/$GH_REPO/releases" --paginate --slurp)
if printf '%s' "$RELEASES" | jq -e --arg tag "$TAG" 'any(.[][]; .tag_name == $tag)' >/dev/null; then
  echo "Release $TAG already exists; preserving its published notes."
else
  gh release create "$TAG" --verify-tag --title "$TAG" --notes-file "$NOTES" "${FLAGS[@]}"
fi
URL=$(gh release view "$TAG" --json url --jq .url)
printf 'release-url=%s\n' "$URL" >> "$GITHUB_OUTPUT"