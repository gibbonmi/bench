#!/usr/bin/env bash
# Assemble one target archive for every row in the canonical platform matrix.
set -euo pipefail

npm_artifacts="${1:?usage: build-offline-archives.sh <npm-artifact-dir> <output-dir>}"
output="${2:?usage: build-offline-archives.sh <npm-artifact-dir> <output-dir>}"
root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
parent="$(dirname "$output")"
mkdir -p "$parent"
stage="$(mktemp -d "$parent/.bench-offline-archives.XXXXXX")"
backup=""

cleanup() {
  if [[ -n "$backup" && -e "$backup" && ! -e "$output" ]]; then
    mv "$backup" "$output"
  fi
  [[ -z "$backup" || ! -e "$backup" ]] || rm -rf "$backup"
  rm -rf "$stage"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM HUP

matrix="$stage/platform-matrix.tsv"
node -e 'for (const p of require(process.argv[1])) console.log([p.os,p.arch].join("\t"))' \
  "$root/scripts/platforms.json" > "$matrix"
matrix_count="$(wc -l < "$matrix" | tr -d '[:space:]')"
[[ "$matrix_count" == 4 ]] || {
  printf 'bench offline archives: canonical matrix has %s rows, want exactly four\n' "$matrix_count" >&2
  exit 1
}

version="$(node -p 'require(process.argv[1]).version' "$root/package.json")"
mkdir "$stage/output"

while IFS=$'\t' read -r os arch; do
  archive_root="redbench-${version}-${os}-${arch}"
  archive_dir="$stage/$archive_root"
  wrapper="$npm_artifacts/redbench-${version}.tgz"
  native="$npm_artifacts/redbench-${os}-${arch}-${version}.tgz"
  [[ -s "$wrapper" ]] || { printf 'bench offline archives: missing wrapper %s\n' "$wrapper" >&2; exit 1; }
  [[ -s "$native" ]] || { printf 'bench offline archives: missing platform package %s\n' "$native" >&2; exit 1; }

  mkdir -p "$archive_dir/bin" "$archive_dir/packages"
  cp "$wrapper" "$archive_dir/packages/$(basename "$wrapper")"
  cp "$native" "$archive_dir/packages/$(basename "$native")"
  tar -xzf "$native" -C "$archive_dir" package/bin/bench
  mv "$archive_dir/package/bin/bench" "$archive_dir/bin/bench"
  rm -rf "$archive_dir/package"
  tar -C "$stage" -czf "$stage/output/${archive_root}.tar.gz" "$archive_root"
  rm -rf "$archive_dir"
done < "$matrix"

actual="$(find "$stage/output" -maxdepth 1 -type f -name '*.tar.gz' -print | wc -l | tr -d '[:space:]')"
[[ "$actual" == "$matrix_count" ]] || {
  printf 'bench offline archives: emitted %s archives, expected %s\n' "$actual" "$matrix_count" >&2
  exit 1
}

if [[ -e "$output" ]]; then
  backup="$(mktemp -d "$parent/.bench-offline-archives.previous.XXXXXX")"
  rmdir "$backup"
  mv "$output" "$backup"
fi
mv "$stage/output" "$output"
[[ -z "$backup" || ! -e "$backup" ]] || rm -rf "$backup"
backup=""
trap - EXIT INT TERM HUP
rm -rf "$stage"
