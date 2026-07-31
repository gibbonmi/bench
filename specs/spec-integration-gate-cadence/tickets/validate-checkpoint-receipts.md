# Validate checkpoint receipts

Blocked by: Start runs and assign frontier tickets

Ownership fence: `internal/specbuild`
Assumptions: Start and Assign persist the exact candidate and assignment facts

## What to build

Accept one canonical, bounded coordinator receipt only after every charged row,
focused check, independent probe, live assignment tree, ticket digest, ownership
path, and declared assumption agrees. Create and retain one attributed checkpoint
without moving the candidate or calling the gate. Split the existing foundation
tests by durable responsibility so this package returns to its structural baseline.

## Acceptance

- [ ] [R10] A complete receipt creates one attributed checkpoint while candidate and gate-call count stay unchanged.
- [ ] [R11] Delegate-produced, stale, other-assignment, other-tree, and inside-worktree probes refuse before checkpoint mutation.
- [ ] [R12] Passed, already-covered, and not-TDD-able row outcomes remain honest admissible classifications.
- [ ] [R13] Omission or failure of each charged row, focused check, or probe field refuses with state, refs, and worktree bytes unchanged.
- [ ] [R14] Changed base/tree/ticket digest, outside-fence patch, unexplained path, and assumption drift each refuse before commit creation.
- [ ] [R15] Empty, oversized, malformed, unreadable, FIFO, device, socket, regular symlink, and dangling-symlink receipts refuse without blocking.
- [ ] [R54] Final-newline and no-final-newline records have one deterministic documented posture; malformed framing has one stable error.

