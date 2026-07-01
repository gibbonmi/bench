#!/usr/bin/env bash
# Stop hook: the completion oracle. When a shift is armed (BENCH_SHIFT=1), the
# agent is not allowed to stop while the gate is red. This is what makes "the
# gate is the oracle, you never grade your own work" enforceable rather than
# aspirational — the agent cannot declare done by fiat.
#
# Wire it in .claude/settings.json under hooks.Stop. Claude Code feeds the hook
# JSON on stdin; exit 2 blocks the stop and returns stderr to the agent.
set -euo pipefail

# Only enforce inside an armed shift — never block ordinary conversational stops.
[[ "${BENCH_SHIFT:-0}" == "1" ]] || exit 0

# Read the hook JSON once. When Claude Code has already looped on a prior block it
# re-invokes with stop_hook_active=true; honor it (allow the stop) so a permanently
# red gate cannot trap the agent forever. Empty or invalid stdin => not active.
input="$(cat 2>/dev/null || true)"
stop_active="$(printf '%s' "$input" | python3 -c 'import sys,json; print(json.load(sys.stdin).get("stop_hook_active", False))' 2>/dev/null || true)"
[[ "$stop_active" == "True" ]] && exit 0

# Record the gate verdict for `bench status` to read (the ambient surface never runs
# the gate cold). The cache lives in the git dir, so it is never tracked or committed:
#   <status> <HEAD sha> <iso8601>
record_gate() {
  local gitdir
  gitdir="$(git rev-parse --absolute-git-dir 2>/dev/null)" || return 0
  printf '%s %s %s\n' "$1" "$(git rev-parse HEAD 2>/dev/null || echo none)" \
    "$(date -u +%Y-%m-%dT%H:%M:%SZ)" > "$gitdir/bench-last-gate"
}

bench_cmd() {
  local root candidate
  root="$(git rev-parse --show-toplevel 2>/dev/null || true)"
  if [[ -n "$root" ]]; then
    for candidate in "$root/.bench/bin/bench.sh" "$root/bin/bench.sh"; do
      [[ -x "$candidate" ]] && { printf '%s\n' "$candidate"; return 0; }
    done
  fi
  command -v bench 2>/dev/null || return 1
}

# A findable CLI is what gives the gate authority here. If none exists, fail open:
# allow the stop with a loud warning and write no gate cache — never forge a verdict
# the ambient dashboard would then report as real.
if ! cmd="$(bench_cmd)"; then
  root="$(git rev-parse --show-toplevel 2>/dev/null || true)"
  {
    echo "WARNING: bench CLI not found — allowing this stop without a gate verdict."
    echo "  searched: ${root:+$root/.bench/bin/bench.sh, $root/bin/bench.sh, }bench on PATH"
  } >&2
  exit 0
fi

gate_out="${TMPDIR:-/tmp}/bench-gate.$$"
trap 'rm -f "$gate_out"' EXIT

if "$cmd" gate >"$gate_out" 2>&1; then
  record_gate green
  exit 0   # green — allow the agent to stop
fi
record_gate red

{
  echo "BLOCKED: the gate is red, so this shift is not done."
  echo "Fix the failing checks at the seam — do not weaken or skip a check —"
  echo "then stop again. Gate output:"
  if [[ -s "$gate_out" ]]; then
    tail -n 30 "$gate_out"
  else
    echo "could not locate bench; expected .bench/bin/bench.sh, bin/bench.sh, or bench on PATH"
  fi
} >&2
exit 2
