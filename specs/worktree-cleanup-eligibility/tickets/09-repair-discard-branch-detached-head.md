# Repair: restore the detached-HEAD guard around DiscardBranch

Blocked by: 08-contract-eligibility-policy.md
Writes: internal/worktree/subshell.go, internal/worktree/clean_branch_test.go, specs/worktree-cleanup-eligibility/tickets/09-repair-discard-branch-detached-head.md

Source: reviews/worktree-cleanup-eligibility.md, Coverage axis finding 1.

## What happened

Ticket 03's extraction moved the `options.DiscardBranch` override so it now applies
outside the `headRef != "detached"` guard it lived inside at base. On a registered
worktree with a detached HEAD, `headRef` is the literal string `"detached"`, so
`--discard-branch` now sets `deleteBranch=true` and `branchRef="detached"`,
changing the fingerprint and, on apply, attempting to delete a branch literally
named `"detached"`. Reachable via `bench worktree clean --discard-branch <path>`
and `--landed --discard-branch`. No existing fixture exercises `DiscardBranch`
together with a detached HEAD.

## What to build

Restore the derived-after `DiscardBranch` application to only apply when HEAD is
not detached, matching base behavior exactly (verify against `868a4e4e` directly,
don't work from this description alone). Add one regression case proving a
detached-HEAD, `DiscardBranch: true` fixture leaves `plan.deleteBranch == false`
and `plan.branchRef == ""`.

## Acceptance

- [ ] `--discard-branch` on a detached-HEAD worktree does not set deleteBranch/branchRef, matching base.
- [ ] A new test in clean_branch_test.go (alongside `TestDiscardBranchLeavesTheDerivedClassificationUnchanged`/`TestDiscardBranchNeverBypassesARefusal`) pins the detached+DiscardBranch case.
- [ ] All pre-existing DiscardBranch/detached tests remain green and unmodified.
