# Thread the verb's run through one landing and its landed-set cleanup

Blocked by: 02-thread-the-run-through-the-ledger-and-store.md
Writes: internal/gittest/gittest.go, internal/worktree/clean_landed.go, internal/worktree/clean_landed_plan.go (new), internal/worktree/clean_landed_apply_test.go, internal/worktree/clean_set_apply.go, internal/worktree/clean_set_apply_test.go, internal/worktree/clean_set_refusal_test.go, internal/worktree/identity_component.go, internal/worktree/joins.go, internal/worktree/land.go, internal/worktree/land_effects.go, internal/worktree/land_empty_sibling_test.go, internal/worktree/land_identity.go, internal/worktree/land_journey_test.go, internal/worktree/land_lease.go, internal/worktree/land_resume.go, internal/worktree/land_run_test.go (new), internal/worktree/land_stack.go, internal/worktree/land_stack_preflight.go, internal/worktree/land_stack_resume.go, internal/worktree/list.go, internal/worktree/merge.go, internal/worktree/ownership.go, internal/worktree/parallel_census_test.go, internal/worktree/path.go, internal/worktree/resume.go, internal/worktree/worktree.go, internal/worktree/journey_test.go, internal/worktree/commitment_landing_test.go, internal/worktree/land_effects_test.go
Covers: GQ34, GQ41

## What to build

One landing passes the verb's run, `a.run` from `newAmbient`, to each store, ledger read, and helper that it reaches. Ticket 02 supplies `intent.In` and `Store.Run`.
Build each landing store literal with `Run: a.run`. The window store at `land.go:198` reads `(commitrepo.Store{Root: root, Run: a.run}).AdmitWindow(publication)`.
The completion-plan probe of this chunk swaps that exact text. The stores at `land.go:206`, `land_lease.go:14`, and `land_stack.go:97` and `:106` take the run too.
The stores at `land_stack_preflight.go:26` and `:114`, `joins.go:77`, and `land_effects.go:180` take the run when the landing reaches them.

Give `assignmentForRequest` and `unmatchedRequestRecovery` (`worktree.go:516`, `:531`) the run. Their callers are `land_identity.go:148`, `ownership.go:258`, and `land_resume.go:191`.
Give `identityBundleRefusal` (`identity_component.go:95`) the run, and remove its own construction at `:102`. Its callers are `land_identity.go:156`, `land_resume.go:204`, and `merge.go:180` and `:302`.
The resume path passes `a.run`, not nil and not a new run. A resume and a first landing then share one cache.

Each ancestry check of `landingSourceProofs` (`land.go:273`, `:278`) calls `run.IsAncestor` and reads an error as no, as `git.OK` does today.
Add no yes-or-no form to `git.Run`.

The landing runs the landed-set cleanup on each landing (`land_effects.go:71`), so this ticket threads that group too.
Give the run to `planLandedSet`, `selectLandedCleanupRow`, `planLandedAssignment`, `requalifyLandedRow`, `replanLandedCleanupRow`, and `assignmentLeaseStates`.
Remove the construction at `clean_landed.go:145`. The callers outside `clean_landed.go` are `land_effects.go:214`, `worktree.go:326`, `ownership.go:343`, `resume.go:468`, `list.go:118`, `path.go:118`, and `clean_set_apply.go:233`.
Each test caller passes a run or nil: `clean_landed_apply_test.go:157`, `clean_set_apply_test.go:153`, `clean_set_refusal_test.go:27`, `:93`, `:118`, `:147`, `:156`, and `land_empty_sibling_test.go:54`, `:72`, `:85`, `:88`.

Each helper takes the run as its first parameter after any `joins` value. Each collapsed pair keeps one form, and no `XIn` wrapper is left.
`land.go` has 6 lines of headroom, `land_stack.go` 7, and `clean_landed.go` none. A change that grows one of them past 400 lines moves a function to a new file in this ticket.
For `clean_landed.go`, that destination is `clean_landed_plan.go`. For `land.go` and `land_stack.go`, the orchestrator widens `Writes:` for a new file first.
`land_journey_test.go` and `parallel_census_test.go` are over their budgets. Change them in place and add no line.

Add the serial test `TestLandingReadsTheDestinationCommonDirOnce` in the new file `land_run_test.go`. It drives the in-process landing of `land_journey_test.go` with the `log` stub on `PATH`.
The test pins the destination common-directory starts with `==` at the threaded count that it observes. Its doc records the larger count at the spec base as context only.
Tickets 04 and 05 set this pin to the new observed count when their change lowers it. This ticket owns the row.
Raise `worktreeSerialCeiling` from 21 to 22 and `worktreeTestCount` from 878 to 879.

After this chunk checkpoint, the coordinator times the package five times with the spec command (GQ41). A gain inside the 7.6 s baseline spread returns the build to the spec.

## Acceptance

- [ ] One in-process landing under the `log` stub starts the destination common-directory query exactly the pinned threaded count (GQ34).
- [ ] A landing whose window store has no run starts more destination common-directory queries than the pin, so `TestLandingReadsTheDestinationCommonDirOnce` fails.
- [ ] No landing-path or landed-set helper in `internal/worktree` builds its own run.
- [ ] `worktreeSerialCeiling` reads 22 and `worktreeTestCount` reads 879.
- [ ] After the checkpoint, five timed runs give a smallest gain outside the 7.6 s baseline spread (GQ41, review-owned).

## Checkpoint verification

- `bench test --package ./internal/worktree`
- `bench test --package ./internal/gittest`
- `bench test --package ./internal/conformance`

Run `bench consumers` on each changed unexported function before the edit. Record the red of the probe mutation, then the green.
