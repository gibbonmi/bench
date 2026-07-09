#!/usr/bin/env bash
# Canary fixture: the git guard shim's inlined resolve_wrapper() has its two
# wrapper candidates SWAPPED (kit wrapper searched before the repo wrapper).
# The guard-resolver-order-drift conformance check must go red comparing this
# against .bench/lib/resolve-bench.sh's canonical repo-then-kit order.
set -uo pipefail

resolve_wrapper() {
  local root candidate
  root="$(git rev-parse --show-toplevel 2>/dev/null || true)"
  if [[ -n "$root" ]]; then
    for candidate in "$root/bin/bench.sh" "$root/.bench/bin/bench.sh"; do
      [[ -x "$candidate" ]] && { printf '%s\n' "$candidate"; return 0; }
    done
  fi
  command -v bench 2>/dev/null || return 1
}

input="$(cat)"
cmd="$(resolve_wrapper)" || exit 0
printf '%s' "$input" | "$cmd" guard-git
