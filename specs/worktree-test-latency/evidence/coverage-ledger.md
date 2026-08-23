# Coverage ledger — ticket 06 contraction

This ledger maps every repository-backed test function that ticket 06 removes
to one retained pure, adapter, or public-journey test. It is written before
the deletions. The verification log at the bottom records the base-to-tip
deleted-test comparison and the proof-inventory red/green demonstrations.
(Coverage rows CV1, RJ1, FA1.)

Base: e00b74e1. Package under contraction: `internal/worktree`.

## Removed test functions and their surviving disposition

| removed test (file) | disposition | surviving test |
|---|---|---|
| TestExplicitEligibilityOutcomeMatrix (eligibility_test.go) | pure table | lifecyclepolicy TestExplicitDecisionTable, with the plan projection held by TestEligibilityVerdictProjectsWithoutSecondDecision; its combined-fact cases "ignored/declaration-overrides-marker-malformed", "lease/live-overrides-lock-mismatch", and "ignored/undeclared-overrides-live-lease" are the precedence survivors of the EX3, EX5, and EX6 subtests |
| TestAutomaticEligibilityOutcomeMatrix (eligibility_test.go) | pure table | lifecyclepolicy TestAutomaticDecisionTable, with fact translation held by the lifecycle fact adapters; its combined-fact cases "lease/live-overrides-everything" and "eligibility/retain-landed-swap" are the precedence survivors |
| TestReleaseMalformedLeaseRetainsAsUncertain (lifecycle_policy_test.go) | adapter | TestLifecycleLeaseFactAdapterTranslatesRealLeases (malformed lease reads LeaseUnknown), verdict row in lifecyclepolicy TestExplicitDecisionTable |
| TestReleaseNumericLeaseWithBadTimestampAndExtraFieldsRetainsAsUncertain (lifecycle_policy_test.go) | pure table | lifecyclepolicy TestLeaseOwnerPIDTable content partitions |
| TestReleasePartialNumericLeaseRetainsAsUncertain (lifecycle_policy_test.go) | pure table | lifecyclepolicy TestLeaseOwnerPIDTable content partitions |
| TestPoolKeyPredicateRetainsAKeyWhoseExistenceItCannotProve (pool_reclaim_test.go) | adapter | TestReclaimFactAdapterTranslatesDeniedReads, uncertainty rows in reclaimpolicy TestClassifyKeyTable |
| TestPoolKeyPredicateRetainsSymlinksUnfollowed (pool_reclaim_test.go) | adapter | TestReclaimFactAdapterTranslatesHostileShapesUnopened (symlink shape undescended), symlink rows in reclaimpolicy TestClassifyKeyTable |
| TestPoolKeyRetainsAChildWhoseGitdirTargetIsRelative (pool_reclaim_test.go) | pure table | reclaimpolicy TestClassifyKeyTable relative-pointer row and TestGitdirTargetTable |
| TestPoolKeyRetainsAChildWhoseTargetEndsInASpace (pool_reclaim_test.go) | pure table | reclaimpolicy TestGitdirTargetTable "trailing space kept" row |
| TestPoolKeyRetainsAChildWhoseGitEntryIsAFifo (pool_reclaim_test.go) | adapter | TestReclaimFactAdapterTranslatesHostileShapesUnopened (FIFO pointer unopened), FIFO row in reclaimpolicy TestClassifyKeyTable |

Retained near the removals: the symlink-lease and directory-lease release
tests stay, because the lease fact adapter does not yet cover a non-regular
lease file; their observation has no pure or adapter substitute.

## Removed compatibility helper

`Claim` in `lifecycle.go` was the temporary ambient compatibility form of
`claimAt`. Its only callers were tests; they now call `claimAt` with an
explicit instant, and the form is removed. The other compatibility forms
(`Create`, `Acquire`, `Pool`, `PlanAutomatic`, `ConservativeCleanup`) keep
callers outside `internal/worktree/` or in production commands, so they stay.

## Verification log

Recorded 2026-08-23 after the diff was complete, on the ticket 06 worktree
(base e00b74e1).

**CV1 — base-to-tip deleted-test comparison.** Command:
`git diff e00b74e1 -- 'internal/worktree/*_test.go' | rg '^-func (Test[A-Za-z0-9_]+)'`.
Result: 10 deleted test functions. Every one of the 10 appears in the table
above, and every table row appears in the diff. No unlisted deletion.

**RJ1 — journey proof red/green.** The `markProof(t, "reclaim/journey/deletion")`
line in `TestReclaimApplyRemovesExactlyThePlannedKeys` was removed; the full
package run then printed
`missing required proof reclaim/journey/deletion owner=reclaim class=journey`
and exited non-zero with every test green. Restoring the line returned the
package to green.

**FA1 — adapter proof red/green.** The `markProof(t, "landing/adapter/facts")`
line in `TestLandingResidueFactAdapterTranslatesRealDestination` was removed;
the full package run then printed
`missing required proof landing/adapter/facts owner=landing class=adapter`
and exited non-zero. Restoring the line returned the package to green.

**EI1 — out-of-harness effect red/green.** A probe test with `t.Setenv` and
`os.Chdir` was appended to `snapshot_test.go`; `TestSerialJourneyHarnessCensus`
failed naming both sites and the required harness route. Removing the probe
returned the census to green.

**PR1 — explicit precedence repair red/green.** Recorded 2026-08-23. The
repair added the three combined-fact cases above to
`TestExplicitDecisionTable`. A temporary probe then short-circuited the
ignored/build-output block in `DecideExplicit` to skip a verdict that already
retains. The probe turned "ignored/declaration-overrides-marker-malformed"
and "ignored/undeclared-overrides-live-lease" red with the earlier block's
reason. The restore returned the package to green.

**NP1 — concurrency census.** A source sweep over the changed
`internal/worktree/` tree finds no `t.Parallel` call and no scheduler import
outside the census tests that name the string.
