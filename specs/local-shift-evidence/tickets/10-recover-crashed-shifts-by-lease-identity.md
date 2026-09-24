# 10. Recover crashed shifts by lease identity

Blocked by: 7-retain-the-shift-memory.md, 8-complete-the-shift-intent-on-exit.md, 9-abandon-lease-less-stale-intent.md
Writes: internal/shift/recover.go (new), internal/shift/recover_test.go (new), internal/shift/fault.go, internal/shift/fault_test.go, internal/shift/record.go, internal/shift/record_test.go, internal/shift/pass_test.go, internal/shift/loop.go, internal/worktree/snapshot.go, internal/worktree/snapshot_test.go, internal/worktree/lifecycle.go, internal/worktree/lifecycle_test.go, internal/worktree/subshell.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: LE64, LE65, LE91, LE105, LE66, LE67, LE68, LE69, LE70, LE89, LE103, LE71, LE75, LE90, LE97, LE98, LE99, LE101, LE100, LE104

## What to build

Chunk: LE-C3.

Give the recovery pass its verdict for an open entry with a lease. The pass grades the lease file at the entry's worktree without a follow. The grade comes before the comparison, so a malformed line never counts as another owner:

- If the file is absent, the pass abandons the entry.
- If the file is not regular, is unreadable, or holds a malformed lease line, the pass skips the entry.
- If a well-formed line differs from the recorded line plus one newline and its owner is dead, the pass abandons the entry.
- If a differing line has a live owner, the pass skips the entry, because a concurrent pass or a new acquirer holds the lease.
- If the file holds the recorded line and its owner is alive, the pass skips the entry.
- If the file holds the recorded line and its owner is dead, the pass starts a recovery.

A recovery first retains the memory through ticket 7's store, before any claim. A concurrent acquire of a clean tree runs `git clean -qfdx`, so the notes must leave the tree first. A concede then leaves an orphan memory file, and the memory prune bounds it.

The recovery then takes the lease through the pool's own takeover protocol. Add an identity claim beside `RetainAndLock` in `internal/worktree/snapshot.go`. It runs `claimAt` with the recorded line as the one accepted judgment. That protocol renames the lease to `.stale.<pid>`, compares the moved bytes, and creates the pass's own lease exclusively.

If the judgment becomes a parameter of `claimAt`, edit all three existing call sites in place. `lifecycle.go` and `subshell.go` are over their line budgets and must not grow. If the claim loses, the pass concedes and changes nothing more.

After a won claim, the pass writes the entry: its own lease line and the outcome `recovered`. A dirty worktree gets the `worktree:` pointer, and a clean worktree gets the recovery `none`. Last, the pass acts as `preserveAndRecover` does, per the reviewer's 2026-09-23 decision:

- A worktree with dirty paths beyond the scratch is locked in place. The lock step is skipped when the worktree is already locked, and the pass removes its own lease.
- A clean worktree loses its scratch and is released to the pool.

After the act, the pass writes its entry once more with the outcome `recovered` and its pointer. So a concurrent abandon that landed between the two writes loses, and a dirty tree keeps its live pointer.

A resume arm judges each entry with the outcome `recovered`. When its lease file holds the entry's lease line with a dead owner, the pass runs only the act again. The resume retains no second memory file and writes one recovery span with no memory keys. Every other `recovered` entry stays as it is: an absent lease, a live owner, another identity, or a malformed or special lease file.

Add two steps to the shift fault seam, one at the claim and one between the entry write and the act. The concede rows fault the claim step in process. The resume rows run the first pass in a helper role that this ticket adds beside ticket 5's role. That role arms the act fault from a test-only variable, runs the pass, and exits. The test waits for it, so the first pass's lease line names a dead pid and not a zombie that `pidAlive` reads as alive. The second pass then runs in the test process.

The worktree claim test uses the existing takeover gap seam, and its test names start with `TestClaimRecordedLease`. `snapshot_test.go` routes its repository and process effects through the worktree journey harness. The killed-helper rows re-exec the shift test binary into ticket 5's helper role and kill that process with SIGKILL. The dead-other-owner rows rewrite the lease to a pid that the test started and reaped.

The command registry and its conformance tests join the Writes line through the binding closure only, and the build expects no edit there.

## Acceptance

- [ ] After a helper shift gets SIGKILL during its adapter, the pass writes a `shift.recovery` span with work state `recovered` and the crashed span's intent key.
- [ ] The recovered entry carries the outcome `recovered`, and one memory file holds the crashed shift's notes.
- [ ] When the identity claim concedes, one memory file still holds the crashed shift's notes.
- [ ] A dirty crashed worktree stays locked with its dirty file unchanged, and its entry keeps a `worktree:` pointer and stays live.
- [ ] A crashed worktree with only scratch files loses its lease and scratch files and records cleanup `released`.
- [ ] Another identity with a dead owner abandons the entry and leaves the lease bytes, worktree files, and lock state unchanged.
- [ ] Another identity with a live owner leaves the entry unchanged, with no recovery span.
- [ ] An absent lease abandons the entry and leaves every worktree file unchanged.
- [ ] A live helper shift keeps its entry unchanged, with no recovery span.
- [ ] A malformed lease line and a FIFO lease each leave the entry unchanged, and the pass returns within the test deadline.
- [ ] When another writer replaces the lease in the takeover gap, the identity claim returns false and the other writer's lease stays.
- [ ] A pass whose claim step is faulted leaves the entry, the lease, and the lock state unchanged and writes no recovery span.
- [ ] A helper-role pass exits between its entry write and its act, and the test waits for it.
- [ ] After that exit, a second pass locks the dirty worktree and removes the first pass's lease.
- [ ] After that second pass, exactly one memory file exists for the crashed shift.
- [ ] A second pass over a finished recovery changes no file, lease, lock, or entry and writes no recovery span.
- [ ] A `recovered` entry whose lease holds another identity with a dead owner stays unchanged, and the pass changes no worktree file.
