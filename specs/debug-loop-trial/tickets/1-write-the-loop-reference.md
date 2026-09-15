# Write the loop reference

Blocked by: none
Writes: .agents/skills/bench-craft-tdd/references/loop.md (new), .agents/skills/bench-craft-tdd/SKILL.md, .agents/commands/bench-debug.md, CONTEXT.md, projects/benchkit.md, CHANGELOG.md, internal/anchors/registry_debug_loop.go (new), internal/anchors/registry_data.go, internal/conformance/registry_test.go, tests/canary/workflow-guidance-anchors, tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted, tests/canary/docs-currency-token-diet/signal-vocabulary-drift, tests/canary/guidance-prose-budgets/over-budget-skill, tests/canary/line-routing/line-binding-prose-drift, tests/canary/skill-description-budgets, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: DL1, DL2, DL3, DL4, DL5, DL6, DL7, DL38, DL39, DL42, DL44

## What to build

Write the loop reference with the four moves and the four-property checkbox list. Point `craft-tdd`'s cycle at it in one sentence, and raise the `craft-tdd` budget row to 124 in the same commit. Move the checkbox list out of `/bench-debug` Phase 1. Make Phase 1 state that the command is complete when it meets the four properties in the reference. Keep the one-command rule, the loop-constructions pointer, and the stop-gate sentence in Phase 1.

Add the **loop** entry to `CONTEXT.md`. Create the new anchor registry file, append it to the ordered registry, and list it in the canary family registry. Move the checkbox anchor's file to the reference. Add one fixture per new anchor row, and run the fixture-bite test in the focused checks. Add the `Debug loop trial` heading to `CHANGELOG.md` with this ticket's entry.

## Acceptance

- [ ] `loop.md` states the four moves in order, one sentence each, and a fixture that drops any one move reds `docs-currency-workflow`.
- [ ] `loop.md` carries the checkbox list, and a fixture that restores the list to `bench-debug.md` reds the check.
- [ ] `bench-debug.md` Phase 1 names `references/loop.md`, and a fixture that drops the pointer reds the check.
- [ ] `craft-tdd` points at `references/loop.md` from its cycle, and a fixture that drops the pointer reds the check.
- [ ] `CONTEXT.md` defines **loop**, and a fixture that drops the entry reds the check.
- [ ] Every new fixture bites through its registered owner and restores.
- [ ] `craft-tdd`, `bench-debug.md`, and every other budgeted file pass `guidance-prose-budgets`.
