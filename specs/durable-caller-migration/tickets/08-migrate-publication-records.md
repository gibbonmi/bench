# Migrate publication record replacement

Blocked by: none
Writes: internal/publication/record.go, internal/publication/replacement_test.go (new), internal/conformance/injected_ports_registry_test.go
Covers: D24

## What to build

Use D-A for SaveRecord and keep the state-machine lock, transitions, validation, 0644 mode, and encoded newline.
The record write now synchronizes file and directory.
If migration changes an audited port, update its exact registry disposition with a real-producer witness.

Review chunk: D-C1.
Complete acceptance and landing of durable-file-replacement supplies the shared leaf and its native qualification.
No ticket in this successor starts before that complete prerequisite lands.
Run this checkpoint with successor tickets unbuilt.
Preserve all existing assertions and successful behavior outside the approved delta.

## Acceptance

Drive SaveRecord through its caller encoder and inject a directory-sync fault after rename.
The operation returns the classified failure without changing release-state policy.

- [ ] The publication record caller exposes an injected directory-sync failure. (D24).

## Checkpoint verification

- `bench test --package ./internal/publication`
- `bench test --check injected-port-registry`

## Required omission evidence

Restore the former unsynchronized writer.
The publication record junction test must detect the missing failure.

Record the exact failing assertion, restored source, and passing result in the review pickup.
A shared-owner test alone cannot prove its caller or transaction outcome.
