#!/usr/bin/env bash
set -euo pipefail

TAG="v$VERSION"
if [[ "$DRY_RUN" == true ]]; then
  REVISION=$(git rev-parse HEAD)
else
  REVISION=$(git rev-parse "refs/tags/$TAG^{commit}")
fi
PREVIOUS=""
while IFS= read -r CANDIDATE; do
  if [[ "$CANDIDATE" != "$TAG" && "$CANDIDATE" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    PREVIOUS="$CANDIDATE"
    break
  fi
done < <(git for-each-ref --merged="$REVISION" --sort=-version:refname --format='%(refname:short)' refs/tags)
RANGE="$REVISION"
if [[ -n "$PREVIOUS" ]]; then
  RANGE="$PREVIOUS..$REVISION"
fi
printf 'tag=%s\nrange=%s\nprevious-tag=%s\n' "$TAG" "$RANGE" "$PREVIOUS" >> "$GITHUB_OUTPUT"