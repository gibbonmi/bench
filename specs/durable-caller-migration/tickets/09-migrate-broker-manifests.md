# Migrate broker manifest replacement

Blocked by: none
Writes: internal/brokermanifest/brokermanifest.go, internal/brokermanifest/replacement_test.go (new)
Covers: D25

## What to build

Use D-A for the replaceable broker manifest while retaining validation, tab-framed payload, and 0644 mode.
Executable-plus-seal publication and adoption staging remain excluded.

Review chunk: D-C2.
Complete acceptance and landing of durable-file-replacement supplies the shared leaf and its native qualification.
No ticket in this successor starts before that complete prerequisite lands.
Run this checkpoint with successor tickets unbuilt.
Preserve all existing assertions and successful behavior outside the approved delta.

## Acceptance

Write a valid real manifest, then fail directory sync after rename.
The error preserves uncertainty and read-back keeps the exact manifest fields.

- [ ] The broker manifest caller exposes an injected directory-sync failure. (D25).

## Checkpoint verification

- `bench test --package ./internal/brokermanifest`

## Required omission evidence

Bypass the leaf with the former manifest writer.
The manifest composition test must report the injected missing-failure red.

Record the exact failing assertion, restored source, and passing result in the review pickup.
A shared-owner test alone cannot prove its caller or transaction outcome.
