# Fold main only at chunk boundaries

Blocked by: none
Writes: .agents/commands/bench-implement-spec.md, .agents/skills/bench-craft-delegate/SKILL.md, CHANGELOG.md, projects/benchkit.md, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/anchors/registry_calibration.go, internal/anchors/registry_calibration_test.go, internal/anchors/registry_charge_binding.go, internal/anchors/registry_charge_binding_test.go, internal/anchors/registry_check_floor_test.go, internal/anchors/registry_commitment.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_debug_loop.go, internal/anchors/registry_ft311_preparation.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_retained_workflow.go, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, internal/conformance/docs_workflow_helpers_test.go, internal/conformance/fixture_bite_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/claude-agent-definitions, tests/canary/workflow-guidance-anchors, tests/canary/guidance-prose-budgets, tests/canary/line-routing, tests/canary/skill-description-budgets
Covers: BD1, BD2, BD3, BD4, BD5, BD6, BD60, BD61

## What to build

A build folds `main` only at its sanctioned points, and the guidance says so in one place.
The Land section of `/bench-implement-spec` takes a new paragraph with the five sentences that the spec quotes under "Fold timing".
The paragraph names the three fold points and forbids a fold between two tickets of a chunk.
It also makes the post-ticket preflight inside a chunk pass the chunk's `--base` and the current `--source-tip`.

The craft-delegate skill adds the scope sentence that the spec quotes at the end of line 68.
Its two stale-base sentences stay byte-exact, because `internal/anchors/registry_data.go:69` and `internal/conformance/fixture_bite_test.go:350` quote them.
The pinned fold sentences at implement-spec:56 and review-implementation:98 stay byte-exact.

New Require rows in `registry_chunk_chain.go` pin each new sentence, and independent mirror rules in `registry_chunk_chain_test.go` restate them.
The implement-spec line budget in `projects/benchkit.md` moves from 85 to 87, and the craft-delegate budget stays 126.
One `CHANGELOG.md` entry records the change.

## Acceptance

- [ ] Implement-spec Land holds each of the five fold-paragraph sentences, and TestChunkChainAnchors goes red when any one of them is removed (BD1, BD2, BD3, BD60, BD61).
- [ ] The craft-delegate skill holds the scope sentence, and TestChunkChainAnchors goes red when it is removed (BD4).
- [ ] The two craft-delegate stale-base sentences stay byte-exact, so TestWorkflowCadenceAnchorsRejectDeletionAndSwap stays green (BD5).
- [ ] The implement-spec:56 and review-implementation:98 pinned sentences raise no diagnostic (BD6).
- [ ] Implement-spec has 87 lines, its budget row reads 87, and TestGuidanceProseBudgetsHoldOnTheLiveTree passes.
