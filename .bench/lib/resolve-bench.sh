# shellcheck shell=sh
# This is the shared bench-wrapper resolver, the one source of the wrapper search
# order for every hook shim (.bench/hooks/*.sh) and shift adapter
# (.bench/adapters/*). It carries only the search: repo `.bench/bin/bench.sh`, then
# kit `bin/bench.sh`, then `bench` on PATH. No policy lives here.
#
# Each shim owns its own fail posture when the search comes up empty. The git and
# agent-line guards refuse or warn-and-allow. `stop` and `session-start` warn and
# fail open. The adapters fail closed. No shim sources this blindly: each
# source-guards it under that posture, because a missing lib that let the shim
# error before its rims run would turn a fail-closed guard fail-open.
#
# This file is POSIX sh; the adapters, which must stay POSIX-clean, source it with
# `.`. It is meant to be sourced, never run directly. Its variables are
# underscore-prefixed, so sourcing it cannot clobber an adapter's locals. It defines
# bench_resolve_wrapper: this prints the wrapper path and returns 0, or returns 1
# when none of the three candidates resolves to an executable.
bench_resolve_wrapper() {
  _bench_root=$(git rev-parse --show-toplevel 2>/dev/null) || _bench_root=
  if [ -n "$_bench_root" ]; then
    for _bench_candidate in "$_bench_root/.bench/bin/bench.sh" "$_bench_root/bin/bench.sh"; do
      if [ -x "$_bench_candidate" ]; then
        printf '%s\n' "$_bench_candidate"
        return 0
      fi
    done
  fi
  command -v bench 2>/dev/null || return 1
}
