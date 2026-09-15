# Keep the fix in the loop session

Blocked by: 1-write-the-loop-reference.md
Writes: .agents/commands/bench-debug.md, .agents/skills/bench-craft-delegate/SKILL.md, CHANGELOG.md, internal/anchors/registry_debug_loop.go (new), tests/canary/workflow-guidance-anchors, tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted, tests/canary/claude-agent-definitions, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: DL5, DL6, DL7, DL8, DL9

## What to build

The first ticket creates the anchor registry file, so its marker here satisfies the preflight at the spec tip. This ticket extends that file.

Rewrite the "How it meets the rest of Bench" section of `/bench-debug`. Remove the sentence that routes code authorship through `craft-delegate`. State that the session that owns the loop writes the fix and delegates only a read-only fan-out search. State that a write delegate runs the phase through Phase 6 on a bug inside its fence.

State that on a bug outside its fence it runs Phases 1 to 3 and stops implementation edits. State that it keeps its in-fence work dirty and returns the bounded blocked report. State that the report carries the loop command, the red output digest, the ranked hypotheses, the failing surface, and the in-fence dirty paths. Keep the coordinator's reslice route, the worktree rule, the declared line, the seam ownership, and the gate-replacement rule.

Replace `craft-delegate`'s blocked-delegate sentence with one that names the three debug phases the delegate runs before it stops. Keep `bench-debug.md` inside 170 lines and `craft-delegate` inside 126. Add the anchors and one fixture per anchor row, and add this ticket's entry under the `Debug loop trial` changelog heading.

## Acceptance

- [ ] A fixture that restores the `craft-delegate` fix route to `bench-debug.md` reds `docs-currency-workflow`.
- [ ] A fixture that drops the loop-session-writes-the-fix sentence reds the check.
- [ ] A fixture that drops the in-fence permission reds the check.
- [ ] A fixture that drops the out-of-fence limit or the loop command from the blocked report reds the check.
- [ ] A fixture that drops the three debug phases from `craft-delegate`'s blocked-delegate sentence reds the check.
- [ ] Every new fixture bites through its registered owner and restores.
- [ ] `bench-debug.md` and `craft-delegate` pass `guidance-prose-budgets`.
