# 12. Document the local evidence contract

Blocked by: 3-rotate-and-retain-the-record.md, 7-retain-the-shift-memory.md, 9-recover-crashed-shifts-by-lease-identity.md, 10-record-the-worktree-shell-session.md
Writes: DATA_HANDLING.md, tests/canary/data-handling-derivation/undocumented-passlist-var
Covers: LE86, LE87

## What to build

Chunk: LE-D.

Update the seam record entry and the retention section of `DATA_HANDLING.md` to the decided state. State that local records are mutable evidence inputs and not a tamper-proof central audit system. Name the live and sealed segments, the memory files and their mode, and the three bound entries by name, not by value. Name the three resource keys. Keep the claim that a line carries no environment value, because ticket 2 makes it true.

Write in ASD-STE100 prose, and run `bench gate-prose . -- DATA_HANDLING.md` before the commit.

## Acceptance

- [ ] `DATA_HANDLING.md` states that local records are mutable evidence inputs and not a tamper-proof central audit system.
- [ ] `DATA_HANDLING.md` names the sealed segments, the memory files, `bounds.RecordSegmentLimit`, `bounds.RecordSegmentsRetained`, `bounds.RecordMemoryRetained`, and the three resource keys.
- [ ] `bench gate-prose . -- DATA_HANDLING.md` exits 0.
