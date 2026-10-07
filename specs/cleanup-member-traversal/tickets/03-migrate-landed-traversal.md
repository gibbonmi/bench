# Migrate landed and sibling cleanup traversal

Blocked by: 02-migrate-unrecorded-traversal.md
Writes: internal/worktree/clean_landed.go, internal/worktree/clean_landed_apply_test.go, internal/worktree/clean_set_apply_test.go, internal/worktree/clean_set_refusal_test.go, internal/worktree/clean_set_outcomes_test.go, internal/worktree/land_effects_cleanup_test.go, internal/worktree/cleanup_members_baseline_test.go (new), internal/worktree/testdata/cleanup-member-traversal/ (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: R01, R02, R04, R05, R07, R08, R09, R10, R11, R12, R13, R14, R17, R19, R23, R24, R27, R28

## What to build

Deliver C3: landed command cleanup and folded-sibling cleanup consume the accepted C1 traversal.
Keep applyLandedSet's existing entry point and scope argument.
Remove its original loop with its command and sibling evidence.
C2 is the serial blocker because the baseline and mandatory registry paths are shared.

Keep all-member landed preflight before the first effect.
Keep StepMemberRequalify and the fresh tuple comparison before each removable visit.
Keep the locked planner and terminal callback in applyCleanupTransaction.
They preserve locked drift and the post-settlement fault boundary.
A failure projects landed unreached outcomes and retains every completed effect.

cleanLandedSiblings keeps its narrowed base scope and existing landed adapter call.
Do not edit its caller or broaden its selection.
Use the frozen C1 records for landed command and sibling outcomes.
Keep head, tip, and checked-out-ref preservation scenarios distinct.

The mandatory registry paths preserve current dispatch, help, AXI membership, and routing.
This ticket introduces no command or rendered schema change.

## Acceptance

- [ ] Empty, retained-only, and mixed landed selections preserve order and unchanged retained rows without retained qualification or effects (R01, R02, R04, R05).
- [ ] TestCleanSetPreflightAllRows and TestCleanSetPreflightFaultOnLaterMember refuse before effects and name the actual offender after retained rows (R07, R08).
- [ ] TestCleanLandedApplyReplansEachRowBeforeMutation and TestCleanSetLateDrift retain unsafe members at outer and locked checks (R09, R10).
- [ ] Internal retained and not-attempted drift outcomes preserve their distinction from fault rows (R11, R12, R13).
- [ ] TestCleanLandedApplyStopsAfterCompletedRowFault and TestCleanupLandedUnreachedSuffix preserve completed rows and append the selected unreached suffix (R14, R17, R19).
- [ ] TestLandCleansTheFoldedSibling, TestLandLeavesAPriorLandedAssignment, and TestLandRetainsAnUnprovenSibling retain narrowed scope (R23).
- [ ] Landed stale remedies, fault boundaries, ownership refusals, and durable state match the frozen cases (R24, R27, R28).
- [ ] Skipped preflight, skipped locked qualification, and prefix-only landed projection each produce an observed behavioral red before restoration.

## Checkpoint verification

Run C3's verification commands from the spec's completion plan.
The focused producer checks run at this checkpoint while successors remain unbuilt.
Record the named red, restoration, green, and baseline comparison in the review pickup.
Use bench probe for each behavioral omission or swap after implementation.
A compile failure does not prove a preservation property.

The orchestrator records the green lane commit and freezes C3's delta for independent review.
That review accepts Standards, Spec, and Coverage before the next chunk starts.
