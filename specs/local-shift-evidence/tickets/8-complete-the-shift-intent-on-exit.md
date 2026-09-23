# 8. Complete the shift intent on exit

Blocked by: 4-record-the-shift-boundaries.md
Writes: internal/intent/admissionpolicy/admissionpolicy.go, internal/intent/admissionpolicy/liveness_test.go (new), internal/intent/ledger/ledger.go, internal/shift/loop.go, internal/shift/record_test.go (new)
Covers: LE60, LE61, LE62, LE63

## What to build

Chunk: LE-C1.

First verify the premise against the tree: `admissionpolicy.Live` keeps an entry while its worktree exists, and a shift's pool worktree persists after release. So a shift that exits normally leaves a live entry today.

Change the liveness rule. An entry that holds an outcome and a recovery of empty or `none` is done, so `Live` drops it and `Compact` removes it. An entry with a `worktree:` recovery stays live while its worktree exists, so the status keeps its recovery pointer.

Add an optional `lease` field to the ledger entry. Right after `worktree.Acquire` returns, the shift reads its own lease file and records the line without its final newline. Ticket 10 consumes this field. `admissionpolicy_test.go` is near its line budget, so the new policy tests go in `liveness_test.go`.

Ticket 4 creates `internal/shift/record_test.go`, and this ticket adds its two shift rows there. This ticket and chunk LE-B2 both write `loop.go` and that test file, so the retained author runs this ticket after the LE-B2 checkpoint.

## Acceptance

- [ ] `Live` drops an entry with the outcome `complete` and the recovery `none` whose worktree exists.
- [ ] `Live` keeps an entry with the outcome `failed` and a `worktree:` recovery whose worktree exists.
- [ ] After a green shift exits, `intent.Snapshot` holds no entry for its key.
- [ ] After the acquire, the shift's ledger entry carries a `lease` value equal to the lease line that the acquire wrote.
