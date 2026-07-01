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

structure_check() {
  local mode="$1" scoped_files="$2" root max_lines max_files exts files violations=0
  root="$(repo_root)"
  max_lines="${BENCH_MAX_LINES:-400}"
  max_files="${BENCH_MAX_DIR_FILES:-12}"
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
    local n; n=$(wc -l < "$root/$f" 2>/dev/null || echo 0)
    if (( n > max_lines )); then
      echo "FILE TOO LONG   $n lines (max $max_lines)   $f"; violations=$((violations+1))
    fi
  done <<< "$files"
  local current_dir="" dir="" dir_count=0 saw_dir=0
  while IFS= read -r dir; do
    if (( saw_dir == 0 )) || [[ "$dir" != "$current_dir" ]]; then
      if (( saw_dir == 1 && dir_count > max_files )); then
        echo "DIR CROWDED     $dir_count source files (max $max_files), group into modules   $current_dir/"; violations=$((violations+1))
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
  if (( saw_dir == 1 && dir_count > max_files )); then
    echo "DIR CROWDED     $dir_count source files (max $max_files), group into modules   $current_dir/"; violations=$((violations+1))
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

status() {
  local root head cache; root="$(repo_root)"
  head="$(git -C "$root" rev-parse HEAD 2>/dev/null || echo none)"
  local -a rows=()
  local footer=""

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

  local wtc; wtc="$(git -C "$root" worktree list 2>/dev/null | wc -l | tr -d ' ' || true)"
  if [[ "${wtc:-1}" -gt 1 ]]; then
    rows+=("2|worktree|$((wtc-1)) extra worktree(s)|resume or clean up (bench worktree)")
  fi

  local floor open=0; floor="${BENCH_LEARNINGS_FLOOR:-1}"
  [[ -f "$root/.bench/learnings.md" ]] && open="$(grep -c '\[open\]' "$root/.bench/learnings.md" 2>/dev/null || true)"
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
