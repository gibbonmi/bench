## Standards

Finding count: 0. Worst: none.

No actionable Standards findings.

## Spec

Finding count: 1. Worst: unauthorized behavioral narrowing of the reviewer-resolved admin-entry predicate.

- `ask-user` — The frozen spec and its named decision source require every direct
  entry under `<git-common-dir>/worktrees/<id>/` to be regular or a directory
  (`specs/worktree-enumeration-hang/decisions/worktree-enumeration-hang.md:59`).
  The candidate rewrites that predicate to exempt `bench-lease`
  (`specs/worktree-enumeration-hang/spec.md:61`) and skips that name regardless
  of shape (`internal/git/git.go:244`). A FIFO, symlink, socket, or device at
  `worktrees/<id>/bench-lease` is therefore accepted without a reviewer decision.
  The reviewer must either retain the every-entry rule or explicitly approve and
  cover the lifecycle-control-record exception.

## Coverage

Finding count: 1. Worst: the scan-I/O typed-error partition is not exercised through `git.Worktrees`.

- `auto-fix` — The spec requires every `git.Worktrees` failure except a genuine
  porcelain-child nonzero exit to carry its typed recovery action
  (`specs/worktree-enumeration-hang/spec.md:127`). The scanner has fallible
  `Lstat` and `ReadDir` operations at the root, id, and child levels
  (`internal/git/git.go:205`), but the only `WorktreeScanError` test calls
  `ScanWorktreeAdmin` directly (`internal/git/git_test.go:223`). Add a
  deterministic `git.Worktrees` test for a scan traversal failure that asserts
  the typed path/action and proves porcelain was not launched. This kills an
  implementation that swallows only `WorktreeScanError` while leaving the
  existing shape-error integration tests green.
