# Route phase command repairs to fresh sessions

Blocked by: 2-align-the-delegation-discipline-with-fresh-authors.md
Writes: .agents/commands/bench-review-implementation.md, .agents/commands/bench-implement-spec.md, .agents/commands/bench-debug.md, .agents/commands/bench-final-check.md, internal/anchors/registry_calibration.go, internal/anchors/registry_calibration_test.go, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_debug_loop.go, internal/anchors/registry_front_door.go, internal/anchors/registry_front_door_test.go, internal/anchors/registry_ft311_preparation.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_retained_workflow.go, internal/conformance/retained_workflow_test.go, internal/conformance/docs_workflow_helpers_test.go, tests/canary/workflow-guidance-anchors/, tests/canary/skills-index-command-adapters/, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/docs-currency-token-diet/
Covers: GR33, GR34, GR35, GR36, GR37, GR38, GR39, GR40, GR41, GR42, GR43, GR44, GR116

## What to build

The phase commands send each repair to a fresh session for its ticket. They name the orchestrator as the final performer. They point to `.bench/BENCH.md` for the chunk-review rules. `.agents/commands/bench-implement-spec.md` and `.agents/commands/bench-debug.md` are at their budgets, so their edits stay line-neutral.

Make these changes:

- In `bench-review-implementation.md`, replace the single repair ticket with this rule: a coverage-map amendment updates each affected ticket's `Covers:` line under the plan-expansion policy.
- In `bench-review-implementation.md`, remove the copies of the review-repeat rule, the chunk-review start rule, and the final reconciliation rule. Point to `.bench/BENCH.md`.
- In `bench-implement-spec.md`, reduce the wrong-tier route to a pointer to `craft-line`'s ladder, and remove the review-repeat copy.
- In `bench-debug.md`, route an out-of-fence cause through the plan-expansion policy, or through the reviewer's split for a scope change.
- In `bench-final-check.md`, retain the orchestrator's final `integration-verification` results in the review record.
- In `bench-final-check.md`, send a spec-backed red to a fresh repair session under `.bench/BENCH.md`'s repair rule. Other work keeps the approve-then-fix route.

Replace each Require row whose sentence goes with the Forbid row that the spec names, and add each planned Require row. The chunk-review start rule keeps a Require guard on its owner. Add a Require row for the `.bench/BENCH.md` sentence "Every ticket contribution reaches the integrated chunk tip before that chunk's review begins." Move the `TestRetainedWorkflow` row for the chunk-review start rule to that sentence.

`checkReviewConvergenceContract` requires the three removed copies. Remove those three requirements, because the Forbid rows GR41, GR42, and GR44 and the `.bench/BENCH.md` Require rows now guard them.

Retarget the canaries `review-repair-ticket-owner`, `review-repair-ticket-covers`, and `delegated-chunk-tip-review`. The last one mutates the `.bench/BENCH.md` sentence. Do not remove a canary directory.

## Acceptance

- [ ] The review phase updates each affected ticket's `Covers:` line for a coverage-map amendment (GR34).
- [ ] The final check retains the orchestrator's `integration-verification` results (GR36).
- [ ] The final check sends a spec-backed red to a fresh repair session (GR37).
- [ ] The debug phase routes an out-of-fence cause through plan expansion or the reviewer's split (GR39).
- [ ] The registry requires the `.bench/BENCH.md` chunk-tip sentence (GR116).
- [ ] The phase commands contain none of the sentences that GR33, GR35, GR38, GR40, GR41, GR42, GR43, and GR44 name.
- [ ] `.agents/commands/bench-implement-spec.md` and `.agents/commands/bench-debug.md` stay inside their budgets.
