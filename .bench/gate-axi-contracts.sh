# AXI-conformance contracts for the benchkit gate: the agent-facing query
# subcommands (`bench learnings`, `bench maps`) must emit flat-table TOON on
# stdout, give definitive empty states, escape hostile field values, and use
# honest exit codes (0 ok, 2 usage). Sourced by gate.sh so it shares $root,
# $gate_dir, err(), and $fail with the other fragments. Fixtures are throwaway
# repos; the commands read files on disk, so no commits are needed. Run the
# real CLI in a fixture and assert stdout shape + exit code — never internals.
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
