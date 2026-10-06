# Migrate preflight and conformance fixtures

Blocked by: 05b-migrate-commit-gate-fixtures.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/build_contracts_test.go, internal/conformance/checks_test.go, internal/conformance/docs_workflow_checks_test.go, internal/conformance/fixture_bite_test.go, internal/conformance/gate_entry_test.go, internal/conformance/git_plumbing_owner_test.go, internal/conformance/handoff_single_source_test.go, internal/conformance/harness_test.go, internal/conformance/package_core_checks_test.go, internal/conformance/package_core_diagnostics_test.go, internal/conformance/subcommand_routing_table_test.go, internal/conformance/validity_checks_test.go, internal/env/git_policy_test.go, internal/env/kit_run_test.go, internal/preflight/preflighttest/fixture.go, internal/preflight/preflighttest/reviewfiles.go
Covers: GF16, GF19, GF29, GF30

## What to build

Create preflight and conformance fixtures through the shared Git owner. Preserve fast-import input, timestamps, explicit child environments, and expected failures.

Consume the accepted GF-C1 helper contract. Remove generic execution and identity copies in this package group.
Keep meaningful fixture composition and policy-free public fixture facades. Classify each specialized probe by its actual operation before an edit.
String literals containing example Go code remain data. No later ticket supplies behavior needed at this checkpoint.

## Acceptance

- [ ] Generic helper callers retain their fixture results (GF16).
- [ ] Specialized Git probes retain their input, environment, and exit contracts (GF19).
- [ ] The same fixture input family has equivalent repository facts before and after migration (GF29).
- [ ] Each migration chunk passes without its successor implementation (GF30).

## Checkpoint verification

- `bench test --package ./internal/conformance`
- `bench test --package ./internal/env`
- `bench test --package ./internal/preflight/preflighttest`
- `bench diff`

Run the named fixture family on the frozen baseline and candidate. Compare repository facts and explicit failure results.
Exclude timestamps and object IDs unless the fixture grades them. Record each retained specialized contract in the review pickup.
