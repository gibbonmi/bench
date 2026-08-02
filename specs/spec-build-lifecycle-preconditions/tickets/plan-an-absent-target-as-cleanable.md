# Plan an absent target as cleanable

Blocked by: nothing
Ownership fence: `internal/worktree/path.go`, `internal/worktree` plan and
resume tests
Assumptions: `planAbandon` short-circuits on a prior cleanup receipt one branch
*earlier* than its absent-path refusal, which is why an interrupted release
already resumes cleanly. `PlanExplicit` enumerates `git.Worktrees(root)` rather
than the filesystem, so a worktree whose directory was removed but whose
registration survives is still discoverable there. Verify that claim against the
tree before relying on it — if the registration is also gone, this ticket's
shape changes and the coordinator wants to hear so before you build.

## What to build

An assignment whose worktree directory no longer exists on disk is cleanable,
not an ownership fault. Today `planAbandon` disagrees: it `Lstat`s the target
and, on `ErrNotExist`, returns the same
`abandon request, assignment, or path mismatch; checkout retained` refusal it
gives a genuinely foreign checkout. That refusal is the second of two sites
that make `bench spec build abandon` unable to retire a run whose worktrees are
already gone — the exact state an escape hatch exists for.

Generalize the receipt escape that already sits one branch above it: an absent
target with **no** prior cleanup receipt reaches a cleanable plan the same way
an interrupted release with a receipt does. This does not replace the receipt
path; the receipt path keeps its own behavior unchanged.

**Absent and foreign are different facts and must stay different.** A path that
does not exist softens. A path that exists but belongs to another repository
still refuses — conflating the two lets cleanup act on a stranger's checkout.
Whatever distinguishes them must be a positive check on the path's own
identity, not the absence of an error from a probe that fails for both reasons.

**Plan is not the whole contract.** The consumer of this plan calls
`ApplyAbandon` immediately after `PlanAbandon`, so a plan that succeeds and an
apply that then refuses leaves the caller worse off than the refusal it
replaced. Both halves are this ticket's, and both are asserted here rather than
left to the consumer.

## Acceptance

- [ ] WT1 — `PlanAbandon` returns a cleanable fingerprint for a registered owned worktree whose directory has been removed and that has no prior cleanup receipt.
- [ ] WT2 — `ApplyAbandon` completes for that same removed-directory assignment using the fingerprint WT1 returned.
- [ ] WT3 — a target path that exists but whose git common directory belongs to another repository still refuses.
- [ ] WT4 — an interrupted release that did leave a cleanup receipt still resumes through the receipt path, with its plan bytes unchanged from today.
- [ ] WT5 — a removed-directory target whose path contains a space and a glob metacharacter plans and applies unchanged.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| WT1 | keep the `ErrNotExist` refusal and soften only the non-`ErrNotExist` `Lstat` error | `TestPlanAbandonAcceptsRemovedDirectory` | create an owned worktree, `os.RemoveAll` its directory, call `PlanAbandon`, expect a non-empty fingerprint and no error |
| WT2 | soften the plan only, leaving apply to fail on the missing directory | `TestApplyAbandonCompletesForRemovedDirectory` | as WT1, then call `ApplyAbandon` with the returned fingerprint, expect no error and the registration gone |
| WT3 | treat every `Lstat` failure and every ownership-probe failure as absent | `TestPlanAbandonRefusesForeignCheckout` | create a second repository at the target path, call `PlanAbandon`, expect the mismatch refusal |
| WT4 | route the receipt case through the new absent-path branch | `TestInterruptedReleaseStillResumesThroughReceipt` | drive a release to an in-flight receipt, remove the directory, call `PlanAbandon`, assert the plan equals the receipt-derived plan |
| WT5 | canonicalize the target through a shell-expanding or glob-expanding path helper | `TestRemovedDirectoryWithHostilePathPlansAndApplies` | create an owned worktree at a path containing `a b*c`, remove it, plan and apply, expect success and no collateral removal |
