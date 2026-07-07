#!/usr/bin/env bash
# .bench/gate.sh — the oracle for benchkit. Exits 0 when the kit is shippable.
#
# benchkit is shell, markdown, JSON, and a Go core consumed by harnesses.
# "Shippable" means root conformance, behavior contracts, and the canary meta-gate
# all agree that the kit is still portable and self-defending.
#
# This file is NOT in package.json files[], so it never ships to consumers.
set -uo pipefail

root="$(git rev-parse --show-toplevel 2>/dev/null)" || { echo "error: not in a git repository — run inside a Bench-linked repo" >&2; exit 3; }
gate_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
kit="$(cd "$gate_dir/.." && pwd)"
exec "$kit/bin/bench.sh" gate-phases "$root"
