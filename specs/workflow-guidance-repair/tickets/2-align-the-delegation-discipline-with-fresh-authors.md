# Align the delegation discipline with fresh authors

Blocked by: 1-align-the-delegation-skill-with-fresh-authors.md
Writes: .agents/skills/bench-craft-delegate/references/delegation-discipline.md, .agents/skills/bench-craft-line/references/bounded-repair-policy.md, internal/anchors/registry_calibration.go, internal/anchors/registry_calibration_test.go, internal/anchors/registry_charge_binding.go, internal/anchors/registry_charge_binding_test.go, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_ft311_preparation.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_retained_workflow.go, internal/conformance/implementation_continuation_test.go, internal/conformance/retained_workflow_test.go, tests/canary/workflow-guidance-anchors/, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: GR17, GR18, GR19, GR20, GR21, GR22, GR23, GR24, GR25, GR26, GR27, GR28, GR29, GR30, GR31, GR32, GR115

## What to build

The delegation discipline and the bounded repair policy state the fresh-author model, the probe verb, and reachable recovery steps. The registry refuses each retired sentence. Ticket 1 supplies the skill wording that this reference points back to.

Make these changes in `.agents/skills/bench-craft-delegate/references/delegation-discipline.md`:

- Add one bullet under "Delegated author transfer": outside `--delegate`, a model change under any trigger is a tier move, and it asks the reviewer first. Keep the four trigger sentences byte for byte, so their canaries stay valid.
- Change "A user-directed write delegate treats `Writes:` as an expectation." so that every write delegate treats `Writes:` as an expectation.
- Replace the chunk-union repair fence with the affected ticket's `Writes:` line, widened only through `.bench/BENCH.md`'s plan-expansion policy.
- Remove the copy-aside probe rule. Keep the bytes "Before the coordinator reads a probe verdict, the coordinator confirms the", and tie that confirmation to the `restored` cell that `bench probe` reports.
- Replace the `cmp` restore rule for a shared-worktree reader with a `bench probe` rule under "Read-only returns". The reader probes only through `bench probe` and reads its `restored` cell.
- Remove the main-checkout build exception.
- Replace `bench resume-clean` with the plan and apply forms of `bench worktree clean`.
- Change the digest rule so that a `bench worktree merge`, not a rebase, triggers repeated verification.
- Remove the after-commit preflight copy and the attempt-definition copy, and point to their owners.

In `.agents/skills/bench-craft-line/references/bounded-repair-policy.md`, apply the allowance to every implementation run, the light path included. Remove the retired mode list.

Replace each Require row whose sentence goes with the Forbid row that the spec names, and add each planned Require row. Update the `TestRetainedWorkflow` row for `Writes:` and the `TestImplementationContinuation` row for the mode list. Retarget each canary whose bytes change.

## Acceptance

- [ ] The discipline states that a model change under a trigger outside `--delegate` asks the reviewer first (GR17).
- [ ] The discipline makes every write delegate treat `Writes:` as an expectation (GR19).
- [ ] The discipline makes the repair fence the affected ticket's `Writes:` line (GR21).
- [ ] The discipline repeats verification after a `bench worktree merge` changes the source digest (GR28).
- [ ] The discipline contains none of the sentences that GR18, GR20, GR22, GR23, GR24, GR25, GR26, GR27, GR29, and GR30 name.
- [ ] The discipline makes a shared-worktree reader probe only through `bench probe` and read its `restored` cell (GR115).
- [ ] The bounded repair policy applies to every implementation run (GR32), and its retired mode list is absent (GR31).
- [ ] The four transfer-trigger canaries and `delegate-probe-mutated-bytes` stay red under their mutations.
