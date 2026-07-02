# Agent-facing query subcommands for bench.sh (AXI-conformant): a shared
# flat-table TOON emitter plus `bench learnings` and `bench maps`. These are the
# machine-queryable views of state bench already computes, so agents stop
# re-deriving it with hand-assembled greps. Sourced by bin/bench.sh like the
# other siblings; all functions are defined before the dispatcher runs, so the
# parsers shared with bench-status.sh's status() are available at call time.
#
# AXI contract (see the craft-cli skill): TOON on stdout, minimal schemas,
# definitive empty states, structured stdout errors + exit 1 for real failures,
# usage on stdout + exit 2 for unknown arguments. All three are read-only.

# ---- shared TOON emitter (flat tables only) ---------------------------------
# A field value containing a comma, double-quote, or newline — or one carrying
# leading or trailing whitespace, which a bare field would lose on parse — is
# double-quoted with inner quotes doubled; anything else is emitted verbatim.
toon_escape() {
  local v="$1"
  case "$v" in
    *,* | *'"'* | *$'\n'* | [[:space:]]* | *[[:space:]])
      v="${v//\"/\"\"}"
      printf '"%s"' "$v"
      ;;
    *) printf '%s' "$v" ;;
  esac
}

# toon_table <name> <fields_csv>
#   Reads records from stdin, one per line, fields separated by a TAB. Emits a
#   flat-table TOON block: `name[N]{fields}:` header (N = record count) then one
#   two-space-indented, comma-joined, escaped row per record. Empty stdin yields
#   the definitive empty table `name[0]{fields}:`. Flat tables only — the general
#   TOON format is intentionally not implemented (see decision #5).
toon_table() {
  local name="$1" fields="$2"
  local -a fnames=() rows=()
  IFS=',' read -r -a fnames <<< "$fields"
  local ncols="${#fnames[@]}"
  local line
  while IFS= read -r line || [[ -n "$line" ]]; do
    rows+=("$line")
  done
  printf '%s[%d]{%s}:\n' "$name" "${#rows[@]}" "$fields"
  local rec rest out i
  if [[ ${#rows[@]} -gt 0 ]]; then
    for rec in "${rows[@]}"; do
      local -a cells=()
      rest="$rec"
      for ((i = 1; i < ncols; i++)); do
        cells+=("${rest%%$'\t'*}")
        rest="${rest#*$'\t'}"
      done
      cells+=("$rest")
      out=""
      for ((i = 0; i < ncols; i++)); do
        [[ $i -gt 0 ]] && out+=","
        out+="$(toon_escape "${cells[$i]}")"
      done
      printf '  %s\n' "$out"
    done
  fi
}

# ---- structured error / usage helpers ---------------------------------------
axi_error() { printf 'error: %s — %s\n' "$1" "$2"; }
axi_usage() { printf 'usage: %s (unknown argument: %s)\n' "$1" "$2"; }

# ---- open-learnings parser (shared with status) -----------------------------
# The single source for "which journal headings are open". status() and
# bench learnings both read through here, so the counter regex lives once — this
# is the two-derivations bug class the state surface exists to end.
LEARNINGS_OPEN_RE='^## [0-9]{4}-[0-9]{2}-[0-9]{2}.*\[open\]'

learnings_open_headings() {
  local file="$1"
  [[ -f "$file" ]] || return 0
  grep -E "$LEARNINGS_OPEN_RE" "$file" 2>/dev/null || true
}

# Count of open learnings — the value status() surfaces, read through the shared
# LEARNINGS_OPEN_RE above so status and `bench learnings` count by one rule.
learnings_open_count() {
  local file="$1"
  [[ -f "$file" ]] || { echo 0; return 0; }
  grep -cE "$LEARNINGS_OPEN_RE" "$file" 2>/dev/null || true
}

# Parse open headings (`## <date> — <title>  [open]`) into `date<TAB>title` rows.
# The title is the heading minus its `## <date>` prefix and the separator run that
# follows — spaces, an ASCII hyphen, or an em-dash — so an ASCII-hyphen or
# separator-less heading yields a clean title instead of leaking the date prefix.
# A trailing CR from a CRLF journal is stripped first (the ${x%$'\r'} posture the
# guards use), so no carriage return rides into a field.
learnings_rows() {
  local hd date title
  while IFS= read -r hd; do
    hd="${hd%$'\r'}"
    [[ -n "$hd" ]] || continue
    date="${hd#\#\# }"
    date="${date%%[[:space:]]*}"
    title="${hd#\#\# }"
    title="${title#"$date"}"
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
  case "${1-}" in
    "") ;;
    -h | --help) echo "usage: bench learnings"; return 0 ;;
    *) axi_usage "bench learnings" "$1"; return 2 ;;
  esac
  local root
  root="$(git rev-parse --show-toplevel 2>/dev/null)" || {
    axi_error "not in a git repository" "run inside a Bench-linked repo"
    return 1
  }
  learnings_open_headings "$root/.bench/learnings.md" | learnings_rows \
    | toon_table learnings date,title
}

# ---- unresolved decision-map tickets ----------------------------------------
# A map is unresolved when its body still carries an `— (open` / `— (deferred`
# placeholder or a `GRILL DEFERRED` banner. The detection rules — strip a trailing
# CR, skip ``` fenced blocks (so a fenced placeholder example is not a real marker),
# and match a placeholder/banner only at line start — live once, in this awk prelude,
# so `bench maps` (which lists tickets) and status (which counts files) can never
# drift on what "unresolved" means. This is the two-derivations bug class the state
# surface exists to end. marker() returns the unresolved kind, or "" for a normal line.
maps_awk_prelude='
  function marker() {
    if ($0 ~ /^— \(open/)       return "open"
    if ($0 ~ /^— \(deferred/)   return "deferred"
    if ($0 ~ /^GRILL DEFERRED/) return "grill-deferred"
    return ""
  }
  { sub(/\r$/, "") }
  substr($0, 1, 3) == "```" { in_fence = !in_fence; next }
  in_fence { next }
'

# Emits `map<TAB>ticket<TAB>type<TAB>state` per unresolved ticket. The `state` field
# (open|deferred|grill-deferred) is kept because the placeholder kind is the actionable
# distinction for an agent scanning the frontier; `type` carries the ticket's `Type:`
# value (Grill|Research|Prototype), or `unknown` when the ticket has no Type line.
maps_rows() {
  local root="$1" f base
  [[ -d "$root/decisions" ]] || return 0
  for f in "$root"/decisions/*.md; do
    [[ -f "$f" ]] || continue
    base="$(basename "$f" .md)"
    awk -v map="$base" "$maps_awk_prelude"'
      function flush(   t) {
        if (num != "" && state != "") {
          t = (type == "" ? "unknown" : type)
          printf "%s\t%s\t%s\t%s\n", map, num, t, state
        }
      }
      /^## #[0-9]+:/ {
        flush()
        num = $0; sub(/^## #/, "", num); sub(/:.*/, "", num)
        type = ""; state = ""
        next
      }
      /^Type:/ { type = $0; sub(/^Type:[ \t]*/, "", type); next }
      num != "" && state == "" { m = marker(); if (m != "") { state = m; next } }
      END { flush() }
    ' "$f"
  done
}

# Count of DISTINCT map FILES that carry at least one unresolved marker — the figure
# `status` surfaces. Shares maps_awk_prelude with maps_rows (one detection core), but
# scans at file scope: a file counts as an unresolved map even if its placeholder is
# not under a `## #` ticket heading, which the per-ticket listing would not emit.
maps_unresolved_count() {
  local root="$1" f n=0
  [[ -d "$root/decisions" ]] || { echo 0; return 0; }
  for f in "$root"/decisions/*.md; do
    [[ -f "$f" ]] || continue
    if awk "$maps_awk_prelude"'
        marker() != "" { found = 1; exit }
        END { exit(found ? 0 : 1) }
      ' "$f"; then
      n=$((n + 1))
    fi
  done
  echo "$n"
}

maps() {
  case "${1-}" in
    "") ;;
    -h | --help) echo "usage: bench maps"; return 0 ;;
    *) axi_usage "bench maps" "$1"; return 2 ;;
  esac
  local root
  root="$(git rev-parse --show-toplevel 2>/dev/null)" || {
    axi_error "not in a git repository" "run inside a Bench-linked repo"
    return 1
  }
  maps_rows "$root" | toon_table maps map,ticket,type,state
}

# ---- guard aggregation ------------------------------------------------------
# A guard is a boundary that can deny. `bench guards` discovers every guard by
# convention — each `.bench/hooks/*.sh`, the adapters' sourced `_line-guard.sh`,
# and the installed git pre-push hook — and reads each one's `--describe` manifest
# so the block surface is learnable without a collision. Manifests are single-
# sourced from the same tables each guard enforces (see decision #6), so the table
# below cannot drift from enforcement.
#
# Discovery rules:
#   - a guard answering `denies: nothing (informational)` (session-start) is not a
#     deny surface and is excluded from the rows;
#   - a script that does not answer --describe (nonzero exit or an unparseable
#     manifest) is never skipped silently — it gets a `no manifest` row;
#   - an absent pre-push hook is a definitive `not installed` row.

# Pull one `key: value` field from a manifest blob ($1) by key ($2); empty if absent.
# Never propagates grep's no-match exit — the caller runs under `set -e` and an empty
# manifest (a stub that ignores --describe) is a valid `no manifest`, not a failure.
guard_manifest_field() {
  printf '%s\n' "$1" | grep -E "^$2: " | head -n1 | sed -E "s/^$2: //" || true
}

# Run a guard's --describe under a time bound so a hook that ignores --describe and
# blocks on stdin (or loops) cannot hang aggregation. Uses coreutils `timeout` when
# present — the same best-effort-optional-tool posture as gate.sh's shellcheck — and a
# background watchdog when it is not. Stdin is /dev/null (a describe reader must not
# swallow ours); returns the guard's exit code, or 124 on timeout.
guard_describe() {
  local path="$1"
  if command -v timeout >/dev/null 2>&1; then
    timeout 5 bash "$path" --describe </dev/null 2>/dev/null
    return $?
  fi
  local tmp rc=0 pid watchdog
  tmp="$(mktemp)"
  bash "$path" --describe </dev/null >"$tmp" 2>/dev/null &
  pid=$!
  ( sleep 5; kill "$pid" 2>/dev/null ) >/dev/null 2>&1 &
  watchdog=$!
  wait "$pid" 2>/dev/null || rc=$?
  # Watchdog already gone ⇒ it fired the kill ⇒ the guard overran ⇒ report a timeout.
  if kill -0 "$watchdog" 2>/dev/null; then
    kill "$watchdog" 2>/dev/null; wait "$watchdog" 2>/dev/null
  else
    rc=124
  fi
  cat "$tmp"; rm -f "$tmp"
  return "$rc"
}

# Emit a `guard<TAB>boundary<TAB>denies` row for one discovered guard, or nothing
# when the guard is informational. $1 = path to run with --describe; $2 = fallback
# display name used when no manifest name is parseable.
guard_row() {
  local path="$1" fallback="$2" out rc name boundary denies
  out="$(guard_describe "$path")" && rc=0 || rc=$?
  if [[ "$rc" -eq 124 ]]; then
    printf '%s\t\tno manifest (timed out)\n' "$fallback"
    return
  fi
  if [[ "$rc" -ne 0 ]]; then
    printf '%s\t\tno manifest\n' "$fallback"
    return
  fi
  name="$(guard_manifest_field "$out" name)"
  boundary="$(guard_manifest_field "$out" boundary)"
  denies="$(guard_manifest_field "$out" denies)"
  # A parseable manifest needs all three table fields; a missing one is treated as
  # "no manifest" rather than a silently thinned row.
  if [[ -z "$name" || -z "$boundary" || -z "$denies" ]]; then
    printf '%s\t\tno manifest\n' "$fallback"
    return
  fi
  [[ "$denies" == "nothing (informational)" ]] && return
  printf '%s\t%s\t%s\n' "$name" "$boundary" "$denies"
}

# Emit guard rows (TAB-separated guard/boundary/denies) for every discovered guard.
guards_rows() {
  local root="$1" f lg hooks_git prepush
  for f in "$root"/.bench/hooks/*.sh; do
    [[ -e "$f" ]] || continue
    guard_row "$f" "$(basename "$f" .sh)"
  done
  lg="$root/.bench/adapters/_line-guard.sh"
  if [[ -f "$lg" ]]; then guard_row "$lg" "_line-guard"; fi
  hooks_git="$(git -C "$root" rev-parse --git-path hooks 2>/dev/null)" || hooks_git=""
  if [[ -n "$hooks_git" && "$hooks_git" != /* ]]; then hooks_git="$root/$hooks_git"; fi
  prepush="$hooks_git/pre-push"
  if [[ -n "$hooks_git" && -f "$prepush" ]]; then
    if grep -q 'bench:managed-pre-push' "$prepush" 2>/dev/null; then
      guard_row "$prepush" "pre-push"
    else
      # A foreign pre-push is never executed for its manifest: running an unknown
      # hook's body just to read --describe is the collision this surface avoids.
      printf 'pre-push\t\tunmanaged (no manifest)\n'
    fi
  else
    printf 'pre-push\t\tnot installed\n'
  fi
}

guards() {
  local brief=0
  case "${1-}" in
    "") ;;
    --brief) brief=1 ;;
    -h | --help) echo "usage: bench guards [--brief]"; return 0 ;;
    *) axi_usage "bench guards" "$1"; return 2 ;;
  esac
  local root
  root="$(git rev-parse --show-toplevel 2>/dev/null)" || {
    axi_error "not in a git repository" "run inside a Bench-linked repo"
    return 1
  }
  if [[ "$brief" -eq 1 ]]; then
    # One plain line per deny-capable guard plus exactly one footer pointing at the
    # full command. This is the surface session-start.sh injects. Split on the first
    # and last TAB by parameter expansion, not `read` — an empty boundary field would
    # collapse under whitespace-IFS and mis-pair the denies clause.
    local line g d
    while IFS= read -r line; do
      [[ -n "$line" ]] || continue
      g="${line%%$'\t'*}"
      d="${line##*$'\t'}"
      printf '%s: %s\n' "$g" "$d"
    done < <(guards_rows "$root")
    printf 'full manifests: bench guards\n'
    return 0
  fi
  guards_rows "$root" | toon_table guards guard,boundary,denies
}

# ---- review-base + changed files --------------------------------------------
# `bench diff` is the single source of review-base truth (decision map
# second-wave-parsers #2/#3): the recorded `branch.<name>.benchBase` key wins
# when it names a reachable commit; otherwise merge-base with the default
# branch. The preamble names which method resolved so a review agent can see a
# fallback happen. Named bench_diff, not diff — the CLI must never shadow the
# system diff inside this script.
bench_diff() {
  case "${1-}" in
    "") ;;
    -h | --help) echo "usage: bench diff"; return 0 ;;
    *) axi_usage "bench diff" "$1"; return 2 ;;
  esac
  git rev-parse --show-toplevel >/dev/null 2>&1 || {
    axi_error "not in a git repository" "run inside a Bench-linked repo"
    return 1
  }
  local branch key base="" method="" def
  branch="$(git symbolic-ref --quiet --short HEAD 2>/dev/null || true)"
  if [[ -n "$branch" ]]; then
    key="$(git config "branch.$branch.benchBase" 2>/dev/null || true)"
    if [[ -n "$key" ]]; then
      if git cat-file -e "$key^{commit}" 2>/dev/null; then
        base="$key" method="recorded"
      else
        method="merge-base (recorded sha unreachable)"
      fi
    fi
  fi
  if [[ -z "$base" ]]; then
    def="$(default_branch)"
    base="$(git merge-base "$def" HEAD 2>/dev/null)" || {
      axi_error "cannot resolve a review base" "no reachable branch.<name>.benchBase and no merge-base with '$def'"
      return 1
    }
    [[ -n "$method" ]] || method="merge-base"
  fi
  printf 'branch: %s\n' "${branch:-(detached)}"
  printf 'base: %s\n' "$base"
  printf 'method: %s\n' "$method"
  # --no-renames keeps every row a flat status/path pair (no two-path R rows);
  # three-dot matches the review phase's convention and equals two-dot here
  # because the resolved base is always an ancestor of HEAD.
  git diff --name-status --no-renames "$base...HEAD" | toon_table files status,path
}

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
