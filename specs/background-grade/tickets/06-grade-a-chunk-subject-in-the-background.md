# Grade a chunk subject in the background

Blocked by: 02-derive-the-run-binary-source-once.md, 03-read-the-baseline-members-from-commits.md, 04-add-the-background-width-input.md, 05-record-background-claims-and-outcomes.md
Writes: internal/gate/background_grade.go (new), internal/gate/background_grade_test.go (new), internal/gate/run_transaction.go, internal/gate/checkpoint.go, internal/gate/backgroundstate (new), cmd/bench/main.go, cmd/bench/command_registry.go, cmd/bench/main_test.go, bin/bench.sh, internal/conformance/subcommand_routing_table_test.go, internal/conformance/entry_point_parity_test.go, CHANGELOG.md
Covers: BG1, BG23, BG24, BG26, BG27, BG28, BG29, BG41, BG42, BG106, BG107, BG114

## What to build

Add one plumbing entry that runs one background grade of one worktree, in the foreground of its own process.
Ticket 08 launches it detached. This ticket drives it directly.

For a kit worktree with an open build of a staged spec, the run grades the chunk subject:

- The subject is the HEAD tree of the tip without the review record of the spec, as `executeChunkCheckpoint` computes it (`internal/gate/checkpoint.go:286-306`).
- The chunk id comes from the latest chunk entry of the review record at the tip. When no entry resolves, the run takes the whole schedule. Ticket 07 adds the delta decision.
- The run grades the subject alone, with no evidence pass and no record prose grade, as `gradesSubjectAlone` does today (`internal/gate/checkpoint.go:84-88`).
- The run executes as the executable that the wrapper resolves for the target worktree. It clears any inherited `BENCH_RUN_BINARY`, `BENCH_KIT`, and `BENCH_WRAPPER`, so the wrapper resolves fresh (`bin/bench.sh:293-317`).
- The run sets the width input of ticket 04.

The run holds a claim of ticket 05 on its grade key while it runs.
It writes its green to the evidence store, and it records `green` or `red` in the background record of its worktree.
It never writes a gate cache, and it takes no checkout lock (`internal/gate/run_transaction.go:85-104`).
A wrong chunk id only costs a lookup miss, because the run records its green under the key that it ran.

Register the entry in the internal command inventory.
`cmd/bench/main.go` has 415 lines, which is over its budget, and `cmd/bench/command_registry.go` has 1 line of headroom.
`bin/bench.sh` has 538 lines, so the wrapper reuses an existing route line.
`run_transaction.go` has 27 lines of headroom, which ticket 10 also uses, so put the run in `background_grade.go`.

## Acceptance

- [ ] A chunk checkpoint after a background chunk green of its tip prints the reuse line and runs no oracle.
- [ ] The background green is under the key of the HEAD tree without the review record.
- [ ] With no chunk entry in the record, the run takes the whole schedule.
- [ ] A record whose evidence the checkpoint refuses still yields a background green.
- [ ] The grade key frames the digest of the executable that the wrapper resolves for the target worktree.
- [ ] The run leaves the gate cache byte for byte unchanged.
- [ ] A foreground `bench gate` in the worktree takes the checkout lock while the run holds its claim.
- [ ] A run with a wrong chunk id writes its green only under the key that it ran.
- [ ] A run across an edit of a worktree Go file records its green.
- [ ] A run that ends green records `green`, and a run that ends red records `red`.
- [ ] `bench structure --growth` reports no growth of `cmd/bench/main.go`, `cmd/bench/command_registry.go`, or `bin/bench.sh`.
