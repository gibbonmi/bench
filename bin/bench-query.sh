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
# A field value containing a comma, double-quote, or newline is double-quoted
# with inner quotes doubled; anything else is emitted verbatim.
toon_escape() {
  local v="$1"
  case "$v" in
    *,* | *'"'* | *$'\n'*)
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

# Count of open learnings — the value status() surfaces (was an inline grep).
learnings_open_count() {
  local file="$1"
  [[ -f "$file" ]] || { echo 0; return 0; }
  grep -cE "$LEARNINGS_OPEN_RE" "$file" 2>/dev/null || true
}

# Parse open headings (`## <date> — <title>  [open]`) into `date<TAB>title` rows.
learnings_rows() {
  local hd date rest title
  while IFS= read -r hd; do
    [[ -n "$hd" ]] || continue
    date="${hd#\#\# }"
    date="${date%% *}"
    rest="${hd#* — }"
    title="${rest%\[open\]*}"
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
# A ticket is unresolved when its body still carries an `— (open` / `— (deferred`
# placeholder or a `GRILL DEFERRED` banner — the same markers status() detects.
# Emits `map<TAB>ticket<TAB>type<TAB>state` per unresolved ticket. The `state`
# field (open|deferred|grill-deferred) is kept because the placeholder kind is
# the actionable distinction for an agent scanning the frontier; `type` carries
# the ticket's `Type:` value (Grill|Research|Prototype).
maps_rows() {
  local root="$1" f base
  [[ -d "$root/decisions" ]] || return 0
  for f in "$root"/decisions/*.md; do
    [[ -f "$f" ]] || continue
    base="$(basename "$f" .md)"
    awk -v map="$base" '
      function flush() {
        if (num != "" && state != "")
          printf "%s\t%s\t%s\t%s\n", map, num, type, state
      }
      /^## #[0-9]+:/ {
        flush()
        num = $0; sub(/^## #/, "", num); sub(/:.*/, "", num)
        type = ""; state = ""
        next
      }
      /^Type:/ { type = $0; sub(/^Type:[ \t]*/, "", type); next }
      num != "" && state == "" && /^— \(open/     { state = "open"; next }
      num != "" && state == "" && /^— \(deferred/ { state = "deferred"; next }
      num != "" && state == "" && /GRILL DEFERRED/ { state = "grill-deferred"; next }
      END { flush() }
    ' "$f"
  done
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

# Emit a `guard<TAB>boundary<TAB>denies` row for one discovered guard, or nothing
# when the guard is informational. $1 = path to run with --describe; $2 = fallback
# display name used when no manifest name is parseable.
guard_row() {
  local path="$1" fallback="$2" out rc name boundary denies
  # Redirect stdin from /dev/null: a script that ignores --describe and reads stdin
  # (a stub, or a hook that never short-circuits) must not hang or swallow ours.
  out="$(bash "$path" --describe </dev/null 2>/dev/null)" && rc=0 || rc=$?
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
    guard_row "$prepush" "pre-push"
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
