# Canary overlay: a bench-query.sh whose open-learnings regex is MANGLED so it never
# matches a real `## <date> … [open]` heading. `bench learnings` then reports zero
# entries for a two-entry journal, so the AXI learnings two-entry contract must go red
# with "AXI learnings two-entry contract failed" — proving that check still bites if
# the journal parser stops recognizing open entries.
toon_escape() {
  local v="$1"
  case "$v" in
    *,* | *'"'* | *$'\n'* | [[:space:]]* | *[[:space:]])
      v="${v//\"/\"\"}"; printf '"%s"' "$v" ;;
    *) printf '%s' "$v" ;;
  esac
}
toon_table() {
  local name="$1" fields="$2"
  local -a fnames=() rows=()
  IFS=',' read -r -a fnames <<< "$fields"
  local ncols="${#fnames[@]}" line
  while IFS= read -r line || [[ -n "$line" ]]; do rows+=("$line"); done
  printf '%s[%d]{%s}:\n' "$name" "${#rows[@]}" "$fields"
  local rec rest out i
  if [[ ${#rows[@]} -gt 0 ]]; then
    for rec in "${rows[@]}"; do
      local -a cells=()
      rest="$rec"
      for ((i = 1; i < ncols; i++)); do cells+=("${rest%%$'\t'*}"); rest="${rest#*$'\t'}"; done
      cells+=("$rest")
      out=""
      for ((i = 0; i < ncols; i++)); do [[ $i -gt 0 ]] && out+=","; out+="$(toon_escape "${cells[$i]}")"; done
      printf '  %s\n' "$out"
    done
  fi
}

# BROKEN: this regex matches a literal "[OPEN-DISABLED]" tag no real heading carries,
# so learnings_open_headings returns nothing for a genuine journal.
LEARNINGS_OPEN_RE='^## NOPE-[0-9]+ \[OPEN-DISABLED\]'
learnings_open_headings() {
  local file="$1"; [[ -f "$file" ]] || return 0
  grep -E "$LEARNINGS_OPEN_RE" "$file" 2>/dev/null || true
}
learnings_rows() {
  local hd date title
  while IFS= read -r hd; do
    hd="${hd%$'\r'}"; [[ -n "$hd" ]] || continue
    date="${hd#\#\# }"; date="${date%%[[:space:]]*}"
    title="${hd#\#\# }"; title="${title#"$date"}"
    while true; do
      case "$title" in
        ' '*) title="${title# }" ;;
        '-'*) title="${title#-}" ;;
        '—'*) title="${title#—}" ;;
        *) break ;;
      esac
    done
    title="${title%\[open\]*}"
    while [[ "$title" == *[[:space:]] ]]; do title="${title%[[:space:]]}"; done
    printf '%s\t%s\n' "$date" "$title"
  done
}
learnings() {
  local root; root="$(git rev-parse --show-toplevel 2>/dev/null)" || return 1
  learnings_open_headings "$root/.bench/learnings.md" | learnings_rows \
    | toon_table learnings date,title
}
