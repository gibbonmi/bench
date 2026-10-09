# Sync the adoption record once per transaction step

Blocked by: none
Writes: internal/adopt/transaction/transaction.go, internal/adopt/transaction/publication.go, internal/adopt/transaction/transaction_test.go
Covers: none

## What to build

`Store.Apply` saves the recovery record once for each transaction step, not once for each entry. Today Apply saves and syncs the full record before each entry's rename, so a transaction with N changed entries saves the record N+2 times. Apply first stages every replacement. Then it records every staged identity in one record save. Then it publishes each replacement and saves the applied state. Thus Apply saves the record three times for each transaction, whatever the entry count.

Recovery does not change. The record that names the staged identities is durable before the first rename. An interruption at any point leaves each destination at its preimage or at its recorded replacement, and `Undo` restores every destination. If staging fails, Apply removes every staged temporary file and changes no destination. `Store.Undo` and the record format stay as they are.

This change is the largest test-time lever for the adopt package: every real Link, Setup, Upgrade, Unlink, and compatibility repair calls Apply.

## Acceptance

- [ ] An Apply over one changed entry and an Apply over many changed entries save the record the same number of times. A test counts the record-directory syncs through the `SyncDirectory` hook and fails when the count grows with the entry count.
- [ ] If a publication fails after some destinations changed, `Undo` with the returned identifier restores every destination to its preimage.
- [ ] If staging fails for a later entry, no destination changes and no staged temporary file stays in any destination directory.
- [ ] The existing transaction tests stay green and unchanged.
- [ ] The author reports the adopt package elapsed time on the base commit and on the ticket tip.
