# Grade the merge lane from the target checkout for every caller root

Blocked by: none
Writes: internal/landing/merge.go, internal/worktree/merge_caller_root_test.go (new)
Covers: none

## What to build

The `bench worktree merge` verb resolves its fast lane against the target worktree.
The lane checks name the target worktree as their file root.
Then the landing merge grades the composed tree with the caller's root as the lane root.
From the primary checkout, the lane does not move the target anchor into the private
checkout of the composed tree.
Thus the prose check reads the incoming Markdown from the target checkout.
The merge did not put that file there, and the prose check refuses the merge.

Make the landing merge grade the composed tree from the target worktree, as `bench commit`
in that worktree does.
The Git operations of the merge stay on the caller's root.

## Acceptance

- [ ] A sibling adds one Markdown file, and a merge from the primary checkout passes its
  lane. The prose check reads that file from the composed tree.
- [ ] A regression test runs the merge from a caller root that is not the target
  worktree. The test fails when the lane root is the caller's root.
