# Repair: follow-up review hardening

Blocked by: 11-repair-comment-cleanup-and-dedup.md
Writes: internal/worktree/subshell.go, internal/worktree/clean_branch_test.go, specs/worktree-cleanup-eligibility/tickets/12-repair-followup-hardening.md

Source: follow-up review after repairs 09-11 (Standards finding 1, Coverage finding 1).

## What to build

Two trivial, no-behavior-change hardenings found by the follow-up review:

1. **Standards:** the restored `if options.DiscardBranch && !facts.headDetached`
   guard in subshell.go has no comment explaining the detached-HEAD conjunct
   specifically — a future editor could drop it without knowing why. Add one
   clause to the existing comment block naming the invariant: a detached HEAD
   has no branch name for the operator override to authorize deleting.
2. **Coverage:** `TestDiscardBranchLeavesADetachedHeadUnaffected` (clean_branch_test.go)
   only asserts the two zero-value fields (`deleteBranch`, `branchRef`). It
   doesn't anchor that the fixture actually reaches a removing action, so a
   future change that silently reclassified the fixture to a retain verdict
   would leave the test green for the wrong reason. Add an assertion that
   `plan.Action` is a removing action (e.g. `ActionRecoverRemove`, matching what
   this detached+dirty-or-recoverable fixture actually resolves to — verify
   against the live fixture, don't guess) alongside the existing two checks.

## Acceptance

- [ ] subshell.go's DiscardBranch comment names the detached-HEAD invariant.
- [ ] The regression test asserts a removing Action in addition to the zero-value checks.
- [ ] All existing tests remain green with unmodified assertions elsewhere.
