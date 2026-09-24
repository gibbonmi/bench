# 9. Abandon lease-less stale intent

Blocked by: 4-record-the-shift-boundaries.md, 8-complete-the-shift-intent-on-exit.md
Writes: internal/shift/recover.go (new), internal/shift/recover_test.go (new), internal/shift/loop.go, internal/intent/intent.go, internal/intent/ledger/ledger.go, internal/sessioninspect/sessioninspect.go, internal/sessioninspect/sessioninspect_test.go, internal/otelrecord/attributes.go, internal/otelrecord/registry.go, internal/worktree/lifecycle.go, internal/worktree/lifecycle_test.go
Covers: LE72, LE73, LE74, LE76, LE77, LE78, LE79

## What to build

Chunk: LE-C2.

Add the recovery pass to the shift package in a new file, `recover.go`. The pass judges each shift entry that has no outcome. This ticket gives a verdict only for an entry with no lease. It skips each entry with a lease, and ticket 10 adds that verdict.

Add a key parse to the intent package beside `NewEntry`, the one producer of the key format. The parse returns the owner process of a shift key. A dead owner abandons the entry. A live owner or an unparsable key skips it. An abandon changes no worktree file, lease, or lock.

An abandon sets the entry outcome to `abandoned`, the work-state word. It writes one `shift.recovery` span with the intent key, the work state, and the cleanup decision `none`. Add the seam to the registry. Update the ledger comment on the outcome field, because the field now also holds the two recovery words.

Run the pass in a new session-inspect phase after the resume phase, and in `bench shift` before its acquire. `worktree.Acquire` has that one production caller. The pass prints `bench shift recovery: recovered <n>, abandoned <m>` only when it acted.

## Acceptance

- [ ] An entry with no lease and no worktree whose key names a dead process gets the outcome `abandoned` and one `shift.recovery` span.
- [ ] An entry with no lease whose key names the live test process stays unchanged.
- [ ] An entry whose key does not parse as a shift key stays unchanged.
- [ ] `Inspect` runs the pass after the resume phase and abandons a seeded dead-key entry.
- [ ] `bench shift` runs the pass before its acquire and abandons a seeded dead-key entry.
- [ ] One abandon prints exactly `bench shift recovery: recovered 0, abandoned 1`, and a pass that acted on nothing prints nothing.
