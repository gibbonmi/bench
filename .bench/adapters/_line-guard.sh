#!/usr/bin/env bash
# Shared model guard for the Bench shift adapters (claude, codex, opencode).
# Invariant #2 says a long run declares its line — model, effort, token cap — with
# no silent escalation. Headless `bench shift` has no interactive Agent-line hook,
# so this guard makes the model half enforceable at the adapter boundary: in a
# routed repo the delegation must carry a BENCH_MODEL naming a bound tier, or the
# adapter refuses to launch the harness.
#
# Sourced, not run. `bench_resolve_model` sets BENCH_RESOLVED_MODEL to the model to
# pass via the harness's model flag (empty => pass no flag), or exits 1 with a
# legible stderr error. The binding lives in .bench/lines.env
# (BENCH_TIER_TOP/MID/CHEAP); the narrative lives in projects/<name>.md "Lines" and
# BENCH.md "Harness adapter for the shift loop".
#
# Rules:
#   - No lines.env at the repo root (or not in a repo): unrouted. An explicit
#     BENCH_MODEL still passes through; an absent one means no model flag — today's
#     behavior.
#   - lines.env present but any tier empty/missing: incomplete binding — warn and
#     fall back to the unrouted passthrough rather than deny against a partial oracle.
#   - lines.env present and complete: BENCH_MODEL must name exactly one bound tier.
#     Unset or unbound => stderr error naming the three models, exit 1, harness never runs.

# Read a tier value from lines.env by grep, not by sourcing: lines.env is a repo
# file and a blind `source` would execute anything in it. Take the last assignment
# for the key, strip surrounding quotes and any trailing carriage return.
_bench_tier_value() {
  local key="$1" file="$2" line
  line="$(grep -E "^[[:space:]]*${key}=" "$file" 2>/dev/null | tail -n1)" || true
  line="${line#*=}"
  line="${line%$'\r'}"
  line="${line%"${line##*[![:space:]]}"}"; line="${line#"${line%%[![:space:]]*}"}"
  line="${line#\"}"; line="${line%\"}"
  line="${line#\'}"; line="${line%\'}"
  printf '%s' "$line"
}

bench_resolve_model() {
  BENCH_RESOLVED_MODEL=""

  local root lines_env
  root="$(git rev-parse --show-toplevel 2>/dev/null)" || true
  lines_env="${root:+$root/}.bench/lines.env"

  # Unrouted repo (or not a repo): explicit beats absent.
  if [ -z "$root" ] || [ ! -f "$lines_env" ]; then
    BENCH_RESOLVED_MODEL="${BENCH_MODEL:-}"
    return 0
  fi

  local top mid cheap
  top="$(_bench_tier_value BENCH_TIER_TOP "$lines_env")"
  mid="$(_bench_tier_value BENCH_TIER_MID "$lines_env")"
  cheap="$(_bench_tier_value BENCH_TIER_CHEAP "$lines_env")"

  # Incomplete binding: treat as absent, warn, fall back to passthrough.
  if [ -z "$top" ] || [ -z "$mid" ] || [ -z "$cheap" ]; then
    printf '%s\n' "WARNING: bench adapter: a BENCH_TIER_* value is unset or empty in $lines_env — ignoring the binding and falling back to BENCH_MODEL." >&2
    BENCH_RESOLVED_MODEL="${BENCH_MODEL:-}"
    return 0
  fi

  # Routed repo: the line must be declared and bound.
  if [ -z "${BENCH_MODEL:-}" ]; then
    printf '%s\n' "bench shift in a routed repo requires a declared line: set BENCH_MODEL to one of top=$top mid=$mid cheap=$cheap" >&2
    exit 1
  fi

  if [ "$BENCH_MODEL" != "$top" ] && [ "$BENCH_MODEL" != "$mid" ] && [ "$BENCH_MODEL" != "$cheap" ]; then
    printf '%s\n' "bench shift: BENCH_MODEL='$BENCH_MODEL' is not a bound model; set it to one of top=$top mid=$mid cheap=$cheap" >&2
    exit 1
  fi

  BENCH_RESOLVED_MODEL="$BENCH_MODEL"
  return 0
}
