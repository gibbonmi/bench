# Review repairs

Blocked by: 07-document-the-bounded-gate-output.md
Writes: internal/gate/report.go, internal/gate/report_test.go, internal/gate/capability_skips.go, internal/gate/capability_skips_test.go, internal/gate/manifest.go, internal/gate/manifest_test.go, internal/gate/runner.go, internal/sanitize/sanitize.go, internal/sanitize/sanitize_test.go, internal/testlines/testlines.go, internal/testlines/testlines_test.go, internal/gate/phase_stream.go, internal/gate/phases_test.go, internal/gate/run_log.go, internal/gate/run_log_prune_test.go, internal/gate/run_stream_test.go

## What to build

The initial `/bench-review-implementation` pass over tickets 01–07 found
four accepted repairs. This ticket retroactively names them, so the
acceptance coverage map's new rows trace to a ticket file. The repairs
already landed on this integration source before this ticket file was
written; it documents the four fixes rather than specifying new work.

**BG01's fix.** The red table's `capability-skips` totals line printed on
every run, red or green. This violated BG01's "and nothing else." The
totals line now prints only on the green branch.

**BG37.** A red phase can carry a `StartErr`: a bad working directory, an
empty argv, an exec failure, or a scheduler deadlock. Such a phase now
renders that error's text as its row, instead of the generic
`exit <code> with no output` text.

**BG38.** A `WARNING: DATA RACE` line now opens a failure block, the same
way `--- FAIL:` and `# <package>` do. A race report's stack now survives as
rows ahead of its `--- FAIL:` block.

**BG39.** A manifest phase named `capability` is now refused before the run
starts. That name is reserved for the skip reporter's filed rows.

Two more repairs closed Standards findings with no new acceptance row.
`sanitize.Strip` composes `toon.Representable` instead of re-deriving its
control-byte threshold. The production wiring line that retains a run's
`.out` stream now has an end-to-end test.

## Acceptance

- [ ] A run with one red phase and zero capability skips prints
      `failures[N]{phase,line}` and `gate: red` on stdout and nothing else. (BG01)
- [ ] A red phase with `StartErr` set renders that error's text as its row. (BG37)
- [ ] A red Go phase whose stream carries a `WARNING: DATA RACE` block
      followed by a `--- FAIL:` block yields rows for both. (BG38)
- [ ] A manifest phase named `capability` is refused before the run starts. (BG39)
