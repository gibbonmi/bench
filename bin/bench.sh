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
AGENT="${BENCH_AGENT:-claude}"          # headless agent command; override per harness
MAX_ITERS="${BENCH_MAX_ITERS:-12}"
MAX_TOKENS="${BENCH_MAX_TOKENS:-4000000}"

repo_root() { git rev-parse --show-toplevel 2>/dev/null || { echo "not in a git repo" >&2; exit 1; }; }
default_branch() { git symbolic-ref --quiet --short refs/remotes/origin/HEAD 2>/dev/null | sed 's@^origin/@@' || echo main; }
shift_scratch_status() {
  git -C "$1" status --porcelain | grep -vE '^[ MADRCU?!]{2} \.bench-(objective|notes\.md)$' || true
}
cleanup_shift_scratch() {
  rm -f "$1/.bench-objective" "$1/.bench-notes.md"
}

worktree_pool() {
  local root key
  root="$1"
  key="$(basename "$root")-$(echo "$root" | cksum | cut -d' ' -f1)"
  printf '%s\n' "$BENCH_HOME/worktrees/$key"
}

worktree_lease_file() {
  git -C "$1" rev-parse --git-path bench-lease
}

# ---- gate: the oracle -------------------------------------------------------
run_gate() {
  local root; root="$(repo_root)"
  # Run, do not exec: the shift loop calls `if run_gate` inline, so an exec here would
  # replace the bench process at the first gate check and the loop would never iterate.
  if [[ -x "$root/.bench/gate.sh" ]]; then "$root/.bench/gate.sh"; return $?; fi
  if [[ -n "${BENCH_GATE:-}" ]]; then bash -c "$BENCH_GATE"; return $?; fi
  # auto-detect — best-effort defaults; a project should ship .bench/gate.sh instead
  if [[ -f "$root/pnpm-lock.yaml" ]]; then
    ( cd "$root" && pnpm -s typecheck && pnpm -s test && pnpm -s lint ); return $?
  elif [[ -f "$root/package.json" ]]; then
    ( cd "$root" && npm run -s typecheck && npm test --silent && npm run -s lint ); return $?
  elif [[ -f "$root/pyproject.toml" ]]; then
    ( cd "$root" && mypy . && pytest -q && ruff check . ); return $?
  elif [[ -f "$root/Cargo.toml" ]]; then
    ( cd "$root" && cargo test --quiet && cargo clippy -q -- -D warnings ); return $?
  fi
  echo "no gate found: add an executable .bench/gate.sh or set BENCH_GATE" >&2
  return 3
}

# ---- worktree: warm, isolated, reusable -------------------------------------
worktree_acquire() {
  local root reset_ref reset_mode pool wt
  root="$1"
  reset_ref="${2:-}"
  reset_mode="${3:-hard}"
  pool="$(worktree_pool "$root")"; mkdir -p "$pool"
  git -C "$root" fetch -q origin 2>/dev/null || true
  for d in "$pool"/*/; do
    [[ -d "$d/.git" || -f "$d/.git" ]] || continue
    [[ -f "$(worktree_lease_file "$d")" ]] && continue
    if [[ -z "$(git -C "$d" status --porcelain 2>/dev/null)" ]]; then wt="$d"; break; fi
  done
  if [[ -z "${wt:-}" ]]; then
    wt="$pool/$(date +%s)-$$"; git -C "$root" worktree add -q --detach "$wt" "origin/$(default_branch)" 2>/dev/null \
      || git -C "$root" worktree add -q --detach "$wt"
  fi
  git -C "$wt" switch -q --detach >/dev/null 2>&1 || true
  if [[ -n "$reset_ref" ]]; then
    if ! git -C "$wt" reset -q --hard "$reset_ref"; then
      [[ "$reset_mode" == soft ]] || return 1
      git -C "$wt" reset -q --hard >/dev/null 2>&1 || true
    fi
  else
    git -C "$wt" reset -q --hard >/dev/null 2>&1 || true
  fi
  git -C "$wt" clean -qfd
  : > "$(worktree_lease_file "$wt")"
  printf '%s\n' "${wt%/}"
}

worktree_release() {
  local wt="$1"
  [[ -n "$wt" ]] || return 0
  rm -f "$(worktree_lease_file "$wt")"
  git -C "$wt" switch -q --detach >/dev/null 2>&1 || true
  git -C "$wt" reset -q --hard >/dev/null 2>&1 || true
  git -C "$wt" clean -qfd >/dev/null 2>&1 || true
  rm -f "$(worktree_lease_file "$wt")"
}

worktree() {
  local root target wt
  root="$(repo_root)"
  target="origin/$(default_branch)"
  wt="$(worktree_acquire "$root" "$target" soft)"
  echo "🪵 worktree: $wt  (exit to release)" >&2
  ( cd "$wt" && "${SHELL:-bash}" ) || true
  worktree_release "$wt"
  echo "🪵 released" >&2
}

# ---- shift: the gated loop --------------------------------------------------
shift_cleanup() {
  local wt="${1:-}"
  [[ -n "$wt" ]] || return 0
  cleanup_shift_scratch "$wt"
  worktree_release "$wt"
}

shift_loop() {
  local objective="$1" main_root root wt branch base i started tokens=0 committed=0
  main_root="$(repo_root)"
  [[ -z "$(git -C "$main_root" status --porcelain)" ]] || { echo "working tree not clean; commit or stash first" >&2; exit 1; }
  base="$(git -C "$main_root" rev-parse HEAD)"
  wt="$(worktree_acquire "$main_root" "$base" hard)"
  root="$wt"
  branch="bench/shift-$(date +%Y%m%d-%H%M%S)"
  git -C "$root" switch -q -c "$branch"
  printf '%s\n' "$objective" > "$root/.bench-objective"
  : > "$root/.bench-notes.md"   # carried between iterations so each learns from the last
  BENCH_SHIFT_ROOT="$root"
  trap 'shift_cleanup "$BENCH_SHIFT_ROOT"' EXIT
  trap 'shift_cleanup "$BENCH_SHIFT_ROOT"; exit 130' INT TERM
  echo "▶ shift on $branch — objective: $objective"
  echo "  worktree: $root"
  echo "  caps: $MAX_ITERS iterations, ~$MAX_TOKENS tokens. Ctrl-C to pull the line."
  started=$(date +%s)
  for ((i=1; i<=MAX_ITERS; i++)); do
    echo "── iteration $i/$MAX_ITERS ──"
    # one bounded iteration: agent makes one small change toward the objective.
    # BENCH_SHIFT=1 arms the Stop hook so the agent cannot declare done on red.
    ( cd "$root" && BENCH_SHIFT=1 "$AGENT" -p "$(iteration_prompt "$objective")" ) || true
    if ( cd "$root" && run_gate ); then
      git -C "$root" add -A -- ':!.bench-objective' ':!.bench-notes.md'
      if git -C "$root" diff --cached --quiet; then
        echo "  gate green, no change this iteration — objective likely met."; break
      fi
      git -C "$root" commit -q -m "shift: iteration $i — $objective"
      committed=$((committed+1))
      echo "  ✓ green — committed iteration $i"
      if ( cd "$root" && objective_met "$objective" ); then echo "  objective met."; break; fi
    else
      echo "  ✗ red gate — rolling back iteration $i, retrying"
      git -C "$root" reset -q --hard; git -C "$root" clean -qfd
    fi
  done
  # Implementation loop is done. Only NOW pay down structural debt, if the work
  # pushed past the budget. Refactor at green — never mid-implementation: splitting
  # before the feature's shape has settled produces premature, bad module boundaries.
  # This is "only trigger a refactor once it's over the threshold", at the loop edge.
  if ! ( cd "$root" && structure_touched_since "$base" >/dev/null 2>&1 ); then
    echo "▶ structure over budget — refactor phase (split at green, not before)"
    local rcap="${BENCH_REFACTOR_ITERS:-4}" r attempted=0
    for ((r=1; r<=rcap; r++)); do
      ( cd "$root" && structure_touched_since "$base" >/dev/null 2>&1 ) && break
      attempted="$r"
      echo "── refactor $r/$rcap ──"
      ( cd "$root" && BENCH_SHIFT=1 "$AGENT" -p "$(refactor_prompt)" ) || true
      if ( cd "$root" && run_gate ); then
        git -C "$root" add -A -- ':!.bench-objective' ':!.bench-notes.md'
        if git -C "$root" diff --cached --quiet; then
          echo "  gate green, refactor $r made no staged change - stopping refactor phase"
          break
        fi
        git -C "$root" commit -q -m "refactor: reduce structural debt"
        echo "  ✓ tests green - refactor $r committed"
      else
        echo "  ✗ refactor broke the gate — rolling back"
        git -C "$root" reset -q --hard; git -C "$root" clean -qfd
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
  cat <<'EOF'
The implementation is complete and tests are green, but the structure budget is
exceeded. Run `bench structure` to see the flagged files and directories. Fix them
by splitting along responsibility, using the deletion test from the craft-seams skill: lift
a cluster out only if extracting it *concentrates* complexity behind a real interface
rather than just moving it. Never fragment a cohesive file to beat the line count — if
a file is genuinely one deep module, leave it and say so. Group a crowded directory
into a package with a clear entry point. Keep every test green; change structure, not
behavior. Make one split, then stop — the loop re-checks and continues.
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

init() {
  local root; root="$(repo_root)"; mkdir -p "$root/.bench"
  if [[ ! -e "$root/.bench/gate.sh" ]]; then
    cat > "$root/.bench/gate.sh" <<'EOF'
#!/usr/bin/env bash
# The external oracle for this repo — correctness only. Exit 0 = done is allowed.
set -euo pipefail
# Stack checks — uncomment what fits:
#   mypy . && pytest -q && ruff check .
#   pnpm -s typecheck && pnpm -s test && pnpm -s lint
#
# Structural debt is NOT checked here. `bench shift` runs `bench structure` once the
# implementation loop finishes and triggers a refactor pass only if a file or dir is
# over budget — so splits happen at green, not mid-iteration. Uncomment the next line
# only if you also want structure hard-blocked at the PR boundary (every commit):
#   bench structure
echo "edit .bench/gate.sh to run this project's checks" >&2; exit 3
EOF
    chmod +x "$root/.bench/gate.sh"
    echo "scaffolded .bench/gate.sh — edit it to run your real checks"
  fi
  if [[ ! -e "$root/.bench/learnings.md" ]]; then
    cat > "$root/.bench/learnings.md" <<'EOF'
# Learnings — usage journal

Append one entry when you deviate from the workflow, make a judgment call you're
unsure about, or catch a should-have-asked in hindsight. You capture; the reviewer
decides. `/bench-integrate-learnings` reviews the open entries, promotes the generalizable ones
into the kit with sign-off, and marks them resolved. Never rewrite a kit rule
yourself — that is the whole point of capturing here instead.

Format per entry:

## <date> — <short title>  [open|resolved]
- **What happened:** …
- **Right behavior:** …
- **Proposed rule change:** … (or "none")

An entry only becomes `[resolved]` via /bench-integrate-learnings.

<!-- entries below -->
EOF
    echo "scaffolded .bench/learnings.md — the self-learning journal"
  fi
  echo "see projects/<name>.md in the Bench kit for the profile template"
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

BENCH_BIN_DIR="$(dirname "$(resolve_script_path)")"
# shellcheck source=/dev/null
. "$BENCH_BIN_DIR/bench-link.sh"
# shellcheck source=/dev/null
. "$BENCH_BIN_DIR/bench-status.sh"

case "${1:-help}" in
  gate)     run_gate ;;
  worktree) worktree ;;
  shift)    shift; shift_loop "${*:-improve the codebase}" ;;
  link)     shift; link "${1:-copy}" ;;
  models)   models ;;
  structure) structure ;;
  init)     init ;;
  idea)     shift; idea "$@" ;;
  roadmap)  roadmap ;;
  status)   status ;;
  *) cat <<EOF
bench — Pocock pipeline meets Kun Chen substrate, gated by your invariants.
  bench link [copy|symlink]  safely wire the kit into this repo for every harness
  bench init                 scaffold .bench/gate.sh in the current repo
  bench models               discover models available in this harness (for the lines)
  bench structure            flag oversized files + crowded dirs (wire into the gate)
  bench idea "<text>"        park an out-of-scope idea in ROADMAP.md (commit to nothing)
  bench roadmap              list parked ideas
  bench status               ambient dashboard: what needs attention + the next action
  bench gate                 run the project gate (the oracle)
  bench worktree             warm, isolated worktree subshell
  bench shift "<objective>"  gated loop in a pooled worktree; commit on green
EOF
  ;;
esac
