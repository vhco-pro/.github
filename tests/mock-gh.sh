#!/usr/bin/env bash
set -euo pipefail

printf '%s\n' "$*" >> "$GH_LOG"
if [[ "$1" == api ]]; then
  case "$RELEASE_CASE" in
    existing) printf '[[{"tag_name":"v2.0.0"}],[]]\n' ;;
    absent) printf '[[]]\n' ;;
    failure) exit 1 ;;
    *) exit 1 ;;
  esac
elif [[ "$1 $2" == 'release view' ]]; then
  printf 'https://github.com/example/project/releases/tag/%s\n' "$TAG"
fi