#!/usr/bin/env bash
# bench — the operational substrate for the Bench workflow.
# Fuses a warm worktree pool (treehouse-lite) with a gated loop (gnhf-lite),
# where the gate is the external oracle: a shift only commits on green.
#
#   bench link                 safely wire the kit into this repo
#   bench init                 scaffold .bench/gate.sh and .bench/learnings.md
#   bench models               discover model ids for line binding
#   bench structure            flag oversized files and crowded directories
#   bench idea "<text>"        park an out-of-scope idea
#   bench roadmap              list parked ideas
#   bench status               print the ambient dashboard
#   bench gate                 run the project gate; exit code is the verdict
#   bench worktree             drop into a warm, isolated worktree subshell
#   bench shift "<objective>"  run the gated loop in a pooled worktree
#
# Config resolution for the gate, in order:
#   1. ./.bench/gate.sh        (executable in the repo — preferred)
#   2. $BENCH_GATE             (a command string)
#   3. auto-detect             (pnpm / npm / pyproject / cargo)
set -euo pipefail

BENCH_HOME="${BENCH_HOME:-$HOME/.bench}"
AGENT="${BENCH_AGENT:-}"                # harness adapter executable; no default — see .bench/adapters/
MAX_ITERS="${BENCH_MAX_ITERS:-12}"

repo_root() { git rev-parse --show-toplevel 2>/dev/null || { echo "not in a git repo" >&2; exit 1; }; }
default_branch() { git symbolic-ref --quiet --short refs/remotes/origin/HEAD 2>/dev/null | sed 's@^origin/@@' || echo main; }
shift_scratch_status() {
  git -C "$1" status --porcelain | grep -vE '^[ MADRCU?!]{2} \.bench-(objective|notes\.md)$' || true
}
shift_dirty_paths() {
  # Sorted newline list of dirty paths, scratch excluded — comm(1) input for the
  # touched-path diff in shift_loop. -z, because plain porcelain C-quotes paths
  # with spaces or specials and the quoted form never matches a pathspec; a path
  # containing a real newline is the one shape this still misreads.
  local entry p
  git -C "$1" status --porcelain -z --no-renames 2>/dev/null | while IFS= read -r -d '' entry; do
    p="${entry:3}"
    case "$p" in .bench-objective|.bench-notes.md) continue ;; esac
    printf '%s\n' "$p"
  done | sort
}
cleanup_shift_scratch() {
  rm -f "$1/.bench-objective" "$1/.bench-notes.md"
}

# ---- gate: the oracle -------------------------------------------------------
gate_record() {
  # Record the verdict for `bench status` (same format the Stop hook writes):
  #   <status> <tree hash> <iso8601>
  # Keyed by gate_tree_hash (bench-status.sh) — the content tested, not the commit
  # sha — so commit-on-green does not stale the verdict that authorized it. The
  # cache lives in the git dir, so it is never tracked or committed.
  local root="$1" rc="$2" gitdir verdict=green
  gitdir="$(git -C "$root" rev-parse --absolute-git-dir 2>/dev/null)" || return 0
  [[ "$rc" -eq 0 ]] || verdict=red
  printf '%s %s %s\n' "$verdict" "$(gate_tree_hash "$root")" \
    "$(date -u +%Y-%m-%dT%H:%M:%SZ)" > "$gitdir/bench-last-gate"
}

run_gate() {
  local root rc=0; root="$(repo_root)"
  # The gate is run from the working tree by design — an agent can edit the file it is
  # graded by. The canary tripwire, not this call site, is what keeps that safe; see
  # docs/adr/0001-working-tree-gate-tripwire.md.
  # Run, do not exec: the shift loop calls `if run_gate` inline, so an exec here would
  # replace the bench process at the first gate check and the loop would never iterate.
  # The `|| rc=$?` form matters under set -e: a bare failing subshell would exit the
  # script before the verdict could be recorded.
  if [[ -x "$root/.bench/gate.sh" ]]; then
    ( cd "$root" && "$root/.bench/gate.sh" ) || rc=$?
  elif [[ -n "${BENCH_GATE:-}" ]]; then
    ( cd "$root" && bash -c "$BENCH_GATE" ) || rc=$?
  # auto-detect — best-effort defaults; a project should ship .bench/gate.sh instead
  elif [[ -f "$root/pnpm-lock.yaml" ]]; then
    ( cd "$root" && pnpm -s typecheck && pnpm -s test && pnpm -s lint ) || rc=$?
  elif [[ -f "$root/package.json" ]]; then
    ( cd "$root" && npm run -s typecheck && npm test --silent && npm run -s lint ) || rc=$?
  elif [[ -f "$root/pyproject.toml" ]]; then
    ( cd "$root" && mypy . && pytest -q && ruff check . ) || rc=$?
  elif [[ -f "$root/Cargo.toml" ]]; then
    ( cd "$root" && cargo test --quiet && cargo clippy -q -- -D warnings ) || rc=$?
  else
    echo "no gate found: add an executable .bench/gate.sh or set BENCH_GATE" >&2
    return 3
  fi
  # A run that happened is a verdict either way; the no-gate path records nothing.
  gate_record "$root" "$rc"
  return "$rc"
}

# ---- shift: the gated loop --------------------------------------------------
# The adapter is the harness seam: shift passes the generated prompt as the
# adapter's single positional argument, so harness-specific flags live in the
# adapter (e.g. .bench/adapters/claude), never in the loop. Arguments belong in
# a wrapper script — a multi-word BENCH_AGENT resolves as one executable name.
require_adapter() {
  if [[ -z "$AGENT" ]]; then
    echo "no harness adapter configured: set BENCH_AGENT to an adapter executable (references in .bench/adapters/)" >&2
    exit 1
  fi
  # type -P, not command -v: the adapter must resolve to an executable file, so a
  # value colliding with a shell keyword or builtin is rejected here instead of
  # exec-failing (or silently no-opping) inside the loop.
  if [[ -z "$(type -P "$AGENT")" ]]; then
    echo "adapter not executable: BENCH_AGENT='$AGENT' is neither an executable file nor a command on PATH" >&2
    exit 1
  fi
}

shift_cleanup() {
  local wt="${1:-}"
  [[ -n "$wt" ]] || return 0
  cleanup_shift_scratch "$wt"
  worktree_release "$wt"
}

shift_stage_touched() {
  # Stage exactly what the agent touched: dirty after it ran ($3) minus dirty
  # before it ran ($2). Snapshotted before the gate runs, so gate byproducts
  # (unignored build artifacts, caches) never ride into an iteration commit —
  # the same sweep class as a blanket `git add -A`, previously contained only
  # by worktree isolation. :(literal) keeps glob characters in a path from
  # being read as a pathspec pattern.
  local root="$1" pre="$2" post="$3" p
  comm -13 <(printf '%s\n' "$pre") <(printf '%s\n' "$post") | while IFS= read -r p; do
    [[ -n "$p" ]] || continue
    git -C "$root" add -A -- ":(literal)$p" || true
  done
}

shift_loop() {
  local objective="$1" main_root root wt branch base i started committed=0 pre post
  require_adapter
  main_root="$(repo_root)"
  [[ -z "$(git -C "$main_root" status --porcelain)" ]] || { echo "working tree not clean; commit or stash first" >&2; exit 1; }
  base="$(git -C "$main_root" rev-parse HEAD)"
  wt="$(worktree_acquire "$main_root" "$base" hard)"
  root="$wt"
  branch="bench/shift-$(date +%Y%m%d-%H%M%S)"
  git -C "$root" switch -q -c "$branch"
  # The true review base for this branch. `bench diff` resolves it from here;
  # worktrees share repo config, so the key is visible wherever review runs.
  git -C "$root" config "branch.$branch.benchBase" "$base"
  printf '%s\n' "$objective" > "$root/.bench-objective"
  : > "$root/.bench-notes.md"   # carried between iterations so each learns from the last
  BENCH_SHIFT_ROOT="$root"
  trap 'shift_cleanup "$BENCH_SHIFT_ROOT"' EXIT
  trap 'shift_cleanup "$BENCH_SHIFT_ROOT"; exit 130' INT TERM
  echo "▶ shift on $branch — objective: $objective"
  echo "  worktree: $root"
  echo "  cap: $MAX_ITERS iterations. Ctrl-C to pull the line."
  started=$(date +%s)
  for ((i=1; i<=MAX_ITERS; i++)); do
    echo "── iteration $i/$MAX_ITERS ──"
    # one bounded iteration: agent makes one small change toward the objective.
    # BENCH_SHIFT=1 arms the Stop hook so the agent cannot declare done on red.
    pre="$(shift_dirty_paths "$root")"
    ( cd "$root" && BENCH_SHIFT=1 "$AGENT" "$(iteration_prompt "$objective")" ) || true
    post="$(shift_dirty_paths "$root")"
    if ( cd "$root" && run_gate ); then
      shift_stage_touched "$root" "$pre" "$post"
      if git -C "$root" diff --cached --quiet; then
        echo "  gate green, no change this iteration — objective likely met."; break
      fi
      git -C "$root" commit -q -m "shift: iteration $i — $objective"
      committed=$((committed+1))
      echo "  ✓ green — committed iteration $i"
      if ( cd "$root" && objective_met "$objective" ); then echo "  objective met."; break; fi
    else
      echo "  ✗ red gate — rolling back iteration $i, retrying"
      git -C "$root" reset -q --hard; git -C "$root" clean -qfdx -e .bench-objective -e .bench-notes.md
    fi
  done
  # Implementation loop is done. Only NOW pay down structural debt, if the work
  # pushed past the budget. Refactor at green — never mid-implementation: splitting
  # before the feature's shape has settled produces premature, bad module boundaries.
  # This is "only trigger a refactor once it's over the threshold", at the loop edge.
  if ! ( cd "$root" && structure_touched_since "$base" >/dev/null 2>&1 ); then
    echo "▶ structure over budget — refactor phase (split at green, not before)"
    local rcap="${BENCH_REFACTOR_ITERS:-4}" r attempted=0 flagged
    for ((r=1; r<=rcap; r++)); do
      ( cd "$root" && structure_touched_since "$base" >/dev/null 2>&1 ) && break
      attempted="$r"
      echo "── refactor $r/$rcap ──"
      # Scope the prompt to the files this shift flagged — never repo-wide debt.
      flagged="$( (cd "$root" && structure_touched_since "$base") 2>&1 || true)"
      pre="$(shift_dirty_paths "$root")"
      ( cd "$root" && BENCH_SHIFT=1 "$AGENT" "$(refactor_prompt "$flagged")" ) || true
      post="$(shift_dirty_paths "$root")"
      if ( cd "$root" && run_gate ); then
        shift_stage_touched "$root" "$pre" "$post"
        if git -C "$root" diff --cached --quiet; then
          echo "  gate green, refactor $r made no staged change - stopping refactor phase"
          break
        fi
        git -C "$root" commit -q -m "refactor: reduce structural debt"
        echo "  ✓ tests green - refactor $r committed"
      else
        echo "  ✗ refactor broke the gate — rolling back"
        git -C "$root" reset -q --hard; git -C "$root" clean -qfdx -e .bench-objective -e .bench-notes.md
      fi
    done
    if ( cd "$root" && structure_touched_since "$base" >/dev/null 2>&1 ); then echo "  structure back under budget."
    else echo "  ⚠ still over budget after ${attempted:-$rcap} refactor pass(es) - review manually, or run a Bench deep pass with bench structure and craft-seams."; fi
  fi
  trap - EXIT INT TERM
  cleanup_shift_scratch "$root"
  git -C "$root" switch -q --detach >/dev/null 2>&1 || true
  worktree_release "$root"
  echo "■ shift done: $branch, $committed committed iteration(s), $(( ($(date +%s)-started)/60 ))m elapsed"
  echo "  review: git -C $main_root log --oneline ${base}..$branch"
  echo "  the merge is yours."
}

refactor_prompt() {
  cat <<EOF
The implementation is complete and tests are green, but the structure budget is
exceeded. These are the flagged files and directories this shift touched — fix only
these, nothing else:

$1

Fix them by splitting along responsibility, using the deletion test from the craft-seams
skill: lift a cluster out only if extracting it *concentrates* complexity behind a real
interface rather than just moving it. Never fragment a cohesive file to beat the line
count — if a file is genuinely one deep module, leave it and say so. Group a crowded
directory into a package with a clear entry point. Keep every test green; change
structure, not behavior. Make one split, then stop — the loop re-checks and continues.
EOF
}

iteration_prompt() {
  cat <<EOF
You are one iteration of a Bench shift. Objective: $1
First read .bench-notes.md for what prior iterations learned, did, and left
unfinished. Then make ONE small, self-contained change toward the objective, at
the pre-agreed seams. Read the spec under specs/ and projects/ if present. Do not
try to finish everything; advance it by one honest step. Do not weaken or skip any
gate check. Before you stop, append 2–4 lines to .bench-notes.md: what you changed,
what you learned, and the next step you'd take. Then stop — the gate, not you,
decides if it counts.
EOF
}

# Override per project: an executable .bench/done.sh that exits 0 when the objective
# is complete (e.g. all spec stories covered). Absent => run to the iteration cap.
objective_met() {
  local root; root="$(repo_root)"
  [[ -x "$root/.bench/done.sh" ]] && "$root/.bench/done.sh" "$1"
}

# Resolve where the canonical kit lives (parent of this script's bin/), following
# symlinks without relying on GNU-only `readlink -f`.
resolve_script_path() {
  local source="${BASH_SOURCE[0]:-$0}" dir target
  while [[ -L "$source" ]]; do
    dir="$(cd -P "$(dirname "$source")" >/dev/null 2>&1 && pwd)"
    target="$(readlink "$source")"
    [[ "$target" == /* ]] && source="$target" || source="$dir/$target"
  done
  dir="$(cd -P "$(dirname "$source")" >/dev/null 2>&1 && pwd)"
  printf '%s/%s\n' "$dir" "$(basename "$source")"
}

kit_dir() {
  local script
  script="$(resolve_script_path)"
  cd "$(dirname "$script")/.." && pwd
}

# ---- strangler router: send a ported subcommand to the Go binary ------------
# The one seam that grows across the port: later slices add subcommand names to the
# dispatch, never a second resolver. platform_pkg maps this host to its
# @benchkit/<os>-<arch> package (npm os/cpu spelling); an off-matrix host returns
# non-zero so the caller can emit the "unsupported platform" error instead of naming
# a package that does not exist.
platform_pkg() {
  local os arch
  case "$(uname -s)" in
    Darwin) os=darwin ;;
    Linux)  os=linux ;;
    *) return 1 ;;
  esac
  case "$(uname -m)" in
    arm64|aarch64) arch=arm64 ;;
    x86_64|amd64)  arch=x64 ;;
    *) return 1 ;;
  esac
  printf '@benchkit/%s-%s' "$os" "$arch"
}

# route_binary <subcommand> [args...] — resolve and exec the Go binary, passing the
# whole argv through. Resolution order: (1) repo-local dev build (kit checkout),
# (2) the platform package bundled under the wrapper's node_modules, (3) the hoisted
# sibling npm produces for global installs. First executable, non-empty match wins;
# a present-but-empty or non-executable file is treated as missing (never exec'd), so
# a torn build falls through to the named-package error rather than exec-failing.
route_binary() {
  local kit pkg c
  kit="$(kit_dir)"
  if ! pkg="$(platform_pkg)"; then
    echo "bench: unsupported platform: $(uname -s | tr '[:upper:]' '[:lower:]')/$(uname -m)" >&2
    exit 2
  fi
  for c in "$kit/dist/bench" "$kit/node_modules/$pkg/bin/bench" "$kit/../$pkg/bin/bench"; do
    [[ -x "$c" && -s "$c" ]] && exec "$c" "$@"
  done
  echo "bench: no binary for this platform — install $pkg (npm install $pkg)" >&2
  exit 127
}

BENCH_BIN_DIR="$(dirname "$(resolve_script_path)")"
# shellcheck source=/dev/null
. "$BENCH_BIN_DIR/bench-link.sh"
# shellcheck source=/dev/null
. "$BENCH_BIN_DIR/bench-status.sh"
# shellcheck source=/dev/null
. "$BENCH_BIN_DIR/bench-worktree.sh"
# shellcheck source=/dev/null
. "$BENCH_BIN_DIR/bench-query.sh"
# shellcheck source=/dev/null
. "$BENCH_BIN_DIR/bench-diff.sh"
# shellcheck source=/dev/null
. "$BENCH_BIN_DIR/bench-coverage.sh"
# shellcheck source=/dev/null
. "$BENCH_BIN_DIR/bench-init.sh"
# shellcheck source=/dev/null
. "$BENCH_BIN_DIR/bench-doctor.sh"

case "${1:-help}" in
  version)  route_binary "$@" ;;
  gate)     run_gate ;;
  doctor)   shift; doctor "$@" ;;
  worktree) worktree ;;
  shift)    shift; shift_loop "${*:-improve the codebase}" ;;
  link)     shift; link "${1:-copy}" ;;
  models)   models ;;
  structure) structure ;;
  init)     init ;;
  idea)     shift; idea "$@" ;;
  roadmap)  roadmap ;;
  status)   status ;;
  learnings) shift; learnings "$@" ;;
  maps)     shift; maps "$@" ;;
  guards)   shift; guards "$@" ;;
  diff)     shift; bench_diff "$@" ;;
  coverage) shift; coverage "$@" ;;
  *) cat <<EOF
bench — Pocock pipeline meets Kun Chen substrate, gated by your invariants.
  bench link [copy|symlink]  safely wire the kit into this repo for every harness
  bench init                 scaffold .bench/gate.sh in the current repo
  bench models               discover models available in this harness (for the lines)
  bench structure            flag oversized files + crowded dirs (wire into the gate)
  bench idea "<text>"        park an out-of-scope idea in ROADMAP.md (commit to nothing)
  bench roadmap              list parked ideas
  bench status               ambient dashboard: what needs attention + the next action
  bench learnings            open journal entries as a TOON table (date, title)
  bench maps                 unresolved decision-map tickets as TOON (map, ticket, type, state)
  bench guards               every guard's deny surface as TOON (guard, boundary, denies)
  bench diff                 review base (recorded or merge-base) + changed files as TOON
  bench coverage <spec>      acceptance-coverage state and rows as TOON (--check to validate)
  bench doctor [--fix]       report (and repair) the PATH shim under a node version manager
  bench gate                 run the project gate (the oracle)
  bench worktree             warm, isolated worktree subshell
  bench shift "<objective>"  gated loop in a pooled worktree; commit on green
  bench version              print the installed benchkit version (os/arch)
EOF
  ;;
esac
