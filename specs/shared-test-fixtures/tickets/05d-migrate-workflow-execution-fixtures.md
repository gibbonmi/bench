# Migrate workflow execution fixtures

Blocked by: 05c-migrate-preflight-conformance-fixtures.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/intent/assignment_lookup_test.go, internal/intent/intent_test.go, internal/intent/worktree_owner_test.go, internal/landing/close_test.go, internal/landing/composition_test.go, internal/landing/landing_helpers_test.go, internal/responsebound/responseboundtest/checkout.go, internal/sessioninspect/sessioninspect_test.go, internal/shift/fault_test.go, internal/shift/refresh_test.go, internal/shift/shift_test.go, internal/testreport/check_test.go, internal/testreport/selection_test.go
Covers: GF16, GF19, GF29, GF30

## What to build

Create intent, landing, shift, response, session, and test-report fixtures through the shared Git owner. Preserve each workflow refusal and explicit environment.

Consume the accepted GF-C1 helper contract. Remove generic execution and identity copies in this package group.
Keep meaningful fixture composition and policy-free public fixture facades. Classify each specialized probe by its actual operation before an edit.
String literals containing example Go code remain data. No later ticket supplies behavior needed at this checkpoint.

## Acceptance

- [ ] Generic helper callers retain their fixture results (GF16).
- [ ] Specialized Git probes retain their input, environment, and exit contracts (GF19).
- [ ] The same fixture input family has equivalent repository facts before and after migration (GF29).
- [ ] Each migration chunk passes without its successor implementation (GF30).

## Checkpoint verification

- `bench test --package ./internal/intent`
- `bench test --package ./internal/landing`
- `bench test --package ./internal/responsebound/responseboundtest`
- `bench test --package ./internal/sessioninspect`
- `bench test --package ./internal/shift`
- `bench test --package ./internal/testreport`
- `bench diff`

Run the named fixture family on the frozen baseline and candidate. Compare repository facts and explicit failure results.
Exclude timestamps and object IDs unless the fixture grades them. Record each retained specialized contract in the review pickup.
