#!/usr/bin/env bash
# name: check-agent-line
# boundary: PreToolUse:Agent
# denies: Agent delegation off the bound tier
# why: invariant #2 forbids silent escalation; a delegate runs on a bound tier or not at all
# PreToolUse guard: a delegation must run on a bound tier. Invariant #2 says the
# line (model, effort, token cap) is declared before a long run and there is no
# silent escalation; this makes the model half enforceable rather than
# aspirational — an Agent-tool call on a model that isn't one of the three bound
# tiers is denied at the boundary.
#
# Threat model: honest-mistake layer, like the git guard. It stops a well-meaning
# agent from delegating onto an unbound model; it is not evasion-resistant. The
# binding lives in .bench/lines.env (BENCH_TIER_TOP/MID/CHEAP); the narrative
# lives in projects/<name>.md "Lines" and the craft-line skill.
#
# This is a thin shim over the Go core: it resolves the bench wrapper via the shared
# resolver, pipes the Agent envelope to `bench check-agent-line`, and passes the
# verdict exit code back. All binding logic (parse model, read binding, verdict, the
# DENIED message and every warn-and-allow branch) lives in internal/lines. The shim
# owns only its fail-open rim: a broken guard must never brick
# delegation, so a missing lib, wrapper, or core warns and ALLOWS (exit 0). Only a
# present model matching no bound tier is denied (exit 2), by the core.
#
# Wire under hooks.PreToolUse with matcher "Agent". Exit 2 denies and returns the
# message to the agent.
set -uo pipefail

# Intentionally mirrors internal/lines' warn() format: the shim emits this only for
# its rims (missing lib/wrapper/core) where the Go core is by definition unreachable,
# so the shared prefix cannot be single-sourced across the shell/Go boundary.
warn() { echo "WARNING: check-agent-line: $1 — allowing delegation." >&2; }

# The shared wrapper resolver is resolved relative to THIS script (pwd -P survives the
# .claude/hooks symlink), not the cwd repo. A hook copied without the lib fails open
# like every other broken-hook case — a broken hook never bricks delegation.
lib="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)/../lib/resolve-bench.sh"

# Rim: a missing shared resolver lib is treated exactly like a missing core — warn
# and allow, so a hook copied without its lib never bricks delegation.
[[ -f "$lib" ]] || { warn "wrapper resolver missing at $lib"; exit 0; }
# shellcheck source=../lib/resolve-bench.sh
. "$lib"

input="$(cat 2>/dev/null || true)"

# Rim: no reachable core → fail open. A broken guard must never brick delegation.
cmd="$(bench_resolve_wrapper)" || { warn "bench core not found"; exit 0; }

# Hand the envelope to the core. It parses the model, reads the binding, and owns the
# verdict: exit 0 allow (or any degraded warn-and-allow, with its WARNING on stderr),
# exit 2 deny (with the DENIED message on stderr). The wrapper exits 127 when no
# binary is installed for this platform.
rc=0
printf '%s' "$input" | "$cmd" check-agent-line || rc=$?
case "$rc" in
  0 | 2) exit "$rc" ;;                       # allow/degraded (0) or deny (2) — the core owns it
  *) warn "bench core errored (exit $rc)"; exit 0 ;;  # 127 missing binary / 3 panic → fail open
esac
