# 3. Align the phase commands and docs

Blocked by: 2-align-the-line-and-delegation-skills.md
Writes: .agents/commands/bench-write-spec.md, .agents/commands/bench-review-implementation.md, .agents/commands/bench-final-check.md, README.md, docs/field-guide.html, docs/adr/0021-benchmark-workflow-orchestration.md, docs/adr/0023-each-ticket-gets-a-fresh-author.md (new), internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_data.go, internal/conformance/retained_workflow_test.go, internal/conformance/docs_workflow_helpers_test.go, tests/canary/workflow-guidance-anchors/, internal/anchors/registry_calibration.go, internal/anchors/registry_calibration_test.go, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/anchors/registry_data_test.go, internal/anchors/registry_decision_maps.go, internal/anchors/registry_decision_maps_test.go, internal/anchors/registry_ft311_preparation.go, internal/anchors/registry_retained_workflow.go, tests/canary/docs-currency-token-diet/, tests/canary/load-validity-metadata/, tests/canary/skills-index-command-adapters/, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, projects/benchkit.md, CONTEXT.md, tests/canary/guidance-prose-budgets/, tests/canary/line-routing/, tests/canary/skill-description-budgets/
Covers: FA20, FA21, FA22, FA23, FA24, FA25, FA29, FA31

## What to build

Chunk: FA.

Align the phase commands and the public docs to the owner rule in `.bench/BENCH.md`.

- `bench-write-spec`: the exit handoff recommends the line for fresh ticket authors on one integration source.
- `bench-review-implementation`: accepted findings go to a fresh repair author, and the orchestrator performs the final reconciliation.
- `bench-final-check`: the retro reports how each ticket author performed against the planned boundaries.
- `README.md` and `docs/field-guide.html`: state that each ticket gets a fresh author session. Replace the field-guide sentence "One implementation session retains authorship through the approved ticket graph and its chunk reviews." and its pin in `internal/anchors/registry_data.go`.
- ADR 0023 records the decision as the current state. ADR 0021's first, second, third, and fifth consequences change to the decided state and point to ADR 0023.

The implement-now light path in `bench-drain` stays in the main session, so its anchor keeps its bytes. Add the rows of FA20 to FA23, change each pin of a changed sentence, and update each canary whose mutation names a changed sentence.

## Acceptance

- [ ] `bench test --check docs-currency-workflow` passes with the FA20 to FA23 needles in place and the FA25 needle unchanged.
- [ ] ADR 0023 exists, and ADR 0021 no longer states the single-author contract.
- [ ] Each changed `write-spec` canary still reds its mutation.
