# The acceptance-coverage-map parser for bin/bench.sh: `bench coverage <spec>`
# extraction and the `--check` validation mode the gate consumes (spec
# second-wave-parsers story 7 — one parser for the convention). Sourced by
# bin/bench.sh after bench-query.sh; composes its TOON emitter and error helpers.

# ---- acceptance-coverage-map parser ------------------------------------------
# One parser for the convention the gate enforces and the review phase audits
# (decision map second-wave-parsers #7): `bench coverage <spec>` extracts state
# and rows; `--check` is the gate's validation mode. The rules replicate the
# retired embedded validator exactly — canonical heading, five-cell header,
# non-empty cells, story references within `## User stories` numbering,
# `<!-- coverage-map: historical -->` opt-out — and the error phrasings are
# load-bearing: canary EXPECT files match them by substring.
#
# coverage_parse <spec> <mode>: mode=state|rows|check. state prints one word
# (mapped|historical|no-map); rows prints story<TAB>seam<TAB>red_signal per data
# row; check prints one violation message per line (empty output = valid).
coverage_parse() {
  awk -v mode="$2" '
    function trim(s) { sub(/^[ \t]+/, "", s); sub(/[ \t]+$/, "", s); return s }
    { sub(/\r$/, "") }
    /^<!-- coverage-map: historical -->$/ { historical = 1 }
    /^## User stories$/ { in_stories = 1; next }
    in_stories && /^## / { in_stories = 0 }
    in_stories && /^[0-9]+\. / { n = $0; sub(/\..*$/, "", n); if (n + 0 > maxstory) maxstory = n + 0 }
    $0 == "### Acceptance coverage map" && !seen { seen = 1; in_map = 1; next }
    in_map && /^##+ / { in_map = 0 }
    in_map {
      line = trim($0)
      if (line !~ /^\|/) next
      sub(/^\|/, "", line); sub(/\|$/, "", line)
      gsub(/\\\|/, SUBSEP, line)
      ncells = split(line, c, "|")
      allsep = 1
      for (i = 1; i <= ncells; i++) {
        gsub(SUBSEP, "\\|", c[i])
        c[i] = trim(c[i])
        if (c[i] !~ /^-*$/) allsep = 0
      }
      if (!got_header) {
        hdr = tolower(c[1])
        for (i = 2; i <= ncells; i++) hdr = hdr "|" tolower(c[i])
        got_header = 1
        header_ok = (hdr == "story|behavior|seam|red signal|why it catches the failure")
        next
      }
      if (allsep) next
      ndata++
      row_story[ndata] = c[1]; row_seam[ndata] = c[3]; row_red[ndata] = c[4]
      row_cells[ndata] = ncells
      for (i = 1; i <= 5 && i <= ncells; i++) row_all[ndata "," i] = c[i]
    }
    END {
      state = !seen ? "no-map" : (historical ? "historical" : "mapped")
      if (mode == "state") { print state; exit }
      if (state != "mapped") exit
      if (mode == "rows") {
        for (r = 1; r <= ndata; r++) printf "%s\t%s\t%s\n", row_story[r], row_seam[r], row_red[r]
        exit
      }
      # mode == check
      split("story|behavior|seam|red signal|why it catches the failure", fname, "|")
      if (!header_ok) { print "coverage map missing the canonical header"; exit }
      if (ndata == 0) { print "coverage map has no data rows"; exit }
      for (r = 1; r <= ndata; r++) {
        if (row_cells[r] != 5) {
          printf "coverage map row %d has %d cells (want 5)\n", r, row_cells[r]
          continue
        }
        for (i = 1; i <= 5; i++) if (row_all[r "," i] == "")
          printf "coverage map row %d has an empty %s%s%s cell\n", r, "\047", fname[i], "\047"
        story = row_all[r "," 1]
        sub(/[ \t]*\(.*\)$/, "", story)
        if (story == "" || story ~ /^[Ee][Dd][Gg][Ee]/) continue
        if (story ~ /^[0-9]+([ \t]*(–|-)[ \t]*[0-9]+)?([ \t]*,[ \t]*[0-9]+([ \t]*(–|-)[ \t]*[0-9]+)?)*$/) {
          nums = story
          gsub(/[^0-9]+/, " ", nums)
          split(nums, nn, " ")
          for (i in nn) if (nn[i] != "" && nn[i] + 0 > maxstory)
            printf "coverage map row %d references story %s but the spec numbers only %d\n", r, nn[i], maxstory
        } else {
          printf "coverage map row %d has an unrecognized story reference %s%s%s\n", r, "\047", story, "\047"
        }
      }
    }
  ' "$1"
}

coverage() {
  local check=0 spec=""
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --check) check=1 ;;
      -h | --help) echo "usage: bench coverage [--check] <spec.md>"; return 0 ;;
      -*) axi_usage "bench coverage" "$1"; return 2 ;;
      *)
        [[ -z "$spec" ]] || { axi_usage "bench coverage" "$1"; return 2; }
        spec="$1"
        ;;
    esac
    shift
  done
  [[ -n "$spec" ]] || { axi_usage "bench coverage" "<spec.md> is required"; return 2; }
  [[ -f "$spec" ]] || {
    axi_error "spec not found: $spec" "pass a path to a spec markdown file"
    return 1
  }
  if [[ "$check" -eq 1 ]]; then
    local violations
    violations="$(coverage_parse "$spec" check)"
    [[ -z "$violations" ]] && return 0
    while IFS= read -r line; do
      printf 'error: %s %s — fix the map or mark it <!-- coverage-map: historical -->\n' "$spec" "$line"
    done <<<"$violations"
    return 1
  fi
  local state
  state="$(coverage_parse "$spec" state)"
  printf 'spec: %s\n' "$spec"
  printf 'state: %s\n' "$state"
  coverage_parse "$spec" rows | toon_table rows story,seam,red_signal
}
