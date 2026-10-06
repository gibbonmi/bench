# 4. State the light-path scope in the guide, the drain, and the implementation command

Blocked by: 3-admit-light-path-landing.md
Writes: .bench/BENCH.md, .agents/commands/bench-drain.md, .agents/commands/bench-implement-spec.md, projects/benchkit.md, CHANGELOG.md, internal/anchors/registry_commitment.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/anchor_harness_diagnostics_test.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_retained_workflow.go, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/anchors/registry_debug_loop.go, internal/anchors/registry_ft311_preparation.go, internal/conformance/commitment_guidance_test.go, internal/conformance/recurrence_maintenance_contract_test.go, tests/canary/workflow-guidance-anchors/, tests/canary/docs-currency-token-diet/benchref-imported, tests/canary/docs-currency-token-diet/benchref-pointer-dropped, tests/canary/docs-currency-token-diet/benchref-section-duplicated, tests/canary/docs-currency-token-diet/dogfood-referent-shipped, tests/canary/docs-currency-token-diet/missing-cli-inventory, tests/canary/docs-currency-token-diet/stale-cli-doc-reference, tests/canary/load-validity-metadata/readme-shared-rule-drift, tests/canary/load-validity-metadata/shared-rule-drift, tests/canary/skills-index-command-adapters/adapter-inert-invocation-key, tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy, tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted, tests/canary/row-next-grammar/token-table-lacks-kit-edit, tests/canary/guidance-prose-budgets/over-budget-skill, tests/canary/line-routing/line-binding-prose-drift, tests/canary/skill-description-budgets/budget-table-missing, tests/canary/skill-description-budgets/description-folded, tests/canary/skill-description-budgets/description-missing, tests/canary/skill-description-budgets/over-budget-command, tests/canary/skill-description-budgets/over-budget-description, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: LP34, LP35, LP36, LP37, LP38, LP39, LP40, LP41, LP42, LP43, LP44, LP45, LP46, LP54, LP55

## What to build

This ticket is the first ticket of review chunk LP-C2. It starts after the
LP-C1 checkpoint, so the guidance never states an exemption that the gate does
not yet admit. It delivers the guidance text of the spec section
`Guidance text`, word for word, and the anchor rows that pin it.

Edit `.bench/BENCH.md`. Change the two commitment-paragraph sentences, and add
the new light-path paragraph after that paragraph. Rewrite the learning
paragraph under the heading `**Delegate each light-path fix.**`. Keep its two
anchored dispatch sentences and its done-claim sentence. Drop its last
sentence.

Change the second sentence of the fix paragraph, and add the one
capture sentence. Keep the light-path table row verbatim.

Edit step 5, the learning sentences, and `## Delegate the evidence` in
`.agents/commands/bench-drain.md`, as the spec states. Add the one-sentence
paragraph after the entry paragraph of
`.agents/commands/bench-implement-spec.md`.

`.bench/BENCH.md` and `.agents/commands/bench-implement-spec.md` are at their
line budgets in `projects/benchkit.md`. Keep each file inside its budget if the
text allows it. Otherwise raise that budget row in this ticket.

In `internal/anchors/registry_commitment.go` and
`internal/anchors/registry_data.go`, require the twelve sentences of rows LP34
to LP40, LP43 to LP46, and LP55, one row each. Forbid the three retired
fragments with the diagnostics that the spec names. Retire the two forbid rows
that the spec names, and keep the forbid row
`a light-path fix that needs no reviewer decision`. Change the retained
implement-now row to the new delegate sentence and the diagnostic
`.agents/commands/bench-drain.md dropped the delegated implement-now light-path route`.
Update `EXPECT` in the canary `drain-implement-now-route` to that diagnostic.

`TestCommitmentGuidance` and `TestRecurrenceMaintenanceContractCheckBites` take
the same sentences as independent copies. Re-anchor the `unadmitted learning fix`
case on `**Delegate each light-path fix.**`. Re-anchor the
`default implementation` and `declined row` cases on
`The coordinator verifies the done-claim against the ticket's acceptance rows and the gate.`

Run one `rg` for each old needle across the tree, and update each hit in this
ticket's `Writes:` paths. Most entries of that line are closure files that build
preflight requires: the fixture canaries that pin each guidance file, the anchor
registry files that name it, and the command-registry files. Edit one only when
its check reds.

In `CHANGELOG.md`, add one entry under `[Unreleased]` for the light-path commit
and landing, and for the drain's delegate route. Correct the `/bench-drain`
entry that says the drain no longer builds a light-path item.

## Acceptance

- [ ] `TestCommitmentGuidance` in `internal/conformance/commitment_guidance_test.go` shows that each guide deletion raises its diagnostic: the start route (LP34) and the intake (LP35).
- [ ] The same test shows the same for the observable (LP36), the `Writes:` boundary (LP37), the row rule (LP38), and the drain dispatch (LP39).
- [ ] The same test shows that each deletion raises its diagnostic: the drain's delegate route (LP40) and the drain's spec-intake sentence (LP43).
- [ ] The same test shows that the deletion of the implementation command's light-path sentence raises its diagnostic (LP44).
- [ ] The same test shows that the restore of `light path and fixes included` (LP41), of the drain's old admission sentence (LP42), and of `a light-path fix that the active committed outcome needs` (LP54) raises its diagnostic.
- [ ] `TestEveryRetainedFixtureBitesThroughRegisteredOwner` in `internal/conformance/fixture_bite_test.go` shows that the canary `drain-implement-now-route` raises `.agents/commands/bench-drain.md dropped the delegated implement-now light-path route` (LP45).
- [ ] `TestRecurrenceMaintenanceContractCheckBites` in `internal/conformance/recurrence_maintenance_contract_test.go` shows that a swap of `Implement-now delegates may run while other reads continue.` (LP46) and of the `Route their line through` sentence (LP55) raises its diagnostic.
- [ ] `bench test --package ./internal/conformance`, `bench test --package ./internal/anchors`, and `bench test --check guidance-prose-budgets` pass.
