# Plan and approve exact commitment changes

Blocked by: none
Writes: internal/jsonfile, internal/commitment (new), internal/intent, internal/roadmap, cmd/bench, internal/tickets/registry_data.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/conformance/help_inventory_single_source_test.go
Covers: DC1, DC2, DC3, DC4, DC5, DC7, DC8, DC9, DC55, DC56, DC57, DC58, DC59, DC60, DC62, DC71

## What to build

Deliver `commitment show`, `inventory`, `plan`, and `approve` through the ordinary CLI. A reviewer can inspect and approve one exact policy transition.
Create the commitment owner and its repository adapter. Keep policy decisions independent of their command consumers.
Reuse the existing strict document reader and intent transaction. Use one canonical projection for policy order and approval effects.

Approval binds the predecessor policy, proposed bytes, affected source identities, and exact delayed and removed sets. An identical replay changes nothing.
Use immutable milestone and outcome identities, criterion identities, and source bindings. Reject ambiguous input before any receipt write.
A planned milestone changes no active selection. An unrelated default-branch commit preserves a receipt when all bound facts remain equal.

This checkpoint stages policy and the derived sequence in a planning assignment. Commit and publication consumers arrive in later tickets.
The usable checkpoint is the proposal and approval command journey, including refusal and replay. It must not claim that all delivery routes are enforced yet.
Add the command family to the existing registry and bind its package in the ticket registry.
Register the commitment form projection in the existing help conformance owner.

Read the spec contracts, `internal/jsonfile`, intent transaction and assignment records, roadmap identity and sequence readers, and the command registry. The read surface is these owners, not every command package. Review this owner seam before ticket 02 starts.

## Acceptance

- [ ] Activating M1 displays exactly A then B; planning M2 leaves that selection intact (DC1, DC2).
- [ ] Insertion, removal, and switching report the exact displacement sets; missing operands refuse approval (DC3 to DC5).
- [ ] Changed predecessors refuse, identical approval replays preserve bytes, and unrelated advances remain usable (DC7 to DC9, DC62).
- [ ] Empty, malformed, special-file, control-bearing, and ambiguous operands refuse before state changes (DC55 to DC60).
- [ ] A literal path and a final input line without a newline produce the same valid proposal (DC58, DC71).
- [ ] The new command appears in help, AXI inventory, and routing checks. Its owner has no import edge to worktree, landing, status, or dashboard.

## Checkpoint verification

Run `bench test --package ./internal/jsonfile`, `bench test --package ./internal/commitment`, `bench test --package ./internal/intent`, `bench test --package ./internal/roadmap`, and `bench test --package ./cmd/bench`. Record a stale-plan or missing-effect red before the corresponding production change.

The JSON owner also validates exact field names for policy input.
A case alias must refuse before a receipt write.
