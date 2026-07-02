# AXI contracts for the wave-2 query subcommands: `bench diff` (review-base
# resolution + changed-file table) and `bench coverage` (acceptance-coverage-map
# extraction + --check validation). Split from gate-axi-contracts.sh to respect
# the structure budget; the gate sources it in its own shell (shares $root,
# $gate_dir, err(), fail). Spec: specs/second-wave-parsers.md.
[ -f "$root/bin/bench.sh" ] || { err "bench CLI missing (wave-2 AXI contracts skipped)"; return 0 2>/dev/null || exit 0; }

gciw() { git -c user.email=bench@local -c user.name=bench "$@"; }

# ---- diff: recorded base beats merge-base on a stacked branch (story 2) ------
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  printf 'r\n' > README.md; gciw add -A; gciw commit -q -m c1
  gciw branch -m main
  gciw switch -qc feature
  printf 'f\n' > f.txt; gciw add -A; gciw commit -q -m c2
  c2="$(gciw rev-parse HEAD)"
  gciw switch -qc bench/shift-stacked
  gciw config branch.bench/shift-stacked.benchBase "$c2"
  printf 'w\n' > work.txt; gciw add -A; gciw commit -q -m c3
  out="$(bash "$root/bin/bench.sh" diff)"; rc=$?
  [ "$rc" = "0" ] || { echo "diff on a stacked branch did not exit 0 (exit $rc): $out"; exit 1; }
  grep -qxF "base: $c2" <<<"$out" || { echo "diff did not prefer the recorded benchBase (want base $c2): $out"; exit 1; }
  grep -qxF 'method: recorded' <<<"$out" || { echo "diff method did not say recorded: $out"; exit 1; }
  grep -qxF 'files[1]{status,path}:' <<<"$out" || { echo "diff files table not count-1 TOON: $out"; exit 1; }
  grep -qxF '  A,work.txt' <<<"$out" || { echo "diff changed-file row missing: $out"; exit 1; }
  # A recorded key that is reachable but NOT an ancestor of HEAD must fall back
  # loudly — three-dot against a divergent base would silently diff from the
  # merge-base while the preamble claims the recorded sha.
  gciw switch -q main; gciw switch -qc other
  printf 'o\n' > o.txt; gciw add -A; gciw commit -q -m o1
  o1="$(gciw rev-parse HEAD)"
  gciw switch -q bench/shift-stacked
  gciw config branch.bench/shift-stacked.benchBase "$o1"
  out="$(bash "$root/bin/bench.sh" diff)"
  grep -qxF 'method: merge-base (recorded sha not an ancestor)' <<<"$out" \
    || { echo "divergent recorded sha not reported as fallback: $out"; exit 1; }
  grep -qxF "base: $(gciw merge-base main HEAD)" <<<"$out" \
    || { echo "divergent-key fallback base is not the merge-base: $out"; exit 1; }
) || err "AXI diff recorded-base contract failed"
rm -rf "$tmp"

# ---- diff: merge-base fallback, unreachable-sha fallback, spaces, subdir,
# ---- empty diff (stories 2, 3) -----------------------------------------------
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  printf 'r\n' > README.md; gciw add -A; gciw commit -q -m c1
  gciw branch -m main
  c1="$(gciw rev-parse HEAD)"
  gciw switch -qc feature
  mkdir -p sub/deeper
  printf 'f\n' > 'a b.txt'; printf 'u\n' > 'café.txt'; printf 'q\n' > 'a"q.txt'
  gciw add -A; gciw commit -q -m c2
  # No recorded key → merge-base with the default branch.
  out="$(cd sub/deeper && bash "$root/bin/bench.sh" diff)"
  grep -qxF "base: $c1" <<<"$out" || { echo "diff fallback base is not the merge-base with main: $out"; exit 1; }
  grep -qxF 'method: merge-base' <<<"$out" || { echo "diff fallback method not named: $out"; exit 1; }
  grep -qF '  A,a b.txt' <<<"$out" || { echo "path with a space did not round-trip: $out"; exit 1; }
  # Paths git would C-quote (non-ASCII, embedded quote) must arrive raw and be
  # TOON-escaped exactly once — never git-quoted and then quoted again.
  grep -qxF '  A,café.txt' <<<"$out" || { echo "non-ASCII path was git-quoted or mangled: $out"; exit 1; }
  grep -qxF '  A,"a""q.txt"' <<<"$out" || { echo "quote-bearing path not single-layer TOON-escaped: $out"; exit 1; }
  # Recorded key pointing at an unreachable sha → loud fallback.
  gciw config branch.feature.benchBase 0123456789abcdef0123456789abcdef01234567
  out="$(bash "$root/bin/bench.sh" diff)"
  grep -qxF "base: $c1" <<<"$out" || { echo "unreachable-sha fallback base wrong: $out"; exit 1; }
  grep -qxF 'method: merge-base (recorded sha unreachable)' <<<"$out" \
    || { echo "unreachable recorded sha not reported: $out"; exit 1; }
  # A branch with no commits since base → definitive empty table, exit 0.
  gciw switch -q main; gciw switch -qc idle
  out="$(bash "$root/bin/bench.sh" diff)"; rc=$?
  [ "$rc" = "0" ] || { echo "empty diff did not exit 0 (exit $rc)"; exit 1; }
  grep -qxF 'files[0]{status,path}:' <<<"$out" || { echo "empty diff not definitive: $out"; exit 1; }
) || err "AXI diff fallback/shape contract failed"
rm -rf "$tmp"

# ---- diff: error posture (story 4) -------------------------------------------
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"
  out="$(bash "$root/bin/bench.sh" diff 2>/dev/null)" && rc=0 || rc=$?
  [ "$rc" = "1" ] || { echo "diff outside a repo did not exit 1 (exit $rc)"; exit 1; }
  grep -q '^error: ' <<<"$out" || { echo "diff outside a repo not a structured error: $out"; exit 1; }
  git init -q
  printf 'r\n' > README.md; gciw add -A; gciw commit -q -m c1
  gciw branch -m trunk   # no 'main', no origin/HEAD → base unresolvable
  out="$(bash "$root/bin/bench.sh" diff)" && rc=0 || rc=$?
  [ "$rc" = "1" ] || { echo "unresolvable base did not exit 1 (exit $rc): $out"; exit 1; }
  grep -q '^error: cannot resolve a review base' <<<"$out" || { echo "unresolvable base error not named: $out"; exit 1; }
  out="$(bash "$root/bin/bench.sh" diff bogusarg)" && rc=0 || rc=$?
  [ "$rc" = "2" ] || { echo "diff unknown arg exit not 2 (got $rc)"; exit 1; }
  grep -qi 'usage' <<<"$out" || { echo "diff unknown arg did not print usage: $out"; exit 1; }
) || err "AXI diff error-posture contract failed"
rm -rf "$tmp"

# ---- coverage: mapped extraction with hostile cells (story 5) ----------------
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir specs
  cat > specs/t.md <<'SPEC'
# t

## User stories
1. As a, I want b, so c.
2. As d, I want e, so f.
3. As g, I want h, so i.

## Testing decisions

### Acceptance coverage map
| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 2–3 | does x \| y | cli seam | cmd fails, loudly | catches z |
| edge of 1 | edge case | gate | already covered | catches w |

## Out of scope
SPEC
  out="$(bash "$root/bin/bench.sh" coverage specs/t.md)"; rc=$?
  [ "$rc" = "0" ] || { echo "coverage on a mapped spec did not exit 0 (exit $rc): $out"; exit 1; }
  grep -qxF 'spec: specs/t.md' <<<"$out" || { echo "coverage spec preamble missing: $out"; exit 1; }
  grep -qxF 'state: mapped' <<<"$out" || { echo "coverage state not mapped: $out"; exit 1; }
  grep -qxF 'rows[2]{story,seam,red_signal}:' <<<"$out" || { echo "coverage rows header wrong: $out"; exit 1; }
  grep -qxF '  2–3,cli seam,"cmd fails, loudly"' <<<"$out" || { echo "range/comma row misparsed: $out"; exit 1; }
  grep -qxF '  edge of 1,gate,already covered' <<<"$out" || { echo "edge row missing: $out"; exit 1; }
  # CRLF endings must parse the same (hand-edited tables).
  sed -i 's/$/\r/' specs/t.md
  out="$(bash "$root/bin/bench.sh" coverage specs/t.md)"; rc=$?
  [ "$rc" = "0" ] || { echo "CRLF spec did not exit 0 (exit $rc): $out"; exit 1; }
  grep -qxF 'state: mapped' <<<"$out" || { echo "CRLF spec state not mapped: $out"; exit 1; }
  grep -qxF 'rows[2]{story,seam,red_signal}:' <<<"$out" || { echo "CRLF rows misparsed: $out"; exit 1; }
) || err "AXI coverage extraction contract failed"
rm -rf "$tmp"

# ---- coverage: historical / no-map / error posture (stories 5, 6) ------------
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir specs
  printf '# h\n\n<!-- coverage-map: historical -->\n\n### Acceptance coverage map\n| a |\n' > specs/h.md
  out="$(bash "$root/bin/bench.sh" coverage specs/h.md)"; rc=$?
  [ "$rc" = "0" ] || { echo "historical spec did not exit 0 (exit $rc): $out"; exit 1; }
  grep -qxF 'state: historical' <<<"$out" || { echo "historical state missing: $out"; exit 1; }
  grep -qxF 'rows[0]{story,seam,red_signal}:' <<<"$out" || { echo "historical rows not definitive empty: $out"; exit 1; }
  printf '# n\nno map here\n' > specs/n.md
  out="$(bash "$root/bin/bench.sh" coverage specs/n.md)"; rc=$?
  [ "$rc" = "0" ] || { echo "no-map spec did not exit 0 (exit $rc): $out"; exit 1; }
  grep -qxF 'state: no-map' <<<"$out" || { echo "no-map state missing: $out"; exit 1; }
  out="$(bash "$root/bin/bench.sh" coverage)" && rc=0 || rc=$?
  [ "$rc" = "2" ] || { echo "missing spec arg exit not 2 (got $rc)"; exit 1; }
  grep -qi 'usage' <<<"$out" || { echo "missing spec arg did not print usage: $out"; exit 1; }
  out="$(bash "$root/bin/bench.sh" coverage --bogus specs/n.md)" && rc=0 || rc=$?
  [ "$rc" = "2" ] || { echo "unknown flag exit not 2 (got $rc)"; exit 1; }
  out="$(bash "$root/bin/bench.sh" coverage specs/absent.md)" && rc=0 || rc=$?
  [ "$rc" = "1" ] || { echo "nonexistent spec exit not 1 (got $rc)"; exit 1; }
  grep -q '^error: ' <<<"$out" || { echo "nonexistent spec not a structured error: $out"; exit 1; }
) || err "AXI coverage state/error contract failed"
rm -rf "$tmp"

# ---- coverage --check: every validation rule survives the port (story 7) -----
tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"; git init -q
  mkdir specs
  stories='## User stories
1. As a, I want b, so c.
2. As d, I want e, so f.
3. As g, I want h, so i.
'
  hdr='| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|'
  # Valid map (incl. list+range story ref) → silent, exit 0.
  printf '# v\n\n%s\n### Acceptance coverage map\n%s\n| 1, 2–3 | b | s | r | w |\n' "$stories" "$hdr" > specs/v.md
  out="$(bash "$root/bin/bench.sh" coverage --check specs/v.md)"; rc=$?
  [ "$rc" = "0" ] || { echo "valid map --check did not exit 0 (exit $rc): $out"; exit 1; }
  [ -z "$out" ] || { echo "valid map --check was not silent: $out"; exit 1; }
  # Historical and no-map → silent, exit 0.
  printf '# h\n<!-- coverage-map: historical -->\n### Acceptance coverage map\n|bad|\n' > specs/h.md
  bash "$root/bin/bench.sh" coverage --check specs/h.md >/dev/null || { echo "historical --check not silent-green"; exit 1; }
  printf '# n\nprose only\n' > specs/n.md
  bash "$root/bin/bench.sh" coverage --check specs/n.md >/dev/null || { echo "no-map --check not silent-green"; exit 1; }
  # Missing canonical header.
  printf '# b\n\n%s\n### Acceptance coverage map\n| a | b |\n|---|---|\n| 1 | x |\n' "$stories" > specs/b1.md
  out="$(bash "$root/bin/bench.sh" coverage --check specs/b1.md)" && rc=0 || rc=$?
  [ "$rc" = "1" ] || { echo "missing header --check exit not 1 (got $rc)"; exit 1; }
  grep -qF 'coverage map missing the canonical header' <<<"$out" || { echo "missing-header phrasing lost: $out"; exit 1; }
  # No data rows.
  printf '# b\n\n%s\n### Acceptance coverage map\n%s\n' "$stories" "$hdr" > specs/b2.md
  out="$(bash "$root/bin/bench.sh" coverage --check specs/b2.md)" && rc=0 || rc=$?
  [ "$rc" = "1" ] || { echo "no-data-rows --check exit not 1 (got $rc)"; exit 1; }
  grep -qF 'coverage map has no data rows' <<<"$out" || { echo "no-data-rows phrasing lost: $out"; exit 1; }
  # Wrong cell count.
  printf '# b\n\n%s\n### Acceptance coverage map\n%s\n| 1 | b | s | r |\n' "$stories" "$hdr" > specs/b3.md
  out="$(bash "$root/bin/bench.sh" coverage --check specs/b3.md)" && rc=0 || rc=$?
  [ "$rc" = "1" ] || { echo "cell-count --check exit not 1 (got $rc)"; exit 1; }
  grep -qF 'coverage map row 1 has 4 cells (want 5)' <<<"$out" || { echo "cell-count phrasing lost: $out"; exit 1; }
  # Empty cell.
  printf '# b\n\n%s\n### Acceptance coverage map\n%s\n| 1 | b |  | r | w |\n' "$stories" "$hdr" > specs/b4.md
  out="$(bash "$root/bin/bench.sh" coverage --check specs/b4.md)" && rc=0 || rc=$?
  [ "$rc" = "1" ] || { echo "empty-cell --check exit not 1 (got $rc)"; exit 1; }
  grep -qF "coverage map row 1 has an empty 'seam' cell" <<<"$out" || { echo "empty-cell phrasing lost: $out"; exit 1; }
  # Story reference beyond the spec's numbering.
  printf '# b\n\n%s\n### Acceptance coverage map\n%s\n| 9 | b | s | r | w |\n' "$stories" "$hdr" > specs/b5.md
  out="$(bash "$root/bin/bench.sh" coverage --check specs/b5.md)" && rc=0 || rc=$?
  [ "$rc" = "1" ] || { echo "out-of-range --check exit not 1 (got $rc)"; exit 1; }
  grep -qF 'references story 9 but the spec numbers only 3' <<<"$out" || { echo "out-of-range phrasing lost: $out"; exit 1; }
  # Unrecognized story reference.
  printf '# b\n\n%s\n### Acceptance coverage map\n%s\n| x | b | s | r | w |\n' "$stories" "$hdr" > specs/b6.md
  out="$(bash "$root/bin/bench.sh" coverage --check specs/b6.md)" && rc=0 || rc=$?
  [ "$rc" = "1" ] || { echo "unrecognized-ref --check exit not 1 (got $rc)"; exit 1; }
  grep -qF "has an unrecognized story reference 'x'" <<<"$out" || { echo "unrecognized-ref phrasing lost: $out"; exit 1; }
) || err "AXI coverage --check validation contract failed"
rm -rf "$tmp"
