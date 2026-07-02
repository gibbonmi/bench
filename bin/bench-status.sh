# Status, roadmap, model, and structure support for bench.sh.

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

structure() {
  structure_check all ""
}

structure_touched_since() {
  local root base files
  root="$(repo_root)"
  base="$1"
  files="$(git -C "$root" diff --name-only --diff-filter=ACMR "$base"..HEAD || true)"
  structure_check touched "$files"
}

# Reviewer-owned per-path overrides from .bench/structure.budgets: `<path> <n>`
# per line, trailing `/` on the path means a directory file-count budget, `#`
# comments and blanks ignored. The value replaces the global cap for that path
# (lower as well as higher). Paths are matched exactly — no globs, a grant is a
# named decision — and cannot contain spaces (the format splits on whitespace).
structure_budgets_load() {
  local file="$1" line p b
  [[ -f "$file" ]] || return 0
  while IFS= read -r line || [[ -n "$line" ]]; do
    line="${line%%#*}"
    [[ -z "${line//[[:space:]]/}" ]] && continue
    read -r p b _ <<<"$line"
    if [[ ! "${b:-}" =~ ^[0-9]+$ ]]; then
      echo "structure.budgets: ignoring malformed line: $line" >&2
      continue
    fi
    printf '%s %s\n' "$p" "$b"
  done < "$file"
}

structure_budget_for() { # budgets, path, fallback
  local hit
  hit="$(printf '%s' "$1" | awk -v p="$2" '$1 == p { print $2; exit }')"
  printf '%s\n' "${hit:-$3}"
}

structure_check() {
  local mode="$1" scoped_files="$2" root max_lines max_files exts files violations=0 budgets
  root="$(repo_root)"
  max_lines="${BENCH_MAX_LINES:-400}"
  max_files="${BENCH_MAX_DIR_FILES:-12}"
  budgets="$(structure_budgets_load "$root/.bench/structure.budgets")"
  exts='py|ts|tsx|js|jsx|go|rs|java|rb|kt|scala|cs|cpp|cc|c|h|hpp|sh'
  if [[ "$mode" == all ]]; then
    files=$(git -C "$root" ls-files 2>/dev/null | grep -E "\.($exts)\$" || true)
  else
    files=$(printf '%s\n' "$scoped_files" | grep -E "\.($exts)\$" || true)
  fi
  [[ -z "$files" ]] && { echo "structure: no tracked source files to check"; return 0; }
  while IFS= read -r f; do
    [[ -z "$f" ]] && continue
    [[ -f "$root/$f" ]] || continue
    local n cap; n=$(wc -l < "$root/$f" 2>/dev/null || echo 0)
    cap="$(structure_budget_for "$budgets" "$f" "$max_lines")"
    if (( n > cap )); then
      echo "FILE TOO LONG   $n lines (max $cap)   $f"; violations=$((violations+1))
    fi
  done <<< "$files"
  local current_dir="" dir="" dir_count=0 saw_dir=0 dir_cap
  while IFS= read -r dir; do
    if (( saw_dir == 0 )) || [[ "$dir" != "$current_dir" ]]; then
      dir_cap="$(structure_budget_for "$budgets" "$current_dir/" "$max_files")"
      if (( saw_dir == 1 && dir_count > dir_cap )); then
        echo "DIR CROWDED     $dir_count source files (max $dir_cap), group into modules   $current_dir/"; violations=$((violations+1))
      fi
      current_dir="$dir"
      dir_count=1
      saw_dir=1
    else
      dir_count=$((dir_count+1))
    fi
  done < <(
    while IFS= read -r f; do
      [[ -z "$f" ]] && continue
      dir="${f%/*}"
      [[ "$dir" == "$f" ]] && dir="."
      printf '%s\n' "$dir"
    done <<< "$files" | sort
  )
  dir_cap="$(structure_budget_for "$budgets" "$current_dir/" "$max_files")"
  if (( saw_dir == 1 && dir_count > dir_cap )); then
    echo "DIR CROWDED     $dir_count source files (max $dir_cap), group into modules   $current_dir/"; violations=$((violations+1))
  fi
  if (( violations > 0 )); then
    echo "structural debt: $violations issue(s). Split along responsibility (see the craft-seams skill); don't fragment to beat the number." >&2
    return 1
  fi
  echo "structure ok (≤$max_lines lines/file, ≤$max_files source files/dir)"
}

idea() {
  local root text file
  root="$(repo_root)"
  text="$*"
  if [[ -z "${text//[[:space:]]/}" ]]; then
    echo 'usage: bench idea "<text>"' >&2
    return 2
  fi
  file="$root/ROADMAP.md"
  # Normalize a missing trailing newline before appending, or a hand-edited file
  # whose last line lacks one would swallow this entry onto the same physical line
  # (and break the `^- ` count that `roadmap`/`status` rely on).
  [[ -s "$file" && -n "$(tail -c1 "$file")" ]] && printf '\n' >> "$file"
  printf '%s  %s\n' "- $(date +%F)" "$text" >> "$file"
  echo "parked: $text"
}

roadmap() {
  local root file
  root="$(repo_root)"
  file="$root/ROADMAP.md"
  [[ -s "$file" ]] || { echo "roadmap empty"; return 0; }
  cat "$file"
}

gate_tree_hash() {
  # Hash of the content actually on disk (tracked + untracked-unignored), computed
  # through a throwaway index so the real index is never touched. This is the gate
  # cache key: the verdict is content-addressed, so a commit that doesn't change
  # the tree keeps the verdict fresh, and a verdict from a dirty tree never gets
  # attributed to a commit whose content was not tested. Mirrored in
  # .bench/hooks/stop.sh, which cannot source this file.
  local root="$1" idx hash
  idx="${TMPDIR:-/tmp}/bench-tree-idx.$$"
  hash="$(
    cd "$root" || exit 1
    export GIT_INDEX_FILE="$idx"
    git read-tree HEAD 2>/dev/null || git read-tree --empty
    git add -A 2>/dev/null
    git write-tree
  )" || hash=""
  rm -f "$idx"
  printf '%s\n' "${hash:-none}"
}

status() {
  local root cache; root="$(repo_root)"
  local -a rows=()
  local footer=""

  cache="$(git -C "$root" rev-parse --absolute-git-dir 2>/dev/null)/bench-last-gate"
  if [[ -f "$cache" ]]; then
    local cstatus ctree _crest tree
    read -r cstatus ctree _crest < "$cache" || true
    tree="$(gate_tree_hash "$root")"
    if [[ "$ctree" != "$tree" ]]; then
      rows+=("6|gate|stale (gated tree ${ctree:0:7}, work tree ${tree:0:7})|re-run the gate")
    elif [[ "$cstatus" == "red" ]]; then
      rows+=("0|gate|red|fix before commit")
    fi
  fi

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

  # Warm pooled worktrees (released, no lease) are the pool doing its job, not a
  # signal. Count only leased pool entries and worktrees outside the pool.
  local wtc=0 pool wpath
  pool="$(worktree_pool "$root")"
  while IFS= read -r wpath; do
    [[ "$wpath" == worktree\ * ]] || continue
    wpath="${wpath#worktree }"
    [[ "$wpath" == "$root" ]] && continue
    [[ "$wpath" == "$pool"/* && ! -f "$(worktree_lease_file "$wpath")" ]] && continue
    wtc=$((wtc+1))
  done < <(git -C "$root" worktree list --porcelain 2>/dev/null || true)
  if (( wtc > 0 )); then
    rows+=("2|worktree|$wtc active worktree(s)|resume or clean up (bench worktree)")
  fi

  local floor open=0; floor="${BENCH_LEARNINGS_FLOOR:-1}"
  [[ -f "$root/.bench/learnings.md" ]] && open="$(grep -cE '^## [0-9]{4}-[0-9]{2}-[0-9]{2}.*\[open\]' "$root/.bench/learnings.md" 2>/dev/null || true)"
  if [[ "${open:-0}" -ge "$floor" && "${open:-0}" -gt 0 ]]; then
    rows+=("3|learnings|$open open|/bench-integrate-learnings")
  fi

  local sviol; sviol="$(structure 2>/dev/null | grep -cE '^(FILE TOO LONG|DIR CROWDED)' || true)"
  if [[ "${sviol:-0}" -gt 0 ]]; then
    rows+=("4|structure|$sviol issue(s)|split (craft-seams)")
  fi

  local maps=0
  if [[ -d "$root/decisions" ]]; then
    maps="$(grep -lE '^— \((open|deferred)|GRILL DEFERRED' "$root"/decisions/*.md 2>/dev/null | wc -l | tr -d ' ' || true)"
  fi
  if [[ "${maps:-0}" -gt 0 ]]; then
    rows+=("5|decisions|$maps unresolved map(s)|craft-grill → /bench-write-spec")
  fi

  if [[ -s "$root/ROADMAP.md" ]]; then
    local n; n="$(grep -c '^- ' "$root/ROADMAP.md" 2>/dev/null || true)"
    [[ "${n:-0}" -gt 0 ]] && footer="$n idea(s) parked — bench roadmap"
  fi

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
