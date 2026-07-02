#!/usr/bin/env bash
# SessionStart hook: print the ambient dashboard when a session opens cold.
# A thin wrapper over `bench status` — the single renderer the user also runs on demand
# (one source of truth). Never blocks the session: outside a repo, or on any error, it
# prints nothing and exits 0. Shared across harnesses; the .claude adapter wires it under
# hooks.SessionStart, and any AGENTS.md harness can point its own start hook here.
set -uo pipefail

# `--describe` (first arg) answers the guard-manifest protocol so `bench guards`
# can classify this hook. SessionStart denies nothing — it is informational — and
# the `nothing (informational)` denies clause is how the aggregator excludes it
# from the guard rows.
if [[ "${1:-}" == "--describe" ]]; then
  printf 'name: session-start\n'
  printf 'boundary: SessionStart\n'
  printf 'denies: nothing (informational)\n'
  printf 'why: prints the ambient dashboard and guard brief on session open; never blocks\n'
  exit 0
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
"$cmd" status 2>/dev/null || true
# The guard brief: one line per deny-capable guard plus a pointer. Never blocks —
# any failure is swallowed so the session opens regardless.
"$cmd" guards --brief 2>/dev/null || true
