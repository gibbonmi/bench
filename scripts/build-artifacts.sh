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
repro_backup=""
repro_file="$parent/reproducibility.json"
second_parent=""
commit_complete=0
artifacts_promoted=0
reproducibility_promoted=0
promote_reproducibility=0
# npm writes pack bookkeeping even for local, script-free packs. Keep that
# mutable state inside the private generation so HOME/cache/locale perturbations
# cannot affect release bytes or require a trusted ambient cache.
export npm_config_cache="$stage/npm-cache"
export GOCACHE="$stage/go-cache"
cleanup() {
  if [[ "$commit_complete" != 1 ]]; then
    if [[ "$artifacts_promoted" == 1 ]]; then
      [[ ! -e "$output" ]] || rm -rf "$output"
    fi
    if [[ -n "$backup" && -e "$backup" ]]; then mv "$backup" "$output"; fi
    if [[ -n "$repro_backup" && -e "$repro_backup" ]]; then
      if [[ "$reproducibility_promoted" == 1 ]]; then
        [[ ! -e "$repro_file" ]] || rm -f "$repro_file"
      fi
      mv "$repro_backup" "$repro_file"
    fi
  fi
  [[ -z "$backup" || ! -e "$backup" ]] || rm -rf "$backup"
  [[ -z "$repro_backup" || ! -e "$repro_backup" ]] || rm -f "$repro_backup"
  [[ -z "$lock" || ! -d "$lock" ]] || rmdir "$lock"
  [[ -z "$second_parent" || ! -d "$second_parent" ]] || rm -rf "$second_parent"
  rm -rf "$stage"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM HUP
lock_path="${output}.lock"
if ! mkdir "$lock_path" 2>/dev/null; then
  printf 'bench artifacts: another build owns output %s\n' "$output" >&2
  exit 1
fi
lock="$lock_path"
node "$source_root/scripts/build-release-evidence.mjs" --validate-required-sources "$source_root"
wrapper="$stage/wrapper"
packages="$stage/packages"
artifacts="$stage/artifacts"
mkdir -p "$wrapper" "$packages" "$artifacts"

if [[ -n "${BENCH_TEST_PREPARED_ARTIFACTS:-}" ]]; then
  [[ -n "${BENCH_TEST_PROMOTION_READY_FILE:-}" ]] || { printf 'bench artifacts: prepared artifacts require the promotion test seam\n' >&2; exit 1; }
  [[ -d "$BENCH_TEST_PREPARED_ARTIFACTS" && ! -L "$BENCH_TEST_PREPARED_ARTIFACTS" ]] || { printf 'bench artifacts: prepared test artifacts are not a directory\n' >&2; exit 1; }
  cp -a "$BENCH_TEST_PREPARED_ARTIFACTS/." "$artifacts/"
else
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
  "$source_root/scripts/build-offline-archives.sh" "$artifacts" "$artifacts"
  expected="$((matrix_count * 2 + 1))"
  actual="$(find "$artifacts" -maxdepth 1 -type f \( -name '*.tgz' -o -name '*.tar.gz' \) -print | wc -l | tr -d ' ')"
  [[ "$actual" == "$expected" ]] || { printf 'bench artifacts: emitted %s tarballs, expected %s\n' "$actual" "$expected" >&2; exit 1; }

  if [[ "${BENCH_REPRO_BUILD:-0}" != 1 ]]; then
    second_parent="$(mktemp -d "$parent/.bench-repro-build.XXXXXX")"
    second_output="$second_parent/artifacts"
    if ! env -u BENCH_TEST_PROMOTION_READY_FILE BENCH_REPRO_BUILD=1 bash "$source_root/scripts/build-artifacts.sh" "$source_root" "$second_output"; then
      printf 'bench artifacts: independent reproducibility build failed\n' >&2
      exit 1
    fi
    bash "$source_root/scripts/compare-artifacts.sh" "$artifacts" "$second_output" "$stage/reproducibility.json"
    rm -rf "$second_parent"
    second_parent=""
    promote_reproducibility=1
  fi
fi

if [[ -n "${BENCH_TEST_PROMOTION_READY_FILE:-}" ]]; then
  : > "$BENCH_TEST_PROMOTION_READY_FILE"
  while [[ -e "$BENCH_TEST_PROMOTION_READY_FILE" ]]; do sleep 0.05; done
fi
if [[ -e "$output" ]]; then
  backup="$(mktemp -d "$parent/.bench-artifacts.previous.XXXXXX")"
  rmdir "$backup"
  mv "$output" "$backup"
fi
if [[ "$promote_reproducibility" == 1 && -e "$repro_file" ]]; then
  repro_backup="$(mktemp "$parent/.bench-reproducibility.previous.XXXXXX")"
  mv "$repro_file" "$repro_backup"
fi
mv "$artifacts" "$output"
artifacts_promoted=1
if [[ "$promote_reproducibility" == 1 ]]; then
  mv "$stage/reproducibility.json" "$repro_file"
  reproducibility_promoted=1
fi
commit_complete=1
[[ -z "$backup" || ! -e "$backup" ]] || rm -rf "$backup"
backup=""
[[ -z "$repro_backup" || ! -e "$repro_backup" ]] || rm -f "$repro_backup"
repro_backup=""
rmdir "$lock"
lock=""
trap - EXIT INT TERM HUP
rm -rf "$stage"
