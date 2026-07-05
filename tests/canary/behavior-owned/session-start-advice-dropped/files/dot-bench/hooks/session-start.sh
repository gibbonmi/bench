#!/usr/bin/env bash
# Canary fixture: the by-path advisory branch omits the "run `bench doctor --fix`"
# pointer (story 14), so a cold session where bench resolves only by path gets no
# repair hint. The session-start advisory contract must go red.
set -uo pipefail
if [[ "${1:-}" == "--describe" ]]; then
  printf 'name: session-start\nboundary: SessionStart\ndenies: nothing (informational)\nwhy: x\n'; exit 0
fi
bench_cmd() {
  local root candidate
  root="$(git rev-parse --show-toplevel 2>/dev/null || true)"
  if [[ -n "$root" ]]; then
    for candidate in "$root/.bench/bin/bench.sh" "$root/bin/bench.sh"; do
      [[ -x "$candidate" ]] && { printf '%s\n' "$candidate"; return 0; }
    done
  fi
  command -v bench 2>/dev/null || return 1
}
cmd="$(bench_cmd)" || exit 0
if command -v bench >/dev/null 2>&1; then
  printf 'bench CLI: %s (invoke as: bench)\n' "$cmd"
else
  printf 'bench CLI: %s (bench not on PATH; invoke by path)\n' "$cmd"   # no doctor pointer — the regression
fi
"$cmd" status 2>/dev/null || true
