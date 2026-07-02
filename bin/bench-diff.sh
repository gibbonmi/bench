# The review-base parser for bin/bench.sh: `bench diff`, the single source of
# review-base truth for /bench-review-implementation (spec second-wave-parsers).
# Sourced by bin/bench.sh after bench-query.sh — it composes the shared TOON
# emitter and error helpers defined there, and default_branch() from bench.sh.

# ---- review-base + changed files --------------------------------------------
# `bench diff` is the single source of review-base truth (decision map
# second-wave-parsers #2/#3): the recorded `branch.<name>.benchBase` key wins
# when it names a reachable commit; otherwise merge-base with the default
# branch. The preamble names which method resolved so a review agent can see a
# fallback happen. Named bench_diff, not diff — the CLI must never shadow the
# system diff inside this script.
bench_diff() {
  case "${1-}" in
    "") ;;
    -h | --help) echo "usage: bench diff"; return 0 ;;
    *) axi_usage "bench diff" "$1"; return 2 ;;
  esac
  git rev-parse --show-toplevel >/dev/null 2>&1 || {
    axi_error "not in a git repository" "run inside a Bench-linked repo"
    return 1
  }
  local branch key base="" method="" def
  branch="$(git symbolic-ref --quiet --short HEAD 2>/dev/null || true)"
  if [[ -n "$branch" ]]; then
    key="$(git config "branch.$branch.benchBase" 2>/dev/null || true)"
    if [[ -n "$key" ]]; then
      if ! git cat-file -e "$key^{commit}" 2>/dev/null; then
        method="merge-base (recorded sha unreachable)"
      elif ! git merge-base --is-ancestor "$key" HEAD 2>/dev/null; then
        # Reachable but divergent: three-dot against it would silently diff from
        # the merge-base while the preamble claims the recorded sha. Fall back loudly.
        method="merge-base (recorded sha not an ancestor)"
      else
        base="$key" method="recorded"
      fi
    fi
  fi
  if [[ -z "$base" ]]; then
    def="$(default_branch)"
    base="$(git merge-base "$def" HEAD 2>/dev/null)" || {
      axi_error "cannot resolve a review base" "no merge-base with '$def'; record one with: git config branch.<name>.benchBase <sha>"
      return 1
    }
    [[ -n "$method" ]] || method="merge-base"
  fi
  printf 'branch: %s\n' "${branch:-(detached)}"
  printf 'base: %s\n' "$base"
  printf 'method: %s\n' "$method"
  # --no-renames keeps every row a flat status/path pair (no two-path R rows);
  # three-dot matches the review phase's convention and equals two-dot here
  # because the resolved base is always an ancestor of HEAD. -z, because plain
  # output C-quotes non-ASCII or quote-bearing paths and the TOON emitter would
  # quote them a second time (the shift_dirty_paths posture; a path containing
  # a real newline is the one shape this still misreads).
  git diff --name-status --no-renames -z "$base...HEAD" \
    | { while IFS= read -r -d '' st && IFS= read -r -d '' p; do printf '%s\t%s\n' "$st" "$p"; done; } \
    | toon_table files status,path
}
