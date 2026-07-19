#!/usr/bin/env bash
# Rebuild and execute one canonical target on its declared native runner.
set -euo pipefail

artifacts="${1:?usage: native-proof.sh <artifact-dir> <proof-file> <os> <arch> <runner>}"
proof="${2:?usage: native-proof.sh <artifact-dir> <proof-file> <os> <arch> <runner>}"
os_name="${3:?usage: native-proof.sh <artifact-dir> <proof-file> <os> <arch> <runner>}"
arch_name="${4:?usage: native-proof.sh <artifact-dir> <proof-file> <os> <arch> <runner>}"
runner="${5:?usage: native-proof.sh <artifact-dir> <proof-file> <os> <arch> <runner>}"
source_path="${BASH_SOURCE[0]}"
while [[ -L "$source_path" ]]; do
  source_dir="$(cd "$(dirname "$source_path")" && pwd)"
  link_target="$(readlink "$source_path")"
  [[ "$link_target" == /* ]] && source_path="$link_target" || source_path="$source_dir/$link_target"
done
root="$(cd "$(dirname "$source_path")/.." && pwd)"
version="$(node -p 'require(process.argv[1]).version' "$root/package.json")"
target="${os_name}-${arch_name}"
matrix_row="$(node "$root/scripts/release-plan.mjs" "$root" target "$os_name" "$arch_name")" || {
  printf 'native proof: target is not in the canonical platform matrix: %s\n' "$target" >&2
  exit 1
}
IFS=$'\t' read -r _matrix_os _matrix_arch goos goarch matrix_runner <<< "$matrix_row"
[[ "$matrix_runner" == "$runner" ]] || { printf 'native proof: runner does not match canonical matrix for %s\n' "$target" >&2; exit 1; }

native_name="$(node "$root/scripts/release-plan.mjs" "$root" artifact-name "$version" "$target" platform)"
archive_name="$(node "$root/scripts/release-plan.mjs" "$root" artifact-name "$version" "$target" archive)"
native="$artifacts/$native_name"
archive="$artifacts/$archive_name"
[[ -f "$native" && -f "$archive" ]] || { printf 'native proof: target artifacts are missing for %s\n' "$target" >&2; exit 1; }
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT INT TERM HUP
rebuild="$tmp/bench"
export GOCACHE="$tmp/go-cache"
if [[ "$goos" == linux ]]; then
  CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" bash "$root/scripts/go-build.sh" "$root" "$rebuild"
else
  GOOS="$goos" GOARCH="$goarch" bash "$root/scripts/go-build.sh" "$root" "$rebuild"
fi
chmod 0755 "$rebuild"

package_dir="$tmp/package"
archive_dir="$tmp/archive"
mkdir -p "$package_dir" "$archive_dir"
tar -xzf "$native" -C "$package_dir" package/bin/bench
tar -xzf "$archive" -C "$archive_dir"
archive_root="$archive_dir/redbench-${version}-${target}"
cmp -s "$rebuild" "$package_dir/package/bin/bench" || { printf 'native proof: rebuilt binary differs from package for %s\n' "$target" >&2; exit 1; }
cmp -s "$rebuild" "$archive_root/bin/bench" || { printf 'native proof: rebuilt binary differs from offline archive for %s\n' "$target" >&2; exit 1; }
"$rebuild" version >/dev/null
"$rebuild" commands --brief >/dev/null
if [[ "$goos" == linux ]]; then
  file_info="$(file "$rebuild")"
  [[ "$file_info" == *"statically linked"* ]] || { printf 'native proof: Linux binary is not static for %s: %s\n' "$target" "$file_info" >&2; exit 1; }
  if command -v readelf >/dev/null 2>&1; then
    ! readelf -S "$rebuild" | rg -q '\.symtab' || { printf 'native proof: Linux binary is not stripped for %s\n' "$target" >&2; exit 1; }
  else
    printf 'native proof: readelf is required to prove stripped Linux output\n' >&2
    exit 1
  fi
  command -v docker >/dev/null 2>&1 || { printf 'native proof: musl runner is unavailable for %s\n' "$target" >&2; exit 1; }
  docker run --rm --network none -v "$rebuild:/bench:ro" alpine:3.20 /bench version >/dev/null
  musl_status=green
else
  file_info="$(file "$rebuild")"
  [[ "$file_info" == *"Mach-O"* ]] || { printf 'native proof: Darwin binary format is invalid for %s: %s\n' "$target" "$file_info" >&2; exit 1; }
  command -v nm >/dev/null 2>&1 || { printf 'native proof: nm is required to prove stripped Darwin output\n' >&2; exit 1; }
  nm -a "$rebuild" > "$tmp/darwin-symbols" 2>&1 || true
  rg -q 'no symbols|no name list' "$tmp/darwin-symbols" || { printf 'native proof: Darwin binary is not stripped for %s\n' "$target" >&2; exit 1; }
  musl_status=not_applicable
fi

sha256() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{print $1}'
  else
    shasum -a 256 "$1" | awk '{print $1}'
  fi
}
native_digest="$(sha256 "$native")"
archive_digest="$(sha256 "$archive")"
rebuilt_digest="$(sha256 "$rebuild")"
mkdir -p "$(dirname "$proof")"
# shellcheck disable=SC2016 # Node template literals are intentionally literal here.
node -e 'const fs=require("fs"), path=require("path"); const [file,target,runner,rebuilt,packaged,archive,musl,strip]=process.argv.slice(1); const body={schema_version:1,target,runner,status:"green",rebuilt_sha256:rebuilt,binary_sha256:rebuilt,package_sha256:packaged,archive_sha256:archive,musl_status:musl,operations_status:"green",strip_status:strip,tools_status:"green"}; const tmp=`${file}.tmp-${process.pid}`; fs.mkdirSync(path.dirname(file),{recursive:true}); fs.writeFileSync(tmp,JSON.stringify(body)+"\n",{mode:0o644}); fs.renameSync(tmp,file)' "$proof" "$target" "$runner" "$rebuilt_digest" "$native_digest" "$archive_digest" "$musl_status" green
