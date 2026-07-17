#!/usr/bin/env bash
set -euo pipefail
source_path="${BASH_SOURCE[0]}"
while [ -L "$source_path" ]; do
  source_dir="$(cd "$(dirname "$source_path")" && pwd)"
  link_target="$(readlink "$source_path")"
  case "$link_target" in
    /*) source_path="$link_target" ;;
    *) source_path="$source_dir/$link_target" ;;
  esac
done
root="$(cd "$(dirname "$source_path")/.." && pwd)"
binary="$root/dist/bench-preflight"
if [ ! -x "$binary" ]; then
  mkdir -p "$root/dist"
  bash "$root/scripts/go-build.sh" "$root" "$binary"
fi
exec "$binary" release-preflight "$@"
