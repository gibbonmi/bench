# Migrate CLI and adoption fixtures

Blocked by: 04-migrate-workflow-readers.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, cmd/bench/main_test.go, cmd/bench/otel_hook_seams_test.go, internal/adopt/adopt_test.go, internal/adopt/broker_test.go, internal/adopt/link_hook_test.go, internal/adopt/link_transaction_test.go, internal/adopt/repairtest/session_test.go, internal/adopt/setup_prompt_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: GF16, GF19, GF29, GF30

## What to build

Create CLI and adoption fixtures through the shared Git owner. Preserve hook-script input and explicit branch choices.

Coordinate overlapping adoption fixtures with landing-test-efficiency. Prefer its minimal payload migration first, then refresh this ticket's caller inventory and affected fences.
Use the approved amendment route if source drift requires new paths. Shared helpers alone establish no measured speed gain.

Consume the accepted GF-C1 helper contract. Remove generic execution and identity copies in this package group.
Keep meaningful fixture composition and policy-free public fixture facades. Classify each specialized probe by its actual operation before an edit.
String literals containing example Go code remain data. No later ticket supplies behavior needed at this checkpoint.

## Acceptance

- [ ] Generic helper callers retain their fixture results (GF16).
- [ ] Specialized Git probes retain their input, environment, and exit contracts (GF19).
- [ ] The same fixture input family has equivalent repository facts before and after migration (GF29).
- [ ] Each migration chunk passes without its successor implementation (GF30).

## Checkpoint verification

- `bench test --package ./cmd/bench`
- `bench test --package ./internal/adopt`
- `bench test --package ./internal/adopt/repairtest`
- `bench diff`

Run the named fixture family on the frozen baseline and candidate. Compare repository facts and explicit failure results.
Exclude timestamps and object IDs unless the fixture grades them. Record each retained specialized contract in the review pickup.
