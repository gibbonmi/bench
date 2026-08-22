#!/usr/bin/env bash
# This is a compatibility entry point. Artifact construction has one owner. A
# caller that used the old generator now receives the complete, atomically
# promoted tarball set.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
out="${1:?usage: gen-platform-packages.sh <output-dir>}"
exec bash "$here/scripts/build-artifacts.sh" "$here" "$out"
