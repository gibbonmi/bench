# Emit the absolute worktree path

Blocked by: none
Writes: internal/worktree/path.go, internal/worktree tests, cmd/bench/main.go

## What to build

`bench worktree path` prints the resolved absolute path, not the `~`-compacted
form. No shell expands a quoted tilde, so every caller hand-builds the path
today. The other worktree verbs keep their acceptance of the `~` form.

## Acceptance

- [ ] `bench worktree path <target>` prints a path that starts with `/`.
- [ ] The printed path is accepted verbatim by `bench worktree clean` and by `cd`.
- [ ] The help row that `bench help` prints for the verb no longer says "portable".
