# Measure the run-binary cost before exec reuses the worktree binary

Blocked by: 11-record-and-reuse-the-delta-green.md
Writes: specs/chunk-repair-cycle/assets/run-binary-cost.md (new)
Covers: RC76, RC81

## What to build

Review chunk: RC-C7. This ticket is the measurement that the performance premise of RC-C7 needs before its code change.

Measure three values in the integration worktree and record them, with the exact commands and the sample counts, in `specs/chunk-repair-cycle/assets/run-binary-cost.md`:

- the wall time of `freshness.Digest`, which the exec check of ticket 13 pays once for each `bench` child (`internal/freshness/freshness_verify.go:23-40`);
- the wall time of one run-binary selection that a `bench` child pays today, when exec strips `BENCH_RUN_BINARY` (`internal/worktree/exec.go:214-252`);
- the count of exec calls that select a binary in one test-only repair cycle. Walk the test-only repair sequence that `.agents/commands/bench-implement-spec.md` states at this ticket's base once. List the command of each step, and count each `bench worktree exec` call whose child is `bench`.

Record the net saving for each repair: the selection cost times the call count, less the digest cost times the count of `bench` children. If the net saving is not positive, the build stops RC-C7 here and returns to the spec.

Record the answer to the spec's `uncertain` wrapper question. `bin/bench.sh:298-313` checks an inherited `BENCH_RUN_BINARY` for shape only (absolute, regular, executable, not a symbolic link, physical path) and runs it; the Go child's `Inherit` and `Validate` sequence checks the source digest (`internal/runbinary/runbinary.go:155-183`). `internal/systemtest/owner_environment_test.go:36-39` states that the wrapper's own search runs only when the variable is absent. Confirm this route before ticket 13 starts.

The reviewer dropped the exec reuse after this measurement, so this ticket is the last ticket of RC-C7. RC81 is review-owned: the final reconciliation of the build measures one test-only repair cycle against the 45-to-60-minute baseline and records it in the review record. This ticket owns the row so that one ticket holds it; its author records nothing for it.

The asset follows the precedent of `specs/<slug>/assets/` evidence files. Write it in ASD-STE100 prose, with a table for the measured values.

## Acceptance

- [ ] RC76: the asset records the digest cost, the selection cost, and the exec calls of one repair, with the commands that produced them.
- [ ] The asset states the net saving for each repair and the stop decision.
- [ ] The asset states the wrapper route for an inherited `BENCH_RUN_BINARY`, with the lines that it read.
