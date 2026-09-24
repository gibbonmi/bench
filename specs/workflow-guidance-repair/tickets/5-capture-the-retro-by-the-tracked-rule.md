# Capture the retro by the tracked rule

Blocked by: 4-state-the-lane-and-landing-gate-once.md
Writes: .agents/commands/bench-final-check.md, .bench/BENCH-reference.md, internal/anchors/registry_calibration.go, internal/anchors/registry_calibration_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_retained_workflow.go, tests/canary/workflow-guidance-anchors/, tests/canary/docs-currency-token-diet/, tests/canary/skills-index-command-adapters/, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: GR59, GR60, GR61, GR62, GR63, GR64, GR65, GR66, GR67, GR117, GR120

## What to build

The final check writes the retro through `bench retro` and commits it by the capture writer's tracked-or-ignored rule. A tracked retro and its scorecard updates commit with the phase close. An ignored retro stays local until the next reviewer-approved capture drain.

Make these changes in `.agents/commands/bench-final-check.md`:

- Replace the full-rewrite rule. Read `bench retro <slug> --scaffold`, then write the retro once with `bench retro <slug> --body <markdown>`.
- Replace the heading block with a pointer to the scaffold, which owns the headings.
- Remove the calibration column list, and point to the scaffold's table. Keep the duty of one table row per labeled claim, and keep the needle for the Brier mean statement.
- Replace "Do not run another gate or commit just to capture the retro" and the Report copy of that rule with the tracked-or-ignored rule.
- Keep "These files are pending capture for `/bench-drain`, not a second roadmap." The drain still verdicts every retro.
- Remove "The retro leaves through the next reviewer-approved capture drain." A tracked retro enters history at the phase close.

In `.bench/BENCH-reference.md`, keep the needles "/bench-final-check` writes `capture/retros/<spec-slug>.md`" and "/bench-drain` owns their reviewed drain". Replace "and its capture commit" with the tracked-or-ignored rule.

Replace each Require row whose sentence goes with the Forbid row that the spec names, and add each planned Require row. Add a Require row for the one-row-per-claim duty, and retarget the canary `calibration-retro-table-duty` to it.

## Acceptance

- [ ] The final check commits a tracked retro and its scorecard updates with the phase close (GR61).
- [ ] The final check leaves an ignored retro local until the next capture drain (GR62).
- [ ] The final check writes the retro once through `bench retro` after the scaffold (GR65).
- [ ] The guidance contains none of the sentences that GR59, GR60, GR63, GR64, GR66, and GR67 name.
- [ ] The final check keeps one calibration row per labeled claim (GR117).
- [ ] The final check does not contain the drain-only exit sentence that GR120 names.
- [ ] The canary `calibration-retro-table-duty` stays red under its mutation.
