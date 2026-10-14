#!/usr/bin/env bash
set -euo pipefail

NOTES="$RUNNER_TEMP/release-notes.md"
test -f "$CONFIG_FILE"
git-cliff --config "$CONFIG_FILE" --tag "$TAG" "$RANGE" > "$NOTES"
if ! grep -q '[^[:space:]]' "$NOTES"; then
  printf 'Maintenance release.\n' > "$NOTES"
fi
if [[ -n "$PREVIOUS_TAG" ]]; then
  printf '\n**Full diff:** %s/%s/compare/%s...%s\n' \
    "$GITHUB_SERVER_URL" "$GITHUB_REPOSITORY" "$PREVIOUS_TAG" "$TAG" >> "$NOTES"
fi
cat "$NOTES" >> "$GITHUB_STEP_SUMMARY"
printf 'file=%s\n' "$NOTES" >> "$GITHUB_OUTPUT"