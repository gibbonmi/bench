# Migrate repair pilot document replacement

Blocked by: none
Writes: internal/repairpilot/command.go, internal/repairpilot/command_test.go, internal/repairpilot/replacement_test.go (new)
Covers: D29

## What to build

Compose the existing repair-pilot replacement seam into D-A.
Keep 0600 mode, document validation, encoded bytes, newline, and all create, short-write, and rename assertions.

Review chunk: D-C6.
Complete acceptance and landing of durable-file-replacement supplies the shared leaf and its native qualification.
No ticket in this successor starts before that complete prerequisite lands.
Run this checkpoint with successor tickets unbuilt.
Preserve all existing assertions and successful behavior outside the approved delta.

## Acceptance

Drive the real document store through the retained failure seam.
A shared directory-sync failure after rename remains a classified failure with the complete new document visible.

- [ ] The repair pilot document caller exposes an injected directory-sync failure. (D29).

## Checkpoint verification

- `bench test --package ./internal/repairpilot`

## Required omission evidence

Restore the local file-only sequence or weaken the short-write assertion.
The corresponding real Store.Replace witness must turn red.

Record the exact failing assertion, restored source, and passing result in the review pickup.
A shared-owner test alone cannot prove its caller or transaction outcome.
