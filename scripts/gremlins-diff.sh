#!/bin/sh
# This runs mutation testing over the Go packages changed since a base ref.
#
# usage: scripts/gremlins-diff.sh [base-ref]        (default: HEAD~1)
#
# This is the automated referee for a red-mutation-free experiment run
# (BENCH_RED_MUTATIONS_OPTIONAL=1). After a green landing, it selects the Go
# packages changed since base, asks gremlins to grade every covered mutant in each
# selected package, and reports per-package and overall test efficacy — killed over
# killed-plus-lived — for the retro.
#
# This referee is advisory by default. Set BENCH_GREMLINS_THRESHOLD to a percentage
# to make a lower overall efficacy exit 1. A package gremlins cannot grade at all
# exits 2: a broken referee must be loud, never a silent pass.

set -u

base="${1:-HEAD~1}"

if ! command -v gremlins >/dev/null 2>&1; then
  echo "gremlins-diff: gremlins not found on PATH; install with: go install github.com/go-gremlins/gremlins/cmd/gremlins@latest" >&2
  exit 2
fi
if ! git rev-parse --verify --quiet "$base^{commit}" >/dev/null; then
  echo "gremlins-diff: base ref '$base' does not name a commit" >&2
  exit 2
fi
# Multi-package coverage needs the covdata tool, which a trimmed module toolchain
# omits. Fall back to the local toolchain when it has one.
if ! go tool -n covdata >/dev/null 2>&1 && GOTOOLCHAIN=local go tool -n covdata >/dev/null 2>&1; then
  export GOTOOLCHAIN=local
fi

dirs=$(git diff --name-only "$base" -- '*.go' | awk '
  index($0, "/") { sub("/[^/]*$", ""); print; next }
  { print "." }
' | sort -u)

total_killed=0
total_lived=0
failed=0
graded=0

for dir in $dirs; do
  case "$dir" in
  *testdata*) continue ;;
  esac
  [ -d "$dir" ] || continue
  ls "$dir"/*.go >/dev/null 2>&1 || continue
  out=$(gremlins unleash "./$dir" 2>&1)
  summary=$(printf '%s\n' "$out" | grep '^Killed:' | tail -1)
  if [ -z "$summary" ]; then
    echo "package $dir: gremlins produced no verdict" >&2
    printf '%s\n' "$out" | tail -3 >&2
    failed=1
    continue
  fi
  killed=$(printf '%s\n' "$summary" | awk -F'[:,] *' '{print $2}')
  lived=$(printf '%s\n' "$summary" | awk -F'[:,] *' '{print $4}')
  uncovered=$(printf '%s\n' "$summary" | awk -F'[:,] *' '{print $6}')
  echo "package $dir: killed=$killed lived=$lived not-covered=$uncovered"
  total_killed=$((total_killed + killed))
  total_lived=$((total_lived + lived))
  graded=1
done

if [ "$failed" -ne 0 ]; then
  echo "gremlins-diff: at least one changed package could not be graded" >&2
  exit 2
fi
if [ "$graded" -eq 0 ]; then
  echo "overall: no changed Go packages since $base"
  exit 0
fi
if [ $((total_killed + total_lived)) -eq 0 ]; then
  echo "overall: no testable mutants on the lines changed since $base"
  exit 0
fi

efficacy=$(awk "BEGIN {printf \"%.2f\", 100 * $total_killed / ($total_killed + $total_lived)}")
echo "overall: killed=$total_killed lived=$total_lived efficacy=$efficacy%"

if [ -n "${BENCH_GREMLINS_THRESHOLD:-}" ]; then
  below=$(awk "BEGIN {print ($efficacy < $BENCH_GREMLINS_THRESHOLD) ? 1 : 0}")
  if [ "$below" -eq 1 ]; then
    echo "gremlins-diff: efficacy $efficacy% is below the $BENCH_GREMLINS_THRESHOLD% threshold" >&2
    exit 1
  fi
fi
exit 0
