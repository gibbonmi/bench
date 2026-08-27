# Print the gate line and log the footprint

Blocked by: 05-hold-the-cache-lock-and-add-bench-cache-clean.md

Writes: internal/gocache/ (new), internal/gate/report.go, internal/gate/runner.go, internal/gate/run_log.go, internal/gate/report_test.go, internal/systemtest/adoption_test.go, CHANGELOG.md

## What to build

Every gate run shows the footprint it just wrote, without a second command. Add
the gate-tail reporter and the run log event. The `clean` child already exists,
so the over-bound line may advertise `bench cache clean` as the next action.

The reporter is one more report closure beside `reportCapabilitySkips` in the
verdict tail. It reads the directory from its own `GOCACHE` entry, walks once,
and answers no rows, one green line, and never red. It removes no file, because
no gate evicts. The production call site is the closure list that
`internal/gate/runner.go` hands to `aggregateAndReport`.

The green line reads
`go-build-cache: <bytes> bytes in <files> files at <dir> (bound <bound> bytes)`
after the phase table and before `gate: green`. Above the bound the parenthesis
reads `(over bound <bound> bytes, next: bench cache clean)`, and the verdict
stays green. The line prints the directory with a control byte stripped. A red
run prints no line at all.

Every run that reaches a verdict, red or green, logs one `cache.footprint`
event through the inherited run log. The event carries the directory, the
bytes, the files, and the over-bound flag. An interrupted run reaches no
verdict, so it logs no event. The adoption journey in
`internal/systemtest/adoption_test.go` runs a real gate, so it observes the
event after one green run.

The hostile-input rows for the walk land here beside the line. A FIFO and a
dangling symlink each add zero bytes, and neither blocks the walk nor raises an
error.

## Acceptance

- [ ] C12 — After one green gate in the adoption journey, the run's `.jsonl` holds one `cache.footprint` event whose `path` equals the derived dir.
- [ ] R08 — A green gate prints the `go-build-cache:` line after the phase table and before `gate: green`.
- [ ] R09 — Above the bound the line's parenthesis reads `(over bound <bound> bytes, next: bench cache clean)`.
- [ ] R10 — Above the bound the run still prints `gate: green`.
- [ ] R11 — A red gate run logs one `cache.footprint` event.
- [ ] R12 — A FIFO inside the directory adds zero bytes and does not block the walk.
- [ ] R13 — A dangling symlink inside the directory adds zero bytes and no error.
- [ ] R15 — The gate line prints the directory with a control byte stripped.
- [ ] R17 — A red gate run prints no `go-build-cache:` line.
- [ ] L09 — A gate reporter given a footprint above the bound removes no file and answers not red.

Delivered outcome: every gate run that reaches a verdict states or logs the
footprint, and an over-bound run names the remedy without grading the tree.
