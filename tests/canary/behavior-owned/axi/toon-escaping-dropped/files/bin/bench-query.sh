# Canary overlay: a bench-query.sh whose flat-table TOON emitter DROPS escaping —
# toon_escape emits every field verbatim, so a field carrying a comma or quote
# corrupts the row. The AXI field-escaping contract in gate-axi-contracts.sh must go
# red with "AXI TOON field-escaping contract failed", proving that check still bites
# if the shared emitter stops quoting hostile field values.
toon_escape() { printf '%s' "$1"; }   # BROKEN: no quoting of comma/quote/whitespace

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

LEARNINGS_OPEN_RE='^## [0-9]{4}-[0-9]{2}-[0-9]{2}.*\[open\]'
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
  learnings_open_headings "$root/capture/learnings.md" | learnings_rows \
    | toon_table learnings date,title
}
