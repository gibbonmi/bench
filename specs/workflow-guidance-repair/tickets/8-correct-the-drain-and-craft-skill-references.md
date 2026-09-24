# Correct the drain and craft skill references

Blocked by: 7-point-the-guides-at-each-fact-owner.md
Writes: .agents/commands/bench-drain.md, .agents/commands/bench-assess.md, .agents/skills/bench-craft-skills/SKILL.md, .agents/skills/bench-craft-cli/SKILL.md, .agents/skills/bench-craft-spec/SKILL.md, .agents/skills/bench-craft-tdd/references/tests.md, .agents/skills/bench-craft-adr/SKILL.md, internal/conformance/recurrence_maintenance_contract_test.go, internal/anchors/anchor_harness_diagnostics_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_debug_loop.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_retained_workflow.go, tests/canary/workflow-guidance-anchors/, tests/canary/row-next-grammar/, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: GR92, GR93, GR94, GR95, GR96, GR97, GR98, GR99, GR100, GR101, GR102, GR103, GR104, GR105, GR122

## What to build

The drain and assess commands name reachable routes, and each craft skill points to the owner of a fact that another file owns. This ticket is the last one that edits a budgeted file, so it carries the budget row.

Make these changes:

- In `bench-drain.md`, name the `Occurrences:` line in `roadmap/FT<n>.md` for each pending pair.
- In `bench-drain.md`, name `.bench/BENCH.md` as the owner of the batch approval rule.
- In `bench-drain.md`, run `bench handoff` as the last write, and date the `main` section by the handoff file's write time.
- In `bench-assess.md`, park an item with `bench idea` only.
- In `craft-skills`, point each phase adapter trigger to the invocation-policy account under "Harness Invocation" in `.bench/BENCH-reference.md`. Do not claim a table there, because the account is prose.
- In `craft-cli`, add the ambiguous-name re-query disclosure to the `bench consumers` row, and keep the row on one table line.
- In `craft-spec`, replace the literal story count with a pointer to the maximum that `bench coverage --check` prints. Remove the chunk-contract copy after its pointer clause. Keep the bytes that `fixture_bite_test.go` pins in this skill. Run `bench test --package ./internal/conformance --run 'TestWorkflowCadenceAnchorsRejectDeletionAndSwap|TestSpecTicketHandoffWorkflowFixturesAreComplete'` in the focused checks.
- In the tests reference of `craft-tdd`, state that a project's own test-expectation standard in `AGENTS.md` overrides this default.
- In `craft-adr`, replace the no-paths bullet with a pointer to invariant 3 of `.bench/BENCH.md`.

Change the contract strings and the bite table in `checkRecurrenceMaintenanceContract` together with the drain sentences. Add each planned Require and Forbid row in the registry.

## Acceptance

- [ ] The drain names the `Occurrences:` line in `roadmap/FT<n>.md` (GR92).
- [ ] The drain names `bench handoff` as its last write (GR95) and dates the `main` section by write time (GR96).
- [ ] `craft-skills` points each adapter trigger to the invocation-policy account (GR99).
- [ ] `bench test --package ./internal/conformance --run 'TestWorkflowCadenceAnchorsRejectDeletionAndSwap|TestSpecTicketHandoffWorkflowFixturesAreComplete'` passes.
- [ ] `craft-cli` names the ambiguous-name re-query disclosure (GR100).
- [ ] The tests reference defers to the project standard (GR103), and `craft-adr` points to invariant 3 (GR104).
- [ ] The files contain none of the sentences that GR93, GR94, GR97, GR98, GR101, GR102, and GR122 name.
- [ ] Every budgeted file that this spec edits stays inside its budget (GR105).
