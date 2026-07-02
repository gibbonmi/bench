# AXI-conformance contracts for the benchkit gate. The agent-facing query
# subcommands (`bench learnings`, `bench maps`, `bench guards`, plus `guards
# --brief`) must emit flat-table TOON on stdout, give definitive empty states,
# escape hostile field values, and use honest exit codes (0 ok, 2 usage). Beyond
# that baseline this file pins the shared parsers against the edges the
# two-derivations bug class breeds:
#   - guard --describe manifests (all hook guards + the generated pre-push) carry
#     the four-key manifest; session-start classifies informational;
#   - `bench guards` bounds each guard's --describe so a hanging hook cannot stall
#     aggregation (reports `no manifest (timed out)`), never executes an unmanaged
#     (marker-less) pre-push (`unmanaged (no manifest)`), and block-dangerous-git
#     degrades to `manifest unavailable (python3 missing)` when python3 is absent;
#   - the maps parser anchors placeholders / the GRILL DEFERRED banner to line
#     start, skips fenced examples, strips CRLF, and reports a Type-less ticket as
#     `unknown`; maps_unresolved_count gives status the distinct-file figure;
#   - the learnings parser strips the date + separator run (ASCII hyphen or em-dash)
#     and any trailing CR; toon_escape quotes leading/trailing whitespace.
# Sourced by gate.sh so it shares $root, $gate_dir, err(), and $fail with the other
# fragments. Fixtures are throwaway repos; the commands read files on disk, so no
# commits are needed. Run the real CLI in a fixture and assert stdout shape + exit
# code — never internals (the few parser-level checks source bench-query.sh directly,
# which is the seam for a pure helper like toon_escape/maps_unresolved_count).

# ---- guard --describe manifest conformance (stories 5, 6) -------------------
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

# Row 1 — journal with two real entry headings → count-2 TOON table with a row each.
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
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
) || err "AXI learnings two-entry contract failed"
rm -rf "$tmp"

# Row 2 — absent journal → definitive empty exit 0; template-only journal → 0 rows.
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  out="$(bash "$root/bin/bench.sh" learnings)"; rc=$?
  [ "$rc" = "0" ] || { echo "absent journal did not exit 0 (exit $rc)"; exit 1; }
  grep -qxF 'learnings[0]{date,title}:' <<<"$out" || { echo "absent journal not definitive empty: $out"; exit 1; }
  # The scaffold's literal format-example heading must not count as an open entry.
  mkdir -p .bench
  printf '## <date> — <short title>  [open]\n' > .bench/learnings.md
  out="$(bash "$root/bin/bench.sh" learnings)"; rc=$?
  [ "$rc" = "0" ] || { echo "template journal did not exit 0 (exit $rc)"; exit 1; }
  grep -qxF 'learnings[0]{date,title}:' <<<"$out" || { echo "template [open] line counted as an entry: $out"; exit 1; }
) || err "AXI learnings empty/template contract failed"
rm -rf "$tmp"

# Row 3 — maps with open, deferred, and grill-deferred tickets → one row per
# unresolved ticket; a resolved ticket is excluded; absent decisions/ → empty.
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
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
) || err "AXI maps unresolved-ticket contract failed"
rm -rf "$tmp"

# Row 4 — a title carrying a comma and a double-quote is escaped per TOON
# (double-quoted, inner quotes doubled).
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir -p .bench
  printf '## 2026-03-03 — a, "b"  [open]\n' > .bench/learnings.md
  out="$(bash "$root/bin/bench.sh" learnings)"
  grep -qxF '  2026-03-03,"a, ""b"""' <<<"$out" || { echo "comma/quote title not escaped per TOON: $out"; exit 1; }
) || err "AXI TOON field-escaping contract failed"
rm -rf "$tmp"

# Row 5 — an unknown argument prints usage on stdout and exits 2 (usage, not error).
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  out="$(bash "$root/bin/bench.sh" learnings bogusarg)"; rc=$?
  [ "$rc" = "2" ] || { echo "learnings unknown arg exit not 2 (got $rc)"; exit 1; }
  grep -qi 'usage' <<<"$out" || { echo "learnings unknown arg did not print usage on stdout: $out"; exit 1; }
  out="$(bash "$root/bin/bench.sh" maps bogusarg)"; rc=$?
  [ "$rc" = "2" ] || { echo "maps unknown arg exit not 2 (got $rc)"; exit 1; }
  grep -qi 'usage' <<<"$out" || { echo "maps unknown arg did not print usage on stdout: $out"; exit 1; }
) || err "AXI usage/exit-2 contract failed"
rm -rf "$tmp"

# Row 6 — invoked from a subdirectory, the command still resolves the repo root.
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir -p .bench sub/deeper
  printf '## 2026-04-04 — sub check  [open]\n' > .bench/learnings.md
  out="$(cd sub/deeper && bash "$root/bin/bench.sh" learnings)"
  head -1 <<<"$out" | grep -qxF 'learnings[1]{date,title}:' || { echo "subdir invocation lost root resolution: $out"; exit 1; }
  grep -qxF '  2026-04-04,sub check' <<<"$out" || { echo "subdir learnings row missing"; exit 1; }
) || err "AXI subdirectory root-resolution contract failed"
rm -rf "$tmp"

# Row 7 — a fixture repo under a path containing a space works for both commands.
parent="$(mktemp -d)"
tmp="$parent/space dir"
mkdir -p "$tmp"
(
  set -u; cd "$tmp"; git init -q
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
) || err "AXI path-with-spaces contract failed"
rm -rf "$parent"

# ---- shared-parser hardening (maps/learnings/TOON edges) --------------------
# These pin the shared parsers against the over-match, CRLF, separator, escaping,
# and missing-field defects that the two-derivations bug class breeds.

# Fix 3 — placeholder detection is anchored and fence-aware: a resolved ticket whose
# Answer prose only mentions GRILL DEFERRED mid-line is not listed, a fenced
# `— (open)` example is not listed, and genuine line-start placeholders still are.
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
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
) || err "AXI maps over-match anchoring/fence contract failed"
rm -rf "$tmp"

# Fix 4 — a CRLF map file emits no carriage returns in the TOON rows.
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir decisions
  printf '## #1: q?\r\nType: Grill\r\n### Answer\r\n— (open)\r\n' > decisions/c.md
  out="$(bash "$root/bin/bench.sh" maps)"
  grep -qxF '  c,1,Grill,open' <<<"$out" || { echo "CRLF map row missing/misformatted: $out"; exit 1; }
  [ "$(printf '%s' "$out" | grep -c $'\r')" = "0" ] || { echo "CRLF leaked carriage returns into maps output"; exit 1; }
) || err "AXI maps CRLF-stripping contract failed"
rm -rf "$tmp"

# Fix 8 — a ticket with no Type: line reports `unknown`, never an empty type field.
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir decisions
  cat > decisions/n.md <<'MAP'
## #1: typeless?

### Answer
— (open)
MAP
  out="$(bash "$root/bin/bench.sh" maps)"
  grep -qxF '  n,1,unknown,open' <<<"$out" || { echo "no-Type ticket did not report unknown: $out"; exit 1; }
) || err "AXI maps no-Type-ticket contract failed"
rm -rf "$tmp"

# Fix 5 — an ASCII-hyphen (or separator-less) heading yields a clean title with no
# `##`/date prefix, not a title left whole because the split token was the em-dash.
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir -p .bench
  printf '## 2026-01-01 - ascii title  [open]\n' > .bench/learnings.md
  out="$(bash "$root/bin/bench.sh" learnings)"
  grep -qxF '  2026-01-01,ascii title' <<<"$out" || { echo "ascii-hyphen heading title not clean: $out"; exit 1; }
) || err "AXI learnings ascii-separator title contract failed"
rm -rf "$tmp"

# Fix 6 — toon_escape quotes a field with leading or trailing whitespace so padding
# survives the flat table (a bare padded field is ambiguous to a TOON parser).
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"
  # shellcheck source=/dev/null
  . "$root/bin/bench-query.sh"
  [ "$(toon_escape ' padded ')" = '" padded "' ] || { echo "leading/trailing-space field not quoted: [$(toon_escape ' padded ')]"; exit 1; }
  [ "$(toon_escape 'plain')" = 'plain' ] || { echo "plain field wrongly quoted"; exit 1; }
) || err "AXI TOON leading/trailing-space escaping contract failed"
rm -rf "$tmp"

# Fix 2 — status derives its unresolved-maps figure from the shared parser:
# maps_unresolved_count returns DISTINCT MAP FILES (what status counts), not the
# ticket count maps_rows lists. Two unresolved tickets in one file → count 1.
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
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
) || err "AXI maps_unresolved_count distinct-file contract failed"
rm -rf "$tmp"

# ---- bench guards aggregation (story 3) -------------------------------------
# A linked fixture carries all five deny-capable guards; `bench guards` emits one
# TOON row each, excludes the informational session-start hook, reports a stub that
# ignores --describe as `no manifest`, and an absent pre-push as `not installed`.
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  bash "$root/bin/bench.sh" link >/dev/null 2>&1 || { echo "link failed setting up guards fixture"; exit 1; }
  out="$(bash "$root/bin/bench.sh" guards)"
  head -1 <<<"$out" | grep -qxF 'guards[5]{guard,boundary,denies}:' || { echo "guards header not count-5 TOON: $(head -1 <<<"$out")"; exit 1; }
  for g in block-dangerous-git check-agent-line stop _line-guard pre-push; do
    grep -qE "^  $g," <<<"$out" || { echo "guards missing deny-capable row: $g"; exit 1; }
  done
  if grep -qE '^  session-start,' <<<"$out"; then echo "informational session-start leaked into guard rows"; exit 1; fi
  # The generated pre-push hook (story 6) answers --describe with the four keys and
  # exits 0, without reading stdin — so the git-layer guard is part of the surface.
  pp="$(git rev-parse --git-path hooks)/pre-push"
  ppm="$(bash "$pp" --describe </dev/null)"; pprc=$?
  [ "$pprc" = "0" ] || { echo "generated pre-push --describe did not exit 0 (exit $pprc)"; exit 1; }
  for k in name boundary denies why; do
    grep -qE "^$k: " <<<"$ppm" || { echo "generated pre-push --describe manifest missing $k key"; exit 1; }
  done
  # An executable hook that ignores --describe is never skipped silently. Capture
  # first, then grep: `bench guards | grep -q` would SIGPIPE the CLI on an early
  # match and fail under pipefail even though the row is present.
  printf '#!/usr/bin/env bash\ncat >/dev/null\nexit 0\n' > .bench/hooks/extra.sh; chmod +x .bench/hooks/extra.sh
  out="$(bash "$root/bin/bench.sh" guards)"
  grep -qxF '  extra,,no manifest' <<<"$out" || { echo "stub hook not reported as no manifest"; exit 1; }
  rm -f .bench/hooks/extra.sh
  # Absent pre-push is a definitive not-installed row, not an omission.
  rm -f "$(git rev-parse --git-path hooks)/pre-push"
  out="$(bash "$root/bin/bench.sh" guards)"
  grep -qxF '  pre-push,,not installed' <<<"$out" || { echo "absent pre-push not reported as not installed"; exit 1; }
) || err "AXI guards aggregation contract failed"
rm -rf "$tmp"

# guards --brief: one plain line per deny-capable guard plus exactly one footer
# pointing at `bench guards` — the surface session-start injects.
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  bash "$root/bin/bench.sh" link >/dev/null 2>&1
  out="$(bash "$root/bin/bench.sh" guards --brief)"
  [ "$(grep -cF 'full manifests: bench guards' <<<"$out")" = "1" ] || { echo "brief footer not exactly one line"; exit 1; }
  [ "$(grep -c . <<<"$out")" = "6" ] || { echo "brief not five guard lines plus one footer (got $(grep -c . <<<"$out"))"; exit 1; }
  grep -qF 'block-dangerous-git: destructive git' <<<"$out" || { echo "brief lost the git guard deny clause"; exit 1; }
) || err "AXI guards --brief contract failed"
rm -rf "$tmp"

# Unknown argument is a usage error (exit 2); the command resolves the repo root
# from a subdirectory; a fixture under a space-containing path works.
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  bash "$root/bin/bench.sh" link >/dev/null 2>&1
  out="$(bash "$root/bin/bench.sh" guards bogusarg)"; rc=$?
  [ "$rc" = "2" ] || { echo "guards unknown arg exit not 2 (got $rc)"; exit 1; }
  grep -qi 'usage' <<<"$out" || { echo "guards unknown arg did not print usage on stdout: $out"; exit 1; }
  mkdir -p sub/deeper
  sout="$(cd sub/deeper && bash "$root/bin/bench.sh" guards)"
  head -1 <<<"$sout" | grep -qF 'guards[' || { echo "guards from a subdirectory lost root resolution: $sout"; exit 1; }
) || err "AXI guards usage/subdirectory contract failed"
rm -rf "$tmp"

parent="$(mktemp -d)"
tmp="$parent/space dir"
mkdir -p "$tmp"
(
  set -u; cd "$tmp"; git init -q
  bash "$root/bin/bench.sh" link >/dev/null 2>&1
  out="$(bash "$root/bin/bench.sh" guards)"
  head -1 <<<"$out" | grep -qF 'guards[' || { echo "space-path guards failed: $out"; exit 1; }
) || err "AXI guards path-with-spaces contract failed"
rm -rf "$parent"

# Fix 1 — a guard whose --describe blocks (here, sleeps well past the bound) must not
# hang aggregation: `bench guards` bounds each --describe and reports the offender as
# `no manifest (timed out)`, returning promptly.
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir -p .bench/hooks
  printf '#!/usr/bin/env bash\nif [ "${1:-}" = "--describe" ]; then sleep 30; fi\nexit 0\n' > .bench/hooks/slow.sh
  chmod +x .bench/hooks/slow.sh
  start=$(date +%s)
  out="$(bash "$root/bin/bench.sh" guards)"
  elapsed=$(( $(date +%s) - start ))
  [ "$elapsed" -lt 10 ] || { echo "guards did not bound a slow --describe (took ${elapsed}s)"; exit 1; }
  grep -qxF '  slow,,no manifest (timed out)' <<<"$out" || { echo "slow guard not reported as timed out: $out"; exit 1; }
) || err "AXI guards --describe timeout-bound contract failed"
rm -rf "$tmp"

# Fix 7 — an unmanaged (foreign) pre-push must NOT be executed by `bench guards`:
# a marker-less pre-push is reported `unmanaged (no manifest)` and never run (proven
# by a sentinel it would touch on any invocation). Absent pre-push stays not installed.
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  sentinel="$tmp/ran-foreign-prepush"
  hooks="$(git rev-parse --git-path hooks)"
  mkdir -p "$hooks"
  printf '#!/usr/bin/env bash\ntouch %q\nexit 1\n' "$sentinel" > "$hooks/pre-push"
  chmod +x "$hooks/pre-push"
  out="$(bash "$root/bin/bench.sh" guards)"
  grep -qxF '  pre-push,,unmanaged (no manifest)' <<<"$out" || { echo "foreign pre-push not reported unmanaged: $out"; exit 1; }
  [ ! -e "$sentinel" ] || { echo "bench guards executed a foreign pre-push"; exit 1; }
) || err "AXI guards unmanaged-pre-push safety contract failed"
rm -rf "$tmp"

# Fix 9 — block-dangerous-git.sh --describe degrades honestly when python3 is
# unreachable: with PATH stripped to a scratch bin holding only bash + coreutils,
# it prints the `manifest unavailable (python3 missing)` denies line and exits 0.
if [ -f "$root/.bench/hooks/block-dangerous-git.sh" ]; then
  tmp="$(mktemp -d)"
  (
    set -u; cd "$tmp"; mkdir sbin
    for t in bash cat printf env sed grep head; do
      p="$(command -v "$t" 2>/dev/null || true)"; [ -n "$p" ] && ln -sf "$p" "sbin/$t"
    done
    out="$(PATH="$tmp/sbin" bash "$root/.bench/hooks/block-dangerous-git.sh" --describe </dev/null 2>/dev/null)"; rc=$?
    [ "$rc" = "0" ] || { echo "python3-missing --describe did not exit 0 (exit $rc)"; exit 1; }
    grep -qF 'manifest unavailable (python3 missing)' <<<"$out" || { echo "python3-missing denies line absent: $out"; exit 1; }
  ) || err "AXI block-dangerous-git python3-missing manifest contract failed"
  rm -rf "$tmp"
fi

# ---- session-start injects the guard brief (story 7) ------------------------
# In a linked repo session-start runs `bench guards --brief` after the dashboard and
# never blocks; outside a repo it prints nothing and exits 0. PATH excludes any global
# bench so bench_cmd resolves only the repo-local CLI (or nothing outside a repo).
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  bash "$root/bin/bench.sh" link >/dev/null 2>&1
  printf '#!/usr/bin/env bash\nexit 0\n' > .bench/gate.sh; chmod +x .bench/gate.sh
  out="$(PATH=/usr/bin:/bin bash .bench/hooks/session-start.sh 2>&1)"; rc=$?
  [ "$rc" = "0" ] || { echo "session-start blocked in a linked repo (exit $rc)"; exit 1; }
  grep -qF 'full manifests: bench guards' <<<"$out" || { echo "session-start did not inject the guard brief"; exit 1; }
  grep -qE '^bench CLI: .*\.bench/bin/bench\.sh \(bench not on PATH; invoke by path\)$' <<<"$out" \
    || { echo "session-start did not advertise the resolved CLI location: $out"; exit 1; }
) || err "session-start guard-brief injection contract failed"
rm -rf "$tmp"

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"
  out="$(PATH=/usr/bin:/bin bash "$root/.bench/hooks/session-start.sh" 2>&1)"; rc=$?
  [ "$rc" = "0" ] || { echo "session-start blocked outside a repo (exit $rc)"; exit 1; }
  [ -z "$out" ] || { echo "session-start printed outside a repo: $out"; exit 1; }
) || err "session-start never-blocks-outside-repo contract failed"
rm -rf "$tmp"
