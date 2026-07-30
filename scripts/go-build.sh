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
#   Usage: go-build.sh <module-root> <output-path>
#   Cross-compile: set GOOS / GOARCH in the environment.
#
# This file is repo-only: not in package.json files[], so it never ships. Consumers
# get prebuilt binaries via the @redbench/<os>-<arch> platform packages.
set -euo pipefail

modroot="$(cd "${1:?usage: go-build.sh <module-root> <output-path>}" && pwd -P)"
out="${2:?usage: go-build.sh <module-root> <output-path>}"
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
go_build_flags=()
while IFS= read -r arg; do go_build_flags+=("$arg"); done < <(node -e '
  const registry = require(process.argv[1]);
  for (const arg of registry.toolchains.find(tool => tool.name === "go").operations.build) console.log(arg);
' "$modroot/$go_requirements")
for index in "${!go_build_flags[@]}"; do
  go_build_flags[$index]="${go_build_flags[$index]//<package-version>/$version}"
done

# Staging beside the target makes executable promotion an atomic rename.
out_dir="$(dirname "$out")"
mkdir -p "$out_dir"
staged="$(mktemp "$out_dir/.bench.tmp.XXXXXX")"
trap 'rm -f "$staged"' EXIT
go build "${go_build_flags[@]}" -o "$staged" ./cmd/bench
env -u GOOS -u GOARCH go run ./internal/freshness/cmd "$modroot" "$staged" "$out"
trap - EXIT
