# Print the traversal diagnostics when the scan failure is absent

Blocked by: none
Writes: internal/git/worktree_admin_hostile_test.go
Covers: none

## What to build

`TestWorktreesPropagatesScanTraversalFailureBeforePorcelain` in
`internal/git/worktree_admin_hostile_test.go` failed once on CI with the
message `scan traversal failure = <nil>`. The test passes locally and on the
next CI run. The message names no cause. It does not say who the process is,
what mode the directory had at that moment, or whether a second read of the
directory succeeds.

When the `Worktrees` call returns no `WorktreeScanError`, the failure message
carries three more facts: the effective uid of the test process, the mode of
the unreadable directory from a fresh `os.Lstat`, and the error of a fresh
`os.ReadDir` on that directory, or the word `readable` when that read
succeeds. The passing path does not change. The guard that skips a privileged
host does not change. No production code changes.

## Acceptance

- [ ] When the scan error is absent or malformed, the failure message prints the effective uid, the directory mode, and the fresh read result.
- [ ] The test still passes as a non-root user, and the privileged-host skip is unchanged.
- [ ] `go test ./internal/git -run 'TestWorktreesPropagatesScanTraversalFailureBeforePorcelain$' -count=1` exits 0.
