#!/usr/bin/env bash
# PreToolUse guard: the agent has no destructive git authority. This makes
# invariant #4 ("you assist, you don't decide where a decision is mine")
# enforceable for the operations that can silently destroy a shift's work or
# bypass the merge — not just aspirational.
#
# Threat model: this is an honest-mistake layer, not an evasion-resistant
# boundary. It stops a well-meaning agent from reflexively running destructive
# git; it does not try to survive deliberate evasion. Wrapper scanning goes
# exactly one level deep by design (see internal/gitguard). The backstops for a
# misaligned agent are the git pre-push hook and bench's pooled-worktree
# isolation, not this script.
#
# Note the boundary: this intercepts the AGENT's Bash tool calls. Bench's own
# controlled rollback inside `bench shift` runs in-process (not through the
# agent's shell), so the harness can still reset/clean a failed iteration while
# the agent itself cannot. That asymmetry is the point.
#
# This is a thin shim over the Go core: it resolves the bench wrapper, pipes the
# PreToolUse envelope to `bench guard-git`, and passes the verdict through. All
# classification (tokenize, scan, verdict, the BLOCKED message) lives in
# internal/gitguard. The shim owns exactly two fail-closed rims — core
# unresolvable/missing, and core errored — plus `--describe`. Wire under
# hooks.PreToolUse with matcher "Bash". Exit 2 blocks and returns the message to
# the agent.
set -uo pipefail

# resolve_wrapper echoes the bench.sh wrapper path (repo-local first, then a
# global `bench`), or fails when none is reachable. The ~8-line search is inlined
# rather than shared with .bench/hooks/stop.sh: sourcing a shared lib would give
# this hook a new fail-OPEN mode (missing lib → the shim errors before its rims
# run, and a non-2 PreToolUse exit is a non-blocking error that silently grants).
# The conformance check in internal/conformance (checkGuardResolverOrderDrift)
# reds if this inline's search order ever drifts from .bench/lib/resolve-bench.sh.
resolve_wrapper() {
  local root candidate
  root="$(git rev-parse --show-toplevel 2>/dev/null || true)"
  if [[ -n "$root" ]]; then
    for candidate in "$root/.bench/bin/bench.sh" "$root/bin/bench.sh"; do
      [[ -x "$candidate" ]] && { printf '%s\n' "$candidate"; return 0; }
    done
  fi
  command -v bench 2>/dev/null || return 1
}

# `--describe` prints the guard manifest and exits 0 without reading stdin, so
# `bench guards` can aggregate the deny surface without a collision. The denies
# clause is generated from the core's `guard-git --describe-classes` — the same
# deny table classification uses, so the advertisement cannot drift. Core
# unreachable → the manifest degrades honestly rather than lying.
if [[ "${1:-}" == "--describe" ]]; then
  printf 'name: block-dangerous-git\n'
  printf 'boundary: PreToolUse:Bash\n'
  if cmd="$(resolve_wrapper)" && classes="$("$cmd" guard-git --describe-classes 2>/dev/null)" && [[ -n "$classes" ]]; then
    printf 'denies: destructive git — %s\n' "$classes"
  else
    printf 'denies: manifest unavailable (analyzer missing)\n'
  fi
  printf 'why: agents lack destructive-git authority; merge and history rewrites belong to the reviewer\n'
  exit 0
fi

input="$(cat)"

# fail_closed_git_shaped is the "cannot classify" rim: with no reachable core the
# shim can't classify at all, so it refuses anything git-shaped (the destructive
# surface) while leaving non-git commands runnable so the shell stays usable. The
# raw envelope carries the command text, so a substring test catches the honest
# mistake without parsing.
fail_closed_git_shaped() {
  case "$input" in
    *git*) echo "BLOCKED: guard degraded (bench core missing) — can't classify commands, refusing anything git-shaped. Restore the bench core (bench link) or hand back." >&2; exit 2 ;;
    *) exit 0 ;;
  esac
}

# Rim 1: core unresolvable. No wrapper on disk or PATH → cannot classify.
cmd="$(resolve_wrapper)" || fail_closed_git_shaped

# Hand the envelope to the core. The core writes its own BLOCKED message to stderr
# and exits 2 on a block, 0 on allow; the wrapper exits 127 when no binary is
# installed for this platform.
rc=0
printf '%s' "$input" | "$cmd" guard-git || rc=$?
case "$rc" in
  0 | 2) exit "$rc" ;;                 # allow / block — the core owns the verdict + message
  127) fail_closed_git_shaped ;;       # rim 1: binary missing for this platform
  *) echo "BLOCKED: guard analyzer error — failing closed; rephrase the command." >&2; exit 2 ;;  # rim 2: core errored
esac
