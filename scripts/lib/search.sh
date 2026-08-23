# This is a shared search shim for the build and release scripts.
#
# ripgrep is preferred where it is installed, but these scripts run as child
# processes, and an `rg` that exists only as an interactive shell function is not
# on their PATH. Calling `rg` unconditionally makes every such script die with
# "command not found", so this shim selects the tool at call time, with POSIX grep
# as the fallback. The notice fires at most once per script run: a missing optional
# accelerator is an operator hint, never a build failure.
#
# There are two entry points because the dialects differ. Fixed-string matching
# takes identical flags in both tools. ripgrep's default pattern syntax is extended
# regular expressions, which is `grep -E`, not plain `grep`. Calling the wrong one
# silently changes what an alternation matches, so the caller names the choice at
# the call site rather than this shim guessing it.

bench_search_have_rg() { command -v rg >/dev/null 2>&1; }

bench_search_notice() {
  [ -n "${_bench_search_notified:-}" ] && return 0
  _bench_search_notified=1
  printf 'notice: ripgrep (rg) is not installed; using grep instead. Install ripgrep for faster searches: https://github.com/BurntSushi/ripgrep#installation\n' >&2
  return 0
}

# This is fixed-string matching. The caller supplies the remaining flags, the
# pattern, and any files. With no file, both tools read standard input.
bench_search_fixed() {
  if bench_search_have_rg; then
    rg -F "$@"
  else
    bench_search_notice
    grep -F "$@"
  fi
}

# This is extended-regex matching.
bench_search_ere() {
  if bench_search_have_rg; then
    rg "$@"
  else
    bench_search_notice
    grep -E "$@"
  fi
}
