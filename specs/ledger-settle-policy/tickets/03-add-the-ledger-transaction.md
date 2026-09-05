# Add the ledger transaction

Blocked by: 02-extract-the-leaf-ledger-package.md
Writes: internal/intent/transaction.go (new), internal/intent/transaction_test.go (new), internal/intent/intent.go, internal/intent/intent_test.go
Covers: LS3, LS4, LS9

## What to build

Verify the premise first. Read `Address`, `Read`, `readPath`, `writePath`,
`acquire`, `staleLock`, and `Upsert` in internal/intent/intent.go. Read
`PurgeAssignments` in internal/intent/assignment.go for its tolerant read and
its pre-lock absent-file answer. Read
`TestConcurrentWritersKeepEveryEntryAndStaleLockReclaims` and
`TestLedgerCommonDirectoryAndSchemaUpsert` in internal/intent/intent_test.go.
Read the leaf declarations ticket
02-extract-the-leaf-ledger-package.md moves.

Add one exported transaction in the new file internal/intent/transaction.go.
`internal/intent/intent.go` holds 395 lines against the 400-line budget, so the
transaction takes its own file. Every sibling ticket calls the transaction, so
fix one exported name and one signature here.

The transaction resolves the ledger address and acquires the lock file. It then
reads the ledger in the caller's read mode, calls the caller's closure with the
read ledger, and writes on a reported change. The closure returns the next
ledger, a changed report, and an error. The transaction releases the lock on
every exit, including every error exit. It performs the same atomic replacement
`writePath` performs today.

The transaction takes a read mode. The strict mode is today's `readPath`, which
refuses a ledger it cannot account for. The tolerant mode is today's purge
read. It decodes each assignment as a raw value and drops each record it cannot
decode or validate. It reads the entries and the cleanup receipts separately.
Move the tolerant read out of `PurgeAssignments` into the transaction, and
leave `PurgeAssignments` on its current envelope until ticket
05-migrate-the-assignment-mutators.md moves it.

The transaction takes an optional compensation step. The step runs when the
write fails. A nil step runs nothing, so no caller pays for one it does not
supply.

Move `Upsert` onto the transaction as the first caller. Its exported signature
and its external behavior do not change. Keep the identical-entry early return,
which reports no change and writes nothing.

Write the new tests in internal/intent/transaction_test.go. Name the first
`TestLedgerTransactionRefusalPersistsNothingAndReleasesTheLock`. It drives a
real ledger file with a closure that returns an error. Name the second
`TestLedgerTransactionTerminalWriteFailureKeepsThePreviousBytes`. It reports a
change and makes the ledger directory read-only. Each test observes the ledger
file bytes and the lock file.

## Acceptance

- [ ] With a closure that returns an error, the transaction returns that error.
- [ ] With a closure that returns an error, the ledger file bytes are unchanged and the lock file is absent.
- [ ] With a reported change and a read-only ledger directory, the transaction returns a write error and the ledger file bytes are unchanged.
- [ ] With a reported change and a read-only ledger directory, the lock file is absent after the return.
- [ ] Two concurrent writers each keep their entry, and a lock file a dead process left is reclaimed.
- [ ] `Upsert` of an identical entry writes nothing, and the file modification time is unchanged.
- [ ] The pre-existing `internal/intent` suite passes with its test logic unchanged.
- [ ] Self-probe: write before the closure returns, and report the refusal test red.
