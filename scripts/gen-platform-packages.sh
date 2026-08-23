#!/usr/bin/env bash
# This compatibility entry point routes artifact construction to its one owner.
# That owner atomically promotes the complete tarball set.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
out="${1:?usage: gen-platform-packages.sh <output-dir>}"
exec bash "$here/scripts/build-artifacts.sh" "$here" "$out"
