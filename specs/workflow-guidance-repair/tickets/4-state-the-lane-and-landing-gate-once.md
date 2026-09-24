# State the lane and the landing gate once

Blocked by: 3-route-phase-command-repairs-to-fresh-sessions.md
Writes: .agents/skills/bench-craft-synthesis/SKILL.md, .agents/commands/bench-final-check.md, .bench/BENCH-reference.md, docs/adr/0014-main-receives-writes-only-through-landings.md, cmd/bench/main.go, cmd/bench/help_inventory_test.go, internal/commit/commit.go, internal/commit/dry_run_test.go, internal/anchors/registry_calibration.go, internal/anchors/registry_calibration_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_retained_workflow.go, tests/canary/workflow-guidance-anchors/, tests/canary/docs-currency-token-diet/, tests/canary/skills-index-command-adapters/, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/package-core-guard/
Covers: GR45, GR46, GR47, GR48, GR49, GR50, GR51, GR52, GR53, GR54, GR55, GR56, GR57, GR58, GR118, GR119

## What to build

Every guidance file and both help strings say one thing: `bench commit` runs the declared lane, and `bench worktree land` runs the whole-project gate. The landing enforcement, the landing rebuild, and the conflict repair read as the tree applies them.

Make these changes:

- In `craft-synthesis`, take the prose-only green verdict from the whole-project gate. That gate is the landing's, or `bench worktree exec <target> -- bench gate` for a batch that waits for approval.
- In `bench-final-check.md`, commit in the Bench worktree with `bench commit` and land through `bench worktree land`. Keep the needle "bench commit -m". Remove the lane-as-gate and one-command claims, the "gate-then-commit path" clause, and the clause "then runs the gate and commits only on green".
- In `bench-final-check.md`, run `bench spec retire <slug>` in a Bench worktree and land its `spec-retire: <slug>` commit.
- In `.bench/BENCH-reference.md`, replace "the rule is guidance, not a hook" with a pointer to `.bench/BENCH.md`'s enforcement sentence. Keep the needle "The spec is optional on the landing and on its resume".
- In `.bench/BENCH-reference.md`, remove the landing rebuild sentence, and name the reviewer as the person who runs the raw merge in the conflict repair.
- In ADR 0014's fifth paragraph, state that the commit verb refuses the primary checkout. Also state that a file-write guard refuses an agent's write to a tracked path there.
- In the `bench help` row for `commit`, name the declared lane, and the gate only when the project declares no lane. Update `TestHelpInventoryIsComplete`.
- In the `bench commit --help` `--dry-run` line, name the declared lane. Extend `TestHelpAdvertisesDryRun` to refuse the old gate wording and to require the lane.

Add each planned Require and Forbid row in the registry.

## Acceptance

- [ ] `craft-synthesis` takes its green verdict from the whole-project gate (GR46).
- [ ] The final check says that the landing runs the whole-project gate on lane-pass commits (GR50).
- [ ] The guidance contains none of the sentences that GR45, GR47, GR48, GR49, GR54, GR56, GR118, and GR119 name.
- [ ] `bench help` describes `commit` with the declared lane (GR51).
- [ ] `bench commit --help` drops "gate the exact composed snapshot" (GR52) and names the declared lane (GR53).
- [ ] ADR 0014 states the primary-checkout refusal and the file-write guard (GR55).
- [ ] The reference names the reviewer as the raw-merge runner (GR57).
- [ ] The final check retires a spec in a Bench worktree and lands the retirement (GR58).
