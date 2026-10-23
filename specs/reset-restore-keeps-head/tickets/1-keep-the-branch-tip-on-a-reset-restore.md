# Keep the branch tip on a reset restore

Blocked by: none
Writes: internal/worktree/reset.go, internal/worktree/reset_restore.go, internal/worktree/reset_restore_later_test.go, internal/worktree/parallel_census_test.go, .bench/BENCH-reference.md, CHANGELOG.md
Covers: none

## What to build

`bench worktree reset --restore <ref> <target>` must not drop a commit on the assignment branch.
When the branch tip descends from the envelope's tip, the restore keeps the branch tip and HEAD.
It applies the three-way merge of each preserved layer onto that tip, with the envelope's base as the merge base.
A staged layer that merges to the tip tree drops, as a capture drops it.
A layer that does not merge refuses the plan before any write.
When the branch tip does not descend from the envelope's tip and the envelope's tip does not contain the branch tip, the plan refuses before any write.
A restore onto a branch tip that the envelope's tip contains keeps its current behavior.

## Acceptance

- [ ] A restore after two later commits exits 0, keeps the branch tip and HEAD attached, and restores the staged, unstaged, and untracked work.
- [ ] A staged edit that a later commit carries restores as a clean index.
- [ ] A preserved edit that conflicts with a later commit refuses the plan with `envelope payload conflicts with the branch tip`, and the checkout does not move.
- [ ] A branch that diverged from the envelope's tip refuses the plan with `restore would drop commits on the assignment branch`, and the branch does not move.
