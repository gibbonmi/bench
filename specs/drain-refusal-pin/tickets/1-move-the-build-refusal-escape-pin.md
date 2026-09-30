# Move the build refusal escape pin to the shared printer test

Blocked by: none
Writes: internal/treetarget/build_test.go, internal/worktree/target_refusal_test.go
Covers: none

## What to build

Two tests pin the same escaped lines of the shared worktree refusal printer. The row `a control byte on either line escapes` in `internal/treetarget` calls `worktree.PrintTreeBuildRefusal`. The table `TestTargetRefusalEscapesOnlyAnUnsafeLine` in `internal/worktree` calls the printer directly. So one expectation has two sources.

The `internal/treetarget` row also pins one fact that the table does not pin: `PrintTreeBuildRefusal` puts the label into the `next=` line through the build grammar. Move the row into the `internal/worktree` table as a row that calls `PrintTreeBuildRefusal`. Then remove the row from `internal/treetarget`. Do not remove the `--in` seam rows of `internal/treetarget`, such as the backslash row.

## Acceptance

- [ ] Exactly one test pins the escaped lines of a build refusal with a control byte in the detail and the label.
- [ ] That test calls `PrintTreeBuildRefusal`, so it also pins the label in the `next=` line.
- [ ] A probe that always escapes and a probe that never escapes each turn the moved row or its table red.
- [ ] The packages `internal/worktree` and `internal/treetarget` pass.
