# AXI-conformance contracts for the benchkit gate: guard --describe manifest
# conformance plus the query subcommands `bench learnings` and `bench maps` and
# the shared parsers behind them. The commands must emit flat-table TOON on
# stdout, give definitive empty states, escape hostile field values, and use
# honest exit codes (0 ok, 2 usage). Beyond that baseline this file pins the
# shared parsers against the edges the two-derivations bug class breeds:
#   - the maps parser anchors placeholders / the GRILL DEFERRED banner to line
#     start, skips fenced examples, strips CRLF, and reports a Type-less ticket
#     as `unknown`; maps_unresolved_count gives status the distinct-file figure;
#   - the learnings parser strips the date + separator run (ASCII hyphen or
#     em-dash) and any trailing CR; toon_escape quotes leading/trailing
#     whitespace.
# The `bench guards` aggregation and session-start contracts live in
# gate-axi-guards-contracts.sh. Sourced by gate.sh so it shares $root, err(),
# and $fail with the other fragments; fixture provisioning and cleanup are the
# contract runner's (gate-contract-runner.sh). Run the real CLI in a fixture and
# assert stdout shape + exit code — never internals (the few parser-level checks
# source bin/bench-query.sh directly, which is the seam for a pure helper like
# toon_escape/maps_unresolved_count).

# ---- guard --describe manifest conformance ----------------------------------
# Every guard must answer --describe with the four-key manifest (name, boundary,
# denies, why) and exit 0, so the advertised deny surface cannot drift from the
# enforcement and `bench guards` can aggregate it. This runs BEFORE the CLI early
# return below on purpose: a minimal canary fixture (a single broken guard, no CLI)
# must still trip it. session-start answers too, but must classify as informational
# (`denies: nothing`) — the clause the aggregator uses to exclude it from guard rows.
_axi_guard_manifest() {
  local gp="$1" gname="$2" out rc k
  out="$(bash "$gp" --describe </dev/null 2>/dev/null)" && rc=0 || rc=$?
  if [ "$rc" -ne 0 ]; then
    err "guard $gname --describe did not exit 0 (exit $rc)"
    return
  fi
  for k in name boundary denies why; do
    printf '%s\n' "$out" | grep -qE "^$k: " \
      || err "guard $gname --describe manifest missing $k key"
  done
}
for _g in block-dangerous-git check-agent-line stop session-start; do
  [ -f "$root/.bench/hooks/$_g.sh" ] && _axi_guard_manifest "$root/.bench/hooks/$_g.sh" "$_g"
done
[ -f "$root/.bench/adapters/_line-guard.sh" ] \
  && _axi_guard_manifest "$root/.bench/adapters/_line-guard.sh" "_line-guard"
if [ -f "$root/.bench/hooks/session-start.sh" ]; then
  bash "$root/.bench/hooks/session-start.sh" --describe </dev/null 2>/dev/null \
    | grep -qxF 'denies: nothing (informational)' \
    || err "session-start --describe is not classified informational (denies: nothing)"
fi

[ -f "$root/bin/bench.sh" ] || { err "bench CLI missing (AXI contracts skipped)"; return 0 2>/dev/null || exit 0; }

# A journal with two real entry headings → count-2 TOON table with a row each.
contract "AXI learnings two-entry contract" <<'BODY'
  mkdir -p .bench
  {
    printf '## 2026-01-01 — first learning  [open]\n- body\n'
    printf '## 2026-02-02 — second learning  [open]\n- body\n'
  } > .bench/learnings.md
  out="$(bash "$root/bin/bench.sh" learnings)"
  head -1 <<<"$out" | grep -qxF 'learnings[2]{date,title}:' || { echo "learnings header not count-2 TOON: $(head -1 <<<"$out")"; exit 1; }
  [ "$(grep -cE '^  ' <<<"$out")" = "2" ] || { echo "learnings did not emit two rows"; exit 1; }
  grep -qxF '  2026-01-01,first learning' <<<"$out" || { echo "first learning row missing/misformatted"; exit 1; }
  grep -qxF '  2026-02-02,second learning' <<<"$out" || { echo "second learning row missing"; exit 1; }
BODY

# An absent journal → definitive empty state, exit 0; a template-only journal
# (the scaffold's literal format-example heading) → 0 rows, not a phantom entry.
contract "AXI learnings empty/template contract" <<'BODY'
  out="$(bash "$root/bin/bench.sh" learnings)"; rc=$?
  [ "$rc" = "0" ] || { echo "absent journal did not exit 0 (exit $rc)"; exit 1; }
  grep -qxF 'learnings[0]{date,title}:' <<<"$out" || { echo "absent journal not definitive empty: $out"; exit 1; }
  mkdir -p .bench
  printf '## <date> — <short title>  [open]\n' > .bench/learnings.md
  out="$(bash "$root/bin/bench.sh" learnings)"; rc=$?
  [ "$rc" = "0" ] || { echo "template journal did not exit 0 (exit $rc)"; exit 1; }
  grep -qxF 'learnings[0]{date,title}:' <<<"$out" || { echo "template [open] line counted as an entry: $out"; exit 1; }
BODY

# A map with open, deferred, and grill-deferred tickets → one row per unresolved
# ticket; a resolved ticket is excluded; an absent decisions/ dir → definitive empty.
contract "AXI maps unresolved-ticket contract" <<'BODY'
  out="$(bash "$root/bin/bench.sh" maps)"; rc=$?
  [ "$rc" = "0" ] || { echo "absent decisions dir did not exit 0 (exit $rc)"; exit 1; }
  grep -qxF 'maps[0]{map,ticket,type,state}:' <<<"$out" || { echo "absent decisions not definitive empty: $out"; exit 1; }
  mkdir decisions
  cat > decisions/m.md <<'MAP'
# A map

## #1: first?

Type: Grill

### Answer
— (open)

## #2: second?

Type: Research

### Answer
— (deferred)

## #3: third?

Type: Grill

### Answer
GRILL DEFERRED — waiting on X

## #4: resolved?

Type: Grill

### Answer
Decided: yes, do the thing.
MAP
  out="$(bash "$root/bin/bench.sh" maps)"
  head -1 <<<"$out" | grep -qxF 'maps[3]{map,ticket,type,state}:' || { echo "maps header not count-3 TOON: $(head -1 <<<"$out")"; exit 1; }
  grep -qxF '  m,1,Grill,open' <<<"$out" || { echo "open ticket row missing"; exit 1; }
  grep -qxF '  m,2,Research,deferred' <<<"$out" || { echo "deferred ticket row missing"; exit 1; }
  grep -qxF '  m,3,Grill,grill-deferred' <<<"$out" || { echo "grill-deferred ticket row missing"; exit 1; }
  if grep -qE '^  m,4,' <<<"$out"; then echo "resolved ticket #4 leaked into unresolved list"; exit 1; fi
BODY

# A title carrying a comma and a double-quote is escaped per TOON
# (double-quoted, inner quotes doubled).
contract "AXI TOON field-escaping contract" <<'BODY'
  mkdir -p .bench
  printf '## 2026-03-03 — a, "b"  [open]\n' > .bench/learnings.md
  out="$(bash "$root/bin/bench.sh" learnings)"
  grep -qxF '  2026-03-03,"a, ""b"""' <<<"$out" || { echo "comma/quote title not escaped per TOON: $out"; exit 1; }
BODY

# An unknown argument prints usage on stdout and exits 2 (usage, not error).
contract "AXI usage/exit-2 contract" <<'BODY'
  out="$(bash "$root/bin/bench.sh" learnings bogusarg)"; rc=$?
  [ "$rc" = "2" ] || { echo "learnings unknown arg exit not 2 (got $rc)"; exit 1; }
  grep -qi 'usage' <<<"$out" || { echo "learnings unknown arg did not print usage on stdout: $out"; exit 1; }
  out="$(bash "$root/bin/bench.sh" maps bogusarg)"; rc=$?
  [ "$rc" = "2" ] || { echo "maps unknown arg exit not 2 (got $rc)"; exit 1; }
  grep -qi 'usage' <<<"$out" || { echo "maps unknown arg did not print usage on stdout: $out"; exit 1; }
BODY

# Invoked from a subdirectory, the command still resolves the repo root.
contract "AXI subdirectory root-resolution contract" <<'BODY'
  mkdir -p .bench sub/deeper
  printf '## 2026-04-04 — sub check  [open]\n' > .bench/learnings.md
  out="$(cd sub/deeper && bash "$root/bin/bench.sh" learnings)"
  head -1 <<<"$out" | grep -qxF 'learnings[1]{date,title}:' || { echo "subdir invocation lost root resolution: $out"; exit 1; }
  grep -qxF '  2026-04-04,sub check' <<<"$out" || { echo "subdir learnings row missing"; exit 1; }
BODY

# A fixture repo under a path containing a space works for both commands.
contract "AXI path-with-spaces contract" --space-path <<'BODY'
  mkdir -p .bench decisions
  printf '## 2026-05-05 — spaced  [open]\n' > .bench/learnings.md
  cat > decisions/s.md <<'MAP'
## #1: q?

Type: Grill

### Answer
— (open)
MAP
  lout="$(bash "$root/bin/bench.sh" learnings)"
  grep -qxF '  2026-05-05,spaced' <<<"$lout" || { echo "space-path learnings failed: $lout"; exit 1; }
  mout="$(bash "$root/bin/bench.sh" maps)"
  grep -qxF '  s,1,Grill,open' <<<"$mout" || { echo "space-path maps failed: $mout"; exit 1; }
BODY

# ---- shared-parser hardening (maps/learnings/TOON edges) --------------------
# These pin the shared parsers against the over-match, CRLF, separator, escaping,
# and missing-field defects that the two-derivations bug class breeds.

# Placeholder detection is anchored and fence-aware: a resolved ticket whose
# Answer prose only mentions GRILL DEFERRED mid-line is not listed, a fenced
# `— (open)` example is not listed, and genuine line-start placeholders still are.
contract "AXI maps over-match anchoring/fence contract" <<'BODY'
  mkdir decisions
  cat > decisions/o.md <<'MAP'
# Over-match map

## #1: genuine?

Type: Grill

### Answer
— (open)

## #2: prose mentions?

Type: Grill

### Answer
Decided: a mid-line GRILL DEFERRED mention is not an unresolved banner.

## #3: fenced example?

Type: Grill

### Answer
Decided: the placeholder looks like this:

```
— (open)
```

so authors recognize it.
MAP
  out="$(bash "$root/bin/bench.sh" maps)"
  head -1 <<<"$out" | grep -qxF 'maps[1]{map,ticket,type,state}:' || { echo "over-match: expected exactly one unresolved ticket: $(head -1 <<<"$out")"; exit 1; }
  grep -qxF '  o,1,Grill,open' <<<"$out" || { echo "over-match: genuine line-start placeholder dropped: $out"; exit 1; }
  if grep -qE '^  o,2,' <<<"$out"; then echo "over-match: mid-line GRILL DEFERRED prose leaked"; exit 1; fi
  if grep -qE '^  o,3,' <<<"$out"; then echo "over-match: fenced placeholder example leaked"; exit 1; fi
BODY

# A CRLF map file emits no carriage returns in the TOON rows.
contract "AXI maps CRLF-stripping contract" <<'BODY'
  mkdir decisions
  printf '## #1: q?\r\nType: Grill\r\n### Answer\r\n— (open)\r\n' > decisions/c.md
  out="$(bash "$root/bin/bench.sh" maps)"
  grep -qxF '  c,1,Grill,open' <<<"$out" || { echo "CRLF map row missing/misformatted: $out"; exit 1; }
  [ "$(printf '%s' "$out" | grep -c $'\r')" = "0" ] || { echo "CRLF leaked carriage returns into maps output"; exit 1; }
BODY

# A ticket with no Type: line reports `unknown`, never an empty type field.
contract "AXI maps no-Type-ticket contract" <<'BODY'
  mkdir decisions
  cat > decisions/n.md <<'MAP'
## #1: typeless?

### Answer
— (open)
MAP
  out="$(bash "$root/bin/bench.sh" maps)"
  grep -qxF '  n,1,unknown,open' <<<"$out" || { echo "no-Type ticket did not report unknown: $out"; exit 1; }
BODY

# An ASCII-hyphen (or separator-less) heading yields a clean title with no
# `##`/date prefix, not a title left whole because the split token was the em-dash.
contract "AXI learnings ascii-separator title contract" <<'BODY'
  mkdir -p .bench
  printf '## 2026-01-01 - ascii title  [open]\n' > .bench/learnings.md
  out="$(bash "$root/bin/bench.sh" learnings)"
  grep -qxF '  2026-01-01,ascii title' <<<"$out" || { echo "ascii-hyphen heading title not clean: $out"; exit 1; }
BODY

# toon_escape quotes a field with leading or trailing whitespace so padding
# survives the flat table (a bare padded field is ambiguous to a TOON parser).
contract "AXI TOON leading/trailing-space escaping contract" <<'BODY'
  # shellcheck source=/dev/null
  . "$root/bin/bench-query.sh"
  [ "$(toon_escape ' padded ')" = '" padded "' ] || { echo "leading/trailing-space field not quoted: [$(toon_escape ' padded ')]"; exit 1; }
  [ "$(toon_escape 'plain')" = 'plain' ] || { echo "plain field wrongly quoted"; exit 1; }
BODY

# status derives its unresolved-maps figure from the shared parser:
# maps_unresolved_count returns DISTINCT MAP FILES (what status counts), not the
# ticket count maps_rows lists. Two unresolved tickets in one file → count 1.
contract "AXI maps_unresolved_count distinct-file contract" <<'BODY'
  # shellcheck source=/dev/null
  . "$root/bin/bench-query.sh"
  mkdir decisions
  cat > decisions/multi.md <<'MAP'
## #1: a?

Type: Grill

### Answer
— (open)

## #2: b?

Type: Grill

### Answer
— (deferred)
MAP
  cat > decisions/solo.md <<'MAP'
## #1: c?

Type: Grill

### Answer
— (open)
MAP
  [ "$(maps_rows "$PWD" | grep -c .)" = "3" ] || { echo "expected 3 unresolved tickets, got $(maps_rows "$PWD" | grep -c .)"; exit 1; }
  [ "$(maps_unresolved_count "$PWD")" = "2" ] || { echo "distinct-file count not 2 (got $(maps_unresolved_count "$PWD"))"; exit 1; }
BODY
