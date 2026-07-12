#!/usr/bin/env bash
# SessionStart hook: print the ambient dashboard when a session opens cold.
# A thin wrapper over `bench status` — the single renderer the user also runs on demand
# (one source of truth). It never blocks the session and stays silent outside a repo;
# an in-repo resume failure remains visible as a preservation warning. Shared across
# harnesses; the .claude adapter wires it under hooks.SessionStart, and any AGENTS.md
# harness can point its own start hook here.
set -uo pipefail

# `--describe` (first arg) answers the guard-manifest protocol so `bench guards`
# can classify this hook. SessionStart denies nothing — it is informational — and
# the `nothing (informational)` denies clause is how the aggregator excludes it
# from the guard rows.
if [[ "${1:-}" == "--describe" ]]; then
  printf 'name: session-start\n'
  printf 'boundary: SessionStart\n'
  printf 'denies: nothing (informational)\n'
  printf 'why: prints the CLI location, ambient dashboard, and guard brief on session open; never blocks\n'
  exit 0
fi

# Resolve the wrapper via the shared resolver (one source for the search order),
# relative to this script (pwd -P survives the .claude/hooks symlink). A missing lib
# is failure like any other: print nothing, exit 0 — never block a session.
lib="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)/../lib/resolve-bench.sh"
[[ -f "$lib" ]] || exit 0
# shellcheck source=../lib/resolve-bench.sh
. "$lib"

cmd="$(bench_resolve_wrapper)" || exit 0
# Advertise the CLI location from the same resolution that runs below, so a cold
# session is told how to invoke bench here and the advice cannot drift from reality.
if command -v bench >/dev/null 2>&1; then
  printf 'bench CLI: %s (invoke as: bench)\n' "$cmd"
else
  printf 'bench CLI: %s (bench not on PATH; invoke by path — run `bench doctor --fix` to install a stable-PATH shim)\n' "$cmd"
fi
if ! "$cmd" resume-clean; then
  printf 'warning: bench session-start: resume-clean failed; inspect retained worktree state\n' >&2
fi
"$cmd" status 2>/dev/null || true
# The guard brief: one line per deny-capable guard plus a pointer. Never blocks —
# any failure is swallowed so the session opens regardless.
"$cmd" guards --brief 2>/dev/null || true
