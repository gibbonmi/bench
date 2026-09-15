# Add the debug trial flag

Blocked by: 1-write-the-loop-reference.md, 2-keep-the-fix-in-the-loop-session.md
Writes: .agents/skills/bench-craft-tdd/references/loop-trial.md (new), .agents/commands/bench-implement-spec.md, .agents/commands/bench-write-spec.md, .agents/commands/bench-review-implementation.md, projects/benchkit.md, CHANGELOG.md, internal/anchors/registry_debug_loop.go (new), tests/canary/workflow-guidance-anchors, tests/canary/skills-index-command-adapters/adapter-inert-invocation-key, tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy, tests/canary/guidance-prose-budgets/over-budget-skill, tests/canary/line-routing/line-binding-prose-drift, tests/canary/skill-description-budgets, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: DL14, DL15, DL16, DL17, DL18, DL19, DL20, DL21, DL22, DL23, DL24, DL25, DL26, DL27, DL28, DL29, DL30, DL31, DL32, DL33, DL34, DL35, DL36, DL37, DL40, DL41, DL43

## What to build

The first ticket creates the anchor registry file, so its marker here satisfies the preflight at the spec tip. This ticket extends that file.

Write the trial reference with one sentence per rule, and open it with the off arm. Then state the three on arms, the record contract, the quality measures, and the tolerance rule. Add the plan template as one fenced JSON block, the compare command, and the two exit sentences. Add one backticked `` `--debug-trial <on|off>` `` section to each build phase. Its four lines are blank, heading, blank, and one line with the grammar sentence and the pointer sentence. Raise the `bench-implement-spec.md` budget row to 84 and the `bench-write-spec.md` row to 77 in the same commit.

Add the Require anchors for the three pointer sections and the trial reference's sentences. Add the Forbid anchors that keep `references/loop.md` out of the three phase files. Add one fixture per anchor row. Add this ticket's entry under the `Debug loop trial` changelog heading.

At the chunk review, copy the plan template to a scratch path and fill it. Import two synthetic runs on one revision, and run `bench assessment compare` against the plan. Record the report in the review pickup.

## Acceptance

- [ ] Each build phase carries the flag section, and a fixture that drops a pointer reds `docs-currency-workflow`.
- [ ] Each flag section states that an unrecognized or missing value stops the phase with the grammar.
- [ ] A fixture that writes `references/loop.md` into a build phase reds the check.
- [ ] A fixture that drops any one arm sentence, from DL17 to DL26 or DL43, reds the check.
- [ ] A fixture that drops any one record, measure, tolerance, unknown, or exit sentence, from DL27 to DL37 or DL40 to DL41, reds the check.
- [ ] The plan template validates, and two synthetic runs compare with `descriptive evidence only; not default-change evidence` as the report's first limit line.
- [ ] Every new fixture bites through its registered owner and restores.
- [ ] `bench-implement-spec.md` and `bench-write-spec.md` pass `guidance-prose-budgets` at their raised rows.
