# 7. Keep the started phase line after a kill

Blocked by: record-the-gate-run-and-its-phases.md
Line: opus / medium
Rows: OT25
Writes: internal/systemtest/otel_crash_test.go (new), internal/otelrecord/processor.go

## What to build

A system test runs `bench gate` against the built binary with a temporary
`BENCH_HOME`. The test waits for a phase span start line, then kills the
process mid-phase. The record keeps the started phase line, so the recovery
evidence survives the crash. A kill loses at most the line in flight.

This ticket is the FT71 condition. If the test cannot pass, the build stops and
reports a material acceptance shortfall. The build does not land a silent
partial, and it does not weaken the row. FT71 then stays a separate ledger, and
its recommended dependency on this row is withdrawn.

This ticket lands before the lane, commit, landing, worktree, and hook seams. A
failure therefore stops the build before the wide instrumentation lands.

The gate runs the `system` phase only when the graded root is the kit checkout.
The ticket-time observation is therefore the focused hand-run `go test
-tags=system ./internal/systemtest`, with `BENCH_KIT` and `BENCH_RUN_BINARY`
set.

## Acceptance

- [ ] OT25: a gate killed mid-phase leaves the started phase span line in the record.
- [ ] every other line the killed run wrote still parses.
- [ ] the test kills the process by signal, and it reads the record file back from disk.
