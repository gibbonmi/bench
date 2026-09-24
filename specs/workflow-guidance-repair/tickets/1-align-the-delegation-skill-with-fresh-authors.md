# Align the delegation skill with fresh authors

Blocked by: none
Writes: .agents/skills/bench-craft-delegate/SKILL.md, .claude/agents/bench-writer.md, internal/anchors/registry_calibration.go, internal/anchors/registry_calibration_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_debug_loop.go, internal/anchors/registry_ft311_preparation.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_retained_workflow.go, internal/conformance/implementation_continuation_test.go, tests/canary/workflow-guidance-anchors/, tests/canary/claude-agent-definitions/, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: GR1, GR2, GR3, GR4, GR5, GR6, GR7, GR8, GR9, GR10, GR11, GR12, GR13, GR14, GR15, GR16, GR109, GR110, GR111, GR112, GR113

## What to build

`craft-delegate` states the fresh-author model, and the anchor registry refuses each retired sentence. The skill keeps its 126-line budget, so each edit stays line-neutral or shrinks the file.

Make these changes in `.agents/skills/bench-craft-delegate/SKILL.md`:

- Every write delegate runs as `bench-writer`: a fresh ticket author, a repair session, and a user-directed write delegation.
- Replace "A change of implementation author or session requires user direction." with a pointer that names `craft-line` as the owner of a model or session change.
- Scope the coverage-row charge to every write charge from a spec, not only a user-directed one.
- State that a ticket author commits its ticket on a lane pass. Remove the coordinator `bench commit` clause and the clause "a write-delegate stops at diff-ready".
- Re-key the charge example to one ticket of its spec, because a ticketless charge returns an uncommitted diff. Replace "Stop at diff ready;" with a commit on a lane pass. Keep " and the mutation's kind (omission or swap)".
- Replace the dependency-only share rule with the rule that a build's serial ticket authors share its one integration source. Keep the needles "reviewed dependent tickets share one retained integration source" and "each charge names its root and expected tip".
- Scope the returned-tree coordinator probe to a user-directed write delegate. Keep the needle "The coordinator probe's mutation kind differs from the delegate author's mutation kind."
- Keep the `git stash` ban, and state that the guard refuses only `git stash drop` and `git stash clear`.
- Remove the second copies "A diagnostic helper can inspect evidence, but it receives no implementation or repair assignment." and "A spec-backed ticket goes to a fresh author session on its integration source, in `Blocked by:` order."

Make these changes in `.claude/agents/bench-writer.md`:

- Name a fresh ticket author and a repair session beside a user-directed write delegation in the description.
- Replace "The user directed this delegation" with a premise that fits all three roles.
- Replace the stop at an uncommitted diff with two rules. A ticket author or a repair author commits its ticket on a lane pass. A user-directed delegate with no ticket returns an uncommitted diff.
- Keep the rule that no role runs `bench worktree land`.

`internal/conformance/fixture_bite_test.go` pins exact bytes of this skill, with their line breaks. Keep the bytes that its `probe kind`, `stale-base coordinator ownership`, `delegate self-probe`, `probe site differs`, `probe kind vocabulary`, `registry tracing`, and `backup isolation` rows name. Run `bench test --package ./internal/conformance --run 'TestWorkflowCadenceAnchorsRejectDeletionAndSwap|TestSpecTicketHandoffWorkflowFixturesAreComplete'` in the focused checks.

In the registries, replace each Require row whose sentence goes with the Forbid row that the spec names. Add each planned Require row. Update the `TestImplementationContinuation` row for the diagnostic helper. Retarget the canaries `delegate-coverage-row-charge` and `delegated-per-ticket-author`. Each one names the new Require bytes or plants the forbidden bytes. Do not remove a canary directory.

Before this ticket starts, the reviewer grants a one-time permission rule for `.claude/agents/bench-writer.md`.

## Acceptance

- [ ] `craft-delegate` names `bench-writer` as the type of a fresh ticket author and a repair session (GR1), and the retired sentence is absent (GR2).
- [ ] The `bench-writer` description names a fresh ticket author and a repair session (GR3).
- [ ] The skill contains none of the retired or copied sentences that GR4, GR7, GR8, GR10, GR12, GR13, GR15, and GR16 name.
- [ ] The skill names `craft-line` as the owner of a model or session change (GR5).
- [ ] The skill makes every write charge from a spec carry its coverage rows (GR6).
- [ ] The skill says that a ticket author commits on a lane pass (GR9).
- [ ] The skill scopes the returned-tree probe to a user-directed write delegate (GR11).
- [ ] The skill states the guard's real `git stash` deny surface (GR14).
- [ ] The agent body tells a ticket or repair author to commit on a lane pass and not to land (GR109).
- [ ] The agent body and the skill contain none of the clauses that GR110, GR111, GR112, and GR113 name.
- [ ] The charge example names one ticket and commits on a lane pass (GR113).
- [ ] `bench test --package ./internal/conformance --run 'TestWorkflowCadenceAnchorsRejectDeletionAndSwap|TestSpecTicketHandoffWorkflowFixturesAreComplete'` passes.
- [ ] The `claude-agent-definitions` check stays green, and `.agents/skills/bench-craft-delegate/SKILL.md` stays inside its budget.
