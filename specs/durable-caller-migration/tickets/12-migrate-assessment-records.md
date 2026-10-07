# Migrate assessment record replacement

Blocked by: 11-preserve-probe-recovery.md, 08-migrate-publication-records.md
Writes: internal/assessment/store.go, internal/assessment/record_test.go, internal/assessment/replacement_test.go (new), internal/conformance/injected_ports_registry_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: D28

## What to build

Compose the existing assessment FileOps seam into D-A rather than copy the ordered writer.
Keep validation, 0600 mode, indented record bytes, newline, and old failure assertions.
Update only the affected audited port disposition with its real Record producer.

Review chunk: D-C5.
Complete acceptance and landing of durable-file-replacement supplies the shared leaf and its native qualification.
No ticket in this successor starts before that complete prerequisite lands.
Run this checkpoint with successor tickets unbuilt.
Preserve all existing assertions and successful behavior outside the approved delta.

## Acceptance

Record validated assessment input, then fail shared directory sync after rename.
The real Record operation retains the classified cause and unchanged encoded bytes.

- [ ] The assessment record caller exposes an injected directory-sync failure. (D28).

## Checkpoint verification

- `bench test --package ./internal/assessment`
- `bench test --check injected-port-registry`
- `bench test --package ./cmd/bench --run 'CommandRegistry'`
- `bench test --check axi-query-registry`
- `bench test --check subcommand-routing`

## Required omission evidence

Retain the FileOps local replacement copy or bypass shared sync.
The Record failure-composition witness must turn red.

Record the exact failing assertion, restored source, and passing result in the review pickup.
A shared-owner test alone cannot prove its caller or transaction outcome.

The co-owned command registries satisfy the existing package binding.
Keep command membership, routes, help grammar, and declared inventories unchanged.
Modify a registry only when the approved caller change requires its exact disposition.
