# Hold light-path fixes while a build runs

Blocked by: 01-fold-only-at-chunk-boundaries.md
Writes: .agents/commands/bench-implement-spec.md, .bench/BENCH.md, docs/adr/0023-each-ticket-gets-a-fresh-author.md, CHANGELOG.md, projects/benchkit.md, internal/anchors/anchor_harness_diagnostics_test.go, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/anchors/registry_check_floor_test.go, internal/anchors/registry_commitment.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_debug_loop.go, internal/anchors/registry_ft311_preparation.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_retained_workflow.go, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, internal/conformance/commitment_guidance_test.go, internal/conformance/retained_workflow_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/docs-currency-token-diet, tests/canary/load-validity-metadata, tests/canary/skills-index-command-adapters, tests/canary/workflow-guidance-anchors
Covers: BD7, BD8, BD9, BD10, BD11, BD12, BD13, BD57

## What to build

A build dispatches no light-path fix while it runs, and the operating guide points to `/bench-implement-spec` for that rule.
The Build section of `/bench-implement-spec` takes a new paragraph with the two sentences that the spec quotes under "Light paths during a build".
A fix that the build's own outcome needs still lands as the build's own commit under the fix-don't-park rule.

In `.bench/BENCH.md`, the sentence after the `**Delegate each light-path fix.**` label becomes the "Outside a build" sentence that the spec quotes.
The label stays byte-exact, because `internal/conformance/commitment_guidance_test.go:35` mutates from it.
Every other sentence of that paragraph stays, and the file keeps 192 lines.
ADR 0023 line 19 scopes its light-path clause to a phase outside a build, as the spec states.
It then adds "A build dispatches no light-path fix while it runs."
The `/bench-final-check` sentence at final-check:52 stays, because it governs a landing from another session.

Require and Forbid rows pin each new and retired sentence, with independent mirrors or commitment-guidance cases.
The implement-spec line budget moves from 87 to 89.
One `CHANGELOG.md` entry records the change.

## Acceptance

- [ ] Implement-spec holds the hold sentence and the own-commit sentence, and TestChunkChainAnchors goes red when either is removed (BD7, BD8).
- [ ] The operating guide no longer grants an in-build light-path fix, and TestCommitmentGuidance goes red when the old grant returns (BD9).
- [ ] The operating guide holds the "Outside a build" clause and the pointer clause, and TestCommitmentGuidance goes red when either is removed (BD10, BD57).
- [ ] The `**Delegate each light-path fix.**` label and its pinned neighbour sentences stay, so TestCommitmentGuidance keeps its mutation anchor (BD11).
- [ ] ADR 0023 holds the build sentence and drops the unscoped clause, and TestChunkChainAnchors goes red when either changes back (BD12, BD13).
- [ ] Implement-spec has 89 lines, `.bench/BENCH.md` has 192, and TestGuidanceProseBudgetsHoldOnTheLiveTree passes.
