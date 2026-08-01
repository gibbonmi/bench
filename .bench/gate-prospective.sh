#!/usr/bin/env bash
set -euo pipefail

root="${1:?usage: gate-prospective.sh <root>}"
bash "$root/scripts/go-build.sh" "$root" "$root/dist/bench"
exec "$root/.bench/gate.sh"
