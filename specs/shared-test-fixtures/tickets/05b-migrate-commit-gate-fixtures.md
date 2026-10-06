# Migrate commit and gate fixtures

Blocked by: 05a-migrate-cli-adopt-fixtures.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/commit/chain_grammar_test.go, internal/commit/deletion_test.go, internal/commit/dry_run_test.go, internal/commit/landing_test.go, internal/commit/lane_structure_test.go, internal/commit/lane_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/gate/authorization/lane_test.go, internal/gate/cache_env_test.go, internal/gate/gate_prose_staged_root_test.go, internal/gate/gate_prose_staged_test.go, internal/gate/greenmarker/greenmarker_test.go, internal/gate/lane_record_test.go, internal/gate/lane_run_test.go, internal/gate/lane_select_test.go, internal/gate/prospective_owner_test.go, internal/gate/prospectiveartifact/prospectiveartifact_test.go, internal/gate/run_failure_outcomes_test.go, internal/gate/run_outcomes_test.go, internal/gate/verdict_registry_guard_test.go
Covers: GF16, GF19, GF29, GF30

## What to build

Create commit and gate fixtures through the shared Git owner. Preserve raw index facts and prospective fixture semantics.

Consume the accepted GF-C1 helper contract. Remove generic execution and identity copies in this package group.
Keep meaningful fixture composition and policy-free public fixture facades. Classify each specialized probe by its actual operation before an edit.
String literals containing example Go code remain data. No later ticket supplies behavior needed at this checkpoint.

## Acceptance

- [ ] Generic helper callers retain their fixture results (GF16).
- [ ] Specialized Git probes retain their input, environment, and exit contracts (GF19).
- [ ] The same fixture input family has equivalent repository facts before and after migration (GF29).
- [ ] Each migration chunk passes without its successor implementation (GF30).

## Checkpoint verification

- `bench test --package ./internal/commit`
- `bench test --package ./internal/gate`
- `bench test --package ./internal/gate/authorization`
- `bench test --package ./internal/gate/greenmarker`
- `bench test --package ./internal/gate/prospectiveartifact`
- `bench diff`

Run the named fixture family on the frozen baseline and candidate. Compare repository facts and explicit failure results.
Exclude timestamps and object IDs unless the fixture grades them. Record each retained specialized contract in the review pickup.
