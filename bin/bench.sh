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
  # Implementation loop is done. Only NOW pay down structural debt, if the work
  # pushed past the budget. Refactor at green — never mid-implementation: splitting
  # before the feature's shape has settled produces premature, bad module boundaries.
  # This is "only trigger a refactor once it's over the threshold", at the loop edge.
  if ! structure >/dev/null 2>&1; then
    echo "▶ structure over budget — refactor phase (split at green, not before)"
    local rcap="${BENCH_REFACTOR_ITERS:-4}" r
    for ((r=1; r<=rcap; r++)); do
      structure >/dev/null 2>&1 && break
      echo "── refactor $r/$rcap ──"
      BENCH_SHIFT=1 "$AGENT" -p "$(refactor_prompt)" || true
      if run_gate; then
        git -C "$root" add -A -- ':!.bench-objective' ':!.bench-notes.md'
        git -C "$root" diff --cached --quiet || \
          git -C "$root" commit -q -m "refactor: reduce structural debt"
        echo "  ✓ tests green — refactor $r committed"
      else
        echo "  ✗ refactor broke the gate — rolling back"
        git -C "$root" reset -q --hard; git -C "$root" clean -qfd
      fi
    done
    if structure >/dev/null 2>&1; then echo "  structure back under budget."
    else echo "  ⚠ still over budget after $rcap passes — review manually, or run /improve-codebase-architecture for a deep pass."; fi
  fi
  rm -f "$root/.bench-objective" "$root/.bench-notes.md"
  echo "■ shift done: $branch, $((i-1)) committed iteration(s), $(( ($(date +%s)-started)/60 ))m elapsed"
  echo "  review: git -C $root log --oneline origin/$(default_branch)..$branch"
  echo "  the merge is yours."
}

refactor_prompt() {
  cat <<'EOF'
The implementation is complete and tests are green, but the structure budget is
exceeded. Run `bench structure` to see the flagged files and directories. Fix them
by splitting along responsibility, using the deletion test from the bench-craft-seams skill: lift
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
decides. `/bench-learn` reviews the open entries, promotes the generalizable ones
into the kit with sign-off, and marks them resolved. Never rewrite a kit rule
yourself — that is the whole point of capturing here instead.

Format per entry:

## <date> — <short title>  [open|resolved]
- **What happened:** …
- **Right behavior:** …
- **Proposed rule change:** … (or "none")

An entry only becomes `[resolved]` via /bench-learn.

<!-- entries below -->
EOF
    echo "scaffolded .bench/learnings.md — the self-learning journal"
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
  if [[ -e "$root/.git/hooks/pre-push" ]] && ! grep -q 'bench:managed-pre-push' "$root/.git/hooks/pre-push"; then
    echo "conflict: .git/hooks/pre-push exists and is not Bench-managed" >&2
    return 1
  fi
  cat > "$root/.git/hooks/pre-push" <<EOF
#!/usr/bin/env bash
# bench:managed-pre-push
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

hash_file() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{print $1}'
  else
    shasum -a 256 "$1" | awk '{print $1}'
  fi
}

hash_stdin() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum | awk '{print $1}'
  else
    shasum -a 256 | awk '{print $1}'
  fi
}

fingerprint_path() {
  local path="$1"
  if [[ -L "$path" ]]; then
    printf 'symlink:%s\n' "$(readlink "$path")" | hash_stdin
  elif [[ -f "$path" ]]; then
    hash_file "$path"
  else
    return 1
  fi
}

manifest_hash() {
  local root="$1" rel="$2" manifest="$root/.bench/link-manifest.tsv"
  [[ -f "$manifest" ]] || return 1
  awk -F '\t' -v p="$rel" '$1 == p { h = $2 } END { if (h != "") print h; else exit 1 }' "$manifest"
}

manifest_owned_clean() {
  local root="$1" rel="$2" path="$root/$rel" old now
  old="$(manifest_hash "$root" "$rel" 2>/dev/null)" || return 1
  [[ -e "$path" || -L "$path" ]] || return 0
  now="$(fingerprint_path "$path" 2>/dev/null)" || return 1
  [[ "$old" == "$now" ]]
}

bench_agents_block() {
  cat <<'EOF'
<!-- bench:start -->
## Bench

Bench is installed in this repo.

- The gate is the oracle: run `bench gate`; done means it exits zero.
- Full operating guide: `.bench/BENCH.md`.
- Portable commands: `.agents/commands/`.
- Portable skills: `.agents/skills/`.
- Project profile: `projects/<name>.md` when present.
- The reviewer owns merge decisions.
<!-- bench:end -->
EOF
}

count_marker() {
  local marker="$1" file="$2"
  awk -v marker="$marker" 'index($0, marker) { count++ } END { print count + 0 }' "$file" 2>/dev/null
}

validate_agents_block() {
  local root="$1" agents="$root/AGENTS.md" starts=0 ends=0
  [[ -f "$agents" ]] || return 0
  starts="$(count_marker '<!-- bench:start -->' "$agents")"
  ends="$(count_marker '<!-- bench:end -->' "$agents")"
  if [[ "$starts" == "0" && "$ends" == "0" ]]; then return 0; fi
  if [[ "$starts" == "1" && "$ends" == "1" ]]; then return 0; fi
  echo "conflict: AGENTS.md has malformed Bench managed block markers" >&2
  return 1
}

write_agents_block() {
  local root="$1" agents="$root/AGENTS.md" block tmp starts
  block="$(mktemp)"; tmp="$(mktemp)"
  bench_agents_block > "$block"
  if [[ ! -f "$agents" ]]; then
    cp "$block" "$agents"
  else
    starts="$(count_marker '<!-- bench:start -->' "$agents")"
    if [[ "$starts" == "0" ]]; then
      cat "$agents" > "$tmp"
      printf '\n' >> "$tmp"
      cat "$block" >> "$tmp"
      mv "$tmp" "$agents"
    else
      awk -v block="$block" '
        /<!-- bench:start -->/ {
          if (!done) {
            while ((getline line < block) > 0) print line
            close(block)
            done = 1
          }
          skip = 1
          next
        }
        /<!-- bench:end -->/ { skip = 0; next }
        !skip { print }
      ' "$agents" > "$tmp"
      mv "$tmp" "$agents"
    fi
  fi
  rm -f "$block" "$tmp"
}

append_tree_to_plan() {
  local src_root="$1" dest_root="$2" kind="$3" plan="$4" src rel
  [[ -d "$src_root" ]] || return 0
  while IFS= read -r src; do
    rel="${src#$src_root/}"
    printf '%s\t%s/%s\t%s\n' "$src" "$dest_root" "$rel" "$kind" >> "$plan"
  done < <(find "$src_root" -type f | sort)
}

build_link_plan() {
  local kit="$1" plan="$2"
  : > "$plan"
  printf '%s\t%s\t%s\n' "$kit/.bench/BENCH.md" ".bench/BENCH.md" "file" >> "$plan"
  printf '%s\t%s\t%s\n' "$kit/.claude/README.md" ".claude/README.md" "file" >> "$plan"
  printf '%s\t%s\t%s\n' "$kit/.claude/settings.json" ".claude/settings.json" "file" >> "$plan"
  printf '%s\t%s\t%s\n' "$kit/.codex/hooks.json" ".codex/hooks.json" "file" >> "$plan"
  append_tree_to_plan "$kit/.bench/hooks" ".bench/hooks" "file" "$plan"
  append_tree_to_plan "$kit/.agents/commands" ".agents/commands" "file" "$plan"
  append_tree_to_plan "$kit/.agents/skills" ".agents/skills" "file" "$plan"
  append_tree_to_plan "$kit/.agents/commands" ".claude/commands" "adapter" "$plan"
  append_tree_to_plan "$kit/.agents/skills" ".claude/skills" "adapter" "$plan"
}

adapter_target() {
  local rel="$1" dest_dir source_rel ups="" part
  case "$rel" in
    .claude/commands/*) source_rel=".agents/commands/${rel#.claude/commands/}" ;;
    .claude/skills/*) source_rel=".agents/skills/${rel#.claude/skills/}" ;;
    *) return 1 ;;
  esac
  dest_dir="$(dirname "$rel")"
  IFS='/' read -r -a parts <<< "$dest_dir"
  for part in "${parts[@]}"; do ups="../$ups"; done
  printf '%s%s' "$ups" "$source_rel"
}

has_symlink_parent() {
  local root="$1" rel="$2" dir path part
  dir="$(dirname "$rel")"
  [[ "$dir" == "." ]] && return 1
  path="$root"
  IFS='/' read -r -a parts <<< "$dir"
  for part in "${parts[@]}"; do
    path="$path/$part"
    [[ -L "$path" ]] && return 0
  done
  return 1
}

preflight_link() {
  local root="$1" plan="$2" conflicts=0 src rel kind dest parent old
  validate_agents_block "$root" || conflicts=$((conflicts+1))
  if [[ -e "$root/.git/hooks/pre-push" ]] && ! grep -q 'bench:managed-pre-push' "$root/.git/hooks/pre-push"; then
    echo "conflict: .git/hooks/pre-push exists and is not Bench-managed" >&2
    conflicts=$((conflicts+1))
  fi
  while IFS=$'\t' read -r src rel kind; do
    [[ -n "$rel" ]] || continue
    [[ -f "$src" ]] || { echo "conflict: kit asset missing: $src" >&2; conflicts=$((conflicts+1)); continue; }
    parent="$root/$(dirname "$rel")"
    if has_symlink_parent "$root" "$rel"; then
      echo "conflict: $rel has a symlink parent directory" >&2
      conflicts=$((conflicts+1))
      continue
    fi
    if [[ -e "$parent" && ! -d "$parent" ]]; then
      echo "conflict: parent path for $rel is not a directory" >&2
      conflicts=$((conflicts+1))
      continue
    fi
    dest="$root/$rel"
    if [[ -e "$dest" || -L "$dest" ]]; then
      if ! manifest_owned_clean "$root" "$rel"; then
        old="$(manifest_hash "$root" "$rel" 2>/dev/null || true)"
        if [[ -n "$old" ]]; then
          echo "conflict: modified Bench-managed file: $rel" >&2
        else
          echo "conflict: project-owned file exists: $rel" >&2
        fi
        conflicts=$((conflicts+1))
      fi
    fi
  done < "$plan"
  [[ "$conflicts" -eq 0 ]]
}

install_planned_file() {
  local root="$1" mode="$2" src="$3" rel="$4" kind="$5" dest target
  dest="$root/$rel"
  mkdir -p "$(dirname "$dest")"
  rm -f "$dest"
  if [[ "$kind" == "adapter" ]]; then
    target="$(adapter_target "$rel")"
    ln -s "$target" "$dest"
  elif [[ "$mode" == "symlink" ]]; then
    ln -s "$src" "$dest"
  else
    cp -p "$src" "$dest"
  fi
}

install_plan() {
  local root="$1" mode="$2" plan="$3" manifest tmp_manifest src rel kind fp
  manifest="$root/.bench/link-manifest.tsv"
  tmp_manifest="$manifest.tmp"
  mkdir -p "$root/.bench"
  : > "$tmp_manifest"
  while IFS=$'\t' read -r src rel kind; do
    [[ -n "$rel" ]] || continue
    install_planned_file "$root" "$mode" "$src" "$rel" "$kind"
    fp="$(fingerprint_path "$root/$rel")"
    printf '%s\t%s\n' "$rel" "$fp" >> "$tmp_manifest"
  done < "$plan"
  mv "$tmp_manifest" "$manifest"
}

# Best-effort model discovery. There is no universal cross-harness "list models"
# command, so: query the Anthropic Models API when a key is present (authoritative),
# otherwise point at the harness's own list. The agent/setup bind tiers from this.
models() {
  if [[ -n "${ANTHROPIC_API_KEY:-}" ]]; then
    echo "Available models (Anthropic Models API):"
    if ! curl -fsS https://api.anthropic.com/v1/models \
          -H "x-api-key: $ANTHROPIC_API_KEY" \
          -H "anthropic-version: 2023-06-01" \
        | python3 -c 'import sys,json
for m in json.load(sys.stdin).get("data",[]): print("  "+m["id"])'; then
      echo "  (query failed — check the key, or read your harness model list)"
    fi
  else
    cat <<'EOF'
No ANTHROPIC_API_KEY set, so I can't query the model list directly. Discover from
your harness instead, then bind the tiers (cheap / mid / top) in projects/<name>.md:
  - Claude Code: `claude --help`, or the in-app /model picker
  - Codex:       `codex --help`, or its model config
  - or export ANTHROPIC_API_KEY and re-run `bench models`
EOF
  fi
}

# Deterministic, language-agnostic structural-debt check, for the gate. Flags files
# over a line budget and directories with too many source files (the "30 scripts in
# src/" smell). Thresholds via env. Exit 1 on violations so it gates. The *how* to
# split is the bench-craft-seams skill's job — this only measures.
structure() {
  local root max_lines max_files exts files violations=0
  root="$(repo_root)"
  max_lines="${BENCH_MAX_LINES:-400}"
  max_files="${BENCH_MAX_DIR_FILES:-12}"
  exts='py|ts|tsx|js|jsx|go|rs|java|rb|kt|scala|cs|cpp|cc|c|h|hpp'
  files=$(git -C "$root" ls-files 2>/dev/null | grep -E "\.($exts)\$" || true)
  [[ -z "$files" ]] && { echo "structure: no tracked source files to check"; return 0; }
  while IFS= read -r f; do
    [[ -z "$f" ]] && continue
    local n; n=$(wc -l < "$root/$f" 2>/dev/null || echo 0)
    if (( n > max_lines )); then
      echo "FILE TOO LONG   $n lines (max $max_lines)   $f"; violations=$((violations+1))
    fi
  done <<< "$files"
  while read -r count dir; do
    [[ -z "$count" ]] && continue
    if (( count > max_files )); then
      echo "DIR CROWDED     $count source files (max $max_files), group into modules   $dir/"; violations=$((violations+1))
    fi
  done <<< "$(printf '%s\n' "$files" | xargs -n1 dirname | sort | uniq -c)"
  if (( violations > 0 )); then
    echo "structural debt: $violations issue(s). Split along responsibility (see the bench-craft-seams skill); don't fragment to beat the number." >&2
    return 1
  fi
  echo "structure ok (≤$max_lines lines/file, ≤$max_files source files/dir)"
}

# ---- roadmap: capture an idea without committing to it ----------------------
# `bench idea "<text>"` parks an out-of-scope idea in ROADMAP.md at the repo root and
# exits — no prompt, no workflow, no spec. Append-only, one dated line per entry; the
# file is a dumb sink (no status, no lifecycle). `/bench-ideate` promotes a parked
# idea into a decision map. `bench roadmap` prints the list on demand.
idea() {
  local root text
  root="$(repo_root)"
  text="$*"
  if [[ -z "${text//[[:space:]]/}" ]]; then
    echo 'usage: bench idea "<text>"' >&2
    return 2
  fi
  printf '%s  %s\n' "- $(date +%F)" "$text" >> "$root/ROADMAP.md"
  echo "parked: $text"
}

roadmap() {
  local root file
  root="$(repo_root)"
  file="$root/ROADMAP.md"
  [[ -s "$file" ]] || { echo "roadmap empty"; return 0; }
  cat "$file"
}

# ---- status: the ambient dashboard ------------------------------------------
# `bench status` is the single renderer (SessionStart hook + on-demand). It reads cheap
# repo state and a gate cache, ranks the signals that fire on a fixed severity ladder,
# and prints a one-line lead, up to five rows, and a zero-severity roadmap footer.
# Deterministic plain shell — no model at render time; the agent reading it adds judgment.
# Show-only-on-signal: a signal with nothing to report prints no row.
status() {
  local root head cache; root="$(repo_root)"
  head="$(git -C "$root" rev-parse HEAD 2>/dev/null || echo none)"
  local -a rows=()   # "rank|signal|detail|action"
  local footer=""

  # rank 0 / 6 — gate, from the Stop-hook cache (never a cold run). red iff fresh-red;
  # stale iff the cache sha != HEAD (a stale green is not a clean bill); else silent.
  # The cache lives in the git dir (never tracked), so it can't pollute commits.
  cache="$(git -C "$root" rev-parse --absolute-git-dir 2>/dev/null)/bench-last-gate"
  if [[ -f "$cache" ]]; then
    local cstatus csha _crest
    read -r cstatus csha _crest < "$cache" || true
    if [[ "$csha" != "$head" ]]; then
      rows+=("6|gate|stale (cache ${csha:0:7}, HEAD ${head:0:7})|re-run the gate")
    elif [[ "$cstatus" == "red" ]]; then
      rows+=("0|gate|red|fix before commit")
    fi
  fi

  # rank 1 — uncommitted / unpushed
  local dirty="" ahead=""
  dirty="$(git -C "$root" status --porcelain 2>/dev/null || true)"
  if git -C "$root" rev-parse --abbrev-ref --symbolic-full-name '@{u}' >/dev/null 2>&1; then
    ahead="$(git -C "$root" log --oneline '@{u}..HEAD' 2>/dev/null || true)"
  fi
  if [[ -n "$dirty" || -n "$ahead" ]]; then
    local d="uncommitted + unpushed"
    [[ -n "$dirty" && -z "$ahead" ]] && d="uncommitted changes"
    [[ -z "$dirty" && -n "$ahead" ]] && d="unpushed commits"
    rows+=("1|git|$d|commit on green / push")
  fi

  # rank 2 — stray worktree / active shift
  local wtc; wtc="$(git -C "$root" worktree list 2>/dev/null | wc -l | tr -d ' ' || true)"
  if [[ "${wtc:-1}" -gt 1 ]]; then
    rows+=("2|worktree|$((wtc-1)) extra worktree(s)|resume or clean up (bench worktree)")
  fi

  # rank 3 — open learnings (configurable floor)
  local floor open=0; floor="${BENCH_LEARNINGS_FLOOR:-1}"
  [[ -f "$root/.bench/learnings.md" ]] && open="$(grep -c '\[open\]' "$root/.bench/learnings.md" 2>/dev/null || true)"
  if [[ "${open:-0}" -ge "$floor" && "${open:-0}" -gt 0 ]]; then
    rows+=("3|learnings|$open open|/bench-learn")
  fi

  # rank 4 — structural debt (reuses bench structure)
  local sviol; sviol="$(structure 2>/dev/null | grep -cE '^(FILE TOO LONG|DIR CROWDED)' || true)"
  if [[ "${sviol:-0}" -gt 0 ]]; then
    rows+=("4|structure|$sviol issue(s)|split (bench-craft-seams)")
  fi

  # rank 5 — unresolved decision map (open-ticket marker convention)
  local maps=0
  if [[ -d "$root/decisions" ]]; then
    maps="$(grep -lE '^— \((open|deferred)|GRILL DEFERRED' "$root"/decisions/*.md 2>/dev/null | wc -l | tr -d ' ' || true)"
  fi
  if [[ "${maps:-0}" -gt 0 ]]; then
    rows+=("5|decisions|$maps unresolved map(s)|/bench-craft-grill → /bench-spec")
  fi

  # footer — roadmap count (zero severity: never ranked, never the lead, never budgeted)
  if [[ -s "$root/ROADMAP.md" ]]; then
    local n; n="$(grep -c '^- ' "$root/ROADMAP.md" 2>/dev/null || true)"
    [[ "${n:-0}" -gt 0 ]] && footer="$n idea(s) parked — bench roadmap"
  fi

  # --- render ---
  if [[ ${#rows[@]} -eq 0 ]]; then
    echo "bench: clean — nothing pending"
    [[ -n "$footer" ]] && echo "$footer"
    return 0
  fi
  local sorted; sorted="$(printf '%s\n' "${rows[@]}" | sort -t'|' -k1,1n)"
  local _r lead_signal _d lead_action
  IFS='|' read -r _r lead_signal _d lead_action <<< "$(printf '%s\n' "$sorted" | head -1)"
  echo "▶ $lead_action  ($lead_signal)"
  local total shown=0; total="$(printf '%s\n' "$sorted" | grep -c .)"
  while IFS='|' read -r _r signal detail action; do
    [[ -z "$_r" ]] && continue
    shown=$((shown+1))
    [[ "$shown" -le 5 ]] && printf '  %-10s %-30s → %s\n' "$signal" "$detail" "$action"
  done <<< "$sorted"
  [[ "$total" -gt 5 ]] && echo "  +$((total-5)) more"
  [[ -n "$footer" ]] && echo "$footer"
  return 0
}

# Wire the kit into the current repo for EVERY harness at once. Idempotent.
# After this, Claude Code and any AGENTS.md harness work in the same repo with no
# per-switch reconfiguration — switching harnesses is just running a different agent.
link() {
  local root kit mode="${1:-copy}" plan
  root="$(repo_root)"; kit="${BENCH_KIT:-$(kit_dir)}"
  case "$mode" in copy|symlink) ;; *) echo "usage: bench link [copy|symlink]" >&2; return 2 ;; esac
  # npx/dlx run from an ephemeral cache that gets cleaned up; symlinks into it would
  # dangle. Detect that and copy instead, so a one-shot `npx benchkit link` is durable.
  case "$kit" in
    *_npx*|*/dlx-*|*npm-cache*|*/.npm/_cacache/*)
      if [[ "$mode" == symlink ]]; then mode=copy
        echo "(running from an ephemeral package cache — using copy mode so files don't dangle)"
      fi ;;
  esac
  plan="$(mktemp)"
  build_link_plan "$kit" "$plan"
  preflight_link "$root" "$plan" || { rm -f "$plan"; return 1; }
  write_agents_block "$root"
  [[ -e "$root/CLAUDE.md" ]] || printf '# Bench\n\nCanonical agreement in AGENTS.md.\n\n@AGENTS.md\n' > "$root/CLAUDE.md"
  install_plan "$root" "$mode" "$plan"
  rm -f "$plan"
  install_git_hook "$root" || return 1
  echo "linked Bench into $root (mode: $mode)."
  echo "  instructions: AGENTS.md managed block -> .bench/BENCH.md"
  echo "  portable surface: .agents/{skills,commands}"
  echo "  adapters: .claude/settings.json, .codex/hooks.json, .claude/{skills,commands} -> .agents"
  echo "  enforcement: shared .bench/hooks + git pre-push guard + the bench shift loop"
  echo "Run 'bench init' next to scaffold .bench/gate.sh."
}

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
  bench shift "<objective>"  gated loop: one small change per iteration, commit on green
EOF
  ;;
esac
