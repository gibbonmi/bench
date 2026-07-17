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
cleanup() {
  if [[ -n "$backup" && -e "$backup" && ! -e "$output" ]]; then mv "$backup" "$output"; fi
  [[ -z "$backup" || ! -e "$backup" ]] || rm -rf "$backup"
  rm -rf "$stage"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM HUP
wrapper="$stage/wrapper"
packages="$stage/packages"
artifacts="$stage/artifacts"
mkdir -p "$wrapper" "$packages" "$artifacts"

while IFS=$'\t' read -r os arch _goos _goarch; do
  mkdir -p "$packages/$os-$arch/bin"
done < <(node -e 'for (const p of require(process.argv[1])) console.log([p.os,p.arch,p.goos,p.goarch].join("\t"))' "$source_root/scripts/platforms.json")

while IFS=$'\t' read -r os arch goos goarch; do
  binary="$packages/$os-$arch/bin/bench"
  if [[ "$goos" == linux ]]; then
    CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" bash "$source_root/scripts/go-build.sh" "$source_root" "$binary"
  else
    GOOS="$goos" GOARCH="$goarch" bash "$source_root/scripts/go-build.sh" "$source_root" "$binary"
  fi
  chmod 0755 "$binary"
done < <(node -e 'for (const p of require(process.argv[1])) console.log([p.os,p.arch,p.goos,p.goarch].join("\t"))' "$source_root/scripts/platforms.json")

node "$source_root/scripts/build-release-evidence.mjs" "$source_root" "$wrapper" "$packages"

while IFS=$'\t' read -r os arch _goos _goarch; do
  npm pack "$packages/$os-$arch" --pack-destination "$artifacts" --ignore-scripts --silent >/dev/null
done < <(node -e 'for (const p of require(process.argv[1])) console.log([p.os,p.arch,p.goos,p.goarch].join("\t"))' "$source_root/scripts/platforms.json")

npm pack "$wrapper" --pack-destination "$artifacts" --ignore-scripts --silent >/dev/null
expected="$(node -e 'process.stdout.write(String(require(process.argv[1]).length + 1))' "$source_root/scripts/platforms.json")"
actual="$(find "$artifacts" -maxdepth 1 -type f -name '*.tgz' -print | wc -l | tr -d ' ')"
[[ "$actual" == "$expected" ]] || { printf 'bench artifacts: emitted %s tarballs, expected %s\n' "$actual" "$expected" >&2; exit 1; }

if [[ -e "$output" ]]; then
  backup="$(mktemp -d "$parent/.bench-artifacts.previous.XXXXXX")"
  rmdir "$backup"
  mv "$output" "$backup"
fi
if [[ -n "${BENCH_TEST_PROMOTION_READY_FILE:-}" ]]; then
  : > "$BENCH_TEST_PROMOTION_READY_FILE"
  while [[ -e "$BENCH_TEST_PROMOTION_READY_FILE" ]]; do sleep 0.05; done
fi
mv "$artifacts" "$output"
[[ -z "$backup" || ! -e "$backup" ]] || rm -rf "$backup"
backup=""
trap - EXIT INT TERM HUP
rm -rf "$stage"
