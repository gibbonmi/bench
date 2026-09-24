# 8. Complete the shift intent on exit

Blocked by: 4-record-the-shift-boundaries.md, 5-record-each-pass-under-the-shift.md, 6-record-the-resolved-line.md, 7-retain-the-shift-memory.md
Writes: internal/intent/admissionpolicy/admissionpolicy.go, internal/intent/admissionpolicy/liveness_test.go (new), internal/intent/ledger/ledger.go, internal/intent/ledger_aliases.go, internal/shift/loop.go, internal/shift/result.go, internal/shift/record_test.go (new), internal/shift/pass_test.go (new), internal/shift/session.go, internal/status/status.go, tests/canary/docs-currency-token-diet/signal-vocabulary-drift
Covers: LE60, LE61, LE62, LE63

## What to build

Chunk: LE-C1.

First verify the premise against the tree: `admissionpolicy.Live` keeps an entry while its worktree exists, and a shift's pool worktree persists after release. So a shift that exits normally leaves a live entry today.

Change the liveness rule. An entry that holds an outcome and a recovery of empty or `none` is done, so `Live` drops it and `Compact` removes it. An entry with a `worktree:` recovery stays live while its worktree exists, so the status keeps its recovery pointer. A branch with no commit past the default branch reads as landed, and that reading never retires a recovery entry.

Add an optional `lease` field to the ledger entry. Right after `worktree.Acquire` returns, the shift reads its own lease file and records the line without its final newline. Ticket 10 consumes this field. `admissionpolicy_test.go` is near its line budget, so the new policy tests go in `liveness_test.go`.

Ticket 4 creates `internal/shift/record_test.go`, and this ticket adds its two shift rows there. Ticket 5 also writes that test file and `loop.go`, and tickets 6 and 7 share the pass tests with ticket 5. So the `Blocked by:` line names all three, and a delegated frontier starts this ticket only after the LE-B2 checkpoint.

## Acceptance

- [ ] `Live` drops an entry with the outcome `complete` and the recovery `none` whose worktree exists.
- [ ] `Live` keeps an entry with the outcome `failed` and a `worktree:` recovery whose worktree exists.
- [ ] After a green shift exits, `intent.Snapshot` holds no entry for its key.
- [ ] After the acquire, the shift's ledger entry carries a `lease` value equal to the lease line that the acquire wrote.
