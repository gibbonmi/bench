#!/usr/bin/env bash
# Shared parser for the .bench/lines.env tier binding — the one source both
# enforcement surfaces load (.bench/hooks/check-agent-line.sh via source-or-fail-open,
# .bench/adapters/_line-guard.sh via source-or-fail-closed). The two surfaces must
# agree on what a tier value IS, so the parse lives here and nowhere else.
#
# Sourced, not run. Read a value from lines.env by grep, not by sourcing: lines.env
# is a repo file and a blind `source` would execute anything in it. Take the last
# assignment for the key, strip surrounding quotes and any trailing carriage return.
bench_tier_value() {
  local key="$1" file="$2" line
  line="$(grep -E "^[[:space:]]*${key}=" "$file" 2>/dev/null | tail -n1)" || true
  line="${line#*=}"
  line="${line%$'\r'}"
  line="${line%"${line##*[![:space:]]}"}"; line="${line#"${line%%[![:space:]]*}"}"
  line="${line#\"}"; line="${line%\"}"
  line="${line#\'}"; line="${line%\'}"
  printf '%s' "$line"
}
