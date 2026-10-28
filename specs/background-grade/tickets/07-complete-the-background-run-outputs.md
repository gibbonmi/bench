# Complete the background run outputs

Blocked by: 06-grade-a-chunk-subject-in-the-background.md
Writes: internal/gate/background_grade.go (new), internal/gate/background_outcome_test.go (new), internal/gate/run_transaction.go, internal/gate/run_durability.go, internal/gate/delta_checkpoint.go, internal/gate/reader_hold.go, internal/gate/reader_store.go, internal/gate/backgroundstate (new), CHANGELOG.md, internal/gate/background_priority_linux.go (new), internal/gate/background_priority_other.go (new)
Covers: BG25, BG30, BG31, BG36, BG37, BG38, BG39, BG40

## What to build

Complete what one background run writes and how it runs.

- With the chunk id of ticket 06, the run takes the run decision of the chunk checkpoint (`internal/gate/delta_checkpoint.go:59-86`). A test-only change after a whole green runs the delta set, and its green goes to the delta record of the worktree.
- A whole green writes the reader table, as any whole green does (`internal/gate/run_transaction.go:256-258`).
- A red run never writes a reader hold (`internal/gate/run_transaction.go:249`, `internal/gate/reader_hold.go:19-33`).
- The oracle runs at nice 19 and in the idle I/O class, or the lowest priority that the platform allows.
- A timeout, an interrupt, and an unavailable subject each record `aborted` in the background record.

The run still writes no gate cache and takes no checkout lock.

## Acceptance

- [ ] A test-only change after a whole green runs the delta set and writes the delta record of the worktree.
- [ ] The oracle of a background run runs at nice 19.
- [ ] On Linux, the oracle of a background run runs in the idle I/O class.
- [ ] A red background landing run writes no reader hold.
- [ ] A whole background green writes the reader table of its tree.
- [ ] A run whose oracle passes the gate timeout records `aborted`.
- [ ] An interrupted run records `aborted`.
- [ ] A run whose subject is unavailable records `aborted` and runs no oracle.
