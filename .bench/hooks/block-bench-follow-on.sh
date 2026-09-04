#!/usr/bin/env bash
# name: block-bench-follow-on
# boundary: PreToolUse:Bash
# denies: Bench shell follow-ons and a cd, an assignment, or a git -C into the Bench worktree pool path
# why: Bench responses are bounded, complete, and self-contained evidence
set -uo pipefail
warn() { echo "WARNING: block-bench-follow-on: $1 — allowing Bash." >&2; }

# kit_root names the tree the rebuild runs in. BENCH_KIT is the caller's own answer
# when it has one; otherwise the enclosing checkout is the best available reading.
kit_root() {
  if [[ -n "${BENCH_KIT:-}" ]]; then
    printf '%s' "$BENCH_KIT"
    return 0
  fi
  git rev-parse --show-toplevel 2>/dev/null
}

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
  printf '%s' "$input" | "$cmd" guard-bench-follow-on 2>"$errors" || rc=$?
  answer="$(cat "$errors")"
else
  answer=""
  printf '%s' "$input" | "$cmd" guard-bench-follow-on || rc=$?
fi

# A stale core splits the posture: a Bench call refuses, because no verb may run
# stale, and every other call passes with the rebuild named, because the shell has to
# stay usable while the binary is rebuilt.
if [[ "$rc" == 2 && "$answer" == *"unknown subcommand"* ]]; then
  action="$(bench_rebuild_action "$(kit_root)")"
  # The core owns the envelope read in every ordinary call; the shim reaches for the
  # shared reader only here, where the core cannot answer. An unreadable envelope is
  # not a Bench call, so the degraded posture stays open rather than denying a call it
  # cannot read.
  if command_text="$(bench_envelope_command "$input")" && bench_invokes_bench "$command_text"; then
    echo "BLOCKED: the Bench binary is stale, so no Bench verb can answer. Rebuild it: $action" >&2
    exit 2
  fi
  warn "the Bench binary is stale, so the follow-on guard is off for this call. Rebuild it: $action"
  exit 0
fi

[[ -n "$answer" ]] && printf '%s' "$answer" >&2
case "$rc" in 0|2) exit "$rc" ;; *) warn "bench core errored (exit $rc)"; exit 0 ;; esac
