# Print a short bitten mutant set

Blocked by: none
Writes: internal/probe/set.go, internal/probe/set_report_test.go (new), internal/probe/set_test.go, internal/probe/set_cleanup_test.go, internal/reviewrecord/recordcmd/excerpt_short_test.go (new)
Covers: BC53, BC54, BC55, BC56, BC57

## What to build

A mutant set in which every row bit and restored, with no preserved copy and no cleanup failure, prints its three tables with `trailing_bytes=0`. It prints neither the baseline report nor a `mutant_report` table.
The predicate is the method `bitten` of the set run, and `answer` calls it on the line `if s.bitten() {`, because the probe-09 command swaps that exact line.
Every other set keeps its current answer: a silent, invalid, or unrestored row, a preserved copy, and a cleanup failure each keep every report.

The short answer stays importable. `bench record excerpt --mutant 1` stages a canned short answer, and `bench record ticket` then writes a probe entry that bit and restored.
The test of the bitten case in `set_test.go` moves to the short answer.

## Acceptance

- [ ] A bitten set prints the three tables, `trailing_bytes=0`, and no `mutant_report` table (BC53).
- [ ] A set with one silent row prints the report of each started mutant (BC54).
- [ ] A set with a failed restore keeps the preserved table (BC55).
- [ ] The short answer stages and records as a probe entry that bit and restored (BC56).
- [ ] A bitten set of four mutants prints at most 10 lines (BC57).
