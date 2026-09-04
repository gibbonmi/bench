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
  # This is an auxiliary artifact, not the checkout's published bench, so its manifest
  # goes beside it. The builder's default publishes into the wrapper's bin/, which the
  # land route reads and execs; a manifest bound to this artifact would outlive the next
  # clean and refuse the landing.
  bash "$root/scripts/go-build.sh" --manifest-dir "$root/dist" "$root" "$binary"
fi
exec env BENCH_RUN_BINARY="$binary" "$binary" release-preflight "$@"
