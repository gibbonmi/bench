#!/usr/bin/env bash
# Compatibility entry point. Artifact construction has one owner; callers that used
# the old generator now receive the complete, atomically promoted tarball set.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
out="${1:?usage: gen-platform-packages.sh <output-dir>}"
exec bash "$here/scripts/build-artifacts.sh" "$here" "$out"
