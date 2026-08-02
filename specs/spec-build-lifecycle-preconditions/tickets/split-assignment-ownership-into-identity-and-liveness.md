# Split assignment ownership into identity and liveness

Blocked by: Exempt abandon from the recomposition refusal; Plan an absent target
as cleanable
Ownership fence: `internal/specbuild/precondition.go` (`ownedAssignments` and
its two callers, `preconditions` and `Abandon`),
`internal/specbuild/abandon_test.go`
Assumptions: `ownedAssignments` returns one undifferentiated error string for
every fault it detects; `Abandon` (the plan phase) calls it directly, not
through `preconditions`; `preconditions` already carries an abandon exemption
for a failed ownership probe, but that exemption requires a *prepared* abandon
operation, which only a successful plan produces. Re-derive all three from the
tree at pickup — two earlier tickets in this build have edited this file.

The package's `abandonOwner` test owner is a **counting pass-through to the real
`internal/worktree` `PlanAbandon` and `ApplyAbandon`**, not a stub. That is what
makes LV1 and LV2 assert this fence against the real producer instead of a
fixture, and it is load-bearing: replacing it with a stub to make a criterion
pass would reproduce the per-component-gate-scoping failure, where two correct
disjoint fences were each green alone and red composed because a fixture
accepted what the real counterpart would not. Keep the pass-through. If a
criterion cannot go green through it, the defect is upstream in `Plan an absent
target as cleanable` — say so rather than routing around it.

## What to build

The abandon plan cannot run for the state abandon exists to clean up. When an
assignment's worktree directory is gone, `ownedAssignments` fails, `Abandon`
returns that failure before it ever reaches `abandonmentPlan`, and no prepared
operation is ever written — so the exemption `ApplyAbandon` already carries is
structurally unreachable on a first apply. An interrupted release, which does
leave a cleanup receipt, is already recoverable and must stay that way.

Split the checks the function performs into two classes and let only one of
them soften:

**Identity** — the uniqueness of assignment ID, path, and owner-request across
the record; the owner-request digest; the map key agreeing with the recorded
request; the empty-field rejections; and the resolved assignment's ID and path
agreeing with the record. These stay fatal for every operation, including
abandon. Softening them would let a forged or hand-edited record drive cleanup.

**Liveness** — whether the worktree still exists on disk. For the abandon
mutation only, a failed liveness probe on an unreleased assignment means the
worktree is already gone and there is nothing to release. The recovery refs the
plan enumerates are what preserve the payload, so softening drops no work.

**The probe conflates two facts and must stop.** Today the repository-identity
probe runs `git -C <path> rev-parse --git-common-dir`, which fails both for a
path that does not exist and for a path holding a stranger's checkout. Only the
first may soften. Distinguish them by the path's own existence before the probe
runs, and keep a present-but-foreign path fatal for every operation.

**The exemption is per-mutation.** Widening it to every caller lets `integrate`
write into a vanished worktree. Assert that directly rather than trusting the
call sites.

## Acceptance

- [ ] LV1 — `Abandon` returns a plan for a run whose unreleased assignment worktree directory has been removed, with that assignment present in the plan's worktree list.
- [ ] LV2 — `ApplyAbandon` then completes for that run and leaves it terminal.
- [ ] LV3 — the recovery refs for the removed worktree survive the apply and still name their assignment.
- [ ] LV4 — duplicate assignment ID, duplicate path, duplicate owner-request, and an owner-request digest mismatch each still refuse `Abandon` and `ApplyAbandon`.
- [ ] LV5 — an assignment path that exists but whose git common directory belongs to another repository still refuses `Abandon`.
- [ ] LV6 — every mutation except `abandon` — `start`, `assign`, `checkpoint`, `integrate`, `review`, and `promote` — still refuses on a failed liveness probe.
- [ ] LV7 — an assignment path containing a space and a glob metacharacter is planned and applied unchanged when its directory has been removed.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| LV1 | soften the probe in `preconditions` only, leaving `Abandon`'s direct call fatal | `TestAbandonPlansForRemovedWorktree` | compose a run with an unreleased assignment, `os.RemoveAll` its worktree, call `Abandon`, assert a plan naming that assignment |
| LV2 | return the plan but leave the apply path's ownership call fatal | `TestAbandonAppliesForRemovedWorktree` | as LV1, then `ApplyAbandon` with the plan fingerprint, assert no error and `Terminal` true |
| LV3 | drop recovery refs from the plan when the worktree is absent | `TestRemovedWorktreeRecoveryRefsSurvive` | as LV2, then enumerate `refs/bench/recovery/`, assert the assignment's ref is present and unchanged |
| LV4 | soften the whole ownership function for abandon rather than its liveness half | `TestAbandonStillRefusesForgedAssignmentIdentity` | four subtests, each corrupting one identity fact in the saved record, calling `Abandon` and `ApplyAbandon`, expecting the ownership refusal from both |
| LV5 | treat any probe failure as absence | `TestAbandonRefusesForeignAssignmentCheckout` | initialize an unrelated repository at the assignment path, call `Abandon`, expect the ownership refusal |
| LV6 | key the softening off the assignment's state rather than off the mutation | `TestNonAbandonMutationsStillRefuseAbsentWorktree` | enumerate the `mutation` constants at test time, exclude `abandon` by name, remove the assignment worktree, invoke each remaining one on an unmoved tip, expecting the ownership refusal from every one |
| LV7 | pass the assignment path through a shell or glob expansion before probing it | `TestRemovedWorktreeWithHostilePathIsPlannedAndApplied` | assign at a path containing `a b*c`, remove the directory, plan and apply, assert success and that no sibling worktree was touched |
