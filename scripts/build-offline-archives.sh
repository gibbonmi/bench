#!/usr/bin/env bash
# This builds the four deterministic, self-contained offline archives from the
# exact npm tarballs the artifact builder already produced.
set -euo pipefail

npm_artifacts="${1:?usage: build-offline-archives.sh <npm-artifact-dir> <output-dir>}"
output="${2:?usage: build-offline-archives.sh <npm-artifact-dir> <output-dir>}"
source_path="${BASH_SOURCE[0]}"
while [[ -L "$source_path" ]]; do
  source_dir="$(cd "$(dirname "$source_path")" && pwd)"
  link_target="$(readlink "$source_path")"
  [[ "$link_target" == /* ]] && source_path="$link_target" || source_path="$source_dir/$link_target"
done
root="$(cd "$(dirname "$source_path")/.." && pwd)"

[[ -d "$npm_artifacts" && ! -L "$npm_artifacts" ]] || {
  printf 'bench offline archives: npm artifact input is not a directory: %s\n' "$npm_artifacts" >&2
  exit 1
}
parent="$(dirname "$output")"
mkdir -p "$parent"
stage="$(mktemp -d "$parent/.bench-offline-archives.XXXXXX")"
backup=""
lock=""

cleanup() {
  if [[ -n "$backup" && -e "$backup" && ! -e "$output" ]]; then
    mv "$backup" "$output"
  fi
  [[ -z "$backup" || ! -e "$backup" ]] || rm -rf "$backup"
  [[ -z "$lock" || ! -d "$lock" ]] || rmdir "$lock"
  rm -rf "$stage"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM HUP

if [[ -e "$output" ]]; then
  [[ -d "$output" && ! -L "$output" ]] || {
    printf 'bench offline archives: output is not a real directory: %s\n' "$output" >&2
    exit 1
  }
  # The swap below replaces this directory and deletes what was here. Only the
  # archives this build emits, and the npm tarballs it consumes, are ours to
  # destroy. Anything else means the caller named a directory that is not a build
  # output, so the refusal comes now, not after the archives are built and the
  # tree is already gone.
  while IFS= read -r -d '' entry; do
    case "${entry##*/}" in
      redbench-*.tar.gz|redbench-*.tgz)
        [[ -f "$entry" && ! -L "$entry" ]] && continue
        ;;
    esac
    printf 'bench offline archives: output directory holds an entry this build did not produce: %s\n' "$entry" >&2
    printf 'bench offline archives: refusing to replace %s\n' "$output" >&2
    exit 1
  done < <(find "$output" -mindepth 1 -maxdepth 1 -print0)
fi
mkdir -p "$stage/roots" "$stage/output"
same_output=0
if [[ "$npm_artifacts" -ef "$output" ]]; then
  same_output=1
fi

matrix="$stage/platform-matrix.tsv"
node "$root/scripts/release-plan.mjs" "$root" targets > "$matrix"

version="$(node -p 'require(process.argv[1]).version' "$root/package.json")"
[[ "$version" =~ ^[0-9]+\.[0-9]+\.[0-9]+([+-][0-9A-Za-z.-]+)?$ ]] || {
  printf 'bench offline archives: package version is invalid: %s\n' "$version" >&2
  exit 1
}

validate_npm_archive() {
  local archive="$1" member
  declare -A seen=()
  while IFS= read -r member; do
    case "$member" in
      package/|package/*) ;;
      *) printf 'bench offline archives: npm tarball has an unexpected path: %s\n' "$member" >&2; exit 1 ;;
    esac
    case "$member" in
      *"../"*|*"/.."|/*|*"\\"*|*"//"*) printf 'bench offline archives: npm tarball has an unsafe path: %s\n' "$member" >&2; exit 1 ;;
    esac
    [[ -n "${seen[$member]:-}" ]] && { printf 'bench offline archives: npm tarball has a duplicate path: %s\n' "$member" >&2; exit 1; }
    seen[$member]=1
  done < <(LC_ALL=C tar -tzf "$archive")
}

validate_extracted_package() {
  local directory="$1" entry
  while IFS= read -r -d '' entry; do
    [[ -d "$entry" && ! -L "$entry" ]] && continue
    [[ -f "$entry" && ! -L "$entry" && -s "$entry" ]] || {
      printf 'bench offline archives: npm tarball has an empty or unsafe member: %s\n' "$entry" >&2
      exit 1
    }
  done < <(find "$directory/package" -mindepth 1 -print0)
}

while IFS=$'\t' read -r os arch _goos _goarch _runner; do
  archive_root="redbench-${version}-${os}-${arch}"
  archive_dir="$stage/roots/$archive_root"
  wrapper="$npm_artifacts/redbench-${version}.tgz"
  native="$npm_artifacts/redbench-${os}-${arch}-${version}.tgz"
  for package in "$wrapper" "$native"; do
    [[ -f "$package" && ! -L "$package" && -s "$package" ]] || {
      printf 'bench offline archives: missing or unsafe npm tarball: %s\n' "$package" >&2
      exit 1
    }
    validate_npm_archive "$package"
    if [[ "$same_output" == 1 && ! -e "$stage/output/$(basename "$package")" ]]; then
      cp "$package" "$stage/output/$(basename "$package")"
      chmod 0644 "$stage/output/$(basename "$package")"
    fi
  done

  wrapper_extract="$stage/${os}-${arch}-wrapper"
  native_extract="$stage/${os}-${arch}-native"
  mkdir -p "$wrapper_extract" "$native_extract"
  tar -xzf "$wrapper" -C "$wrapper_extract"
  tar -xzf "$native" -C "$native_extract"
  validate_extracted_package "$wrapper_extract"
  validate_extracted_package "$native_extract"
  binary="$native_extract/package/bin/bench"
  [[ -f "$binary" && ! -L "$binary" && -s "$binary" ]] || {
    printf 'bench offline archives: target package has no regular executable: %s\n' "$native" >&2
    exit 1
  }

  node "$root/scripts/assemble-offline-archive.mjs" "$root" "$npm_artifacts" "$archive_dir" "${os}-${arch}" "$version" "$binary" "$wrapper_extract" "$native_extract"

  node "$root/scripts/write-deterministic-archive.mjs" "$archive_dir" "$stage/output/${archive_root}.tar.gz"
  rm -rf "$archive_dir" "$wrapper_extract" "$native_extract"
done < "$matrix"

actual="$(find "$stage/output" -mindepth 1 -maxdepth 1 -type f -name 'redbench-*.tar.gz' -print | wc -l | tr -d '[:space:]')"
expected="$(wc -l < "$matrix" | tr -d '[:space:]')"
[[ "$actual" == "$expected" ]] || {
  printf 'bench offline archives: emitted %s archives, expected %s from release plan\n' "$actual" "$expected" >&2
  exit 1
}

if [[ -e "${output}.lock" ]] || ! mkdir "${output}.lock" 2>/dev/null; then
  printf 'bench offline archives: another build owns output %s\n' "$output" >&2
  exit 1
fi
lock="${output}.lock"
if [[ -e "$output" ]]; then
  backup="$(mktemp -d "$parent/.bench-offline-archives.previous.XXXXXX")"
  rmdir "$backup"
  mv "$output" "$backup"
fi
mv "$stage/output" "$output"
[[ -z "$backup" || ! -e "$backup" ]] || rm -rf "$backup"
backup=""
rmdir "$lock"
lock=""
trap - EXIT INT TERM HUP
rm -rf "$stage"
