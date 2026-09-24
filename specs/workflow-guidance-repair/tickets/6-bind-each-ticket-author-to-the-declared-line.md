# Bind each ticket author to the declared line

Blocked by: 5-capture-the-retro-by-the-tracked-rule.md
Writes: .agents/skills/bench-craft-line/SKILL.md, .agents/skills/bench-craft-spec/SKILL.md, internal/anchors/anchor_harness_diagnostics_test.go, internal/anchors/registry_calibration.go, internal/anchors/registry_calibration_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_debug_loop.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_retained_workflow.go, tests/canary/workflow-guidance-anchors/, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: GR68, GR69, GR70, GR71, GR72, GR73, GR74, GR75, GR106

## What to build

`craft-line` binds every ticket author to the spec's one declared line, and every ladder step obeys the `--delegate` tier rule. `craft-spec` asks for the implementation line in its approval table. Both skills are at their budgets, so each edit stays line-neutral or shrinks the file.

Make these changes in `.agents/skills/bench-craft-line/SKILL.md`:

- Replace the per-story ceiling and the per-ticket table re-run with one sentence: every ticket author runs on the spec's declared `Line:`.
- Remove the per-story collapse sentence.
- Make step 3 escalate under the step 2 tier-move rule.
- Limit the step 5 top-tier pause to work outside `--delegate`.
- Keep the needle "Outside `--delegate`, a tier move of a fresh ticket author asks the reviewer first." byte for byte.

In `.agents/skills/bench-craft-spec/SKILL.md`, change "stories and their lines" in the Template approval table sentence to name the implementation line. Replace its RequireInSection needle, and keep that row's diagnostic byte for byte. `fixture_bite_test.go` pins that diagnostic. Keep "ownership fences with an explicit reviewer disposition, and out of scope.", which the canary `write-spec-fence-approval` mutates. Run `bench test --package ./internal/conformance --run 'TestWorkflowCadenceAnchorsRejectDeletionAndSwap|TestSpecTicketHandoffWorkflowFixturesAreComplete'` in the focused checks.

Add each planned Require and Forbid row in the registry.

## Acceptance

- [ ] `craft-line` binds every ticket author to the declared line (GR70).
- [ ] `craft-line` step 3 obeys the step 2 tier-move rule (GR74), and step 5 names the `--delegate` exception (GR75).
- [ ] The step 2 needle stays in place (GR106).
- [ ] `craft-spec`'s approval table names the implementation line (GR72).
- [ ] The skills contain none of the sentences that GR68, GR69, GR71, and GR73 name.
- [ ] `bench test --package ./internal/conformance --run 'TestWorkflowCadenceAnchorsRejectDeletionAndSwap|TestSpecTicketHandoffWorkflowFixturesAreComplete'` passes, and the canary `write-spec-fence-approval` stays red under its mutation.
- [ ] Both skills stay inside their budgets.
