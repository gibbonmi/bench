#!/usr/bin/env bash
# SessionStart hook: print the ambient dashboard when a session opens cold.
# A thin wrapper over `bench status` — the single renderer the user also runs on demand
# (one source of truth). Never blocks the session: outside a repo, or on any error, it
# prints nothing and exits 0. Shared across harnesses; the .claude adapter wires it under
# hooks.SessionStart, and any AGENTS.md harness can point its own start hook here.
set -uo pipefail

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
