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
kit="${BENCH_KIT:-$(cd "$gate_dir/.." && pwd)}"
bench="${BENCH_RUN_BINARY:-}"
case "$bench" in
  /*) ;;
  *) echo "error: start the gate through the Bench wrapper: run 'bash bin/bench.sh gate' from the repository root" >&2
     echo "       the wrapper selects the Bench executable and exports BENCH_RUN_BINARY, which is unset or not absolute here" >&2
     exit 1 ;;
esac
if [[ ! -f "$bench" || ! -x "$bench" || -L "$bench" ]]; then
  echo "error: BENCH_RUN_BINARY is not a regular executable" >&2
  exit 1
fi
bench_physical="$(cd -P "$(dirname "$bench")" 2>/dev/null && pwd)/$(basename "$bench")"
if [[ "$bench_physical" != "$bench" ]]; then
  echo "error: BENCH_RUN_BINARY must be a cleaned physical path" >&2
  exit 1
fi
if ! "$bench" freshness-check "$kit"; then
  exit 1
fi
exec env BENCH_KIT="$kit" BENCH_RUN_BINARY="$bench" "$bench" gate-phases "$root"
