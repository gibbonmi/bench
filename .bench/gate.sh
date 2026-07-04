#!/usr/bin/env bash
# .bench/gate.sh — the oracle for benchkit. Exits 0 when the kit is shippable.
#
# benchkit is shell, markdown, JSON, and a Go core consumed by harnesses.
# "Shippable" means root conformance, behavior contracts, and the canary meta-gate
# all agree that the kit is still portable and self-defending.
#
# This file is NOT in package.json files[], so it never ships to consumers.
set -uo pipefail

root="$(git rev-parse --show-toplevel 2>/dev/null)" || { echo "gate: not in a git repo" >&2; exit 3; }
cd "$root"
gate_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

fail=0
err() { echo "gate: $*" >&2; fail=1; }

# The shared fixture harness the contract fragments call (`contract` — one
# source for provision/report/cleanup); must precede every fragment.
# shellcheck source=/dev/null
. "$gate_dir/gate-contract-runner.sh"

# Root-grading conformance lives in Go tests. The real kit checkout owns the test
# code; BENCH_CONFORMANCE_ROOT names the tree under grade, including canary fixtures.
realkit="$(cd "$gate_dir/.." && pwd)"
if ! (cd "$realkit" && BENCH_CONFORMANCE_ROOT="$root" go test -count=1 ./internal/conformance -run '^TestRootConformance$'); then
  fail=1
fi

# shellcheck source=/dev/null
. "$gate_dir/gate-go-contracts.sh"
# shellcheck source=/dev/null
. "$gate_dir/gate-link-contracts.sh"
# shellcheck source=/dev/null
. "$gate_dir/gate-runtime-contracts.sh"
# shellcheck source=/dev/null
. "$gate_dir/gate-runtime-shift-contracts.sh"
# shellcheck source=/dev/null
. "$gate_dir/gate-doctor-contracts.sh"
# shellcheck source=/dev/null
. "$gate_dir/gate-axi-contracts.sh"
# shellcheck source=/dev/null
. "$gate_dir/gate-axi-guards-contracts.sh"
# shellcheck source=/dev/null
. "$gate_dir/gate-axi-wave2-contracts.sh"

# 6. shellcheck — stronger shell lint, best-effort (runs only when installed).
if command -v shellcheck >/dev/null 2>&1; then
  shellcheck -S warning bin/bench.sh .bench/hooks/*.sh .bench/lib/*.sh || err "shellcheck reported issues"
fi

# 7. Canary — prove the gate's own checks still bite, and that the harness itself is
#    present. Outer mode runs the Go canary runner; inner mode skips only the sweep so
#    fixture gates still exercise the remaining shell conformance and behavior
#    fragments without recursing.
if [ "${BENCH_CANARY_INNER:-0}" != "1" ]; then
  bash "$root/bin/bench.sh" canary "$root" || err "canary sweep failed"
fi

if [ "$fail" -eq 0 ]; then echo "gate: green"; else echo "gate: red" >&2; fi
exit "$fail"
