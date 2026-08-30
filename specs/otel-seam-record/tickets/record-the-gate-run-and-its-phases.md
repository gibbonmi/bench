# 6. Record the gate run and its phases

Blocked by: write-the-start-and-end-lines.md
Line: opus / medium
Rows: OT14, OT15, OT16, OT28
Writes: internal/gate/gate.go, internal/gate/engine.go, internal/gate/runner.go, internal/systemtest/otel_gate_test.go (new)

## What to build

`bench gate` starts a root span at its verb boundary. The root span carries the
subject id, the run's mode, and the run's exit. The subject id groups the
iterations of one subject, so a consumer counts a subject's gate spans.

Each executed phase starts a child span that carries the phase name and the
phase exit, so a reader derives the per-phase time. A skipped phase also starts
a span, and that span names the blocker that skipped it. A cascade skip
therefore stays attributed.

A failed record write never changes the verb's outcome. With an unwritable
record directory, `bench gate` keeps its exit code and reports its own verdict.
The record is evidence, never a condition.

The spans carry only the attributes that ticket 3 declares. This ticket
instruments one verb. It does not edit the spec's acceptance rows, the budget
targets, or the ownership fences.

The system tests run the built binary with `BENCH_HOME` set to a temporary
directory, then read the JSON lines back. The gate runs the `system` phase only
when the graded root is the kit checkout. The ticket-time observation is
therefore the focused hand-run `go test -tags=system ./internal/systemtest`,
with `BENCH_KIT` and `BENCH_RUN_BINARY` set.

## Acceptance

- [ ] OT14: `bench gate` writes a root span that carries the subject id, the run's mode, and the run's exit.
- [ ] OT15: each executed phase writes a phase span with its name and its exit.
- [ ] OT16: a phase that a red need skips writes a span that names its blocker.
- [ ] OT28: `bench gate` keeps its exit code when the record directory is unwritable.
- [ ] every attribute on the gate spans and the phase spans sits inside the declared set.
