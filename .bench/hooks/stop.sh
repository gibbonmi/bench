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

# Avoid an infinite loop if the hook itself re-triggers.
[[ "${BENCH_STOP_CHECKED:-0}" == "1" ]] && exit 0
export BENCH_STOP_CHECKED=1

if bench gate >/tmp/bench-gate.out 2>&1; then
  exit 0   # green — allow the agent to stop
fi

{
  echo "BLOCKED: the gate is red, so this shift is not done."
  echo "Fix the failing checks at the seam — do not weaken or skip a check —"
  echo "then stop again. Gate output:"
  tail -n 30 /tmp/bench-gate.out
} >&2
exit 2
