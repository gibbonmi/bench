# shellcheck shell=sh
# Shared bench-wrapper resolver — the ONE source of the wrapper search order for
# every hook shim (.bench/hooks/*.sh) and shift adapter (.bench/adapters/*). It
# carries only the search: repo `.bench/bin/bench.sh` -> kit `bin/bench.sh` ->
# `bench` on PATH. No policy lives here. Each shim owns its OWN fail posture when
# the search comes up empty — the git/agent-line guards refuse or warn-and-allow,
# `stop`/`session-start` warn and fail open, the adapters fail closed. No shim
# blindly sources this: each source-guards it under that posture, because a missing
# lib that let the shim error before its rims run would turn a fail-closed guard
# fail-open (the slice-4 failure mode).
#
# POSIX sh (sourced with `.` by the adapters, which must stay POSIX-clean). Sourced,
# not run. Variables are underscore-prefixed so sourcing an adapter cannot clobber
# its locals. Defines bench_resolve_wrapper: prints the wrapper path and returns 0,
# or returns 1 when none of the three candidates resolves to an executable.
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
