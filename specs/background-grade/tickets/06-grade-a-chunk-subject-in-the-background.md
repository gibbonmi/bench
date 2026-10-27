# Grade a chunk subject in the background

Blocked by: 02-derive-the-run-binary-source-once.md, 03-read-the-baseline-members-from-commits.md, 04-add-the-background-width-input.md, 05-record-background-claims-and-outcomes.md
Writes: internal/gate/background_grade.go (new), internal/gate/background_grade_test.go (new), internal/gate/run_transaction.go, internal/gate/run_durability.go, internal/gate/delta_checkpoint.go, internal/gate/checkpoint.go, internal/gate/reader_hold.go, internal/gate/reader_store.go, internal/gate/backgroundstate (new), cmd/bench/main.go, cmd/bench/command_registry.go, cmd/bench/main_test.go, bin/bench.sh, internal/conformance/subcommand_routing_table_test.go, internal/conformance/entry_point_parity_test.go, CHANGELOG.md
Covers: BG1, BG23, BG24, BG25, BG26, BG27, BG28, BG29, BG30, BG31, BG36, BG37, BG38, BG39, BG40, BG41, BG42

## What to build

Add one plumbing entry that runs one background grade of one worktree, in the foreground of its own process.
Ticket 07 launches it detached. This ticket drives it directly.

For a kit worktree with an open build of a staged spec, the run grades the chunk subject:

- The subject is the HEAD tree of the tip without the review record of the spec, as `executeChunkCheckpoint` computes it (`internal/gate/checkpoint.go:286-306`).
- The chunk id comes from the latest chunk entry of the review record at the tip. When no entry resolves, the run takes the whole schedule.
- With that chunk id, the run takes the run decision of the chunk checkpoint (`internal/gate/delta_checkpoint.go:59-86`). A delta green goes to the delta record of the worktree.
- The run grades the subject alone, with no evidence pass and no record prose grade, as `gradesSubjectAlone` does today (`internal/gate/checkpoint.go:84-88`).
- The run executes as the executable that the wrapper resolves for the target worktree. It clears any inherited `BENCH_RUN_BINARY` and `BENCH_KIT`, so the wrapper resolves fresh (`bin/bench.sh:293-317`).
- The run sets the width input of ticket 04, and it runs at nice 19 and in the idle I/O class.

The run holds a claim of ticket 05 on its grade key while it runs.
It writes four places only: the evidence store, the reader table after a whole green, the delta record, and the background record.
It never writes a reader hold or a gate cache, and it takes no checkout lock (`internal/gate/run_transaction.go:85-104`, `internal/gate/run_transaction.go:241-258`).
A timeout, an interrupt, or an unavailable subject records `aborted`.
A wrong chunk id only costs a lookup miss, because the run records its green under the key that it ran.

Register the entry in the internal command inventory. `cmd/bench/main.go` and `bin/bench.sh` are over their line budget.
So the wrapper reuses an existing route line, and the registry row moves code out of `main.go` in this ticket.

## Acceptance

- [ ] A chunk checkpoint after a background chunk green of its tip prints the reuse line and runs no oracle.
- [ ] The background green is under the key of the HEAD tree without the review record.
- [ ] With no chunk entry in the record, the run takes the whole schedule.
- [ ] A test-only change after a whole green runs the delta set and writes the delta record of the worktree.
- [ ] A record whose evidence the checkpoint refuses still yields a background green.
- [ ] The grade key frames the digest of the executable that the wrapper resolves for the target worktree.
- [ ] The run leaves the gate cache byte for byte unchanged, and a foreground `bench gate` takes the checkout lock during the run.
- [ ] The oracle runs at nice 19 and, on Linux, in the idle I/O class.
- [ ] A red run writes no reader hold, and a whole green writes the reader table.
- [ ] A timeout, an interrupt, and an unavailable subject each record `aborted`.
- [ ] A run with a wrong chunk id writes its green only under the key that it ran.
- [ ] A run across an edit of a worktree Go file records its green.
