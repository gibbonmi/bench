# Name the ticket author as the freezer

Blocked by: 05-record-one-ticket-in-one-write.md
Writes: .agents/commands/bench-implement-spec.md, .agents/commands/bench-review-implementation.md, .bench/BENCH.md, .bench/BENCH-reference.md, .agents/skills/bench-craft-tickets/SKILL.md, docs/known-issues.md, CHANGELOG.md, .bench/commitment.json, projects/benchkit.md, tests/canary/line-routing/line-binding-prose-drift/, tests/canary/skill-description-budgets/, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, internal/anchors/registry_check_floor.go, internal/anchors/registry_check_floor_test.go, internal/anchors/registry_retained_workflow.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_commitment.go, internal/anchors/registry_debug_loop.go, internal/anchors/registry_ft311_preparation.go, internal/anchors/registry_calibration.go, internal/anchors/registry_calibration_test.go, internal/anchors/anchor_harness_diagnostics_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/conformance/retained_workflow_test.go, internal/conformance/docs_workflow_helpers_test.go, internal/conformance/implementation_continuation_test.go, internal/conformance/commitment_guidance_test.go, internal/conformance/checks_test.go, internal/conformance/skills_index_checks_test.go, internal/conformance/help_inventory_single_source_test.go, internal/conformance/validity_checks_test.go, internal/conformance/recurrence_maintenance_contract_test.go, internal/conformance/docs_workflow_checks_test.go, internal/conformance/fixture_bite_test.go, internal/conformance/prose_budget_test.go, internal/conformance/package_shipped_surface_test.go, internal/conformance/prose_mechanics_test.go, tests/canary/workflow-guidance-anchors/, tests/canary/docs-currency-token-diet/, tests/canary/load-validity-metadata/, tests/canary/skills-index-command-adapters/, tests/canary/guidance-prose-budgets/
Covers: RT47, RT48, RT49, RT50, RT51, RT52, RT53, RT54, RT64

## What to build

The guidance names the last ticket's author as the freezer and `bench record ticket` as the record call.

- The implementation phase changes two Land sentences in place: the freeze sentence and the verification sentence name `bench record ticket`. The file has 0 lines of prose budget headroom, so the line count stays the same.
- The operating guide's chunk sentence names the last ticket's author as the freezer, within the 25-word sentence cap.
- The review phase's step 6 names the last ticket's author as the writer of the frozen pair. The step 6 record paragraph holds 6 sentences, so the sentence that lists the excerpt and ticket forms goes in another paragraph of step 6. The sentence "Write each completed chunk, verification, review, and amendment entry with `bench record`, not with hand-built JSON." stays byte for byte.
- The ticket craft and the reference name the last ticket's author as the freezer, and neither file holds `review freezes`.
- The known-issues row that names FT405 goes. The other rows stay. The workaround of the evidence-preflight row keeps each sentence within 20 words (RT-CURRENT Standards S1).
- The changelog gains one entry that names the excerpt form, the ticket form, the retained response, and the `dirty_record` field.

The chunk-chain registry and its twin get the changed and new anchors.
The Forbid needle "review freezes its predecessor tip and current tip" matches no live text today, so the RT52 Forbid rows replace it.
RT53 is a Forbid row on `FT405` in `docs/known-issues.md`, in `registry_ticket_passes.go` and its twin.
`registry_retained_workflow.go` has 2 lines of headroom, and `registry_data.go` and its test are over their budgets, so new rows go in the chunk-chain and ticket-passes pairs.

`.bench/commitment.json` is in this ticket's `Writes:` line for the spec-stage close binding.

## Acceptance

- [ ] The implementation phase's Land section names `bench record ticket` for the freeze and for the verification entries.
- [ ] The operating guide names the last ticket's author as the freezer before the three axes run.
- [ ] The review phase's step 6 names the last ticket's author and `bench record ticket` for the frozen pair, and names the excerpt and ticket forms.
- [ ] The ticket craft and the reference name the last ticket's author as the freezer and hold no `review freezes` text.
- [ ] `docs/known-issues.md` holds no row that names FT405.
- [ ] Each changed sentence has an anchor and an independent twin that turn red when the sentence is absent.
- [ ] `bench test --check guidance-prose-budgets` and `bench test --check prose-mechanics` pass.
