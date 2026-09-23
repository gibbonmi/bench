# 9. Recover crashed shifts by lease identity

Blocked by: 4-record-the-shift-boundaries.md, 7-retain-the-shift-memory.md, 8-complete-the-shift-intent-on-exit.md
Writes: internal/shift/recover.go (new), internal/shift/recover_test.go (new), internal/shift/loop.go, internal/intent/intent.go, internal/intent/ledger/ledger.go, internal/sessioninspect/sessioninspect.go, internal/sessioninspect/sessioninspect_test.go, internal/otelrecord/attributes.go, internal/otelrecord/registry.go
Covers: LE64, LE65, LE91, LE66, LE67, LE68, LE69, LE70, LE89, LE71, LE72, LE73, LE74, LE75, LE90, LE76, LE77, LE78, LE79

## What to build

Chunk: LE-C2.

Add a recovery pass to the shift package. It judges each shift entry that has no outcome.

For an entry with a lease, the pass grades the lease file at the entry's worktree without a follow. The grade comes before the comparison, so a malformed line never counts as another owner:

- If the file is absent, the pass abandons the entry.
- If the file is not regular, is unreadable, or holds a malformed lease line, the pass skips the entry.
- If a well-formed line differs from the recorded line plus one newline, the pass abandons the entry.
- If the file holds the recorded line and its owner is alive, the pass skips the entry.
- If the file holds the recorded line and its owner is dead, the pass recovers the entry.

A recovery first retains the memory through ticket 7's store. A worktree with dirty paths beyond the scratch is retained and locked with `worktree.RetainAndLock`, and the entry gets its `worktree:` pointer. A clean worktree loses its scratch and is released. An abandon changes no worktree file, lease, or lock.

For an entry with no lease, add a key parse to the intent package beside `NewEntry`. A dead owner abandons the entry. A live owner or an unparsable key skips it.

The pass sets the outcome to `recovered` or `abandoned`. It writes one `shift.recovery` span with the intent key, the work state, and the cleanup decision. The span also carries the recovery reference and the memory reference. Add the seam to the registry. The pass prints `bench shift recovery: recovered <n>, abandoned <m>` only when it acted.

Run the pass in a new session-inspect phase after the resume phase, and in `bench shift` before its acquire.

## Acceptance

- [ ] After a helper shift gets SIGKILL during its adapter, the pass writes a `shift.recovery` span with work state `recovered` and the crashed span's intent key.
- [ ] The recovered entry carries the outcome `recovered`.
- [ ] After the recovery of a crashed shift whose adapter appended `MEMMARK` to the notes, one memory file holds the notes bytes.
- [ ] A dirty crashed worktree stays locked with its dirty file unchanged, and its entry keeps a `worktree:` pointer and stays live.
- [ ] A crashed worktree with only scratch files loses its lease and scratch files and records cleanup `released`.
- [ ] Another identity in the lease abandons the entry and leaves the lease bytes, worktree files, and lock state unchanged.
- [ ] An absent lease abandons the entry and leaves every worktree file unchanged.
- [ ] A live helper shift keeps its entry unchanged, with no recovery span.
- [ ] An entry with no lease abandons only for a dead key owner, and a live owner or an unparsable key stays unchanged.
- [ ] A malformed lease line and a FIFO lease each leave the entry unchanged, and the pass returns within the test deadline.
- [ ] `Inspect` and `bench shift` each abandon a seeded dead-key entry.
- [ ] One abandon prints exactly `bench shift recovery: recovered 0, abandoned 1`, and a pass that acted on nothing prints nothing.
