# 3. Guide relevant raw reads

Blocked by: 1-select-worktrees.md, 2-select-histories.md
Writes: .agents/skills/bench-craft-cli/SKILL.md, .agents/commands/bench-debug.md, .agents/commands/bench-drain.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_debug_loop.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_retained_workflow.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted, tests/canary/row-next-grammar/token-table-lacks-kit-edit, tests/canary/workflow-guidance-anchors/debug-archaeology-anchor, tests/canary/workflow-guidance-anchors/debug-phase1-stop-gate-softened, tests/canary/workflow-guidance-anchors/debug-red-commit, tests/canary/workflow-guidance-anchors/debug-reproduction-economics-deleted, tests/canary/workflow-guidance-anchors/dg-1, tests/canary/workflow-guidance-anchors/dg-1-contradiction, tests/canary/workflow-guidance-anchors/dg-2, tests/canary/workflow-guidance-anchors/dg-2-contradiction, tests/canary/workflow-guidance-anchors/dg-3, tests/canary/workflow-guidance-anchors/dg-3-blanket-ban, tests/canary/workflow-guidance-anchors/dg-3-contradiction, tests/canary/workflow-guidance-anchors/dg-4, tests/canary/workflow-guidance-anchors/dg-4-contradiction, tests/canary/workflow-guidance-anchors/dg-5, tests/canary/workflow-guidance-anchors/dg-5-contradiction, tests/canary/workflow-guidance-anchors/dg-5-dirty, tests/canary/workflow-guidance-anchors/dg-5-dirty-contradiction, tests/canary/workflow-guidance-anchors/dg-6, tests/canary/workflow-guidance-anchors/dg-6-command, tests/canary/workflow-guidance-anchors/dg-6-digest, tests/canary/workflow-guidance-anchors/dg-6-dirty, tests/canary/workflow-guidance-anchors/dg-6-surface, tests/canary/workflow-guidance-anchors/drain-anchor, tests/canary/workflow-guidance-anchors/drain-implement-now-commit, tests/canary/workflow-guidance-anchors/drain-implement-now-per-spec-exception, tests/canary/workflow-guidance-anchors/drain-implement-now-route, tests/canary/workflow-guidance-anchors/drain-implement-now-row-fallback, tests/canary/workflow-guidance-anchors/drain-implement-now-second-exception, tests/canary/workflow-guidance-anchors/drain-roadmap-context-anchor, tests/canary/workflow-guidance-anchors/drain-spec-history-anchor, tests/canary/workflow-guidance-anchors/drain-split-board-detail-owner, tests/canary/workflow-guidance-anchors/drain-split-board-retirement-pair, tests/canary/workflow-guidance-anchors/drain-split-board-row-detail-owner, tests/canary/workflow-guidance-anchors/implementation-retro-drain-anchor
Covers: QU11, QU12, QU13

## What to build

Add focused examples for file, Git, test, shell, archive, and log reads.
Use the new selected views where the examples repeat one domain intent.
Keep existing Bench response ownership and authority boundaries unchanged.
The examples must not invent an unsupported CLI flag or a numeric default.
Keep the drain needle ``use `bench spec history <slug>` for the shipped-row check``.

Before the first edit, run `bench anchors` on each of the three guidance files.
The anchor registries and the `dg-*`, `debug-*`, and `drain-*` canaries are in the fence only for an anchor that the edit moves.
Do not weaken an anchor or a canary.

Read the stream reference source for the example wording, and verify it against the current guidance bytes:

- `194b7dba:.agents/commands/bench-debug.md`
- `194b7dba:.agents/commands/bench-drain.md`
- `194b7dba:.agents/skills/bench-craft-cli/SKILL.md`

## Acceptance

- [ ] Each bounded example states its available complete-detail route.
- [ ] Archive and log examples consolidate discovery without concatenating full bodies.
- [ ] Polling and authority-changing operations remain separate.
- [ ] The guidance retains every existing anchored rule.
