# Migrate capture transaction file replacement

Blocked by: none
Writes: internal/capturetx/transaction.go, internal/capturetx/replacement_test.go (new)
Covers: D20

## What to build

Use D-A for replace operations in Append, Abort, recovery, and the manifest writer.
Keep their encoders, 0644 mode, rollback authority, and existing transaction outcomes.
Exclusive writeNew and directory promotion remain with their existing owners.

Review chunk: D-B3.
Complete acceptance and landing of durable-file-replacement supplies the shared leaf and its native qualification.
No ticket in this successor starts before that complete prerequisite lands.
Run this checkpoint with successor tickets unbuilt.
Preserve all existing assertions and successful behavior outside the approved delta.

## Acceptance

Drive the capture operation through its real encoder and fail shared directory synchronization after rename.
The operation retains the classified cause and its existing transaction recovery contract.

- [ ] The capture transaction files caller exposes an injected directory-sync failure. (D20).

## Checkpoint verification

- `bench test --package ./internal/capturetx`

## Required omission evidence

Bypass the shared owner in capture replace.
The capture operation failure witness must turn red.

Record the exact failing assertion, restored source, and passing result in the review pickup.
A shared-owner test alone cannot prove its caller or transaction outcome.
