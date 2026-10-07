# Share recorded explicit member traversal

Blocked by: none
Writes: internal/worktree/cleanup_members.go (new), internal/worktree/cleanup_members_test.go (new), internal/worktree/clean_set.go, internal/worktree/clean_set_apply.go, internal/worktree/clean_set_apply_test.go, internal/worktree/clean_set_refusal_test.go, internal/worktree/clean_set_outcomes_test.go, internal/worktree/clean_set_wiring_test.go, internal/worktree/cleanup_members_baseline_test.go (new), internal/worktree/testdata/cleanup-member-traversal/ (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: R01, R02, R03, R04, R05, R06, R08, R09, R10, R11, R12, R13, R14, R15, R19, R24, R27, R28

## What to build

Deliver C1: recorded explicit cleanup uses one package-local ordered traversal.
Create the seam and its unit table while migrating applyExplicitSet's recorded visits.
Leave applyUnrecordedRows as the existing successor adapter until C2.
Keep the four existing adapter entry points.
Remove the recorded member loop in this ticket.

The shared owner accepts ordered members, retained pass-through, the reached operation, and the adapter's stopped-result projection.
It accumulates completed results and stops at the first failed visit.
It derives no fingerprint, ownership proof, classification, or assignment identity.
A reached operation returns its current outcome and error.
The accepted C1 interface is the contract that C2 through C4 consume.

Keep explicit all-member preflight before either stage performs an effect.
Keep StepMemberRequalify, later requalification, and applyExplicitWith in their current sequence.
The locked, receipted applyCleanupTransaction remains authoritative.
A recorded failure projects both recorded and unrecorded suffix members.
Call the unrecorded adapter only after recorded visits succeed.

Before the first production edit, capture the complete original case manifest for all modes once.
Use the existing fixture owners and the spec's explicit identity and digest mappings.
C1 checks recorded cases while C2 through C4 add their mode comparisons.
No ordinary test restores source or rewrites the frozen records.
The captured baseline remains the predecessor value consumed by later tickets.

The mandatory registry paths preserve current dispatch, help, AXI membership, and routing.
This ticket introduces no command or rendered schema change.

## Acceptance

- [ ] An empty or retained-only recorded selection performs no qualification or removal, and retained rows remain byte-equivalent (R01, R04, R05).
- [ ] TestCleanupTraversalOrder and TestCleanupTraversalStopsAtFailedVisit prove order and the stop rule with exact call traces (R02, R14).
- [ ] TestCleanSetPreflightAllRows and explicit later-offender tests refuse before any removal despite a retained leading row (R06, R08).
- [ ] TestCleanSetLateDrift and TestCleanSetMemberDriftAfterPreflight preserve later and locked refusal, including retained and not-attempted internal outcomes (R09, R10, R12, R13).
- [ ] TestCleanSetPartialApply and TestCleanSetUnstartedOutcomes retain the failed target, completed prefix, and both suffix families (R11, R15, R19).
- [ ] TestDiscardTargetRequalifiesAfterARecordedRemoval still retains the newly unsafe branch through the unchanged unrecorded successor (R03).
- [ ] Explicit stale rendering, existing fault boundaries, and authority refusals match the recorded baseline cases (R24, R27, R28).
- [ ] The production explicit path calls the shared owner before C2 exists, and its existing entry points remain callable.
- [ ] Skipped preflight, later qualification, locked qualification, continued visits, and changed suffix projection each produce an observed behavioral red before restoration.

## Checkpoint verification

Run C1's verification commands from the spec's completion plan.
The focused producer checks run at this checkpoint while successors remain unbuilt.
Record the named red, restoration, green, and baseline comparison in the review pickup.
Use bench probe for each behavioral omission or swap after implementation.
A compile failure does not prove a preservation property.

The orchestrator records the green lane commit and freezes C1's delta for independent review.
That review accepts Standards, Spec, and Coverage before the next chunk starts.
