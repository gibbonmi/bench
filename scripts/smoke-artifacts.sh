#!/usr/bin/env bash
# Execute the host-native tarball through the same installed wrapper users receive.
set -euo pipefail
artifacts="${1:?usage: smoke-artifacts.sh <artifact-dir>}"
root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
version="$(node -p 'require(process.argv[1]).version' "$root/package.json")"
case "$(uname -s):$(uname -m)" in
  Darwin:arm64) target=darwin-arm64 ;;
  Darwin:x86_64) target=darwin-x64 ;;
  Linux:aarch64|Linux:arm64) target=linux-arm64 ;;
  Linux:x86_64|Linux:amd64) target=linux-x64 ;;
  *) printf 'bench artifacts: unsupported smoke host %s/%s\n' "$(uname -s)" "$(uname -m)" >&2; exit 2 ;;
esac
wrapper="$artifacts/redbench-$version.tgz"
native="$artifacts/redbench-$target-$version.tgz"
[[ -f "$wrapper" && -f "$native" ]] || { printf 'bench artifacts: host tarballs missing for %s\n' "$target" >&2; exit 1; }
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
mkdir -p "$tmp/app" "$tmp/home/cache/bin/$version/$target"
mkdir -p "$tmp/native"
tar -xzf "$native" -C "$tmp/native" package/bin/bench
native_binary="$tmp/native/package/bin/bench"
[[ -x "$native_binary" && -s "$native_binary" ]] || { printf 'bench artifacts: %s binary mode or size invalid\n' "$target" >&2; exit 1; }
format="$(file "$native_binary")"
case "$target" in
  darwin-arm64) [[ "$format" == *"Mach-O 64-bit"* && "$format" == *"arm64"* ]] ;;
  darwin-x64) [[ "$format" == *"Mach-O 64-bit"* && "$format" == *"x86_64"* ]] ;;
  linux-arm64) [[ "$format" == *"ELF 64-bit"* && "$format" == *"ARM aarch64"* && "$format" == *"statically linked"* ]] ;;
  linux-x64) [[ "$format" == *"ELF 64-bit"* && "$format" == *"x86-64"* && "$format" == *"statically linked"* ]] ;;
esac || { printf 'bench artifacts: %s format mismatch: %s\n' "$target" "$format" >&2; exit 1; }
printf '#!/bin/sh\nprintf "poisoned cache\\n"\nexit 99\n' > "$tmp/home/cache/bin/$version/$target/bench"
chmod 0755 "$tmp/home/cache/bin/$version/$target/bench"
printf '{"private":true}\n' > "$tmp/app/package.json"
npm install --prefix "$tmp/app" --ignore-scripts --omit=optional "$wrapper" "$native" >/dev/null
installed="$tmp/app/node_modules/redbench/bin/bench.sh"
out="$(BENCH_HOME="$tmp/home" bash "$installed" version)"
case "$out" in
  "benchkit $version "*) ;;
  *) printf 'bench artifacts: native version mismatch: %s\n' "$out" >&2; exit 1 ;;
esac
printf 'bench artifacts: %s selected %s\n' "$target" "$out"
