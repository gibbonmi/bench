# Give the release the landed verdict of the list

Blocked by: none
Writes: internal/worktree/landed.go, internal/worktree/classifier.go, internal/worktree/unlanded_route_test.go, internal/worktree/resume_test.go, internal/worktree/parallel_census_test.go, CHANGELOG.md
Covers: none

## What to build

An assignment with no own commit and uncommitted output is unlanded work. The `landed` cell of `bench worktree list` shows `landed=false` for it, and the list help row names `bench worktree clean --discard-branch <path>`. But `bench worktree release` refuses that assignment as `dirty` and names `bench worktree release --request <request> <path>` again. That route gives the same refusal each time, so the operator cannot leave the loop.

The cause is in the automatic decision. A release first sets the assignment to cleanup-pending. The decision of a cleanup-pending assignment reads the ancestry proof of the explicit plan, and that proof does not apply the own-commit rule. The landed verdict that applies the rule reads only an active assignment.

One shared function in `internal/worktree/landed.go` gives the landedness of an assignment in any state. It returns the ancestry proof, except that an assignment that holds unlanded output has not landed. The landed verdict of an active assignment reads this function. The automatic classifier gives this function to the decision of a cleanup-pending release. Thus the release refuses the assignment as `unmerged`, and its refusal names the same discard clean as the list help row.

An assignment whose own commit `main` holds keeps its landed verdict. Its release with uncommitted output still refuses as `dirty` and names the release route. The automatic-planner test of ordinary dirt now lands an own commit first. Dirt on a branch with no own commit is unlanded work.

## Acceptance

- [ ] The release of an assignment with no own commit and an untracked file refuses as `unmerged`, and its `next=` line names `bench worktree clean --discard-branch <path>`.
- [ ] After that refusal, the `bench worktree list` help row names the same route, and that route plans the recovery removal and removes the tree.
- [ ] The release of an assignment with a landed own commit and an untracked file still refuses as `dirty`. Its `next=` line names the release route.
- [ ] `CHANGELOG.md` records the change.
