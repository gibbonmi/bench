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

# Version is package.json's; absent or unreadable → "dev" (the canary fixtures carry
# a go.mod but no package.json, and an unstamped dev build is a legitimate outcome).
version="$(node -e '
  const fs = require("fs");
  try { process.stdout.write(String(JSON.parse(fs.readFileSync(process.argv[1], "utf8")).version || "dev")); }
  catch { process.stdout.write("dev"); }
' "$modroot/package.json")"

cd "$modroot"
go_build_flags=()
while IFS= read -r arg; do go_build_flags+=("$arg"); done < <(node -e '
  const registry = require(process.argv[1]);
  for (const arg of registry.toolchains.find(tool => tool.name === "go").operations.build) console.log(arg);
' "$modroot/internal/releaseevidence/requirements.json")
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
