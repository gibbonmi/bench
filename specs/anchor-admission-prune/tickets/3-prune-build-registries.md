# Prune the build-path topic registries

Blocked by: 2-prune-core-registry.md
Writes: internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/anchors/registry_charge_binding.go, internal/anchors/registry_charge_binding_test.go, internal/anchors/registry_calibration.go, internal/anchors/registry_calibration_test.go, internal/anchors/registry_check_floor.go, internal/anchors/registry_check_floor_test.go, internal/anchors/registry_ft311_preparation.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_spec_trace.go, internal/anchors/registry_spec_trace_test.go, internal/anchors/registry_data.go, internal/conformance/ft311_preparation_test.go, internal/conformance/charge_evidence_guidance_test.go, internal/conformance/implementation_continuation_test.go, internal/conformance/retained_workflow_test.go, internal/conformance/tier_live_tree_test.go, internal/conformance/fixture_bite_test.go, tests/canary/workflow-guidance-anchors, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: AP02, AP03, AP04, AP10, AP18, AP19, AP20, AP22, AP23, AP46

## What to build

The build-path topic registries keep only their admitted rows. They are the ticket-pass, chunk-chain, charge-binding, calibration, check-floor, FT311 preparation, FT311 review-dispatch, and spec-trace files. The nine step-scoped rows in the chunk-chain and calibration files stay, because the AP-C5 ticket removes them with their kind.

Each removed row leaves with its test rule and its entries in the conformance pin lists. Its fixture directories, which `assets/fixture-census.md` lists for this ticket, leave too. The spec-trace file becomes empty, so delete it and its test, and drop it from the registry append list. Keep each constant that a row in another file still reads, for example the review-repeat and final-reconciliation rules.

Update the count pins in the FT311 preparation test and the bounded-action rows in the charge-evidence test to the rows that stay. Remove a deleted test name from the live-tree test inventory. Re-cut each COMMAND needle in these files that is over the bound under the `internal/prose` word count, and keep each kept fixture biting. No touched file may grow.

## Acceptance

- [ ] Each of these registry files holds only admitted rows and the nine step-scoped rows.
- [ ] The spec-trace registry file and its test are gone, and the ticket record names each deleted file.
- [ ] The ticket record carries a removal census for the removed rows.
- [ ] `bench test --package ./internal/anchors`, `bench test --check docs-currency-workflow`, and `bench test --package ./internal/conformance` pass.
