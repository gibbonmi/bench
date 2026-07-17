#!/usr/bin/env bash
# Build the complete npm artifact set from explicit regular-file inputs. Nothing is
# written to the requested output until every wrapper and platform tarball is ready.
set -euo pipefail

source_root="${1:?usage: build-artifacts.sh <source-root> <output-dir>}"
output="${2:?usage: build-artifacts.sh <source-root> <output-dir>}"
source_root="$(cd "$source_root" && pwd)"
parent="$(dirname "$output")"
mkdir -p "$parent"
stage="$(mktemp -d "$parent/.bench-artifacts.XXXXXX")"
backup=""
lock=""
cleanup() {
  if [[ -n "$backup" && -e "$backup" && ! -e "$output" ]]; then mv "$backup" "$output"; fi
  [[ -z "$backup" || ! -e "$backup" ]] || rm -rf "$backup"
  [[ -z "$lock" || ! -d "$lock" ]] || rmdir "$lock"
  rm -rf "$stage"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM HUP
wrapper="$stage/wrapper"
packages="$stage/packages"
artifacts="$stage/artifacts"
mkdir -p "$wrapper" "$packages" "$artifacts"

matrix_file="$stage/platform-matrix.tsv"
node -e 'for (const p of require(process.argv[1])) console.log([p.os,p.arch,p.goos,p.goarch].join("\t"))' "$source_root/scripts/platforms.json" > "$matrix_file"
matrix_count="$(wc -l < "$matrix_file" | tr -d '[:space:]')"
npm_pack_flags=()
while IFS= read -r arg; do npm_pack_flags+=("$arg"); done < <(node -e 'for (const arg of require(process.argv[1]).toolchains.find(tool => tool.name === "npm").operations.pack) console.log(arg)' "$source_root/internal/releaseevidence/requirements.json")

while IFS=$'\t' read -r os arch _goos _goarch; do
  mkdir -p "$packages/$os-$arch/bin"
done < "$matrix_file"

while IFS=$'\t' read -r os arch goos goarch; do
  binary="$packages/$os-$arch/bin/bench"
  if [[ "$goos" == linux ]]; then
    CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" bash "$source_root/scripts/go-build.sh" "$source_root" "$binary"
  else
    GOOS="$goos" GOARCH="$goarch" bash "$source_root/scripts/go-build.sh" "$source_root" "$binary"
  fi
  chmod 0755 "$binary"
done < "$matrix_file"

node "$source_root/scripts/build-release-evidence.mjs" "$source_root" "$wrapper" "$packages"

while IFS=$'\t' read -r os arch _goos _goarch; do
  npm pack "$packages/$os-$arch" --pack-destination "$artifacts" "${npm_pack_flags[@]}" >/dev/null
done < "$matrix_file"

npm pack "$wrapper" --pack-destination "$artifacts" "${npm_pack_flags[@]}" >/dev/null
expected="$((matrix_count + 1))"
actual="$(find "$artifacts" -maxdepth 1 -type f -name '*.tgz' -print | wc -l | tr -d ' ')"
[[ "$actual" == "$expected" ]] || { printf 'bench artifacts: emitted %s tarballs, expected %s\n' "$actual" "$expected" >&2; exit 1; }

lock_path="${output}.lock"
if ! mkdir "$lock_path" 2>/dev/null; then
  printf 'bench artifacts: another build owns output %s\n' "$output" >&2
  exit 1
fi
lock="$lock_path"
if [[ -n "${BENCH_TEST_PROMOTION_READY_FILE:-}" ]]; then
  : > "$BENCH_TEST_PROMOTION_READY_FILE"
  while [[ -e "$BENCH_TEST_PROMOTION_READY_FILE" ]]; do sleep 0.05; done
fi
if [[ -e "$output" ]]; then
  backup="$(mktemp -d "$parent/.bench-artifacts.previous.XXXXXX")"
  rmdir "$backup"
  mv "$output" "$backup"
fi
mv "$artifacts" "$output"
[[ -z "$backup" || ! -e "$backup" ]] || rm -rf "$backup"
backup=""
rmdir "$lock"
lock=""
trap - EXIT INT TERM HUP
rm -rf "$stage"
