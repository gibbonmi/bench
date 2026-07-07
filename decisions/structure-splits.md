# Structure Splits (FT32)

## #1: Which flagged files split, and along what seam?

Blocked by: —
Type: Grill

### Question
`bench structure` flags four files. The assessment's verdicts: two genuine,
two budget noise. What are the split seams, and what happens to the noise so
the signal stays credible?

### Answer
Split the two genuine ones: `axi_wave2_test.go` along the command boundary
(one test file per covered command family); `line_routing_checks_test.go`
along the static-parse vs subprocess-exec seam. Pure moves — no test logic
changes, same coverage, gate green before and after. The two noise files
(`runtime_status_test.go`, `status.go`) are barely over budget and cohesive;
fragmenting them would hurt.

## #2: How is an accepted over-budget file recorded?

Blocked by: #1
Type: Grill

### Question
If the two noise files stay flagged, the structure signal shows 2 issues
forever — permanent noise erodes the signal's credibility. Accept-list
mechanism, or live with the noise?

### Answer
Add a small accept mechanism to `bench structure`: a `.bench/structure-accept`
file listing `<path> <one-clause reason>` rows. Accepted files are excluded
from the violation count and shown in a separate `accepted:` section of the
full output, so the acceptance is visible, reasoned, and reversible — not a
silent threshold bump. A row whose path no longer exists is reported stale
(keeps the file honest). Seed it with the two noise files. **Flagged for
veto:** this adds a suppression surface to a signal; the guard is that
acceptance is per-file, reasoned on the page, and the count of accepted files
is printed. Rejected: raising the global budget (hides real debt);
leaving permanent noise (trains readers to ignore the row).

## Handoff

1. **Module boundaries.** `internal/structure` owns the accept file parse and
   reporting; the two test-file splits live where they live today
   (`internal/contract/axi`, `internal/conformance`).
2. **Contracts.** `bench structure`: violation count excludes accepted files;
   output gains `accepted: <path> — <reason>` lines and a stale-row warning;
   exit code semantics unchanged (nonzero only on real violations). Accept
   file: one row per line, `#` comments, missing file = empty accept list.
3. **Deep vs thin.** Accept parsing is thin; the splits are mechanical moves.
4. **Black-box assertables.** Temp-repo fixtures: over-budget file +
   accept row → count drops, accepted line printed; stale row → warning;
   the two split files each land under budget with the full suite still
   green.
5. **Gate attachment.** structure's own tests for the accept mechanism; the
   splits are proven by the unchanged gate.
6. **Hostile-input owners.** structure owns malformed accept rows (no
   reason → treated as malformed, reported, not honored — a reason is the
   price of acceptance) and paths with spaces.
7. **Uncertainty flags.** The exact command-family grouping inside
   axi_wave2_test.go — the delegate picks the natural boundaries and reports
   the resulting file list.
8. **Rejected alternatives.** Global budget raise; permanent noise;
   per-file budget overrides in source comments (invisible to the reader of
   the report).
9. **Domain watch-outs.** The false-empty rule ([FT29]) applies: a failed
   read of the accept file must be loud, not an empty accept list at exit 0
   (a vanished accept file changing counts silently is the same defect
   class).

Dependency order: n/a — single spec.
