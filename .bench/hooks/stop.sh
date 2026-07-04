#!/usr/bin/env bash
# Stop hook: the completion oracle. When a shift is armed (BENCH_SHIFT=1), the
# agent is not allowed to stop while the gate is red. This is what makes "the
# gate is the oracle, you never grade your own work" enforceable rather than
# aspirational — the agent cannot declare done by fiat.
#
# This is a thin shim over the Go core: it resolves the bench wrapper via the shared
# resolver and passes it to `bench stop-verdict`, which reads the Stop envelope,
# honors stop_hook_active, enforces only when BENCH_SHIFT=1, runs `<wrapper> gate`,
# truncates the output for the BLOCKED message, and writes the verdict cache (keyed
# to the tree hash from the Go core — never a forged verdict). The shim owns only
# `--describe` and the fail-OPEN rim: a missing lib, wrapper, or core warns loudly,
# writes NO cache, and allows the stop (exit 0). A missing oracle must degrade to
# "no verdict", never a false green — the mirror of the git guard's fail-closed.
#
# Wire it in .claude/settings.json under hooks.Stop (and .codex/hooks.json Stop).
# The harness feeds the hook JSON on stdin; exit 2 blocks the stop and returns
# stderr to the agent. `--describe` (first arg) prints the guard manifest and exits
# 0 without reading stdin.
set -uo pipefail

# `--describe` short-circuits before any stdin read or wrapper resolution. The stop
# manifest's denies clause is fixed (unlike the agent-line guard's live binding), so
# it needs no core call.
if [[ "${1:-}" == "--describe" ]]; then
  printf 'name: stop\n'
  printf 'boundary: Stop\n'
  printf 'denies: stopping an armed shift (BENCH_SHIFT=1) while the gate is red\n'
  printf 'why: the gate is the completion oracle; the agent cannot declare done on red\n'
  exit 0
fi

# Fail open: allow the stop with a loud warning and write no gate cache — never forge
# a verdict the ambient dashboard would then report as real.
warn_open() {
  echo "WARNING: $1 — allowing this stop without a gate verdict (no forged verdict)." >&2
  exit 0
}

# The shared resolver is resolved relative to THIS script (pwd -P survives the
# .claude/hooks symlink). A missing lib is treated exactly like a missing core.
lib="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)/../lib/resolve-bench.sh"
[[ -f "$lib" ]] || warn_open "bench wrapper resolver missing at $lib"
# shellcheck source=../lib/resolve-bench.sh
. "$lib"

input="$(cat 2>/dev/null || true)"

# A findable CLI is what gives the gate authority here. If none exists, fail open.
cmd="$(bench_resolve_wrapper)" || warn_open "bench CLI not found"

# Hand the envelope and the resolved wrapper to the core. It checks stop_hook_active
# and BENCH_SHIFT, runs `<wrapper> gate`, writes the cache, and owns the verdict:
# exit 0 allow, exit 2 block (with the BLOCKED message + gate tail already on stderr).
# The wrapper exits 127 when no binary is installed for this platform.
rc=0
printf '%s' "$input" | "$cmd" stop-verdict "$cmd" || rc=$?
case "$rc" in
  0 | 2) exit "$rc" ;;                                # allow / block — the core owns it
  *) warn_open "bench core errored (exit $rc)" ;;      # 127 missing binary / 3 panic → fail open, no cache
esac
