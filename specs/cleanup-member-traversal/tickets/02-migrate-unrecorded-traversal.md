# Migrate explicit unrecorded branch traversal

Blocked by: 01-share-recorded-traversal.md
Writes: internal/worktree/clean_discard.go, internal/worktree/clean_discard_test.go, internal/worktree/clean_discard_transaction_test.go, internal/worktree/clean_set_apply_test.go, internal/worktree/cleanup_members_baseline_test.go (new), internal/worktree/testdata/cleanup-member-traversal/ (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: R03, R04, R05, R06, R09, R11, R12, R13, R14, R16, R19, R22, R24, R27, R28, R29

## What to build

Deliver C2: explicit unrecorded branches use the accepted C1 traversal after recorded members finish.
Consume C1's unchanged ordered-operation and stop-projection contract.
Keep applyUnrecordedRows callable with its existing completed prefix.
Remove its original loop in this ticket.

Keep requalifyUnrecordedRow after earlier recorded effects.
Keep discardUnrecordedBranch as the effect owner and exact branch deletion as its final operation.
Only classUnique plans carry a recovery ref.
Landed and subsumed plans retain recovery none and create no discarded ref.
Do not add lifecycle fault hooks to this adapter.

A failure retains the completed recorded and branch prefix and projects every unreached member.
Retained rows remain unchanged and skip qualification and effects.
Use the frozen C1 records for mixed explicit and unrecorded cases.
Extend the permanent comparison without changing original baseline facts.

The mandatory registry paths preserve current dispatch, help, AXI membership, and routing.
This ticket introduces no command or rendered schema change.

## Acceptance

- [ ] TestDiscardTargetRequalifiesAfterARecordedRemoval retains a branch whose holder removal changed its qualification (R03, R09).
- [ ] Mixed retained and removable branches preserve preflight, unchanged retained rows, and zero retained visits (R04, R05, R06).
- [ ] TestCleanupTraversalUnrecordedSuffix preserves failed-member identity, drift outcomes, completed effects, and the unreached branch suffix (R11, R12, R13, R14, R16, R19).
- [ ] TestDiscardTargetWritesTheDiscardedRefFirst preserves the exact planned tip under the planned recovery ref before unique-branch deletion (R22).
- [ ] TestDiscardTargetLandedAndSubsumedWriteNoRef preserves recovery none and creates no discarded ref for either classification (R29).
- [ ] Discarded-ref symref, conflicting tip, moved branch tip, and checked-out branch witnesses retain their existing survivors and recovery state.
- [ ] Existing fault windows, stale command remedies, and preservation refusals match the frozen unrecorded and mixed-mode cases (R24, R27, R28).
- [ ] An omitted later requalification and unconditional recovery creation each produce an observed behavioral red before restoration.

## Checkpoint verification

Run C2's verification commands from the spec's completion plan.
The focused producer checks run at this checkpoint while successors remain unbuilt.
Record the named red, restoration, green, and baseline comparison in the review pickup.
Use bench probe for each behavioral omission or swap after implementation.
A compile failure does not prove a preservation property.

The orchestrator records the green lane commit and freezes C2's delta for independent review.
That review accepts Standards, Spec, and Coverage before the next chunk starts.
