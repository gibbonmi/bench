# Print the worktree path on exec failures

Blocked by: pass-argv-stdin-and-env-to-the-exec-child.md
Writes: internal/worktree/exec.go, internal/worktree/exec_test.go

## What to build

`runWorktreeChild` writes the line `worktree: <absolute path>` to stderr on
every failure after the target resolves. The function writes that line after a
start failure's own error line, after a nonzero child exit, and on the cancel
path. A zero exit prints nothing, so a green run stays silent.

The child's stdout, stderr, and exit code still pass through unchanged. A child
that exits 3 exits 3, and a child killed by SIGINT through the cancel path
exits 130. The verb rewrites no byte of the child's stderr, and it appends the
`worktree:` line after that stderr.

File invariant for `internal/worktree/exec.go`: one writer prints the
`worktree:` line, and it prints the line one time for one child run. A second
printer at a call site duplicates the line on a start failure.

## Acceptance

- [ ] X5: a child `sh -c 'echo child-usage >&2; exit 2'` yields exit 2, and
      stderr holds `child-usage` and then the `worktree:` line.
- [ ] F1: `-- no-such-binary-xyz` exits 1 with the stderr line
      `bench worktree exec: exec: ...` and then `worktree: <absolute path>`.
- [ ] F2: `-- sh -c 'exit 3'` exits 3, and stderr ends with the
      `worktree: <absolute path>` line.
- [ ] F3: a child killed by SIGINT through the cancel path exits 130, and
      stderr ends with the `worktree:` line.
- [ ] F4: `-- true` exits 0 with empty stderr.
