# 10. Recover crashed shifts by lease identity

Blocked by: 7-retain-the-shift-memory.md, 8-complete-the-shift-intent-on-exit.md, 9-abandon-lease-less-stale-intent.md
Writes: internal/shift/recover.go (new), internal/shift/recover_test.go (new), internal/shift/fault.go, internal/worktree/snapshot.go, internal/worktree/snapshot_test.go, internal/worktree/lifecycle.go, internal/worktree/subshell.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: LE64, LE65, LE91, LE66, LE67, LE68, LE69, LE70, LE89, LE71, LE75, LE90, LE97, LE98, LE99, LE101, LE100

## What to build

Chunk: LE-C3.

Give the recovery pass its verdict for an entry with a lease. The pass grades the lease file at the entry's worktree without a follow. The grade comes before the comparison, so a malformed line never counts as another owner:

- If the file is absent, the pass abandons the entry.
- If the file is not regular, is unreadable, or holds a malformed lease line, the pass skips the entry.
- If a well-formed line differs from the recorded line plus one newline, the pass abandons the entry.
- If the file holds the recorded line and its owner is alive, the pass skips the entry.
- If the file holds the recorded line and its owner is dead, the pass starts a recovery.

A recovery takes the lease through the pool's own takeover protocol, because a concurrent acquire can claim the same dead lease. Add an identity claim beside `RetainAndLock` in `internal/worktree/snapshot.go`. It runs `claimAt` with the recorded line as the one accepted judgment. That protocol renames the lease to `.stale.<pid>`, compares the moved bytes, and creates the pass's own lease exclusively.

If the judgment becomes a parameter of `claimAt`, edit all three existing call sites in place. `lifecycle.go` and `subshell.go` are over their line budgets and must not grow. If the claim loses, the pass concedes and changes nothing.

After a won claim, the pass retains the memory through ticket 7's store. It then writes the entry: its own lease line, the outcome `recovered`, and the `worktree:` pointer. Last, it acts: it locks the worktree in place, skipping the lock when the worktree is already locked, and removes its own lease. A recovery never releases or resets a worktree, dirty or clean.

A later pass resumes an entry with the outcome `recovered` whose lease file still holds the entry's lease line with a dead owner. It runs only the act again, retains no second memory file, and writes one recovery span with no memory keys. An entry with the outcome `recovered` and no lease file is finished, and the pass leaves it.

Add two steps to the shift fault seam, one at the claim and one between the entry write and the act. The concede row and the resume row use them. The worktree claim test uses the existing takeover gap seam, and its test names start with `TestClaimRecordedLease`. `snapshot_test.go` routes its repository and process effects through the worktree journey harness.

The killed-helper rows re-exec the shift test binary into ticket 5's helper role and kill that process with SIGKILL. The command registry and its conformance tests join the Writes line through the binding closure only, and the build expects no edit there.

## Acceptance

- [ ] After a helper shift gets SIGKILL during its adapter, the pass writes a `shift.recovery` span with work state `recovered` and the crashed span's intent key.
- [ ] The recovered entry carries the outcome `recovered`, and one memory file holds the crashed shift's notes.
- [ ] A dirty crashed worktree stays locked with its dirty file unchanged.
- [ ] A crashed worktree with only scratch files is also locked in place, loses its lease, and records cleanup `retained`.
- [ ] Each recovered entry keeps a `worktree:` pointer and stays live.
- [ ] Another identity in the lease abandons the entry and leaves the lease bytes, worktree files, and lock state unchanged.
- [ ] An absent lease abandons the entry and leaves every worktree file unchanged.
- [ ] A live helper shift keeps its entry unchanged, with no recovery span.
- [ ] A malformed lease line and a FIFO lease each leave the entry unchanged, and the pass returns within the test deadline.
- [ ] When another writer replaces the lease in the takeover gap, the identity claim returns false and the other writer's lease stays.
- [ ] A pass whose claim step is faulted leaves the entry, the lease, and the lock state unchanged and writes no recovery span.
- [ ] After a fault stops a pass between its entry write and its act, a second pass locks the worktree and removes the pass's lease.
- [ ] After that second pass, exactly one memory file exists for the crashed shift.
- [ ] A second pass over a finished recovery changes no file, lease, lock, or entry and writes no recovery span.
