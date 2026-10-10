# Close the last chunk through the complete checkpoint

Blocked by: 02-hold-light-paths-during-a-build.md
Writes: .agents/commands/bench-implement-spec.md, docs/adr/0023-each-ticket-gets-a-fresh-author.md, README.md, CHANGELOG.md, projects/benchkit.md, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/anchors/registry_check_floor.go, internal/anchors/registry_check_floor_test.go, internal/anchors/registry_commitment.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_debug_loop.go, internal/anchors/registry_decision_maps.go, internal/anchors/registry_decision_maps_test.go, internal/anchors/registry_ft311_preparation.go, internal/anchors/registry_retained_workflow.go, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, internal/conformance/docs_workflow_helpers_test.go, internal/conformance/retained_workflow_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/docs-currency-token-diet, tests/canary/load-validity-metadata, tests/canary/workflow-guidance-anchors, tests/canary/guidance-prose-budgets, tests/canary/line-routing, tests/canary/skill-description-budgets, internal/anchors/registry_ft311_review_dispatch.go
Covers: BD14, BD15, BD16, BD17, BD18, BD19, BD20, BD21, BD50, BD51, BD52, BD55, BD58

## What to build

The last ticket's author freezes and records its chunk, and the last chunk closes through the complete checkpoint alone.
Tickets 01 and 02 move the build-base line numbers, so the author finds each target by its section and sentence lead.
In the Land paragraph that starts "Retain author verification", the author-records sentence that the spec quotes replaces the sentence that starts "When the orchestrator freezes".
That paragraph's last sentence becomes the two sentences that run the chunk checkpoint gate only before a successor chunk.
In the Build paragraph, the freeze-owner sentence that the spec quotes follows the pinned sentence that starts "After the last ticket in a chunk, freeze".

The Land paragraph that starts "After the last chunk, reconcile" keeps that first sentence byte-exact.
The two last-chunk sentences that the spec quotes follow that first sentence.
The handoff clause that ends "at each chunk checkpoint." adds "and at the final reconciliation".
ADR 0023 line 15 takes the two sentences that the spec states under "Chunk records and checkpoints".
The README chunk flow moves the "chunk checkpoint" label to the successor edge.

The chunk-chain row that pins the orchestrator-records sentence changes in place, and its mirror in `registry_chunk_chain_test.go` changes with it.
The handoff row in `registry_ft311_preparation.go` changes in place, and a new BD50 mirror rule goes into `registry_chunk_chain_test.go`, because that registry has no harness test.
New Require rows pin each new sentence, and Forbid rows keep each retired sentence out.
Each edit replaces text in place, so implement-spec keeps 89 lines.
One bullet under the `### Build directions` heading of `CHANGELOG.md` records the change.

## Acceptance

- [ ] Implement-spec holds the author-freeze and author-records sentences and drops the orchestrator sentence, and TestChunkChainAnchors goes red on each reversal (BD14, BD16, BD58).
- [ ] The Land verification sentence still raises no diagnostic (BD15).
- [ ] Implement-spec runs the chunk checkpoint gate only before a successor chunk and drops the every-chunk wording, and TestChunkChainAnchors goes red on each reversal (BD17, BD18).
- [ ] Implement-spec holds the two last-chunk sentences, and its reconciliation sentence keeps its lowercase pin in TestReviewConvergenceContractCurrentDocs (BD19, BD20, BD21).
- [ ] Implement-spec and ADR 0023 name the final-reconciliation handoff refresh, and TestChunkChainAnchors goes red when either is removed (BD50, BD52).
- [ ] README shows the chunk checkpoint on the successor edge only, and TestChunkChainAnchors goes red on the old label (BD51, BD55).
- [ ] Implement-spec has 89 lines, and TestGuidanceProseBudgetsHoldOnTheLiveTree passes.
