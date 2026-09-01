#!/usr/bin/env bash
# name: block-bench-follow-on
# boundary: PreToolUse:Bash
# denies: Bench shell follow-ons and a cd, an assignment, or a git -C into the Bench worktree pool path
# why: Bench responses are bounded, complete, and self-contained evidence
set -uo pipefail
warn() { echo "WARNING: block-bench-follow-on: $1 — allowing Bash." >&2; }
lib="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)/../lib/resolve-bench.sh"
[[ -f "$lib" ]] || { warn "wrapper resolver missing"; exit 0; }
# shellcheck source=../lib/resolve-bench.sh
. "$lib"
input="$(cat 2>/dev/null || true)"
cmd="$(bench_resolve_wrapper)" || { warn "bench core not found"; exit 0; }
rc=0
printf '%s' "$input" | "$cmd" guard-bench-follow-on || rc=$?
case "$rc" in 0|2) exit "$rc" ;; *) warn "bench core errored (exit $rc)"; exit 0 ;; esac
