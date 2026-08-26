# Print the green phase table

Blocked by: 03-bound-and-classify-the-red-failure-rows.md
Writes: internal/gate/report.go (new in 02), internal/gate/report_test.go (new in 02), internal/gate/runner.go, internal/gate/capability_skips.go, internal/gate/capability_skips_test.go

## What to build

A green run prints at most ten rows. The report prints
`phases[N]{phase,verdict,elapsed_ms}` through `toon.Table`, in phase-table
order. The verdict cell holds `green` or `skipped`. A skipped optional phase
leaves the run green. The report then prints one `capability-skips` line,
then `gate: green`.

The per-phase summary lines go away, and `phaseSummary`
goes with them. A green six-phase run therefore prints exactly nine stdout
lines and no `[phase]`-prefixed line.

A run of more than seven phases prints one row `phases: N/N green` and no
table.

This contract crosses into tickets 05 and 06:

- `phaseResult.ElapsedMS int64` carries the elapsed time `schedule` already measures.

`schedule` sets `ElapsedMS` from the same value the `phase.finish` log
record takes, so the table and the progress log cannot disagree.

The skip rows collapse into one line. That line carries the total, both
kinds, and each nonzero class, as in `capability-skips: 6 (capability=6
environment=0; fifo=3 privilege=3)`. Two pins in
`internal/gate/capability_skips_test.go` read the old multi-line shape, so
this ticket rewrites them. The reused-verdict line stays exactly as it is;
ticket 06 pins it.

## Acceptance

- [ ] A green run of six phases prints `phases[6]{phase,verdict,elapsed_ms}`, one `capability-skips` line, and `gate: green`. (BG16)
- [ ] A green run of six phases prints exactly nine stdout lines and no line with a `[phase]` prefix. (BG30)
- [ ] A green run of eight phases prints `phases: 8/8 green` and no table. (BG17)
- [ ] A skipped optional phase shows `skipped` in its verdict cell, and the run stays green. (BG18)
- [ ] Three capability skips in two classes print one `capability-skips` line that names both classes. (BG19)
- [ ] A green run prints no `phase <name>: green` line.
