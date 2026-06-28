#!/usr/bin/env bash
# bench — the operational substrate for the Bench workflow.
# Fuses a warm worktree pool (treehouse-lite) with a gated loop (gnhf-lite),
# where the gate is the external oracle: a shift only commits on green.
#
#   bench gate                 run the project gate; exit code is the verdict
#   bench worktree             drop into a warm, isolated worktree subshell
#   bench shift "<objective>"  run the gated loop toward an objective
#   bench init                 scaffold a project profile + default gate
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

# ---- gate: the oracle -------------------------------------------------------
run_gate() {
  local root; root="$(repo_root)"
  if [[ -x "$root/.bench/gate.sh" ]]; then exec "$root/.bench/gate.sh"; fi
  if [[ -x "$root/.bench/gate" ]]; then exec "$root/.bench/gate"; fi  # legacy: pre-.sh name
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
worktree() {
  local root key pool wt
  root="$(repo_root)"
  key="$(basename "$root")-$(echo "$root" | cksum | cut -d' ' -f1)"
  pool="$BENCH_HOME/worktrees/$key"; mkdir -p "$pool"
  git -C "$root" fetch -q origin 2>/dev/null || true
  # find a clean, unleased worktree in the pool
  for d in "$pool"/*/; do
    [[ -d "$d/.git" || -f "$d/.git" ]] || continue
    [[ -f "$d/.lease" ]] && continue
    if [[ -z "$(git -C "$d" status --porcelain 2>/dev/null)" ]]; then wt="$d"; break; fi
  done
  if [[ -z "${wt:-}" ]]; then
    wt="$pool/$(date +%s)-$$"; git -C "$root" worktree add -q --detach "$wt" "origin/$(default_branch)" 2>/dev/null \
      || git -C "$root" worktree add -q --detach "$wt"
  fi
  git -C "$wt" reset -q --hard "origin/$(default_branch)" 2>/dev/null || true
  git -C "$wt" clean -qfd
  : > "$wt/.lease"
  echo "🪵 worktree: $wt  (exit to release)" >&2
  ( cd "$wt" && "${SHELL:-bash}" ) || true
  rm -f "$wt/.lease"
  git -C "$wt" reset -q --hard >/dev/null 2>&1 || true
  git -C "$wt" clean -qfd >/dev/null 2>&1 || true
  echo "🪵 released" >&2
}

# ---- shift: the gated loop --------------------------------------------------
shift_loop() {
  local objective="$1" root branch i started tokens=0
  root="$(repo_root)"
  [[ -z "$(git -C "$root" status --porcelain)" ]] || { echo "working tree not clean; commit or stash first" >&2; exit 1; }
  branch="bench/shift-$(date +%Y%m%d-%H%M%S)"
  git -C "$root" switch -q -c "$branch"
  printf '%s\n' "$objective" > "$root/.bench-objective"
  : > "$root/.bench-notes.md"   # carried between iterations so each learns from the last
  echo "▶ shift on $branch — objective: $objective"
  echo "  caps: $MAX_ITERS iterations, ~$MAX_TOKENS tokens. Ctrl-C to pull the line."
  started=$(date +%s)
  for ((i=1; i<=MAX_ITERS; i++)); do
    echo "── iteration $i/$MAX_ITERS ──"
    # one bounded iteration: agent makes one small change toward the objective.
    # BENCH_SHIFT=1 arms the Stop hook so the agent cannot declare done on red.
    BENCH_SHIFT=1 "$AGENT" -p "$(iteration_prompt "$objective")" || true
    if run_gate; then
      git -C "$root" add -A -- ':!.bench-objective' ':!.bench-notes.md'
      if git -C "$root" diff --cached --quiet; then
        echo "  gate green, no change this iteration — objective likely met."; break
      fi
      git -C "$root" commit -q -m "shift: iteration $i — $objective"
      echo "  ✓ green — committed iteration $i"
      if objective_met "$objective"; then echo "  objective met."; break; fi
    else
      echo "  ✗ red gate — rolling back iteration $i, retrying"
      git -C "$root" reset -q --hard; git -C "$root" clean -qfd
    fi
  done
  rm -f "$root/.bench-objective" "$root/.bench-notes.md"
  echo "■ shift done: $branch, $((i-1)) committed iteration(s), $(( ($(date +%s)-started)/60 ))m elapsed"
  echo "  review: git -C $root log --oneline origin/$(default_branch)..$branch"
  echo "  the merge is yours."
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
  if   [[ -x "$root/.bench/done.sh" ]]; then "$root/.bench/done.sh" "$1"
  elif [[ -x "$root/.bench/done"    ]]; then "$root/.bench/done" "$1"  # legacy: pre-.sh name
  else return 1; fi
}

init() {
  local root; root="$(repo_root)"; mkdir -p "$root/.bench"
  if [[ ! -e "$root/.bench/gate.sh" && ! -e "$root/.bench/gate" ]]; then
    cat > "$root/.bench/gate.sh" <<'EOF'
#!/usr/bin/env bash
# The external oracle for this repo. Exit 0 = done is allowed. Edit to taste.
set -euo pipefail
# Examples — uncomment what fits:
#   mypy . && pytest -q && ruff check .
#   pnpm -s typecheck && pnpm -s test && pnpm -s lint
echo "edit .bench/gate.sh to run this project's checks" >&2; exit 3
EOF
    chmod +x "$root/.bench/gate.sh"
    echo "scaffolded .bench/gate.sh — edit it to run your real checks"
  fi
  echo "see projects/<name>.md in the Bench kit for the profile template"
}

# Resolve where the canonical kit lives (parent of this script's bin/).
kit_dir() { cd "$(dirname "$(readlink -f "$0" 2>/dev/null || echo "$0")")/.." && pwd; }

# Git-native push guard — harness- and human-agnostic. Protects the default
# branch (you own the merge); rejects force-pushes to it. Works regardless of
# which agent (or none) runs git push.
install_git_hook() {
  local root="$1" def; def="$(default_branch)"
  mkdir -p "$root/.git/hooks"
  cat > "$root/.git/hooks/pre-push" <<EOF
#!/usr/bin/env bash
# Installed by 'bench link'. The merge is the human's; agents don't push $def.
while read -r local_ref local_oid remote_ref remote_oid; do
  if [[ "\$remote_ref" == "refs/heads/$def" ]]; then
    echo "blocked: direct push to $def. Open a PR or merge it yourself." >&2
    exit 1
  fi
done
exit 0
EOF
  chmod +x "$root/.git/hooks/pre-push"
}

# Wire the kit into the current repo for EVERY harness at once. Idempotent.
# After this, Claude Code and any AGENTS.md harness work in the same repo with no
# per-switch reconfiguration — switching harnesses is just running a different agent.
link() {
  local root kit mode="${1:-symlink}"; root="$(repo_root)"; kit="${BENCH_KIT:-$(kit_dir)}"
  # npx/dlx run from an ephemeral cache that gets cleaned up; symlinks into it would
  # dangle. Detect that and copy instead, so a one-shot `npx benchkit link` is durable.
  case "$kit" in
    *_npx*|*/dlx-*|*npm-cache*|*/.npm/_cacache/*)
      if [[ "$mode" == symlink ]]; then mode=copy
        echo "(running from an ephemeral package cache — using copy mode so files don't dangle)"
      fi ;;
  esac
  put() { if [[ "$mode" == copy ]]; then rm -rf "$2"; cp -r "$1" "$2"; else ln -sfn "$1" "$2"; fi; }
  # canonical instructions: AGENTS.md is the source of truth; CLAUDE.md is a shim
  put "$kit/AGENTS.md" "$root/AGENTS.md"
  [[ -e "$root/CLAUDE.md" ]] || printf '# Bench\n\nCanonical agreement in AGENTS.md.\n\n@AGENTS.md\n' > "$root/CLAUDE.md"
  # skills in BOTH locations so every harness finds them
  mkdir -p "$root/.claude" "$root/.agents"
  put "$kit/.claude/skills"   "$root/.claude/skills"
  put "$kit/.claude/skills"   "$root/.agents/skills"
  put "$kit/.claude/commands" "$root/.claude/commands"
  put "$kit/.claude/commands" "$root/.agents/commands"
  # Claude Code accelerants (hooks + settings) — ignored harmlessly by other harnesses
  put "$kit/.claude/hooks"         "$root/.claude/hooks"
  put "$kit/.claude/settings.json" "$root/.claude/settings.json"
  # harness-agnostic enforcement backstop
  install_git_hook "$root"
  echo "linked Bench into $root (mode: $mode)."
  echo "  Claude Code: CLAUDE.md -> @AGENTS.md, .claude/{skills,commands,hooks}"
  echo "  AGENTS.md harnesses (Codex/OpenCode/...): AGENTS.md, .agents/{skills,commands}"
  echo "  enforcement: git pre-push guard + the bench shift loop (both harness-independent)"
  echo "Run 'bench init' next to scaffold .bench/gate.sh."
}

case "${1:-help}" in
  gate)     run_gate ;;
  worktree) worktree ;;
  shift)    shift; shift_loop "${*:-improve the codebase}" ;;
  link)     shift; link "${1:-symlink}" ;;
  init)     init ;;
  *) cat <<EOF
bench — Pocock pipeline meets Kun Chen substrate, gated by your invariants.
  bench link [symlink|copy]  wire the kit into this repo for every harness at once
  bench init                 scaffold .bench/gate.sh in the current repo
  bench gate                 run the project gate (the oracle)
  bench worktree             warm, isolated worktree subshell
  bench shift "<objective>"  gated loop: one small change per iteration, commit on green
EOF
  ;;
esac
