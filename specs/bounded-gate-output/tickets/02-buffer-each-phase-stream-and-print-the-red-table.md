# Buffer each phase stream and print the red table

Blocked by: 01-move-the-runner-line-predicate-into-testlines.md
Writes: internal/gate/runner.go, internal/gate/phase_stream.go (new), internal/gate/phase_stream_test.go (new), internal/gate/report.go (new), internal/gate/report_test.go (new), internal/gate/lane_test.go

## What to build

The gate engine stops relaying phase output. A new phase stream buffer in
`internal/gate/phase_stream.go` replaces the engine's use of
`prefixedPhaseWriters`. The fast lane composes `prefixedPhaseWriters` today
through `laneWriters` in `internal/gate/lane.go` and never reaches the
report, so that function stays exactly as it is. A lane test pins that a red
lane check still prints its own lines.

Each phase's stdout and stderr arrive as lines and land in one per-phase
buffer, in arrival order. A green phase prints nothing. A red phase prints
its failure rows through the report.

This contract crosses into tickets 03, 04, and 05:

- `newPhaseStreams(stderr io.Writer) *phaseStreams` builds the buffer for one run.
- `(*phaseStreams).open(phase Phase) (io.Writer, io.Writer, func())` is the engine's writer factory.
- `(*phaseStreams).lines(phase string) []string` returns one phase's buffered lines.
- `aggregateAndReport(results []phaseResult, cancelled bool, streams *phaseStreams, stdout, stderr io.Writer, redReports ...func() bool) int`.

`schedule` keeps its writer-factory parameter; the engine passes
`streams.open` and the lane passes its own. The buffer holds lines, not
bytes. A last line without a newline flushes as one line at close.

On red the report prints one `failures[N]{phase,line}` table through
`toon.Table`, in phase-table order, then `gate: red` on stdout. A Go test
phase is a phase whose argv starts with `go` and then `test`. Its rows come
from `testlines.FailureRows`. When that returns no row for a red Go phase,
the rows are the phase's last twenty non-empty lines. Every other red phase
yields each non-empty buffered line as a row. A red phase with an empty
buffer yields one row `exit <code> with no output`.

The verdict changes stream. `aggregateAndReport` prints `gate: red` on
stdout, not on stderr. No test pins that line on stderr today, so this
ticket adds the pin. The two operational refusals that print `gate: red`
before any phase runs keep stderr, in `internal/gate/runner.go` and
`internal/gate/phases.go`.

A cancelled run keeps its current shape. It prints the stragglers on stderr
and no verdict on either stream. Ticket 03 adds the cap, the skip rows, and
the control-byte filter.

## Acceptance

- [ ] A run with one red phase prints `failures[N]{phase,line}` and `gate: red` on stdout, and nothing else. (BG01)
- [ ] A run with one green phase and one red phase prints no line from the green phase. (BG02)
- [ ] A red `vet` phase with three non-empty lines yields three rows. (BG06)
- [ ] A red phase with an empty stream prints one row `exit 1 with no output`. (BG08)
- [ ] `gate: red` appears on stdout and not on stderr. (BG11)
- [ ] A cancelled run prints the stragglers on stderr and no verdict on either stream. (BG12)
- [ ] A red stream whose last line has no newline yields that line as a row. (BG15)
- [x] A red Go phase whose stream classifies to no row prints its last
      twenty non-empty lines. (BG32, reworded by ticket 08: a
      `WARNING: DATA RACE` stream now classifies, per BG38, so it no longer
      names this fallback's example fixture.)
- [ ] A lane run with a red prose check prints that check's lines as it does today. (BG29)
- [ ] Two red phases print their rows in phase-table order.
- [ ] A phase whose failure arrives on stderr yields its stderr lines as rows.
