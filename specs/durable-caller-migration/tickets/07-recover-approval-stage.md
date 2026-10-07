# Recover failed approval stages and retain published receipts

Blocked by: 06-preserve-ledger-effects.md
Writes: internal/commitment/repository/repository.go, internal/commitment/repository/replacement_test.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/commitment/repository/replacement.go (new)
Covers: D23, D46, D47, D48, D49, D50, D51

## What to build

Migrate policy and board replacement through D-A and the real D-B5 ledger transaction seam.
Snapshot both external states and construct restore before the first policy write.
Restore any failed policy or board stage, including initial policy failure after rename.
Restore once after an unpublished ledger failure.
Retain new policy, projected board, and approved receipt together after a published ledger failure.
Capture restore errors without changing the public Compensation callback signature.

Review chunk: D-B6.
Complete acceptance and landing of durable-file-replacement supplies the shared leaf and its native qualification.
No ticket in this successor starts before that complete prerequisite lands.
Run this checkpoint with successor tickets unbuilt.
Preserve all existing assertions and successful behavior outside the approved delta.

## Acceptance

Fail initial policy sync after rename with distinct old policy and board states.
Successful restore returns their original bytes, modes, or absence and leaves the receipt unapproved.
Fail approved-ledger sync separately and retain all three matching new artifacts.
A second restore fault reports its cause and unrestored path without a rollback promise.
Repeat the applicable cases with no board.

- [ ] The commitment repository caller exposes an injected directory-sync failure. (D23).
- [ ] An initial post-rename policy failure restores the original policy state before approval returns. (D46).
- [ ] A post-rename board failure restores both original file states. (D47).
- [ ] A pre-rename approval-ledger failure restores the original policy and board. (D48).
- [ ] A post-rename approval-ledger failure retains the new policy and board with its approved receipt. (D49).
- [ ] A failed stage or compensation restore reports its cause and unrestored path. (D50).
- [ ] Failed initial policy stage in a board-absent repository restores policy without creating a board. (D51).

## Checkpoint verification

- `bench test --package ./internal/commitment/repository`
- `bench test --package ./internal/commitment`
- `bench test --package ./cmd/bench --run 'CommandRegistry'`
- `bench test --check axi-query-registry`
- `bench test --check subcommand-routing`

## Required omission evidence

Create restore only after policy write, compensate a published ledger, restore only one file, or discard restore failure.
Each mutation must fail the corresponding real approval witness.

Record the exact failing assertion, restored source, and passing result in the review pickup.
A shared-owner test alone cannot prove its caller or transaction outcome.

The co-owned command registries satisfy the existing package binding.
Keep command membership, routes, help grammar, and declared inventories unchanged.
Modify a registry only when the approved caller change requires its exact disposition.

## Source headroom

Move snapshot replacement and restore helpers into replacement.go in this same checkpoint.
Keep repository.go within its file limit while preserving complete stage-before-ledger behavior.
The extraction and new behavior land together; a later ticket cannot pay this headroom debt.
