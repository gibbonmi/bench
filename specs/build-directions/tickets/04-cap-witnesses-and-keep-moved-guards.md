# Cap the witnesses and keep the guards of moved code

Blocked by: 03-close-the-last-chunk-through-the-complete-checkpoint.md
Writes: .agents/commands/bench-implement-spec.md, docs/known-issues.md, CHANGELOG.md, projects/benchkit.md, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, internal/anchors/registry_charge_binding.go, internal/anchors/registry_charge_binding_test.go, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/anchors/registry_check_floor_test.go, internal/anchors/registry_commitment.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_debug_loop.go, internal/anchors/registry_ft311_preparation.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_retained_workflow.go, internal/conformance/charge_evidence_guidance_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/guidance-prose-budgets, tests/canary/line-routing, tests/canary/skill-description-budgets, tests/canary/workflow-guidance-anchors
Covers: BD22, BD23, BD24, BD25, BD26, BD27, BD28, BD29, BD30

## What to build

The author charge keeps the guards of moved code and caps the witnesses, and the cap has one source.
The author-charge paragraph at implement-spec:25 adds "Code that the author moves or merges keeps its error handling and validation." after its reach-list sentence.
That paragraph then holds six sentences, the prose bound.

A new paragraph after it holds the four witness-cap sentences that the spec quotes.
It stays a separate paragraph before the "Build action requires" paragraph, because `charge_evidence_guidance_test.go` pins every sentence of that paragraph.
The `docs/known-issues.md:12` workaround cell becomes the two pointer sentences that the spec quotes, and the row keeps its symptom and FT168 retire cells.
The delegation-discipline one-red rule stays byte-exact.

New Require rows in `registry_ticket_passes.go` pin each new sentence and the known-issues pointer, and a Forbid row keeps the old cap text out of `docs/known-issues.md`.
Independent mirror rules in `registry_ticket_passes_test.go` restate them.
The implement-spec line budget moves from 89 to 91.
One `CHANGELOG.md` entry records the change.

## Acceptance

- [ ] Implement-spec holds the moved-code guard sentence and keeps the reach list, and TestWriteSpecChunkedSlicing goes red when either is removed (BD22, BD27).
- [ ] Implement-spec holds the four witness sentences, and TestWriteSpecChunkedSlicing goes red when any one is removed (BD23, BD24, BD25, BD26).
- [ ] `docs/known-issues.md` holds the owner pointer and no copy of the cap, and TestWriteSpecChunkedSlicing goes red on each reversal (BD28, BD29).
- [ ] The delegation-discipline one-red rule raises no diagnostic in TestChargeBindingAnchors (BD30).
- [ ] Implement-spec has 91 lines, and TestGuidanceProseBudgetsHoldOnTheLiveTree passes.
