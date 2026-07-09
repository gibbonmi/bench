# shellcheck shell=sh
# Canary fixture: the canonical resolver, unchanged, kept alongside the drifted
# shim so the conformance check has a real order to diff against.
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
