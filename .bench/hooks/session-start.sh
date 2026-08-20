#!/usr/bin/env bash
# name: session-start
# boundary: SessionStart
# denies: nothing (informational)
# why: prints the CLI location, ambient dashboard, and guard brief on session open; never blocks
# SessionStart hook: print the ambient dashboard when a session opens cold.
# A thin resolver shim over `bench session-inspect`, whose Go core runs the
# ambient phases under one aggregate deadline. It never blocks the session and stays
# silent outside a repo; an in-repo resume failure remains visible as a preservation warning. Shared across
# harnesses; the .claude adapter wires it under hooks.SessionStart, and any AGENTS.md
# harness can point its own start hook here.
set -uo pipefail

# One git question settles both halves of the hook's contract. Outside a repository
# there is nothing ambient to say and nothing to recover, so the hook exits silent
# before it resolves anything — a session opened somewhere else must not be told about
# this one. Inside one, the root it found is the tree the recovery hint below names.
root="$(git rev-parse --show-toplevel 2>/dev/null)" || exit 0
[[ -n "$root" ]] || exit 0

# Resolve the wrapper via the shared resolver (one source for the search order),
# relative to this script (pwd -P survives the .claude/hooks symlink). A missing lib
# is failure like any other: print nothing, exit 0 — never block a session.
lib="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)/../lib/resolve-bench.sh"
[[ -f "$lib" ]] || exit 0
# shellcheck source=../lib/resolve-bench.sh
. "$lib"

# shell_quote and rebuild_action are the shell spelling of internal/freshness's
# shellQuote and RebuildAction. The Go function is the source of the rebuild invocation;
# this copy exists only because the one moment the invocation is needed is the moment the
# Go binary cannot run, so nothing can ask it for the string. internal/systemtest asserts
# the printed line against freshness.RebuildAction, so a drift between the two reds.
shell_quote() { printf "'%s'" "${1//\'/\'\\\'\'}"; }
rebuild_action() {
  printf 'cd %s && bash scripts/go-build.sh %s %s' "$(shell_quote "$1")" "$(shell_quote "$1")" "$(shell_quote "$1/dist/bench")"
}

# no_core_hint — what a cold session gets when nothing here can render the dashboard. It
# names the build script rather than plain `go build`, which produces a binary carrying
# no version and so fails the freshness check the session would hit next.
no_core_hint() {
  printf 'bench: no Bench core is reachable here, so this session opened without its dashboard. Rebuild it with: %s\n' "$(rebuild_action "$root")"
}

cmd="$(bench_resolve_wrapper)" || { no_core_hint; exit 0; }
# Advertise the CLI location from the same resolution that runs below, so a cold
# session is told how to invoke bench here and the advice cannot drift from reality.
if command -v bench >/dev/null 2>&1; then
  printf 'bench CLI: %s (invoke as: bench)\n' "$cmd"
else
  printf 'bench CLI: %s (bench not on PATH; invoke by path — run `bench doctor --fix` to install a stable-PATH shim)\n' "$cmd"
fi
# Run rather than exec: the wrapper exits non-zero exactly when it could not reach the Go
# core (127 when no binary is present for this platform), and execing handed that code —
# and the session's silence — straight to the harness. Running it keeps the hook's own
# exit at 0, so a broken core costs a hint rather than a session.
"$cmd" session-inspect || no_core_hint
exit 0
