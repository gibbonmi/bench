#!/usr/bin/env bash
# Build the four deterministic, self-contained offline archives from the exact
# npm tarballs already produced by the artifact builder.
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
  while IFS= read -r -d '' entry; do
    info="$(stat -c '%F' "$entry" 2>/dev/null || stat -f '%HT' "$entry")"
    [[ "$info" == "regular file" || "$info" == "Regular File" ]] || {
      printf 'bench offline archives: output contains a non-regular entry: %s\n' "$entry" >&2
      exit 1
    }
    cp "$entry" "$stage/$(basename "$entry")"
  done < <(find "$output" -mindepth 1 -maxdepth 1 -print0)
fi
mkdir -p "$stage/roots" "$stage/output"
for entry in "$stage"/*; do
  [[ "$(basename "$entry")" == "output" || "$(basename "$entry")" == "roots" ]] && continue
  [[ -f "$entry" ]] && cp "$entry" "$stage/output/$(basename "$entry")"
done

matrix="$stage/platform-matrix.tsv"
node -e '
  const rows = require(process.argv[1]);
  if (!Array.isArray(rows) || rows.length !== 4) throw new Error("canonical platform matrix must contain exactly four rows");
  const seen = new Set();
  for (const row of rows) {
    if (!row || !/^(darwin|linux)$/.test(row.os) || !/^(arm64|x64)$/.test(row.arch)) throw new Error("canonical platform matrix has an invalid target");
    const key = `${row.os}-${row.arch}`;
    if (seen.has(key)) throw new Error(`canonical platform matrix repeats ${key}`);
    seen.add(key);
    process.stdout.write(`${row.os}\t${row.arch}\n`);
  }
' "$root/scripts/platforms.json" > "$matrix"

version="$(node -p 'require(process.argv[1]).version' "$root/package.json")"
[[ "$version" =~ ^[0-9]+\.[0-9]+\.[0-9]+([+-][0-9A-Za-z.-]+)?$ ]] || {
  printf 'bench offline archives: package version is invalid: %s\n' "$version" >&2
  exit 1
}

copy_package_file() {
  local package_dir="$1" rel="$2" destination="$3"
  local source="$package_dir/package/$rel"
  [[ -f "$source" && ! -L "$source" && -s "$source" ]] || {
    printf 'bench offline archives: package evidence is missing or unsafe: %s\n' "$rel" >&2
    exit 1
  }
  mkdir -p "$(dirname "$destination")"
  cp "$source" "$destination"
  chmod 0644 "$destination"
}

validate_npm_archive() {
  local archive="$1" member kind
  while IFS= read -r member; do
    case "$member" in
      package/|package/*) ;;
      *) printf 'bench offline archives: npm tarball has an unexpected path: %s\n' "$member" >&2; exit 1 ;;
    esac
    case "$member" in
      *"../"*|*"/.."|/*|*"\\"*|*"//"*) printf 'bench offline archives: npm tarball has an unsafe path: %s\n' "$member" >&2; exit 1 ;;
    esac
  done < <(LC_ALL=C tar -tzf "$archive")
  while IFS= read -r member; do
    kind="${member:0:1}"
    case "$kind" in
      -|d) ;;
      *) printf 'bench offline archives: npm tarball contains a special member: %s\n' "$member" >&2; exit 1 ;;
    esac
  done < <(LC_ALL=C tar -tvzf "$archive")
}

while IFS=$'\t' read -r os arch; do
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
  done

  wrapper_extract="$stage/${os}-${arch}-wrapper"
  native_extract="$stage/${os}-${arch}-native"
  mkdir -p "$wrapper_extract" "$native_extract"
  tar -xzf "$wrapper" -C "$wrapper_extract"
  tar -xzf "$native" -C "$native_extract"
  binary="$native_extract/package/bin/bench"
  [[ -f "$binary" && ! -L "$binary" && -s "$binary" ]] || {
    printf 'bench offline archives: target package has no regular executable: %s\n' "$native" >&2
    exit 1
  }

  mkdir -p "$archive_dir/bin" "$archive_dir/packages" "$archive_dir/evidence/components"
  cp "$binary" "$archive_dir/bin/bench"
  chmod 0755 "$archive_dir/bin/bench"
  cp "$wrapper" "$archive_dir/packages/$(basename "$wrapper")"
  cp "$native" "$archive_dir/packages/$(basename "$native")"
  chmod 0644 "$archive_dir/packages/$(basename "$wrapper")" "$archive_dir/packages/$(basename "$native")"

  while IFS= read -r evidence_path; do
    [[ -n "$evidence_path" ]] || continue
    copy_package_file "$wrapper_extract" "$evidence_path" "$archive_dir/evidence/$evidence_path"
  done < <(node -e 'for (const record of require(process.argv[1]).records) if (record.package_mode) console.log(record.path)' "$root/internal/releaseevidence/requirements.json")
  copy_package_file "$wrapper_extract" component-manifest.json "$archive_dir/evidence/components/wrapper-component-manifest.json"
  copy_package_file "$native_extract" component-manifest.json "$archive_dir/evidence/components/platform-component-manifest.json"

  printf '%s\n' \
    "# Redbench ${version} offline bundle" \
    "" \
    "This archive is for ${os}/${arch}. Verify the externally supplied release-index.json and SHA256SUMS before use." \
    "" \
    "## Direct execution" \
    "" \
    "Run bin/bench version and bin/bench commands --brief with npm, Node, caches, and repair unavailable." \
    "" \
    "## Local npm installation" \
    "" \
    "Install packages/redbench-${version}.tgz and packages/redbench-${os}-${arch}-${version}.tgz with npm --offline and BENCH_NO_REPAIR=1." \
    "" \
    "## Internal registry" \
    "" \
    "Seed the platform tarball first and the wrapper tarball last. Serve the exact bytes at their immutable ${version} versions." \
    "" \
    "## Removal" \
    "" \
    "Uninstall the package and remove the extracted bundle, prefix, cache, and temporary home; no residue is expected." \
    > "$archive_dir/OFFLINE.md"
  chmod 0644 "$archive_dir/OFFLINE.md"

  node "$root/scripts/write-deterministic-archive.mjs" "$archive_dir" "$stage/output/${archive_root}.tar.gz"
  rm -rf "$archive_dir" "$wrapper_extract" "$native_extract"
done < "$matrix"

actual="$(find "$stage/output" -mindepth 1 -maxdepth 1 -type f -name 'redbench-*.tar.gz' -print | wc -l | tr -d '[:space:]')"
[[ "$actual" == 4 ]] || {
  printf 'bench offline archives: emitted %s archives, expected exactly four\n' "$actual" >&2
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
