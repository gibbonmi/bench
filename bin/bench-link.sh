# Link/install support for bench.sh.

# Resolve the effective git hooks directory — honoring core.hooksPath and worktree
# layouts where .git is a file — as an absolute path under $root.
hooks_dir() {
  local root="$1" dir
  dir="$(git -C "$root" rev-parse --git-path hooks 2>/dev/null)" || dir=".git/hooks"
  [[ "$dir" == /* ]] || dir="$root/$dir"
  printf '%s\n' "$dir"
}

install_git_hook() {
  local root="$1" def hooks
  # A remote without origin/HEAD set otherwise resolves to main and guards the wrong
  # branch; nudge it at link time only so default_branch() stays network-free.
  if git -C "$root" remote get-url origin >/dev/null 2>&1 \
     && ! git -C "$root" symbolic-ref --quiet refs/remotes/origin/HEAD >/dev/null 2>&1; then
    git -C "$root" remote set-head origin --auto >/dev/null 2>&1 || true
  fi
  def="$(default_branch)"
  hooks="$(hooks_dir "$root")"
  mkdir -p "$hooks"
  if [[ -e "$hooks/pre-push" ]] && ! grep -q 'bench:managed-pre-push' "$hooks/pre-push"; then
    echo "conflict: $hooks/pre-push exists and is not Bench-managed" >&2
    return 1
  fi
  cat > "$hooks/pre-push" <<EOF
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
  chmod +x "$hooks/pre-push"
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

# Shared awk prelude for every marker path: a line starting with a ``` fence toggles
# fence state; lines inside a fence are printed but never treated as Bench markers, so
# AGENTS.md may document the markers in an example block without a relink eating it.
fence_awk_prelude='{ if (substr($0, 1, 3) == "```") { in_fence = !in_fence; line_is_fence = 1 } else { line_is_fence = 0 } }'

count_marker() {
  local marker="$1" file="$2"
  awk -v marker="$marker" "$fence_awk_prelude"'
    !line_is_fence && !in_fence && index($0, marker) { count++ }
    END { print count + 0 }
  ' "$file" 2>/dev/null
}

marker_line() {
  local marker="$1" file="$2"
  awk -v marker="$marker" "$fence_awk_prelude"'
    !line_is_fence && !in_fence && index($0, marker) { print NR; exit }
  ' "$file" 2>/dev/null
}

fence_balanced() {
  local file="$1"
  [[ "$(awk "$fence_awk_prelude"' END { print in_fence + 0 }' "$file" 2>/dev/null)" == "0" ]]
}

validate_agents_block() {
  local root="$1" agents="$root/AGENTS.md" starts=0 ends=0 start_line end_line
  [[ -f "$agents" ]] || return 0
  if ! fence_balanced "$agents" && grep -qF 'bench:' "$agents"; then
    echo "conflict: AGENTS.md has an unclosed code fence around Bench markers; marker detection cannot be trusted" >&2
    return 1
  fi
  starts="$(count_marker '<!-- bench:start -->' "$agents")"
  ends="$(count_marker '<!-- bench:end -->' "$agents")"
  if [[ "$starts" == "0" && "$ends" == "0" ]]; then return 0; fi
  if [[ "$starts" == "1" && "$ends" == "1" ]]; then
    start_line="$(marker_line '<!-- bench:start -->' "$agents")"
    end_line="$(marker_line '<!-- bench:end -->' "$agents")"
    if [[ -n "$start_line" && -n "$end_line" && "$start_line" -lt "$end_line" ]]; then
      return 0
    fi
  fi
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
      awk -v block="$block" "$fence_awk_prelude"'
        !line_is_fence && !in_fence && /<!-- bench:start -->/ {
          if (!done) {
            while ((getline line < block) > 0) print line
            close(block)
            done = 1
          }
          skip = 1
          next
        }
        !line_is_fence && !in_fence && /<!-- bench:end -->/ { skip = 0; next }
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
    rel="${src#"$src_root"/}"
    printf '%s\t%s/%s\t%s\n' "$src" "$dest_root" "$rel" "$kind" >> "$plan"
  done < <(find "$src_root" -type f | sort)
}

build_link_plan() {
  local kit="$1" plan="$2"
  : > "$plan"
  append_tree_to_plan "$kit/bin" ".bench/bin" "file" "$plan"
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
  local root="$1" plan="$2" conflicts=0 src rel kind dest parent old hooks
  validate_agents_block "$root" || conflicts=$((conflicts+1))
  hooks="$(hooks_dir "$root")"
  if [[ -e "$hooks/pre-push" ]] && ! grep -q 'bench:managed-pre-push' "$hooks/pre-push"; then
    echo "conflict: $hooks/pre-push exists and is not Bench-managed" >&2
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

link() {
  local root kit mode="${1:-copy}" plan
  root="$(repo_root)"; kit="${BENCH_KIT:-$(kit_dir)}"
  case "$mode" in copy|symlink) ;; *) echo "usage: bench link [copy|symlink]" >&2; return 2 ;; esac
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
