# Share cleanup member traversal

Status: staged
Roadmap: FT363
Decision source: ready compiled map [architecture-cleanup](decisions/architecture-cleanup.md), including its resolved decision tickets.
Verification log: 2 iteration(s) to accept — FT363-SPEC-R2 accepted the repaired spec before ticket slicing.
Ticket review: FT363-TICKETS-R1 accepted the four-ticket graph at static confidence 9.

## Problem

Cleanup selections repeat the ordered visit and stop mechanism.
The explicit adapter visits recorded members before its unrecorded branch members.
The landed adapter visits selected assignments.
The unclaimed adapter visits its replanned branch selection.
Their current loops live in [clean_set.go](../../internal/worktree/clean_set.go:325),
[clean_discard.go](../../internal/worktree/clean_discard.go:224),
[clean_landed.go](../../internal/worktree/clean_landed.go:333), and
[clean_unclaimed.go](../../internal/worktree/clean_unclaimed.go:191).

The FT363 premise that every mode reports unreached members is stale.
Unclaimed deletion returns its completed prefix and the failed row.
It does not append unreached rows.
The other adapters report unreached members.
[The resolved decision](decisions/architecture-cleanup/tickets/1.md) preserves this difference.

A common loop must preserve authority and effects.
A recorded member's removal can remove the proof that an unrecorded branch is subsumed.
The later branch must then survive.
[The existing witness](../../internal/worktree/clean_discard_transaction_test.go:165) covers this case.

## Solution

Give the adapters one package-local ordered traversal.
Keep each adapter's qualification, effect, and stopped-result policy.
Keep the current commands, response fields, exit status, fingerprints, and recovery routes.
This refactor makes the shared stop rule change in one place.
It makes no elapsed-time or reduced-test-count promise.

## User stories

Line: gpt-6.1-sol / high.
Implementation-line reason: the difficult work preserves effect order across distinct authority checks and failure adapters.
Harder chunks: C1, C2, and C4.

1. As a cleanup operator, I want an empty selection to finish without effects, so that a no-op remains safe.
2. As a cleanup operator, I want visits in the approved order, so that results describe the same selection.
3. As an explicit cleanup operator, I want recorded members visited before unrecorded members, so that later qualification sees earlier effects.
4. As a cleanup operator, I want retained rows to pass through unchanged, so that their verdict and digest remain valid.
5. As a cleanup operator, I want retained members skipped by qualification and removal, so that retention stays authoritative.
6. As an explicit cleanup operator, I want all removable members preflighted before removal, so that existing drift prevents avoidable effects.
7. As a landed cleanup operator, I want the same preflight guarantee, so that a later invalid member prevents earlier deletion.
8. As a cleanup operator, I want a preflight fault attributed to its actual member, so that a retained leading row does not hide the offender.
9. As a cleanup operator, I want each reached member requalified after earlier effects, so that a newly unsafe member survives.
10. As a recorded cleanup operator, I want qualification checked inside the locked transaction, so that drift after an outer check still refuses removal.
11. As a cleanup operator, I want a fault row to name the failed member, so that failure remains distinguishable from drift.
12. As a cleanup operator, I want drift to a retained verdict reported as retained internally, so that the fresh verdict survives.
13. As a cleanup operator, I want drift that remains removable reported as not attempted internally, so that no removal is claimed.
14. As a cleanup operator, I want processing to stop at the first failed visit, so that later effects do not occur.
15. As an explicit cleanup operator, I want unreached recorded members reported, so that partial completion is clear.
16. As an explicit cleanup operator, I want unreached unrecorded members reported, so that the mixed selection stays complete.
17. As a landed cleanup operator, I want unreached assignments reported, so that partial completion is clear.
18. As an unclaimed cleanup operator, I want only the completed prefix and failed deletion row, so that the existing result contract stays intact.
19. As a cleanup operator, I want completed effects to remain completed after a later fault, so that output never invents rollback.
20. As an unclaimed cleanup operator, I want the complete selection replanned before deletion, so that added or removed refs invalidate approval.
21. As an unclaimed cleanup operator, I want exact object identity checked at deletion, so that a moved branch survives.
22. As a branch cleanup operator, I want recovery for a unique discarded branch, so that its planned tip remains reachable.
23. As a landing operator, I want folded-sibling cleanup to use the same landed adapter, so that non-command cleanup preserves its scope.
24. As a cleanup operator, I want stale output to keep its refusal and replan action, so that I can recover through the existing grammar.
25. As a maintainer, I want one traversal implementation, so that an old loop cannot remain as a second owner.
26. As a maintainer, I want a comparison against captured pre-change behavior, so that a shared implementation cannot silently change a mode.
27. As a maintainer, I want fault tests at each real mode boundary, so that a generic table does not erase distinct authority.
28. As an owner of surviving work, I want all existing ownership and preservation refusals retained, so that refactoring grants no new deletion authority.

29. As a branch cleanup operator, I want landed and subsumed branches removed without recovery refs, so that their existing classification contract stays intact.

## Implementation decisions

The shared traversal owns iteration order, retained-member pass-through, completed-result accumulation, and stopping.
Its input is the adapter's ordered members.
An adapter supplies the operation for a removable member and the stopped-result projection.
A reached operation returns its current outcome and error.
The traversal does not derive an assignment identity, fingerprint, ownership proof, or branch classification.

Keep the four existing function entry points as adapters.
Keep explicit all-row preflight ahead of the recorded and unrecorded visits.
Keep landed all-row preflight ahead of its visits.
Keep unclaimed full fingerprint and length replan before visits.
Keep each tuple equality check at its current position immediately before that member's retained or deletion visit.
Do not add per-member lifecycle fault hooks to modes that lack them today.

The current preflight owners are [preflightExplicitSet and preflightLandedSet](../../internal/worktree/clean_set_apply.go:208).
The unclaimed replan and exact-delete route is [applyUnclaimedAssignmentSet](../../internal/worktree/clean_unclaimed.go:191).

Keep [applyCleanupTransaction](../../internal/worktree/resume.go:104) as the locked, receipted owner for recorded cleanup.
The outer traversal does not replace its locked planner check.
Only classUnique plans name a discarded recovery ref.
Landed and subsumed branch plans carry recovery none and create no discarded ref.
R22 and R29 preserve both sides of this condition.
The current rule is [planUnrecordedRow](../../internal/worktree/clean_discard.go:189), and [writeDiscardedRef](../../internal/worktree/clean_discard.go:274) skips an absent recovery ref.

Unrecorded branches retain [discardUnrecordedBranch](../../internal/worktree/clean_discard.go:249).
Unclaimed refs retain their direct exact-OID deletion operation.
These separate effect owners must not be forced through recorded-worktree authority.

Use the existing outcome helpers in [clean_set_apply.go](../../internal/worktree/clean_set_apply.go:42).
Its internal drift outcomes and command rendering are separate contracts.
[applyOutcomes](../../internal/worktree/clean_set_apply.go:192) replaces stale internal outcomes with the mode's refusal form.
Tests must cover both contracts.

[ADR0005](../../docs/adr/0005-worktree-cleanup-requires-verifiable-ownership.md:3) remains binding.
Automatic cleanup requires matching ownership and a dead or absent proven owner.
Unprovable ownership, unsafe ignored residue, and incomplete preservation evidence retain work.
Explicit cleanup still binds exact selection and ignored deletion through its plan fingerprint.
This spec changes none of those predicates.

## Implementation chunks

Each ticket is a complete green checkpoint on one retained integration source.
C1 creates the shared seam and immediately serves recorded explicit cleanup.
Its independent chunk review must accept the seam before C2 starts.
Each successor waits for its predecessor's green commit and chunk review.
Shared baseline and registry writes require this serial order.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
|---|---|---|---|---|
| C1 / 01-share-recorded-traversal.md | Recorded explicit cleanup uses the shared owner with unchanged authority and results | R01, R02, R03, R04, R05, R06, R08, R09, R10, R11, R12, R13, R14, R15, R19, R24, R27, R28 | Shared traversal table, explicit command and outcome witnesses, recorded baseline cases | yes |
| C2 / 02-migrate-unrecorded-traversal.md | Explicit branch cleanup uses the owner after recorded effects with conditional recovery | R03, R04, R05, R06, R09, R11, R12, R13, R14, R16, R19, R22, R24, R27, R28, R29 | Mixed explicit baseline, discarded-ref and holder-removal witnesses | yes |
| C3 / 03-migrate-landed-traversal.md | Landed command and folded-sibling cleanup use the owner with unchanged scope | R01, R02, R04, R05, R07, R08, R09, R10, R11, R12, R13, R14, R17, R19, R23, R24, R27, R28 | Landed command, locked drift, suffix, sibling, and baseline witnesses | no |
| C4 / 04-migrate-unclaimed-and-contract.md | Unclaimed cleanup uses the owner and the complete tree has no original traversal copies | R01, R02, R04, R05, R14, R18, R19, R20, R21, R24, R25, R26, R27, R28 | Post-replan deletion witnesses, exact-OID race, full baseline comparison, copy-survival census | yes |

The previous prospective C1 splits into C1 and C2.
The previous C2 becomes C3.
The previous C3 and C4 combine into C4.

Each migration removes its own old loop with its preservation evidence.
Do not retain unreachable legacy loops to create a separate contraction ticket.
C4 completes the sole-owner census and reconciles every baseline case.
No chunk may weaken an existing witness to obtain green.

## Completion plan

This version 1 plan records author verification and final integration checks.
The orchestrator amends it to version 2 before the first author dispatch.
The amendment records fresh authors and independent chunk review results.
No implementation evidence or ticket-review acceptance is claimed here.

```bench-completion-plan
{
  "version": 1,
  "chunks": [
    {
      "id": "C1",
      "tickets": [
        "01-share-recorded-traversal.md"
      ],
      "verification": [
        {
          "id": "mode-tests",
          "command": "bench test --package ./internal/worktree --run 'CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval'"
        },
        {
          "id": "command-contract",
          "command": "bench test --package ./cmd/bench --run 'CommandRegistry|HelpInventory|AXI'"
        },
        {
          "id": "axi-membership",
          "command": "bench test --check axi-query-registry"
        },
        {
          "id": "routing",
          "command": "bench test --check subcommand-routing"
        },
        {
          "id": "preservation",
          "command": "bench test --package ./internal/worktree --run 'CleanupTraversal|CleanSet|DiscardTargetRequalifiesAfterARecordedRemoval'",
          "probe": "Compare this mode with frozen baseline records and record each named behavioral red, restoration, and green."
        }
      ]
    },
    {
      "id": "C2",
      "tickets": [
        "02-migrate-unrecorded-traversal.md"
      ],
      "verification": [
        {
          "id": "mode-tests",
          "command": "bench test --package ./internal/worktree --run 'CleanupTraversal|CleanSet|DiscardTarget'"
        },
        {
          "id": "command-contract",
          "command": "bench test --package ./cmd/bench --run 'CommandRegistry|HelpInventory|AXI'"
        },
        {
          "id": "axi-membership",
          "command": "bench test --check axi-query-registry"
        },
        {
          "id": "routing",
          "command": "bench test --check subcommand-routing"
        },
        {
          "id": "preservation",
          "command": "bench test --package ./internal/worktree --run 'CleanupTraversal|CleanSet|DiscardTarget'",
          "probe": "Compare this mode with frozen baseline records and record each named behavioral red, restoration, and green."
        }
      ]
    },
    {
      "id": "C3",
      "tickets": [
        "03-migrate-landed-traversal.md"
      ],
      "verification": [
        {
          "id": "mode-tests",
          "command": "bench test --package ./internal/worktree --run 'CleanupTraversal|CleanSet|CleanupLanded|CleanLanded|LandCleansTheFoldedSibling|LandLeavesAPriorLandedAssignment|LandRetainsAnUnprovenSibling'"
        },
        {
          "id": "command-contract",
          "command": "bench test --package ./cmd/bench --run 'CommandRegistry|HelpInventory|AXI'"
        },
        {
          "id": "axi-membership",
          "command": "bench test --check axi-query-registry"
        },
        {
          "id": "routing",
          "command": "bench test --check subcommand-routing"
        },
        {
          "id": "preservation",
          "command": "bench test --package ./internal/worktree --run 'CleanupTraversal|CleanSet|CleanupLanded|CleanLanded|LandCleansTheFoldedSibling|LandLeavesAPriorLandedAssignment|LandRetainsAnUnprovenSibling'",
          "probe": "Compare this mode with frozen baseline records and record each named behavioral red, restoration, and green."
        }
      ]
    },
    {
      "id": "C4",
      "tickets": [
        "04-migrate-unclaimed-and-contract.md"
      ],
      "verification": [
        {
          "id": "mode-tests",
          "command": "bench test --package ./internal/worktree --run 'CleanupTraversal|CleanupUnclaimed|CleanUnclaimed|CleanSet|DiscardTarget|CleanLanded|LandCleansTheFoldedSibling|LandLeavesAPriorLandedAssignment|LandRetainsAnUnprovenSibling'"
        },
        {
          "id": "command-contract",
          "command": "bench test --package ./cmd/bench --run 'CommandRegistry|HelpInventory|AXI'"
        },
        {
          "id": "axi-membership",
          "command": "bench test --check axi-query-registry"
        },
        {
          "id": "routing",
          "command": "bench test --check subcommand-routing"
        },
        {
          "id": "preservation",
          "command": "bench test --package ./internal/worktree --run 'CleanupTraversal|CleanupUnclaimed|CleanUnclaimed|CleanSet|DiscardTarget|CleanLanded|LandCleansTheFoldedSibling|LandLeavesAPriorLandedAssignment|LandRetainsAnUnprovenSibling'",
          "probe": "Compare this mode with frozen baseline records and record each named behavioral red, restoration, and green."
        }
      ]
    }
  ],
  "final_verification": [
    {
      "id": "all-mode-tests",
      "command": "bench test --package ./internal/worktree"
    },
    {
      "id": "coverage",
      "command": "bench coverage --check cleanup-member-traversal"
    },
    {
      "id": "command-contract",
      "command": "bench test --package ./cmd/bench"
    },
    {
      "id": "axi-membership",
      "command": "bench test --check axi-query-registry"
    },
    {
      "id": "routing",
      "command": "bench test --check subcommand-routing"
    }
  ]
}
```

## Testing decisions

Use the existing cleanup command seam and the internal adapters where command rendering hides drift outcomes.
A unit table exercises the shared traversal's stop rule with injected visits.
It does not replace mode-specific command, authority, or effect tests.

The existing fixtures in clean_set_apply_test.go remain the shared recorded-set fixture owners.
Use the existing unclaimed and discarded-ref fixture helpers.
Do not paste their repository construction into a new table.

Before production edits, capture pre-change responses and durable-state outcomes once.
Store the baseline under internal/worktree/testdata/cleanup-member-traversal/.
A permanent comparison test drives current adapters against those records.
Ordinary tests must never restore old production source or rewrite the baseline.
Normalize only fixture-dependent paths and identities through named mappings.
Keep actions, row order, digests, errors, exit status, recovery refs, and deletion outcomes meaningful.

The baseline case manifest enumerates each applicable mode and each case below.
Do not take an unchecked Cartesian product.
Recorded and landed cases include outer preflight and locked drift.
Unrecorded cases include holder removal and recovery-ref deletion.
Unclaimed cases include full selection drift and exact-OID races.
Common cases include empty, retained-only, mixed retained/removable, all-success, first failure, and later failure.

Each expected result comes from a captured original run or an independently authored expectation with a demonstrated omission red.

The copy-survival check reads the four adapter bodies and proves that each delegates visits to the shared owner.
It must fail if any original member-iteration loop remains.
Demonstrate that red by restoring one old loop at each adapter site during required implementation verification.
A harmless identifier check alone is insufficient.

### Seam diagram

    trigger: cleanup apply command or folded-sibling landing cleanup
        |
        v
    selected members -> [mode preflight and authority adapter]
                                 |
                                 v
                        [shared ordered traversal]
                                 |
                     mode qualification and effect owner
                                 |
                                 v
                     [mode outcome adapter and renderer]
                                 |
                                 v
                       result rows and durable state

Tests drive the public cleanup command, internal adapter results, and the landing sibling caller.
Faults enter through existing injected boundaries.
Tests observe surviving trees and refs as well as reported rows.

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| R01 | 1 | Empty applicable selections produce no effects | worktree cleanup command, planned TestCleanupTraversalEmptyModes in internal/worktree/cleanup_members_baseline_test.go | A spurious visit changes the effect trace |
| R02 | 2 | Visits retain the approved member order | shared traversal table, planned TestCleanupTraversalOrder in internal/worktree/cleanup_members_test.go | Reordering changes the expected trace |
| R03 | 3 | Recorded visits precede unrecorded visits | `internal/worktree/clean_discard_transaction_test.go` (`TestDiscardTargetRequalifiesAfterARecordedRemoval`) | Reversing the stages removes a branch that must survive |
| R04 | 4 | A retained row passes through unchanged | `internal/worktree/clean_set_outcomes_test.go` (`TestCleanSetRetainedMember`) and mode baseline cases | Rewriting its digest or verdict changes the stored row |
| R05 | 5 | Retained rows invoke no qualification or effect | planned TestCleanupTraversalRetainedVisits in internal/worktree/cleanup_members_test.go | A redundant call changes the expected call trace |
| R06 | 6 | Explicit preflight checks every removable member before removal | `internal/worktree/clean_set_apply_test.go` (`TestCleanSetPreflightAllRows`) | Omitting later preflight permits an earlier deletion |
| R07 | 7 | Landed preflight checks every removable member before removal | `internal/worktree/clean_set_apply_test.go` (`TestCleanSetPreflightAllRows`) | A late invalid assignment leaves an earlier effect |
| R08 | 8 | A preflight fault names its actual offender after a retained leading row | `internal/worktree/clean_set_refusal_test.go` (`TestCleanSetPreflightFaultOnLaterMember`) and `internal/worktree/clean_set_refusal_test.go` (`TestCleanSetExplicitPreflightFaultOnLaterMember`) | An offset bug attributes the error to the retained member |
| R09 | 9 | Requalification observes effects from earlier members | `internal/worktree/clean_set_apply_test.go` (`TestCleanSetLateDrift`) and `internal/worktree/clean_discard_transaction_test.go` (`TestDiscardTargetRequalifiesAfterARecordedRemoval`) | Omitting the later check removes newly unsafe work |
| R10 | 10 | The recorded transaction rechecks qualification while locked | `internal/worktree/clean_set_apply_test.go` (`TestCleanSetLateDrift`) | Omitting the locked check removes a changed later member |
| R11 | 11 | A failed visit produces the failed member's fault outcome | `internal/worktree/clean_set_outcomes_test.go` (`TestCleanSetPartialApply`) and traversal fault table | A generic stop loses target identity or fault reason |
| R12 | 12 | Drift to retention preserves the current retained outcome internally | `internal/worktree/clean_set_refusal_test.go` (`TestCleanSetMemberDriftAfterPreflight`) | Turning drift into a fault changes the internal row |
| R13 | 13 | Drift while still removable produces not-attempted internally | `internal/worktree/clean_set_refusal_test.go` (`TestCleanSetMemberDriftAfterPreflight`) | Reporting removed or error changes the internal action |
| R14 | 14 | No visit runs after the first failed visit | planned TestCleanupTraversalStopsAtFailedVisit in internal/worktree/cleanup_members_test.go | A continued loop changes later durable state |
| R15 | 15 | Explicit recorded members after a failure receive unreached outcomes | `internal/worktree/clean_set_outcomes_test.go` (`TestCleanSetUnstartedOutcomes`) | Dropping the suffix loses selected member rows |
| R16 | 16 | Explicit unrecorded members after a failure receive unreached outcomes | planned TestCleanupTraversalUnrecordedSuffix in internal/worktree/cleanup_members_baseline_test.go | A mixed selection loses its branch suffix |
| R17 | 17 | Landed members after a failure receive unreached outcomes | planned TestCleanupLandedUnreachedSuffix in internal/worktree/cleanup_members_baseline_test.go | Prefix-only projection loses the required suffix |
| R18 | 18 | Unclaimed deletion failure returns only its completed prefix and failed row | planned TestCleanupUnclaimedDeletionFailureAfterReplan in internal/worktree/cleanup_members_baseline_test.go | Appending unreached rows changes this mode's contract |
| R19 | 19 | A completed earlier effect remains completed after a later failure | planned TestCleanupTraversalCompletedPrefix in internal/worktree/cleanup_members_baseline_test.go | A rollback claim or row rewrite contradicts durable state |
| R20 | 20 | Unclaimed selection drift refuses before deletion | planned TestCleanupUnclaimedSelectionDrift in internal/worktree/cleanup_members_baseline_test.go | Omitting full replan permits a changed selection |
| R21 | 21 | Unclaimed deletion checks the planned exact OID | planned TestCleanupUnclaimedMovedRefAfterReplan in internal/worktree/cleanup_members_baseline_test.go | Deleting a moved ref changes the survivor check |
| R22 | 22 | Unique unrecorded branch deletion preserves its planned tip under the planned recovery ref | `internal/worktree/clean_discard_test.go` (`TestDiscardTargetWritesTheDiscardedRefFirst`) | Omitting unique-branch recovery loses the preserved tip |
| R23 | 23 | Folded-sibling cleanup retains the narrowed landed scope | `internal/worktree/land_effects_cleanup_test.go` (`TestLandCleansTheFoldedSibling`), `internal/worktree/land_effects_cleanup_test.go` (`TestLandLeavesAPriorLandedAssignment`), `internal/worktree/land_effects_cleanup_test.go` (`TestLandRetainsAnUnprovenSibling`) | Broadening scope removes an unrelated or unproven assignment |
| R24 | 24 | Stale command rendering preserves its refusal and exact replan action | `internal/worktree/clean_set_outcomes_test.go` (`TestCleanSetStaleReplanAction`), `internal/worktree/clean_set_outcomes_test.go` (`TestCleanSetUnclaimedStaleReplanAction`), `internal/worktree/clean_set_outcomes_test.go` (`TestCleanSetApplyTimeStaleRefusal`) | Rendering internal rows directly changes the public refusal |
| R25 | 25 | Every adapter uses the sole member traversal owner | planned TestCleanupTraversalCopySurvival in internal/worktree/cleanup_members_census_test.go | Restoring any original iteration loop turns the census red |
| R26 | 26 | Each enumerated baseline case preserves response and durable-state outcome | planned TestCleanupTraversalBaseline in internal/worktree/cleanup_members_baseline_test.go | A per-mode behavior delta changes the captured result |
| R27 | 27 | Each mode retains its existing injected fault boundaries | planned TestCleanupTraversalModeFaultBoundaries in internal/worktree/cleanup_members_baseline_test.go | Adding or omitting a boundary changes the mode trace |
| R28 | 28 | Existing authority and preservation refusals remain effective | planned TestCleanupTraversalAuthorityPreservation in internal/worktree/cleanup_members_baseline_test.go | A bypass removes work that the refusal fixture retains |
| R29 | 29 | Landed and subsumed unrecorded deletion reports recovery none and writes no discarded ref | `internal/worktree/clean_discard_test.go` (`TestDiscardTargetLandedAndSubsumedWriteNoRef`) | Unconditional recovery creation changes the ref namespace |

Planned tests are future evidence.
No new red, green, or mutation result is claimed at spec time.
Implementation verification must demonstrate each independent expectation's omission red.
Use distinct mutations for skipped outer preflight, skipped later qualification, skipped locked qualification, changed stop projection, continued visits, and restored loop copies.

### Edge inventory

Cover empty and retained-only selections, mixed recorded and unrecorded members, first and later failure, and a retained leading row.
Cover stale selection before entry, before the first effect, after an earlier effect, and inside the recorded lock.
Cover drift to retained and drift that remains removable separately.
Cover unclaimed namespace additions, removals, tuple changes, and a ref moved before exact deletion.
Cover original fault hooks and zero hooks where the current adapter has none.
Cover partial completion without rollback and replay without repeated deletion.

Won't handle: a new uniform failure-row format — all existing cleanup callers keep their mode's current format.
Won't handle: new cleanup authority — ADR0005 and existing eligibility owners remain authoritative.
Won't handle: a new cleanup transaction — recorded cleanup keeps applyCleanupTransaction.
Won't handle: lifecycle planning outside traversal — FT217 remains a separate outcome.
Won't handle: new command grammar or response fields — current clean command callers remain supported.
Won't handle: performance claims — existing recorded timing cannot establish this refactor's savings.

### Unclaimed deletion witnesses after successful replan

R18 uses planned TestCleanupUnclaimedDeletionFailureAfterReplan in internal/worktree/cleanup_members_baseline_test.go.
Run first-failure and later-failure cases through applyUnclaimedAssignmentSet.
Build a stable selection with multiple removable refs and retained rows.
Place a Git ref lock on the selected first or later removable ref.
The lock must leave selection planning unchanged.

Assert that full fingerprint and length replan pass.
Assert that the returned error is a deletion fault rather than errStaleFingerprint.
Assert the completed prefix, exactly one failed row, and no unreached suffix rows.


Assert earlier successful deletions remain complete.
Assert the failed ref and every later ref retain their original OIDs.
Also assert retained prefix rows remain unchanged.
A first-failure fixture proves that no prior deletion occurs.
The later-failure fixture proves that the completed prefix is not rolled back.

R21 uses planned TestCleanupUnclaimedMovedRefAfterReplan in internal/worktree/cleanup_members_baseline_test.go.
Install a real Git reference-transaction hook in the disposable fixture repository.
On committed deletion of the first selected removable ref, move a later selected ref to a different OID.
This event occurs after the complete replan and after the first completed deletion.
The hook must ignore its own nested update and record that it ran.
The adapter then attempts deletion of the later ref with its planned old OID.


Assert a deletion fault, the completed earlier row, and the failed later row.
Assert no suffix rows appear.
Assert the moved ref retains the new OID and each unvisited ref retains its original OID.
Assert the first deleted ref remains absent.
Omitting the old-OID guard must turn this witness red.

No new production lifecycle hook is permitted.
StepUnlockedReplan runs before replan and cannot prove either post-replan property.
Existing stale-plan tests remain witnesses for R20.
They are not substitutes for R18 or R21.

A bounded disposable Git correctness probe ran during spec repair.
Ref locks caused first and later exact deletion failures.
A committed reference-transaction hook moved a later ref after an earlier deletion.
Exact deletion with the old OID then failed.
The failed and later refs survived.
The disposable repositories were removed.


This verifies the Git fixture primitives on this host.
Full Bench replan, output rows, and authority composition remain planned test evidence.

### Baseline identity and digest mapping

Each baseline case declares its original identity and its current fixture identity.
Map repository roots, assignment paths, generated branch names, session IDs, and timestamps through that explicit bijection.
Compare identity relationships before rendering normalized responses.
Normalize a fingerprint only after comparing its complete canonical inputs in that case.
Keep options, selection order, ref OIDs, classification, recovery intent, and relevant file images in those inputs.


Verify that each fingerprint binds the same inputs within its own run.
Preserve equality and inequality relationships between plan and result digests.
A changed digest relation, input omission, ref tip, or recovery destination must fail the comparison.
Do not replace all digests or OIDs with one placeholder.


Case-specific OID mappings may rename independently created commits only after their parent and content relationships match.
Retain lock receipts and completed-effect outcomes as state evidence.
The baseline comparison must also fail when an earlier effect changes a later member's classification or proof.

## Ownership fences

Independent ticket review accepted these implementation planning fences.
Production implementation remains outside this planning phase.
Each planned test path occurs here.
Ticket slicing changes only this spec and its ticket files.

- `internal/worktree/cleanup_members.go` (new)
- `internal/worktree/cleanup_members_test.go` (new)
- `internal/worktree/cleanup_members_census_test.go` (new)
- `internal/worktree/cleanup_members_baseline_test.go` (new)
- `internal/worktree/testdata/cleanup-member-traversal/` (new)
- `internal/worktree/clean_set.go`
- `internal/worktree/clean_landed.go`
- `internal/worktree/clean_discard.go`
- `internal/worktree/clean_unclaimed.go`
- `internal/worktree/clean_set_apply.go`
- `internal/worktree/clean_set_apply_test.go`
- `internal/worktree/clean_set_refusal_test.go`
- `internal/worktree/clean_set_outcomes_test.go`
- `internal/worktree/clean_discard_transaction_test.go`
- `internal/worktree/clean_discard_test.go`
- `internal/worktree/clean_set_wiring_test.go`
- `internal/worktree/clean_landed_apply_test.go`
- `internal/worktree/land_effects_cleanup_test.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `specs/cleanup-member-traversal/`
- `reviews/cleanup-member-traversal.md` (new)

No public symbol is moved or retired.
The unchanged non-command caller is internal/worktree/land_effects.go:134.
Do not replace its landed adapter or widen its scope.
The binding registry requires every ticket to co-name the five command registry paths.
They preserve current dispatch, help, AXI membership, and routing.
Only spec and ticket writes are authorized during this planning phase.

## Out of scope

Uniform output, adoption lifecycle deepening, changed cleanup authority, and performance tuning require separate outcomes.
This delivery spends zero edits and zero gate runs on those cuts.
A future estimate must use its own concrete scope.
The eventual build runs its affected checks and required gate.
This planning phase authorizes ticket files but no implementation.

## Further notes

### Current source and caller verification

Source HEAD: 37ac80b8adeff4ecbd57f3d813ff8937f2bcbc68.
The author inspected the four loops, outcome helpers, ADR0005, and cited witnesses.
Current-code assessments remain claimed, confidence 9.
The disposable Git fixture primitives were verified during repair.
No Bench test execution verified the planned composition.

The changed-function caller census used bench consumers with worktree-qualified symbols.
Its production routes are cleanExplicitSet to applyExplicitSet, applyExplicitSet to applyUnrecordedRows,
cleanLanded and cleanLandedSiblings to applyLandedSet, and cleanUnclaimed to applyUnclaimedAssignmentSet.
The non-command caller is [cleanLandedSiblings](../../internal/worktree/land_effects.go:134).
Existing internal adapter tests remain in the fence.
The author must refresh this census before slice lock if the source advances.

### Consolidation table

| Existing site | Shared rule after migration | Mode rule retained | Rows |
|---|---|---|---|
| applyExplicitSet | ordered visit, retained pass-through, stop | explicit preflight and recorded locked effect | R03–R15 |
| applyUnrecordedRows | ordered visit, retained pass-through, stop | branch requalification and recovery-ref exact deletion | R03, R09, R16, R22 |
| applyLandedSet | ordered visit, retained pass-through, stop | landed preflight, scope, locked planner, unreached suffix | R07–R14, R17, R23 |
| applyUnclaimedAssignmentSet | ordered visit, retained pass-through, stop | complete replan, exact OID, prefix-only failure | R18, R20, R21 |

### Source-sentence-to-row table

| Source clause | Coverage |
|---|---|
| Decision ticket 1: preserve distinct failure rows | R15, R16, R17, R18 |
| Decision ticket 2: ordered visits and retained members | R01, R02, R03, R04, R05, R14 |
| Decision ticket 2: preflight and later requalification | R06, R07, R08, R09, R10 |
| Decision ticket 2: unclaimed replan and exact object check | R20, R21 |
| Decision ticket 2: existing effect owners | R10, R22, R28 |
| Decision ticket 3: command tests and multi-member failures | R11, R14, R15, R16, R17, R18, R19, R24, R27 |
| Decision ticket 3: preflight, drift, partial and unstarted witnesses | R06, R07, R08, R09, R12, R13, R19 |
| FT363: sole traversal owner and shared fault mechanism | R25, R26, R27 |
| ADR0005: ownership and preservation requirements | R28 |
| Unrecorded unique versus landed/subsumed recovery contract | R22, R29 |

### Pre-review proof checklist

- Cited symbols: current adapter and helper declarations were opened.
- Import edges: none added or changed.
- Source-row clauses and occurrences: enumerated in the source table.
- Promised field labels: none new.
- Changed-function callers: production routes and internal tests enumerated above.
- Copy survival: R25 restores each original loop in its omission proof.
- Rendered-message citations: existing public forms are preserved by R24 and R26.
- Pin operators: no new pinned production fact.
- Entry reads: existing adapters keep ambient inputs and their existing injected seams.
- Derived expectations: R26 uses captured baseline records, not current implementation answers.
- Consolidated rules: the table assigns each old traversal site.
- Quantified obligations: every mode gets its applicable baseline and fault cases.
- Workflow-step writes: none added.
- Flagged additions: new private traversal, copy-survival test, and baseline comparison.
- Unread structured sources: none.

### Review checkpoint

| Item | Proposed disposition |
|---|---|
| Implementation line | gpt-6.1-sol / high |
| Seams | Existing command, internal outcome, and landing caller seams |
| Coverage and edges | Preserve distinct modes and all named failure windows |
| Fences | Exact paths and ticket allocations accepted by independent review |
| Scope cuts | No output unification, authority change, adoption planner, or timing claim |
| Next phase | A separately authorized implementation, after its admission checks |

### Ticket review result

Independent Sol 6.1/high review FT363-TICKETS-R1 accepted the graph at `a382f4848d432815968f044820fad593e856e297`.
It found no material Standards, Spec, or Coverage findings.
The review verified complete caller migrations, all twenty-nine acceptance rows, the serial dependencies, and the exact ownership union.

Clean-source preflight passed with fourteen green checks and two not-applicable checks.
All four ticket proposals had no missing paths or ordering requirements.
These planning checks used the reviewed checkpoint as the prospective implementation baseline.
The original main-to-batch delta also contains unrelated maps; it is not a cleanup implementation delta.
No runtime preservation or mutation result is claimed by this review.
