## Standards

Finding count: 2. Worst: exported Go symbols lack doc comments.

- `auto-fix` — Add a current, full-sentence `WorktreeFailure` doc comment.
- `auto-fix` — Add a current, full-sentence `WorktreeScanError` doc comment.

## Spec

Finding count: 1. Worst: a pre-enumeration refusal has the porcelain label.

- `auto-fix` — `PruneLandedBranches` must use the neutral `worktree discovery
  failed` wrapper for every `Worktrees` error, as the approved spec requires.

## Coverage

Finding count: 0. Worst: none.

No coverage repair is required.
