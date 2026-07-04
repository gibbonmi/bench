#!/usr/bin/env bash
# Stop hook: the completion oracle. When a shift is armed (BENCH_SHIFT=1), the
# agent is not allowed to stop while the gate is red. This is what makes "the
# gate is the oracle, you never grade your own work" enforceable rather than
# aspirational — the agent cannot declare done by fiat.
#
# Wire it in .claude/settings.json under hooks.Stop. Claude Code feeds the hook
# JSON on stdin; exit 2 blocks the stop and returns stderr to the agent. `--describe`
# (first arg) prints the guard manifest and exits 0 without reading stdin, so
# `bench guards` can aggregate the deny surface.
set -euo pipefail

# `--describe` short-circuits before any stdin read or shift-arming check.
if [[ "${1:-}" == "--describe" ]]; then
  printf 'name: stop\n'
  printf 'boundary: Stop\n'
  printf 'denies: stopping an armed shift (BENCH_SHIFT=1) while the gate is red\n'
  printf 'why: the gate is the completion oracle; the agent cannot declare done on red\n'
  exit 0
fi

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
#   <status> <tree hash> <iso8601>
# The key is the tree actually tested, computed by the Go core's `tree-hash`
# (git.TreeHash) through the same bench.sh wrapper this hook already resolves ($cmd) —
# so the hash has ONE source, shared with gate_record in bin/bench.sh (the mirror this
# hook used to carry is gone). If the core binary is missing the tree-hash is
# unavailable: skip the write loudly and never record a verdict keyed to a guessed tree,
# so a missing platform binary degrades to "no verdict", never a false green.
record_gate() {
  local gitdir tree
  gitdir="$(git rev-parse --absolute-git-dir 2>/dev/null)" || return 0
  tree="$("$cmd" tree-hash 2>/dev/null)" || tree=""
  if [[ ! "$tree" =~ ^[0-9a-f]+$ ]]; then
    echo "WARNING: bench tree-hash unavailable — not recording a gate verdict (no forged tree)." >&2
    return 0
  fi
  printf '%s %s %s\n' "$1" "$tree" \
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
