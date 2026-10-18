# Delete the proposal and breakdown forms of build preflight

Blocked by: none
Writes: internal/preflight/proposal.go (deleted), internal/preflight/preparation.go, internal/preflight/command.go, internal/preflight/command_build_test.go, internal/preflight/evidencecmd/operations.go, internal/preflight/evidencecmd/bound_test.go, internal/preflight/charge_cells_test.go, internal/preflight/fence_writes_test.go, internal/preflight/proposal_test.go (deleted), internal/preflight/proposal_command_test.go (deleted), internal/preflight/proposal_edges_test.go (deleted), internal/preflight/proposal_readonly_test.go (deleted), internal/preflight/ownership/authority.go (deleted), internal/preflight/ownership/authority_test.go, internal/preflight/ownership/breakdown.go, internal/preflight/ownership/ownership.go, internal/preflight/ownership/proposal_test.go (deleted), internal/preflight/ownership/source_test.go, internal/preflight/ownership/fieldwrites_test.go, internal/preflight/ownership/effects_test.go, internal/preflight/ownership/premises_test.go, internal/preflight/ownership/kitrow_test.go, internal/preflight/ownership/ownershiptest/repo.go, internal/preflight/staleness/identity_test.go, internal/preflight/staleness/drift_test.go, cmd/bench/help_inventory_test.go, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, .agents/commands/bench-write-spec.md, .agents/skills/bench-craft-tickets/references/slicing-checks.md, .agents/skills/bench-craft-tickets/SKILL.md, .agents/skills/bench-craft-spec/references/map-only-route.md, .agents/skills/bench-craft-spec/references/map-discipline.md, docs/known-issues.md, CHANGELOG.md, .bench/commitment.json, tests/canary/skills-index-command-adapters/adapter-inert-invocation-key, tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy, tests/canary/workflow-guidance-anchors/command-handoff-anchor, tests/canary/workflow-guidance-anchors/write-spec-handoff-anchor, tests/canary/workflow-guidance-anchors/write-spec-post-slicing-handoff, tests/canary/workflow-guidance-anchors/write-spec-review-tier-escalated, tests/canary/workflow-guidance-anchors/ticket-template-anchor, tests/canary/workflow-guidance-anchors/reader-sweep-term, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/anchors/registry_data.go, internal/anchors/registry_decision_maps.go, internal/anchors/registry_decision_maps_test.go, internal/anchors/registry_ft311_preparation.go, internal/anchors/registry_retained_workflow.go, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/anchors/registry_data_test.go, internal/preflight/anchor_closure_test.go, internal/preflight/ownership/parameter_grammar_test.go
Covers: CR07, CR08, CR09, CR10, CR80, CR81, CR82, CR83, CR84, CR85, CR91, CR93, CR95

## What to build

`bench preflight build` loses its `--propose-writes` and `--breakdown` forms. The evidence operations table, the command dispatch, the prepared-command path, and the help drop both forms. A call with either flag gets the grammar's usage refusal at exit 2. The proposal renderer, the uncited projection, and `ownership.Unordered` go with their only reader. `ownership.UncitedIDs` stays, because `rows-owned` reads it.

Every test that drove a removed form goes or loses that case. A new test in the build command tests drives both removed flags and asserts the usage refusal.

The guidance drops each step that runs a removed form. The slicing checks lose the whole `Enforced rules` section, because each rule there repeats a check that the ticket parser or build preflight owns. The craft-tickets pointer to that section goes or names the checks themselves. `/bench-write-spec` drops the `--propose-writes` comparison and its require anchor. The map-only route and map discipline drop the proposal and breakdown steps. The known-issues row whose retire condition names the writes proposal goes. The unreleased changelog entry that adds the closure forms gives way to an entry that states their removal.

The ticket adds one forbid anchor for each retired command sentence in the ticket-slicing registry and its mirror test. The `require-in-section` needle on `slicing-checks.md` "The slicer runs `bench preflight build <slug>` after each change..." stays byte-identical.

The spec-stage close commits this spec's commitment binding inside the range that the plan-only preflight grades, so this ticket's `Writes:` line names `.bench/commitment.json`. The build itself does not edit that file.

Contract for ticket 02: this ticket removes every reader of `Facts.KitObligations`, `ownership.OwnedPaths`, `Snapshot.Requirements`, and `Snapshot.Uncited` outside the closure, so ticket 02 can delete them. The uncited call in `ownership.Collect` leaves with the uncited projection.

## Acceptance

- [ ] `bench preflight build s --propose-writes --ticket a.md --base B --source-tip T` exits 2 with a usage refusal.
- [ ] `bench preflight build s --breakdown` exits 2 with a usage refusal.
- [ ] `bench help` prints neither removed form, and the complete help inventory test is green.
- [ ] `slicing-checks.md` holds no `## Enforced rules` heading, no `--propose-writes` text, and no `--breakdown` text, and a forbid anchor reds each one.
- [ ] `bench-write-spec.md`, `map-only-route.md`, and `map-discipline.md` hold no removed-form step, and a forbid anchor reds each one.
- [ ] `docs/known-issues.md` holds no `writes proposal` text, and a forbid anchor reds it.
- [ ] The `[Unreleased]` changelog section holds no entry that adds `--breakdown`.
- [ ] Each edited budgeted guidance file holds a line count at or below its budget row.
