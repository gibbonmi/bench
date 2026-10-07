# Migrate handoff document replacement

Blocked by: none
Writes: internal/handoffdoc/store.go, internal/handoffdoc/replacement_test.go (new)
Covers: D21

## What to build

Use D-A for handoff document replacement with the existing lock, rendered bytes, and 0644 mode.
Directory-open and directory-sync failures now fail the operation instead of being ignored.
Keep caller validation and section ownership.

Review chunk: D-B4.
Complete acceptance and landing of durable-file-replacement supplies the shared leaf and its native qualification.
No ticket in this successor starts before that complete prerequisite lands.
Run this checkpoint with successor tickets unbuilt.
Preserve all existing assertions and successful behavior outside the approved delta.

## Acceptance

Write a real section, then fail directory synchronization after rename.
The returned error retains uncertainty and the section encoder keeps its exact payload.

- [ ] The handoff document caller exposes an injected directory-sync failure. (D21).

## Checkpoint verification

- `bench test --package ./internal/handoffdoc`

## Required omission evidence

Restore ignored directory errors or bypass the leaf.
The section-writing composition test must detect the missing failure.

Record the exact failing assertion, restored source, and passing result in the review pickup.
A shared-owner test alone cannot prove its caller or transaction outcome.
