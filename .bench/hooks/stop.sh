#!/usr/bin/env bash
# name: stop
# boundary: Stop
# denies: stopping an armed shift while the gate is red
# why: the gate is the completion oracle; the agent cannot declare done on red
# This Stop hook is the completion oracle. When a shift is armed (BENCH_SHIFT=1),
# the agent cannot stop while the gate is red. This enforces the rule that the gate
# is the oracle and the agent never grades its own work: the agent cannot declare
# done by fiat.
#
# This is a thin shim over the Go core. It resolves the bench wrapper through the
# shared resolver and passes it to `bench stop-verdict`, which reads the Stop
# envelope, honors stop_hook_active, enforces only when BENCH_SHIFT=1, runs `<wrapper>
# gate`, truncates the output for the BLOCKED message, and writes the verdict cache —
# keyed to the tree hash from the Go core, never a forged verdict.
#
# The shim owns only the fail-open rim: a missing lib, wrapper, or core warns loudly,
# writes no cache, and allows the stop (exit 0). A missing oracle must degrade to "no
# verdict", never to a false green. This is the mirror of the git guard's
# fail-closed.
#
# Wire this hook in .claude/settings.json under hooks.Stop, and in .codex/hooks.json
# under Stop. The harness feeds the hook JSON on stdin. Exit 2 blocks the stop and
# returns stderr to the agent.
set -uo pipefail

# Fail open: allow the stop with a loud warning, and write no gate cache. Never
# forge a verdict the ambient dashboard would then report as real.
warn_open() {
  echo "WARNING: $1 — allowing this stop without a gate verdict (no forged verdict)." >&2
  exit 0
}

# The shared resolver resolves relative to this script; `pwd -P` survives the
# .claude/hooks symlink. A missing lib is treated exactly like a missing core.
lib="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)/../lib/resolve-bench.sh"
[[ -f "$lib" ]] || warn_open "bench wrapper resolver missing at $lib"
# shellcheck source=../lib/resolve-bench.sh
. "$lib"

input="$(cat 2>/dev/null || true)"

# A findable CLI gives the gate its authority here. If none exists, the hook fails
# open.
cmd="$(bench_resolve_wrapper)" || warn_open "bench CLI not found"

# Hand the envelope and the resolved wrapper to the core. It checks stop_hook_active
# and BENCH_SHIFT, runs `<wrapper> gate`, writes the cache, and owns the verdict.
# Exit 0 allows the stop, and exit 2 blocks it, with the BLOCKED message and the gate
# tail already on stderr. The wrapper exits 127 when this platform has no installed
# binary.
rc=0
printf '%s' "$input" | "$cmd" stop-verdict "$cmd" || rc=$?
case "$rc" in
  0 | 2) exit "$rc" ;;                                # allow / block — the core owns it
  *) warn_open "bench core errored (exit $rc)" ;;      # 127 missing binary / 3 panic → fail open, no cache
esac
