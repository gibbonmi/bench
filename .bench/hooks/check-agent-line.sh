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
# Harness aliases: the Claude Code Agent tool addresses models by alias (opus,
# fable, ...), not full id, so lines.env may declare which aliases bind via
# optional BENCH_ALIAS_TOP/MID/CHEAP. An undeclared alias is denied like any
# unbound model — deliberately: leaving BENCH_ALIAS_CHEAP undeclared is how this
# repo keeps bare `sonnet` (which resolves to an excluded model) out.
#
# Wire under hooks.PreToolUse with matcher "Agent". Exit 2 denies and returns the
# message to the agent. `--describe` (first arg) prints the guard manifest and
# exits 0 without reading stdin, so `bench guards` can aggregate the deny surface;
# the denies clause reads the live .bench/lines.env binding it enforces against
# (or `unrouted` when no binding is present), so it cannot drift from enforcement.
set -euo pipefail

warn() { echo "WARNING: check-agent-line: $1 — allowing delegation." >&2; }
deny() { echo "DENIED: $1" >&2; exit 2; }

# Read a tier value from lines.env by grep, not by sourcing (a repo file). Take the
# last assignment for the key, strip surrounding quotes and any trailing carriage
# return. Defined ahead of --describe so both paths read the binding identically.
tier_value() {
  local key="$1" line
  line="$(grep -E "^[[:space:]]*${key}=" "$lines_env" 2>/dev/null | tail -n1)" || true
  line="${line#*=}"
  line="${line%$'\r'}"
  line="${line%"${line##*[![:space:]]}"}"; line="${line#"${line%%[![:space:]]*}"}"
  line="${line#\"}"; line="${line%\"}"
  line="${line#\'}"; line="${line%\'}"
  printf '%s' "$line"
}

if [[ "${1:-}" == "--describe" ]]; then
  root="$(git rev-parse --show-toplevel 2>/dev/null || true)"
  lines_env="${root:+$root/}.bench/lines.env"
  printf 'name: check-agent-line\n'
  printf 'boundary: PreToolUse:Agent\n'
  if [[ -f "$lines_env" ]]; then
    top="$(tier_value BENCH_TIER_TOP)"; mid="$(tier_value BENCH_TIER_MID)"; cheap="$(tier_value BENCH_TIER_CHEAP)"
    printf 'denies: Agent delegation off the bound line (top=%s mid=%s cheap=%s)\n' \
      "${top:--}" "${mid:--}" "${cheap:--}"
  else
    printf 'denies: unrouted (no .bench/lines.env binding)\n'
  fi
  printf 'why: invariant #2 forbids silent escalation; a delegate runs on a bound tier or not at all\n'
  exit 0
fi

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

# Parse the tier values by grep, not by sourcing (tier_value is defined above the
# --describe branch so both paths read the binding the same way).
top="$(tier_value BENCH_TIER_TOP)"
mid="$(tier_value BENCH_TIER_MID)"
cheap="$(tier_value BENCH_TIER_CHEAP)"

# Any unset/empty tier means the binding is incomplete — fail open rather than
# deny against a partial oracle.
if [[ -z "$top" || -z "$mid" || -z "$cheap" ]]; then
  warn "a BENCH_TIER_* value is unset or empty in .bench/lines.env"
  exit 0
fi

# Optional alias declarations (may be absent — absence means no alias binds).
alias_top="$(tier_value BENCH_ALIAS_TOP)"
alias_mid="$(tier_value BENCH_ALIAS_MID)"
alias_cheap="$(tier_value BENCH_ALIAS_CHEAP)"

# The model is present and the binding is complete: it must be exactly one bound
# tier id or one declared alias.
for bound in "$top" "$mid" "$cheap" "$alias_top" "$alias_mid" "$alias_cheap"; do
  [[ -n "$bound" && "$model" == "$bound" ]] && exit 0
done

deny "delegation model '$model' is not a bound tier; bound: top=$top mid=$mid cheap=$cheap aliases: top=${alias_top:--} mid=${alias_mid:--} cheap=${alias_cheap:--} (see .bench/lines.env and the craft-line skill). Re-delegate on a bound tier or update the binding."
