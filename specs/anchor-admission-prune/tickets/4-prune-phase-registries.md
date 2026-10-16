# Prune the phase topic registries

Blocked by: 3-prune-build-registries.md
Writes: internal/anchors/registry_debug_loop.go, internal/anchors/registry_commitment.go, internal/anchors/registry_decision_maps.go, internal/anchors/registry_decision_maps_test.go, internal/anchors/registry_craft_research.go, internal/anchors/registry_craft_research_test.go, internal/anchors/registry_review_rules.go, internal/anchors/registry_review_rules_test.go, internal/anchors/registry_front_door.go, internal/anchors/registry_front_door_test.go, internal/anchors/registry_data.go, internal/conformance/commitment_guidance_test.go, internal/conformance/data_handling_test.go, internal/conformance/tier_live_tree_test.go, tests/canary/workflow-guidance-anchors, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: AP02, AP03, AP04, AP10, AP18, AP19, AP20, AP46

## What to build

The phase topic registries keep only their admitted rows. They are the debug-loop, commitment, decision-map, craft-research, review-rule, and front-door files. Each removed row leaves together with its test rule, its mutation cases, and the fixture directories that `assets/fixture-census.md` lists for this ticket.

The review-rule and front-door files become empty, so delete them and their tests, and drop them from the registry append list. Keep the exported commitment marker and the commitment diagnostic prefix, because conformance tests read them.

Remove the commitment-guidance mutation cases and the data-handling cases that name removed rows. Keep the cases for kept rows. Remove a deleted test name from the live-tree test inventory. Re-cut each COMMAND needle in these files that is over the bound under the `internal/prose` word count. No touched file may grow.

## Acceptance

- [ ] Each of these registry files holds only admitted rows.
- [ ] The review-rule and front-door registry files and their tests are gone, and the ticket record names each deleted file.
- [ ] The ticket record carries a removal census for the removed rows.
- [ ] `bench test --package ./internal/anchors`, `bench test --check docs-currency-workflow`, and `bench test --package ./internal/conformance` pass.
