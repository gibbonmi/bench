#!/usr/bin/env bash
# name: session-start
# boundary: SessionStart
# denies: nothing (informational)
# why: prints the CLI location, ambient dashboard, and guard brief on session open; never blocks
# This SessionStart hook prints the ambient dashboard when a session opens cold. It
# is a thin resolver shim over `bench session-inspect`, whose Go core runs the ambient
# phases under one aggregate deadline. It never blocks the session, and it stays silent
# outside a repo. An in-repo resume failure still shows as a preservation warning.
#
# This hook is shared across harnesses. The .claude adapter wires it under
# hooks.SessionStart, and any AGENTS.md harness can point its own start hook here.
set -uo pipefail

# One git question settles both halves of the hook's contract. Outside a repository
# there is nothing ambient to say and nothing to recover, so the hook exits silently
# before it resolves anything: a session opened somewhere else must not learn about
# this one. Inside a repository, the root it found is the tree the recovery hint
# below names.
root="$(git rev-parse --show-toplevel 2>/dev/null)" || exit 0
[[ -n "$root" ]] || exit 0

# Resolve the wrapper through the shared resolver — the one source for the search
# order — relative to this script; `pwd -P` survives the .claude/hooks symlink. A
# missing lib fails like any other case: print nothing, and exit 0. Never block a
# session.
lib="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)/../lib/resolve-bench.sh"
[[ -f "$lib" ]] || exit 0
# shellcheck source=../lib/resolve-bench.sh
. "$lib"

# no_core_hint prints what a cold session gets when nothing here can render the
# dashboard. It names the build script, not plain `go build`, because `go build`
# produces a binary that carries no version and fails the freshness check the session
# would hit next.
no_core_hint() {
  printf 'bench: no Bench core is reachable here, so this session opened without its dashboard. Rebuild it with: %s\n' "$(bench_rebuild_action "$root")"
}

cmd="$(bench_resolve_wrapper)" || { no_core_hint; exit 0; }
# Advertise the CLI location from the same resolution that runs below, so a cold
# session learns how to invoke bench here, and the advice cannot drift from reality.
if command -v bench >/dev/null 2>&1; then
  printf 'bench CLI: %s (invoke as: bench)\n' "$cmd"
else
  printf 'bench CLI: %s (bench not on PATH; invoke by path — run `bench doctor --fix` to install a stable-PATH shim)\n' "$cmd"
fi
# Run the wrapper rather than exec it. The wrapper exits non-zero exactly when it
# cannot reach the Go core — 127 when no binary exists for this platform — and exec
# would hand that code, and the session's silence, straight to the harness. Running it
# keeps the hook's own exit at 0, so a broken core costs a hint, not a session.
"$cmd" session-inspect || no_core_hint
exit 0
