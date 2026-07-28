#!/usr/bin/env bash
# Build the complete npm artifact set from explicit regular-file inputs. Nothing is
# written to the requested output until every wrapper and platform tarball is ready.
set -euo pipefail

source_root="${1:?usage: build-artifacts.sh <source-root> <output-dir>}"
output="${2:?usage: build-artifacts.sh <source-root> <output-dir>}"
source_root="$(cd "$source_root" && pwd -P)"
node "$source_root/scripts/build-release-evidence.mjs" --validate-required-sources "$source_root"
if ! git -C "$source_root" rev-parse --verify HEAD >/dev/null 2>&1; then
  printf 'bench artifacts: source state has no authenticated HEAD\n' >&2
  exit 1
fi
source_status=""
if ! source_status="$(git -C "$source_root" status --porcelain=v1 --untracked-files=all --ignore-submodules=none)"; then
  printf 'bench artifacts: could not verify source state at HEAD\n' >&2
  exit 1
fi
if [[ -n "$source_status" ]]; then
  printf 'bench artifacts: source state must be clean and tracked at HEAD\n' >&2
  exit 1
fi
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
# Armed before anything below can fail, so every later exit takes the staging
# directory with it. Everything cleanup touches is initialized above.
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
  # Go's isolated module cache makes downloaded modules read-only. Restore
  # owner write/search before removal, including newly-created descendants.
  chmod -R u+rwX "$stage" 2>/dev/null || true
  rm -rf "$stage"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM HUP
# Dev-tier posture opt-in, matched exactly against 1. Absent, empty, and every other
# value resolve hermetic, so no environment accident can hand a release build the
# ambient caches or drop its reproducibility evidence.
shared_build_cache=0
[[ "${BENCH_SHARED_BUILD_CACHE:-}" != 1 ]] || shared_build_cache=1
if [[ "$shared_build_cache" == 1 ]]; then
  # Resolve while HOME is still ambient: go env reports the environment's value when
  # one is set and the computed default otherwise, and that default sits under HOME.
  ambient_go_caches=""
  if ! ambient_go_caches="$(go env GOCACHE GOMODCACHE)"; then
    printf 'bench artifacts: could not read the ambient Go caches\n' >&2
    exit 1
  fi
  if ! { IFS= read -r ambient_gocache && IFS= read -r ambient_gomodcache; } <<<"$ambient_go_caches"; then
    printf 'bench artifacts: ambient Go caches did not resolve to a cache and a module cache\n' >&2
    exit 1
  fi
fi
# npm writes pack bookkeeping even for local, script-free packs. Keep that mutable
# state, TMPDIR, and HOME inside the private generation under both postures so
# HOME/cache/locale perturbations cannot affect release bytes.
export npm_config_cache="$stage/npm-cache"
export TMPDIR="$stage/tmp"
export HOME="$stage/home"
if [[ "$shared_build_cache" == 1 ]]; then
  export GOCACHE="$ambient_gocache"
  export GOMODCACHE="$ambient_gomodcache"
else
  export GOCACHE="$stage/go-cache"
  export GOMODCACHE="$stage/go-mod-cache"
fi
mkdir -p "$TMPDIR" "$HOME" "$GOMODCACHE"
lock_path="${output}.lock"
if ! mkdir "$lock_path" 2>/dev/null; then
  printf 'bench artifacts: another build owns output %s\n' "$output" >&2
  exit 1
fi
lock="$lock_path"
wrapper="$stage/wrapper"
packages="$stage/packages"
artifacts="$stage/artifacts"
npm_artifacts="$stage/npm-artifacts"
offline_archives="$stage/offline-archives"
mkdir -p "$wrapper" "$packages" "$artifacts" "$npm_artifacts"

if [[ -n "${BENCH_TEST_PREPARED_ARTIFACTS:-}" ]]; then
  [[ -n "${BENCH_TEST_PROMOTION_READY_FILE:-}" ]] || { printf 'bench artifacts: prepared artifacts require the promotion test seam\n' >&2; exit 1; }
  [[ -d "$BENCH_TEST_PREPARED_ARTIFACTS" && ! -L "$BENCH_TEST_PREPARED_ARTIFACTS" ]] || { printf 'bench artifacts: prepared test artifacts are not a directory\n' >&2; exit 1; }
  cp -a "$BENCH_TEST_PREPARED_ARTIFACTS/." "$artifacts/"
else
  matrix_file="$stage/platform-matrix.tsv"
  node "$source_root/scripts/release-plan.mjs" "$source_root" targets > "$matrix_file"
  version="$(node -p 'require(process.argv[1]).version' "$source_root/package.json")"
  npm_pack_flags=()
  while IFS= read -r arg; do npm_pack_flags+=("$arg"); done < <(node -e 'for (const arg of require(process.argv[1]).toolchains.find(tool => tool.name === "npm").operations.pack) console.log(arg)' "$source_root/internal/releaseevidence/requirements.json")

  while IFS=$'\t' read -r os arch _goos _goarch _runner; do
    mkdir -p "$packages/$os-$arch/bin"
  done < "$matrix_file"

  while IFS=$'\t' read -r os arch goos goarch _runner; do
    binary="$packages/$os-$arch/bin/bench"
    if [[ "$goos" == linux ]]; then
      CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" bash "$source_root/scripts/go-build.sh" "$source_root" "$binary"
    else
      GOOS="$goos" GOARCH="$goarch" bash "$source_root/scripts/go-build.sh" "$source_root" "$binary"
    fi
    chmod 0755 "$binary"
  done < "$matrix_file"

  node "$source_root/scripts/build-release-evidence.mjs" "$source_root" "$wrapper" "$packages"

  while IFS=$'\t' read -r os arch _goos _goarch _runner; do
    npm pack "$packages/$os-$arch" --pack-destination "$npm_artifacts" "${npm_pack_flags[@]}" >/dev/null
  done < "$matrix_file"

  pin_path="$(node -p 'require(process.argv[1]).binary_pin_manifest.path' "$source_root/internal/releaseevidence/requirements.json")"
  if [[ "${BENCH_TEST_SKIP_PIN_MANIFEST:-0}" != 1 ]]; then
    node "$source_root/scripts/build-binary-pins.mjs" "$source_root" "$wrapper" "$npm_artifacts" "$version"
    # Refresh the wrapper component inventory now that the declared late-bound
    # manifest exists; platform tarballs are already closed and remain unchanged.
    node "$source_root/scripts/build-release-evidence.mjs" "$source_root" "$wrapper" "$packages"
  fi
  [[ -f "$wrapper/$pin_path" && ! -L "$wrapper/$pin_path" && -s "$wrapper/$pin_path" ]] || { printf 'bench artifacts: binary pin manifest is missing or empty\n' >&2; exit 1; }

  npm pack "$wrapper" --pack-destination "$npm_artifacts" "${npm_pack_flags[@]}" >/dev/null
  bash "$source_root/scripts/build-offline-archives.sh" "$npm_artifacts" "$offline_archives"
  cp "$npm_artifacts"/* "$offline_archives"/* "$artifacts/"
  expected="$(node "$source_root/scripts/release-plan.mjs" "$source_root" artifact-names "$version" | wc -l | tr -d '[:space:]')"
  actual="$(find "$artifacts" -maxdepth 1 -type f \( -name '*.tgz' -o -name '*.tar.gz' \) -print | wc -l | tr -d ' ')"
  [[ "$actual" == "$expected" ]] || { printf 'bench artifacts: emitted %s tarballs, expected %s\n' "$actual" "$expected" >&2; exit 1; }

  if [[ "${BENCH_REPRO_BUILD:-0}" != 1 && "$shared_build_cache" != 1 ]]; then
    second_parent="$(mktemp -d "$parent/.bench-repro-build.XXXXXX")"
    second_source="$second_parent/source"
    mkdir -p "$second_parent/tmp"
    git clone -q --no-hardlinks "$source_root" "$second_source"
    second_output="$second_parent/artifacts"
    if ! env -u BENCH_TEST_PROMOTION_READY_FILE BENCH_REPRO_BUILD=1 HOME="$second_parent/home" TMPDIR="$second_parent/tmp" GOCACHE="$second_parent/go-cache" GOMODCACHE="$second_parent/go-mod-cache" npm_config_cache="$second_parent/npm-cache" bash "$second_source/scripts/build-artifacts.sh" "$second_source" "$second_output"; then
      printf 'bench artifacts: independent reproducibility build failed\n' >&2
      exit 1
    fi
    bash "$source_root/scripts/compare-artifacts.sh" "$artifacts" "$second_output" "$stage/reproducibility.json" "$source_root" "$second_source"
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
# A shared-cache build compares nothing, so any record beside the output grades bytes
# it never saw. Either posture moves that record aside with the artifacts it replaces,
# and cleanup returns both together if the promotion does not complete.
if [[ "$promote_reproducibility" == 1 || "$shared_build_cache" == 1 ]] && [[ -e "$repro_file" ]]; then
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
chmod -R u+rwX "$stage" 2>/dev/null || true
rm -rf "$stage"
