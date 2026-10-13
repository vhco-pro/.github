#!/usr/bin/env bash
set -euo pipefail

VERSION=2.14.2
case "$RUNNER_ARCH" in
  X64)
    TARGET=x86_64-unknown-linux-gnu
    CHECKSUM=26d1f7c8ea2400f2ccb0b2e7f321635f600b583bf4cd1f7afaa5e8c24068ad1088a2504e0a4d9f4333d4bb20a80329327d9462634343ee3af91e1a7fbdc6f18a
    ;;
  ARM64)
    TARGET=aarch64-unknown-linux-gnu
    CHECKSUM=fd27768ea76fea8c3ba3f3247fef2f6604988b629f15e48bf40be5eaa08b8fef2e4087cc9c8af16bc522bc19cc369611fa7692101b513f63469ad58522708600
    ;;
  *) echo "::error::Unsupported Linux runner architecture: $RUNNER_ARCH"; exit 1 ;;
esac
DIRECTORY="$RUNNER_TEMP/release-tools"
mkdir -p "$DIRECTORY"
ARCHIVE="$DIRECTORY/git-cliff.tar.gz"
curl -fsSL "https://github.com/orhun/git-cliff/releases/download/v$VERSION/git-cliff-$VERSION-$TARGET.tar.gz" -o "$ARCHIVE"
printf '%s  %s\n' "$CHECKSUM" "$ARCHIVE" | sha512sum --check --status
tar -xzf "$ARCHIVE" -C "$DIRECTORY"
install -m 0755 "$DIRECTORY/git-cliff-$VERSION/git-cliff" "$DIRECTORY/git-cliff"
printf '%s\n' "$DIRECTORY" >> "$GITHUB_PATH"