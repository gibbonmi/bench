# Migrate execution fixtures

Blocked by: 04-migrate-workflow-readers.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, cmd/bench/main_test.go, cmd/bench/otel_hook_seams_test.go, internal/adopt/adopt_test.go, internal/adopt/broker_test.go, internal/adopt/link_hook_test.go, internal/adopt/link_transaction_test.go, internal/adopt/repairtest/session_test.go, internal/adopt/setup_prompt_test.go, internal/commit/chain_grammar_test.go, internal/commit/deletion_test.go, internal/commit/dry_run_test.go, internal/commit/landing_test.go, internal/commit/lane_structure_test.go, internal/commit/lane_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/build_contracts_test.go, internal/conformance/checks_test.go, internal/conformance/docs_workflow_checks_test.go, internal/conformance/fixture_bite_test.go, internal/conformance/gate_entry_test.go, internal/conformance/git_plumbing_owner_test.go, internal/conformance/handoff_single_source_test.go, internal/conformance/harness_test.go, internal/conformance/package_core_checks_test.go, internal/conformance/package_core_diagnostics_test.go, internal/conformance/subcommand_routing_table_test.go, internal/conformance/validity_checks_test.go, internal/env/git_policy_test.go, internal/env/kit_run_test.go, internal/gate/authorization/lane_test.go, internal/gate/cache_env_test.go, internal/gate/gate_prose_staged_root_test.go, internal/gate/gate_prose_staged_test.go, internal/gate/greenmarker/greenmarker_test.go, internal/gate/lane_record_test.go, internal/gate/lane_run_test.go, internal/gate/lane_select_test.go, internal/gate/prospective_owner_test.go, internal/gate/prospectiveartifact/prospectiveartifact_test.go, internal/gate/run_failure_outcomes_test.go, internal/gate/run_outcomes_test.go, internal/gate/verdict_registry_guard_test.go, internal/intent/assignment_lookup_test.go, internal/intent/intent_test.go, internal/intent/worktree_owner_test.go, internal/landing/close_test.go, internal/landing/composition_test.go, internal/landing/landing_helpers_test.go, internal/preflight/preflighttest/fixture.go, internal/preflight/preflighttest/reviewfiles.go, internal/responsebound/responseboundtest/checkout.go, internal/sessioninspect/sessioninspect_test.go, internal/shift/fault_test.go, internal/shift/refresh_test.go, internal/shift/shift_test.go, internal/testreport/check_test.go, internal/testreport/selection_test.go
Covers: GF16, GF19, GF29, GF30

## What to build

Consume the accepted GF-C1 helper contract in execution fixtures. Preserve meaningful fixture composition and explicit child environments. The preflight fast-import input and expected-failure probes retain their transport. Literal Go source inside checker fixtures remains data.

Build the existing execution fixtures through their public fixture APIs. Preserve explicit stdin, timestamps, child environments, and expected nonzero exits.

## Acceptance

- [ ] Generic helper callers retain their fixture results (GF16).
- [ ] Specialized Git probes retain their input, environment, and exit contracts (GF19).
- [ ] The same fixture input family has equivalent repository facts before and after migration (GF29).
- [ ] Each migration chunk passes without its successor implementation (GF30).

## Checkpoint verification

- `bench test --package ./cmd/bench`
- `bench test --package ./internal/adopt`
- `bench test --package ./internal/adopt/repairtest`
- `bench test --package ./internal/commit`
- `bench test --package ./internal/env`
- `bench test --package ./internal/gate`
- `bench test --package ./internal/gate/authorization`
- `bench test --package ./internal/gate/greenmarker`
- `bench test --package ./internal/gate/prospectiveartifact`
- `bench test --package ./internal/intent`
- `bench test --package ./internal/landing`
- `bench test --package ./internal/preflight/preflighttest`
- `bench test --package ./internal/responsebound/responseboundtest`
- `bench test --package ./internal/sessioninspect`
- `bench test --package ./internal/shift`
- `bench test --package ./internal/testreport`
- `bench test --package ./internal/conformance`
- `bench diff`

Use the frozen baseline and the candidate on the same fixture inputs. Exclude commit timestamps and object IDs from comparisons unless the scenario grades them.
Record the observed facts and the named omission or swap in the review pickup. A changed outcome stops the migration.
