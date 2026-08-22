#!/usr/bin/env bash
# name: check-agent-line
# boundary: PreToolUse:Agent
# denies: Agent delegation off the bound tier
# why: invariant #2 forbids silent escalation; a delegate runs on a bound tier or not at all
# This is a PreToolUse guard: a delegation must run on a bound tier. Invariant #2
# requires the line — model, effort, token cap — declared before a long run, with no
# silent escalation. This guard enforces the model half: an Agent-tool call on a model
# outside the three bound tiers is denied at the boundary.
#
# Threat model: an honest-mistake layer, like the git guard. It stops a well-meaning
# agent from delegating onto an unbound model. It is not evasion-resistant. The binding
# lives in .bench/lines.env, the BENCH_<HARNESS>_<TIER> matrix. The narrative lives in
# projects/<name>.md "Lines" and the craft-line skill.
#
# This is a thin shim over the Go core. It resolves the bench wrapper through the
# shared resolver, pipes the Agent envelope to `bench check-agent-line`, and passes
# the verdict exit code back. It names its own harness, so the core resolves the
# Claude column.
#
# All binding logic — parse model, read binding, verdict, the DENIED message, every
# warn-and-allow branch — lives in internal/lines. The shim owns only its fail-open
# rim: a broken guard must never brick delegation, so a missing lib, wrapper, or core
# warns and allows (exit 0). The core alone denies (exit 2), and only when a present
# model matches no bound tier.
#
# Wire this hook under hooks.PreToolUse with matcher "Agent". Exit 2 denies the call
# and returns the message to the agent.
set -uo pipefail

# This deliberately mirrors internal/lines' warn() format. The shim emits this only
# for its rims — a missing lib, wrapper, or core — where the Go core is by definition
# unreachable, so the shared prefix cannot be single-sourced across the shell/Go
# boundary.
warn() { echo "WARNING: check-agent-line: $1 — allowing delegation." >&2; }

# The shared wrapper resolver resolves relative to this script, not the cwd repo;
# `pwd -P` survives the .claude/hooks symlink. A hook copied without the lib fails
# open, like every other broken-hook case: a broken hook never bricks delegation.
lib="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)/../lib/resolve-bench.sh"

# Rim: a missing shared resolver lib is treated exactly like a missing core. The hook
# warns and allows, so a hook copied without its lib never bricks delegation.
[[ -f "$lib" ]] || { warn "wrapper resolver missing at $lib"; exit 0; }
# shellcheck source=../lib/resolve-bench.sh
. "$lib"

input="$(cat 2>/dev/null || true)"

# Rim: when no core is reachable, the hook fails open. A broken guard must never
# brick delegation.
cmd="$(bench_resolve_wrapper)" || { warn "bench core not found"; exit 0; }

# Hand the envelope to the core. It parses the model, reads the binding, and owns the
# verdict. Exit 0 allows the call, including any degraded warn-and-allow with its
# WARNING on stderr. Exit 2 denies the call, with the DENIED message on stderr. The
# wrapper exits 127 when this platform has no installed binary.
rc=0
printf '%s' "$input" | "$cmd" check-agent-line --harness claude || rc=$?
case "$rc" in
  0 | 2) exit "$rc" ;;                       # allow/degraded (0) or deny (2) — the core owns it
  *) warn "bench core errored (exit $rc)"; exit 0 ;;  # 127 missing binary / 3 panic → fail open
esac
