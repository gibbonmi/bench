# Thread the verb's run through the explicit cleanup plan and the release paths

Blocked by: 03-thread-the-landing.md
Writes: internal/conformance/bounds_policy_test.go, internal/conformance/process_lifetime_test.go, internal/worktree/admin_readers_test.go, internal/worktree/classifier.go, internal/worktree/clean_set.go, internal/worktree/clean_set_apply.go, internal/worktree/eligibility_test.go, internal/worktree/exec.go, internal/worktree/identity_component_test.go, internal/worktree/journey_test.go, internal/worktree/land_resume.go, internal/worktree/land_resume_test.go, internal/worktree/land_run_test.go (new), internal/worktree/land_stack_resume_test.go, internal/worktree/lifecycle_acquire_test.go, internal/worktree/lifecycle_facts_test.go, internal/worktree/lifecycle_policy_test.go, internal/worktree/lifecycle_test.go, internal/worktree/orphan_sweep_test.go, internal/worktree/orphan_test.go, internal/worktree/owner_marker.go, internal/worktree/ownership.go, internal/worktree/ownership_test.go, internal/worktree/reauthorize_test.go, internal/worktree/recovery_retry_test.go, internal/worktree/release_inside_test.go, internal/worktree/resume.go, internal/worktree/resume_test.go, internal/worktree/subshell.go, internal/worktree/unlanded_route_test.go, internal/worktree/worktree.go, internal/worktree/worktree_test.go, internal/worktree/identity_component.go, internal/worktree/reauthorize.go
Covers: GQ35, GQ36, GQ37

## What to build

`planExplicitWith` (`subshell.go:211`) reads each location and each ledger answer through the run that it holds.
Its ledger reads become `intent.In(run, root).LifecycleEvidence()` (`:259`) and `intent.In(run, root).Assignments()` (`:280`). The completion-plan probe swaps the first line.
Its admin-directory read (`:251`) and its index-path read (`:332`) go through `run.AdminDir` and `run.AdminPath`.
`foreignRecoveryAssignment` (`classifier.go:317`) takes the run too. Its callers are `subshell.go:304` and `eligibility_test.go:165`.
`PlanExplicit`, `PlanExplicitWithOptions`, and `PlanAutomatic` are exported entries. They keep or gain one construction each (`subshell.go:192`, `:206`, `classifier.go:261`).

Remove three helper constructions. The root-only wrapper `validateOwnerMarker` (`owner_marker.go:22-24`) goes, and its test callers call `validateOwnerMarkerIn` with a run or nil.
Those test callers are `admin_readers_test.go:102`, `eligibility_test.go:137`, and `journey_test.go:254`. The journey test passes a fresh run, because it counts one start.
`planExplicitCleanupRow` (`clean_set.go:191-192`) and `planAutomaticAt` (`classifier.go:266-267`) take the run and build none.
The callers of `planAutomaticAt` are `classifier.go:261`, `ownership.go:320`, and `resume.go:374`, `:378`, and `:417`. Its test caller is `resume_test.go:504`.

The release and cleanup-transaction paths pass the verb's run to `assignmentByID` (`resume.go:59`), `cleanupIdentity` (`resume.go:22`), and `markerPath` (`ownership.go:158`).
Their production callers are in `worktree.go`, `ownership.go`, `resume.go`, `exec.go:95`, and `land_resume.go:123`.
Their test callers pass a run or nil. The test callers of `assignmentByID` are `release_inside_test.go`, `unlanded_route_test.go`, `lifecycle_facts_test.go`, and `recovery_retry_test.go`.
The other `assignmentByID` test callers are `orphan_sweep_test.go`, `lifecycle_test.go`, `lifecycle_policy_test.go`, `orphan_test.go`, `reauthorize_test.go`, `resume_test.go`, and `worktree_test.go`.
The test callers of `cleanupIdentity` are `lifecycle_test.go`, `resume_test.go`, `recovery_retry_test.go`, `worktree_test.go`, `land_stack_resume_test.go`, and `land_resume_test.go`.
The test callers of `markerPath` are `ownership_test.go`, `admin_readers_test.go`, `identity_component_test.go`, and `lifecycle_acquire_test.go`.
`LeaseFile` (`worktree.go:31`) is exported, so it keeps its signature and delegates to a run form. Its 40 test rows do not change.

Each helper takes the run as its first parameter after any `joins` value. No `XIn` wrapper is left for a helper that this ticket changes.
Keep `bounds.AssignmentStale` in `classifier.go` and `subshellAt` in `subshell.go`. The pins at `internal/conformance/bounds_policy_test.go:50` and `internal/conformance/process_lifetime_test.go:51` read them.
`subshell.go` has 14 lines of headroom, `classifier.go` 13, and `clean_set.go` 9. A change that grows one past its budget moves a function to a new file, after a `Writes:` widen.
`ownership_test.go`, `resume_test.go`, and `lifecycle_test.go` are over their budgets. Change them in place and add no line.

In `TestExplicitPlanReadsEachLocationOnce` (`journey_test.go:271`), set the common-directory pin at `:278` from 3 to 1.
Extend the same test with an admin-directory count for the target, pinned at 1 with `==`. This adds no top-level test.
Correct the comments at `journey_test.go:245-249` and `:266-270`, so that they state the run's current cache scope. Keep the two-run check at `:260-263`.
When this ticket lowers the landing count, set the `==` pin of `TestLandingReadsTheDestinationCommonDirOnce` to the new observed count. Ticket 03 owns that row.

## Acceptance

- [ ] One explicit plan for an owned target starts the common-directory query once (GQ35).
- [ ] One explicit plan for an owned target starts the admin-directory query once for the target (GQ36).
- [ ] Two verb runs hold two distinct run values (GQ37).
- [ ] An explicit plan whose lifecycle-evidence read drops the run starts the common-directory query twice, so the plan test fails.
- [ ] No release or cleanup-transaction helper that this ticket changes builds its own run.
- [ ] `TestLandingReadsTheDestinationCommonDirOnce` passes at its pinned count.

## Checkpoint verification

- `bench test --package ./internal/worktree`
- `bench test --package ./internal/conformance`

Run `bench consumers` on each changed unexported function before the edit. Record the red of the probe mutation, then the green.
