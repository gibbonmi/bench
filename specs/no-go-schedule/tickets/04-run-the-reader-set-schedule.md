# Run the reader set schedule and write the reader record

Blocked by: 03-derive-the-reader-table.md
Writes: internal/gate/reader_set.go (new), internal/gate/reader_set_test.go (new), internal/gate/delta_checkpoint.go, internal/gate/report.go, internal/gate/verdict_registry.go, internal/gate/verdict.go, internal/gate/gate.go, internal/gate/verdict_registry_guard_test.go, internal/gate/durable_replacement_test.go, internal/git/git.go
Covers: NG29, NG30, NG31, NG32, NG33, NG34, NG35, NG36

## What to build

The gate owner passes a reader schedule to the phase child in one environment entry, as the delta set passes `BENCH_GATE_DELTA` (`internal/gate/delta_checkpoint.go:29-31`). The schedule names the reader packages and the race flag. `scheduledTable` reads it and builds the reader phases:

- `test` over the reader packages and `./internal/conformance`, each once, with the root-conformance environment
- `race` with its whole table, when the flag is set
- `system`

The reader phases hold no `gofmt`, `vet`, or `shellcheck` phase. A malformed entry refuses before any phase. `greenVerdict` (`internal/gate/report.go:94-101`) prints `gate: green (reader set)` for a reader run.

The reader record is a new class in `verdictRecordClasses`, with its own reuse refusal. It carries the graded tree, the destination tree, the green marker commit, the reader schedule, the run binary, and the record time. The store keeps one record for each graded tree, under a directory name that `internal/git/git.go` declares.

A run function beside `runDelta` runs an eligible reader set. It writes the reader record on green through the durable replacer. It writes no full verdict, no pending record, and no evidence, and a red writes no record. `gate.Result` states the schedule that ran, so a later ticket can attribute it. This ticket adds no caller that selects the reader set; ticket 05 adds it.

Add the class to the class-count pin (`internal/gate/verdict_registry_guard_test.go:30-42`) and the record to the durable-replacement entries (`internal/gate/durable_replacement_test.go:59`).

## Acceptance

- [ ] A phase child with a reader schedule over a fixture module runs `test` over the named packages and `./internal/conformance`, then `system`, and its last line is `gate: green (reader set)`.
- [ ] A reader schedule with the race flag also runs `race` with every race-table sentinel.
- [ ] A reader schedule prints no `gofmt`, `vet`, or `shellcheck` row.
- [ ] A malformed reader schedule refuses before any phase.
- [ ] A green reader run writes one reader record for its tree and leaves the verdict cache and the evidence store byte-identical.
- [ ] A reader record never makes its tree reusable green.
- [ ] A red reader run writes no reader record.
