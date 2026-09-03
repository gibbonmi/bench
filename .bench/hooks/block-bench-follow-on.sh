#!/usr/bin/env bash
# name: block-bench-follow-on
# boundary: PreToolUse:Bash
# denies: Bench shell follow-ons and a cd, an assignment, or a git -C into the Bench worktree pool path
# why: Bench responses are bounded, complete, and self-contained evidence
set -uo pipefail
warn() { echo "WARNING: block-bench-follow-on: $1 — allowing Bash." >&2; }

# rebuild_action composes the one rebuild command for a kit root. It is the shell
# derivation of internal/freshness.RebuildAction, which the shim cannot call because
# the binary that would answer is the stale thing it is reporting. The system suite
# pins the two derivations byte for byte.
shell_quote() { printf "'%s'" "${1//\'/\'\\\'\'}"; }
rebuild_action() {
  printf 'cd %s && bash scripts/go-build.sh %s %s' "$(shell_quote "$1")" "$(shell_quote "$1")" "$(shell_quote "$1/dist/bench")"
}

# kit_root names the tree the rebuild runs in. BENCH_KIT is the caller's own answer
# when it has one; otherwise the enclosing checkout is the best available reading.
kit_root() {
  if [[ -n "${BENCH_KIT:-}" ]]; then
    printf '%s' "$BENCH_KIT"
    return 0
  fi
  git rev-parse --show-toplevel 2>/dev/null
}

# envelope_command reads the Bash command out of the hook envelope. The core owns
# this read in every ordinary call; the shim needs its own only when the core cannot
# answer. An unreadable field yields the empty string, which classifies as non-Bench,
# so the degraded posture stays open rather than denying a call it cannot read.
envelope_command() {
  local text
  text="$(printf '%s' "$1" | sed -n 's/.*"command"[[:space:]]*:[[:space:]]*"\(.*\)".*/\1/p')"
  text="${text//\\n/$'\n'}"
  text="${text//\\\"/\"}"
  printf '%s' "${text//\\\\/\\}"
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
  action="$(rebuild_action "$(kit_root)")"
  if bench_invokes_bench "$(envelope_command "$input")"; then
    echo "BLOCKED: the Bench binary is stale, so no Bench verb can answer. Rebuild it: $action" >&2
    exit 2
  fi
  warn "the Bench binary is stale, so the follow-on guard is off for this call. Rebuild it: $action"
  exit 0
fi

[[ -n "$answer" ]] && printf '%s' "$answer" >&2
case "$rc" in 0|2) exit "$rc" ;; *) warn "bench core errored (exit $rc)"; exit 0 ;; esac
