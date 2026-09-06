#!/usr/bin/env bash
# name: block-primary-file-write
# boundary: PreToolUse:Edit|Write|MultiEdit|NotebookEdit
# denies: a file-tool write to a tracked path in the primary checkout
# why: within Bench, main receives writes only through landings
set -uo pipefail
warn() { echo "WARNING: block-primary-file-write: $1 — allowing the write." >&2; }

lib="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)/../lib/resolve-bench.sh"
[[ -f "$lib" ]] || { warn "wrapper resolver missing"; exit 0; }
# shellcheck source=../lib/resolve-bench.sh
. "$lib"
input="$(cat 2>/dev/null || true)"
cmd="$(bench_resolve_wrapper)" || { warn "bench core not found"; exit 0; }

# The core's stderr is captured rather than forwarded, because a stale binary answers
# `unknown subcommand` at exit 2 and that answer must not read as a refusal.
errors="$(mktemp 2>/dev/null)" || errors=""
rc=0
if [[ -n "$errors" ]]; then
  trap 'rm -f "$errors"' EXIT
  printf '%s' "$input" | "$cmd" guard-file-write 2>"$errors" || rc=$?
  answer="$(cat "$errors")"
else
  answer=""
  printf '%s' "$input" | "$cmd" guard-file-write || rc=$?
fi

# A file edit outside a repository has to stay possible while the binary is rebuilt, so
# a stale core warns rather than refusing every write the harness attempts.
if [[ "$rc" == 2 && "$answer" == *"unknown subcommand"* ]]; then
  warn "the Bench binary is stale, so the primary-checkout guard is off for this call"
  exit 0
fi

[[ -n "$answer" ]] && printf '%s' "$answer" >&2
case "$rc" in 0|2) exit "$rc" ;; *) warn "bench core errored (exit $rc)"; exit 0 ;; esac
