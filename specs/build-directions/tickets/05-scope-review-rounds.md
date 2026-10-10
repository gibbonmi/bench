# Scope review rounds to one combined charge and named findings

Blocked by: 04-cap-witnesses-and-keep-moved-guards.md
Writes: .agents/commands/bench-review-implementation.md, .agents/skills/bench-craft-delegate/references/delegation-discipline.md, CHANGELOG.md, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/anchors/registry_calibration.go, internal/anchors/registry_calibration_test.go, internal/anchors/registry_charge_binding.go, internal/anchors/registry_charge_binding_test.go, internal/anchors/registry_check_floor.go, internal/anchors/registry_check_floor_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_debug_loop.go, internal/anchors/registry_ft311_preparation.go, internal/anchors/registry_retained_workflow.go, internal/conformance/docs_workflow_helpers_test.go, internal/conformance/fixture_bite_test.go, internal/conformance/ft311_preparation_test.go, internal/conformance/implementation_continuation_test.go, internal/conformance/retained_workflow_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/workflow-guidance-anchors
Covers: BD31, BD32, BD33, BD34, BD35, BD36, BD37, BD38, BD54, BD59

## What to build

A review round's axes run in parallel, and a one-ticket round sends every finding in one combined repair charge.
A confirming round examines only its named findings.
A new `RequireInStep` row (Process, step 4) in `registry_chunk_chain.go` pins the existing parallel-dispatch sentence.
The "Review modes" section of `/bench-review-implementation` takes the three combined-charge sentences that the spec quotes as their own paragraph.
When the findings fall on several tickets, each affected ticket keeps its own repair session, and `.bench/BENCH.md`'s repair rule stays unchanged.

The confirming-round sentence at review-implementation:48 becomes the two sentences that the spec quotes.
The Standards first-look sentence at :49 goes away, and the three-axis sentence at :47 stays.
The delegation-discipline bullet at :216-217 becomes the named-findings bullet that the spec quotes.

The row at `registry_ft311_review_dispatch.go:22` changes in place to the first new sentence, and its diagnostic drops the folds clause.
New rows in `registry_chunk_chain.go` pin the second sentence, the combined-charge sentences, and the delegation-discipline bullet.
Forbid rows keep the first-look sentence and the blocking-bar scope out.
Independent mirror rules in `registry_chunk_chain_test.go` restate every new and rewritten row.
One bullet under the `### Build directions` heading of `CHANGELOG.md` records the change.

## Acceptance

- [ ] Review step 4 keeps its parallel-dispatch sentence, and TestChunkChainAnchors goes red when it is removed (BD31).
- [ ] Review-implementation holds the three combined-charge sentences, and TestChunkChainAnchors goes red when any one is removed (BD32, BD33, BD34).
- [ ] Review-implementation holds both confirming-round sentences and keeps the three-axis sentence, and TestChunkChainAnchors goes red on each removal (BD35, BD36, BD59).
- [ ] The first-look sentence and the blocking-bar scope are absent, and TestChunkChainAnchors goes red when either returns (BD37, BD54).
- [ ] Delegation-discipline holds the named-findings clause, and TestChunkChainAnchors goes red when it is removed (BD38).
- [ ] The "prepared review dispatch: " family keeps 17 rows, so TestPreparedReviewGuidance stays green.
