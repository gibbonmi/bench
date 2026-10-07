# Retain private test directories through phase cleanup

Blocked by: 01-own-lifetime-through-bounds.md, 02-protect-nested-bounds-users.md, 03-migrate-gate-process-policy.md, 04-retain-selected-build-binaries.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, cmd/bench/process_resources.go (new), internal/conformance/axi_query_registry_test.go, internal/conformance/injected_ports_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/env/kit_run.go, internal/env/kit_run_lifetime_test.go (new), internal/env/kit_run_test.go, internal/env/protection.go (new), internal/gate/engine.go, internal/gate/lifetime_resources_test.go (new), internal/gate/phases.go
Covers: PL37, PL38

## What to build

A real phase child retains its private home and temporary directory when cleanup cannot prove completion.
Create protection before OpenKitTestRun permits child use.
Entries retains the complete inherited/local vector through its existing environment routing.
KitTestRun.Close returns the protected or uncertain release result instead of removing active files.

Consume Descriptor and Disposition from ticket 2, with the existing private-run identity as owner authority.
Install that owner's identity adapter in the CLI validator table.
The gate phase owner consumes Close failure before reporting success.
It consumes the selected-binary result from ticket 4 through its existing owner path.
Keep the current private environment contents and ordinary completed-run cleanup.

Before implementation starts, confirm complete accepted and landed durable-file-replacement delivery.
Its leaf, review-record caller, and native qualification must be complete.
Durable-caller-migration is not a prerequisite.
Use one retained integration source and a fresh author for this ticket.

Keep each caller's selected normal-exit policy.
A raw child status cannot certify required cleanup or publication.
Retain uncertain resources and every unresolved obligation.
Never signal from recovered records or add a blanket normal-exit kill.

This ticket belongs to PL-C3.
Finish the preceding chunk review before this chunk starts.
Its blockers include every prior overlapping Writes owner.
No successor caller must exist for this ticket's owned predicates.

## Acceptance

- [ ] PL37: KitTestRun.Close retains a protected private directory.
- [ ] PL38: The gate consumes private-directory Close failure.

## Verification

Run the named checks through their production entries.
For each owned row, record a behavioral omission or swap that fails, restoration, and the passing result.
A compilation failure is not a behavioral red.
Preserve existing assertions and fixture mutation purpose.
Co-owned closure holders receive only necessary reference or fixture updates.

- `bench test --package ./internal/env`
- `bench test --package ./internal/gate`

Mutation witness: Drop the private-run descriptor or discard KitTestRun.Close failure. The child-owned file survives while a successful phase result must disappear.
