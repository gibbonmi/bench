#!/usr/bin/env bash
set -euo pipefail
root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
binary="$root/dist/bench-preflight"
if [ ! -x "$binary" ]; then
  mkdir -p "$root/dist"
  bash "$root/scripts/go-build.sh" "$root" "$binary"
fi
exec "$binary" release-preflight "$@"
