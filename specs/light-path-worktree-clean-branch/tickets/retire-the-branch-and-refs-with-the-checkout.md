# Retire the branch and refs with the checkout

Blocked by: none

Ownership fence: `internal/worktree`, `internal/contract/runtime/runtime_worktree_clean_branch_test.go`
Assumptions: `plan.deleteBranch` / `git.DeleteBranchExact` already exist and are exercised for the landed case; the classifier's fail-closed default stays untouched

## What to build

`bench worktree clean` removes a checkout and then leaves its assignment branch and a
fresh recovery ref behind, so cleaning a worktree produces two new pieces of residue
instead of none. After a ten-ticket build this left ten branches, eleven recovery refs,
and eleven stale pool entries that no sanctioned command could retire — the operator had
to hand-write `git branch -D` and `git update-ref -d` loops.

The machinery is already there: the apply path deletes the branch when
`plan.deleteBranch` is set, and the classifier sets it only when `LandedInDefault` proves
the work landed. That proof tries ancestry, then merge detection, then `git cherry`
patch-equivalence — and all three miss a **squash-landing**, where the branch's commits
were composed into one commit on the default branch. So a fully-landed build classifies
as `unmerged`, retains the branch, and mints a recovery ref for work that is already in.

Add an explicit operator override that says "I know this landed; take the branch and the
ref with the checkout." The classifier's automatic path and its fail-closed default do not
change: absent the override, an unprovable payload is still retained. The override is the
operator supplying the proof the tool cannot derive, exactly as `bench gate --fresh` lets
an operator overrule a verdict the closure still calls current.

It must be plan-then-apply like the rest of the command, so the operator sees the branch
and ref named before anything is removed, and it must be refused in the same breath as any
other precondition failure — an override that skips the ownership and identity checks would
be a different and much worse command.

## Acceptance

- [ ] [RW1] With the override, applying a cleanup removes the checkout, deletes its assignment branch, and leaves no recovery ref behind for that assignment.
- [ ] [RW2] The plan names the branch and the ref it will delete before any removal happens.
- [ ] [RW3] Without the override, behavior is byte-for-byte what it is today, including retaining an unlanded payload.
- [ ] [RW4] The override never bypasses ownership, identity, or path-safety refusals — a foreign or mismatched target is still refused with the same reason.
