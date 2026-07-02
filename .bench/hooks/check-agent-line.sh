#!/usr/bin/env bash
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
# Fail-open is deliberate: a broken hook must never brick delegation. Missing
# lines.env, an unset/empty tier, unparseable stdin, or no model field all allow
# the call with a one-line stderr warning. Only a present model that matches none
# of the three bound tiers (exact string compare) is denied.
#
# Wire under hooks.PreToolUse with matcher "Agent". Exit 2 denies and returns the
# message to the agent.
set -euo pipefail

warn() { echo "WARNING: check-agent-line: $1 — allowing delegation." >&2; }
deny() { echo "DENIED: $1" >&2; exit 2; }

input="$(cat 2>/dev/null || true)"

# Pull the delegation's target model: resolvedModel wins, model is the fallback.
# python exits nonzero on unparseable JSON; an empty result means neither field
# was present (or both were empty) — both fail open.
if ! model="$(printf '%s' "$input" | python3 -c 'import sys,json; d=json.load(sys.stdin).get("tool_input",{}); print(d.get("resolvedModel") or d.get("model") or "")' 2>/dev/null)"; then
  warn "stdin is not parseable as JSON"
  exit 0
fi
[[ -z "$model" ]] && { warn "no resolvedModel/model field in tool_input"; exit 0; }

# Read the binding from the repo's .bench/lines.env. Locate the root the same way
# the sibling hooks do (git rev-parse --show-toplevel); no root => no binding file.
root="$(git rev-parse --show-toplevel 2>/dev/null || true)"
lines_env="${root:+$root/}.bench/lines.env"
[[ -f "$lines_env" ]] || { warn "no .bench/lines.env at repo root"; exit 0; }

# Parse the tier values by grep, not by sourcing: lines.env is a repo file and a
# blind `source` would execute anything in it. Take the last assignment for each
# key, strip surrounding quotes and any trailing carriage return.
tier_value() {
  local key="$1" line
  line="$(grep -E "^[[:space:]]*${key}=" "$lines_env" 2>/dev/null | tail -n1)" || true
  line="${line#*=}"
  line="${line%$'\r'}"
  line="${line#\"}"; line="${line%\"}"
  line="${line#\'}"; line="${line%\'}"
  printf '%s' "$line"
}

top="$(tier_value BENCH_TIER_TOP)"
mid="$(tier_value BENCH_TIER_MID)"
cheap="$(tier_value BENCH_TIER_CHEAP)"

# Any unset/empty tier means the binding is incomplete — fail open rather than
# deny against a partial oracle.
if [[ -z "$top" || -z "$mid" || -z "$cheap" ]]; then
  warn "a BENCH_TIER_* value is unset or empty in .bench/lines.env"
  exit 0
fi

# The model is present and the binding is complete: it must be exactly one tier.
if [[ "$model" == "$top" || "$model" == "$mid" || "$model" == "$cheap" ]]; then
  exit 0
fi

deny "delegation model '$model' is not a bound tier; bound: top=$top mid=$mid cheap=$cheap (see .bench/lines.env and the craft-line skill). Re-delegate on a bound tier or update the binding."
