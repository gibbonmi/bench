# Close the last chunk through the complete checkpoint

Blocked by: 02-hold-light-paths-during-a-build.md
Writes: .agents/commands/bench-implement-spec.md, docs/adr/0023-each-ticket-gets-a-fresh-author.md, README.md, CHANGELOG.md, projects/benchkit.md, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/anchors/registry_check_floor.go, internal/anchors/registry_check_floor_test.go, internal/anchors/registry_commitment.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_debug_loop.go, internal/anchors/registry_decision_maps.go, internal/anchors/registry_decision_maps_test.go, internal/anchors/registry_ft311_preparation.go, internal/anchors/registry_retained_workflow.go, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, internal/conformance/docs_workflow_helpers_test.go, internal/conformance/retained_workflow_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/docs-currency-token-diet, tests/canary/load-validity-metadata, tests/canary/workflow-guidance-anchors
Covers: BD14, BD15, BD16, BD17, BD18, BD19, BD20, BD21, BD50, BD51, BD52, BD55, BD58

## What to build

The last ticket's author freezes and records its chunk, and the last chunk closes through the complete checkpoint alone.
The Land paragraph at implement-spec:58 replaces its orchestrator sentence with the author-records sentence that the spec quotes.
Its last sentence becomes the two sentences that run the chunk checkpoint gate only before a successor chunk.
The Build paragraph at implement-spec:50 keeps its pinned freeze sentence byte-exact and adds "The last ticket's author performs that freeze, as the Land section states."

The paragraph at implement-spec:62 keeps its first sentence byte-exact and adds the two last-chunk sentences that the spec quotes.
The handoff clause at implement-spec:31 adds "and at the final reconciliation", and ADR 0023 line 15 takes the same change.
The README chunk flow moves the "chunk checkpoint" label to the successor edge.

The rows at `registry_chunk_chain.go:20` and `registry_ft311_preparation.go:43` change in place, and their mirrors change with them.
New Require rows pin each new sentence, and Forbid rows keep each retired sentence out.
Each edit replaces text in place, so implement-spec keeps 89 lines.
One `CHANGELOG.md` entry records the change.

## Acceptance

- [ ] Implement-spec holds the author-freeze and author-records sentences and drops the orchestrator sentence, and TestChunkChainAnchors goes red on each reversal (BD14, BD16, BD58).
- [ ] The verification sentence at implement-spec:58 still raises no diagnostic (BD15).
- [ ] Implement-spec runs the chunk checkpoint gate only before a successor chunk and drops the every-chunk wording, and TestChunkChainAnchors goes red on each reversal (BD17, BD18).
- [ ] Implement-spec holds the two last-chunk sentences, and its reconciliation sentence keeps its lowercase pin in TestReviewConvergenceContractCurrentDocs (BD19, BD20, BD21).
- [ ] Implement-spec and ADR 0023 name the final-reconciliation handoff refresh, and TestChunkChainAnchors goes red when either is removed (BD50, BD52).
- [ ] README shows the chunk checkpoint on the successor edge only, and TestChunkChainAnchors goes red on the old label (BD51, BD55).
- [ ] Implement-spec has 89 lines, and TestGuidanceProseBudgetsHoldOnTheLiveTree passes.
