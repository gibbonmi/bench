#!/usr/bin/env bash
# The one source of Go build flags. Both callers use it: the gate (host dist/ build)
# and the release workflow (per-target cross-compile). Kept as one file so the
# reproducibility flags cannot drift between "what the gate proved" and "what ships".
#
# -trimpath + the pinned toolchain (go.mod) = reproducible builds, while
# -buildvcs=false prevents checkout topology from changing embedded module/VCS facts.
# The version stamp comes from package.json — the one canonical version — so an
# unstamped build (prints "dev") never masquerades as a release.
#
#   Usage: go-build.sh [--mode artifact] <module-root> <output-path>
#
# This file is repo-only: not in package.json files[], so it never ships. Consumers
# get prebuilt binaries via the @redbench/<os>-<arch> platform packages.
set -euo pipefail

usage() {
  printf 'usage: go-build.sh [--mode artifact] <module-root> <output-path>\n' >&2
  exit 2
}

mode=subject
case "$#" in
  2)
    module_root="$1"
    out="$2"
    ;;
  4)
    [[ "$1" == --mode && "$2" == artifact ]] || usage
    mode=artifact
    module_root="$3"
    out="$4"
    ;;
  *) usage ;;
esac
modroot="$(cd "$module_root" && pwd -P)"
build_script=
package_version=
go_requirements=
while IFS='=' read -r key value; do
  case "$key" in
    build_script) build_script="$value" ;;
    package_version) package_version="$value" ;;
    go_requirements) go_requirements="$value" ;;
    "") ;;
    *) printf 'unknown Go build input key: %s\n' "$key" >&2; exit 1 ;;
  esac
done < "$modroot/scripts/go-build.inputs"
: "${build_script:?go-build.inputs missing build_script}"
: "${package_version:?go-build.inputs missing package_version}"
: "${go_requirements:?go-build.inputs missing go_requirements}"
for input in "$build_script" "$package_version" "$go_requirements"; do
  test -f "$modroot/$input" || { printf 'missing Go build input: %s\n' "$input" >&2; exit 1; }
done

# Version comes from the package-version input named by the build owner manifest.
version="$(node -e '
  const fs = require("fs");
  try { process.stdout.write(String(JSON.parse(fs.readFileSync(process.argv[1], "utf8")).version || "dev")); }
  catch { process.stdout.write("dev"); }
' "$modroot/$package_version")"

cd "$modroot"
if [[ "$out" != /* ]]; then out="./$out"; fi
go_build_flags=()
while IFS= read -r arg; do go_build_flags+=("$arg"); done < <(node -e '
  const registry = require(process.argv[1]);
  for (const arg of registry.toolchains.find(tool => tool.name === "go").operations.build) console.log(arg);
' "$modroot/$go_requirements")
for index in "${!go_build_flags[@]}"; do
  go_build_flags[$index]="${go_build_flags[$index]//<package-version>/$version}"
done

# Staging beside the target makes executable promotion an atomic rename.
out_dir="$(dirname -- "$out")"
refuse_output() {
  printf 'invalid Go build output: %s\n' "$out" >&2
  exit 1
}
if [[ -L "$out" || (-e "$out" && ! -f "$out") || (-e "$out" && ! -w "$out") ]]; then
  refuse_output
fi
if [[ -L "$out.seal" || (-e "$out.seal" && ! -f "$out.seal") || (-e "$out.seal" && ! -w "$out.seal") ]]; then
  refuse_output
fi
component="$out_dir"
while [[ "$component" != . && "$component" != / ]]; do
  if [[ -L "$component" || (-e "$component" && ! -d "$component") ]]; then
    refuse_output
  fi
  component="$(dirname -- "$component")"
done
mkdir -p -- "$out_dir"
[[ -d "$out_dir" && ! -L "$out_dir" && -w "$out_dir" ]] || refuse_output
staged="$(mktemp "$out_dir/.bench.tmp.XXXXXX")"
# The builder holds no rollback state. Each mode installs by replacing the destination
# in one atomic rename, so every way this script can die short of that rename leaves
# the prior output exactly as it found it, and the sealed mode's own restore belongs to
# the publisher that owns the executable-plus-seal pair.
cleanup() {
  # The rename consumes the staged entry, so its absence is the record that an artifact
  # install promoted, and a promoted install still owes the destination's stale seal — it
  # describes the bytes the rename retired. Finishing the install here rather than after
  # the rename is what carries every handled exit, signalled or not, to the unsealed
  # artifact the mode promises; a signal caught between two directory entries would
  # otherwise strand the new executable beside the retired subject's seal.
  [[ "$mode" != artifact || -e "$staged" ]] || rm -f -- "$out.seal"
  rm -f -- "$staged"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM HUP

if [[ "$mode" == artifact ]]; then
  go build "${go_build_flags[@]}" -o "$staged" ./cmd/bench
  [[ -f "$staged" && -x "$staged" ]] || { printf 'artifact build did not produce an executable\n' >&2; exit 1; }
  # Artifact mode must never execute what it just built, so it cannot delegate to the
  # publisher and instead replaces the destination itself. The rename is the whole
  # install: the prior output is present until the instant the new one takes its place,
  # and the cleanup trap that owns every exit path retires the destination's stale seal.
  mv -- "$staged" "$out"
else
  env -u GOOS -u GOARCH go build "${go_build_flags[@]}" -o "$staged" ./cmd/bench
  env -u GOOS -u GOARCH "$staged" freshness-publish "$modroot" "$out"
fi
