# Verify milestone criteria before completion

Blocked by: 07-close-light-delivery.md
Writes: internal/commitment (new), internal/intent, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: DC43, DC44, DC45, DC46, DC70

## What to build

Deliver `commitment verify` and explicit completion approval for the active milestone.
Ticket 07 supplies published delivery facts for both supported deliverable shapes. Ticket 01 supplies exact proposal and approval identity.
Bind every criterion result to its immutable criterion identity and current published revision. Accept only complete native evidence references with explicit reviewer outcome assessment.

A green gate remains necessary delivery evidence, not proof of arbitrary natural-language criteria. Empty roadmap rows cannot complete a milestone.
Reject missing, duplicate, blocked, unknown, stale, or unmet criterion results. Successful completion clears only the completed active milestone.
The next planned milestone stays inactive until a separate explicit activation.

Read the approved criterion and completion contracts plus the commitment command and receipt owner. Ticket 07 is also the serial write blocker for the shared policy files; this ticket adds a separate useful operator command.

## Acceptance

- [ ] Empty roadmap rows without criterion evidence cannot complete M1 (DC43).
- [ ] An unmet criterion refuses completion despite a green gate (DC44).
- [ ] Complete current evidence and explicit approval complete M1 and leave M2 inactive (DC45).
- [ ] Changed criteria or revision invalidate evidence (DC46).
- [ ] Missing, duplicate, blocked, and unknown criterion results each refuse without changing milestone state (DC70).

## Checkpoint verification

Run `bench test --package ./internal/commitment` and `bench test --package ./internal/intent`. Drive verification and approval commands over published-fact fixtures produced by the delivery contract.
