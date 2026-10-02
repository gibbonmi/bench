# Derive each worktree command name from its usage line

Blocked by: none
Writes: internal/usage/parse.go, internal/usage/command_name_test.go (new), internal/worktree/show.go, internal/worktree/build.go, internal/worktree/merge.go, internal/worktree/reset.go, internal/worktree/exec.go, internal/worktree/land.go
Covers: none

## What to build

The `usage` package exports `CommandName`. This function returns the leading
command words of a usage line. It stops at the first operand, flag, or group
token.

The show, build, merge, reset, exec, and land grammars of
`internal/worktree` set their `Cmd` from `CommandName` and the `usage`
constant of the verb. Each target refusal in those files reads the name from
its grammar. Thus the `usage` constant is the one source of each command name.
The printed text of each grammar refusal and target refusal does not change.

## Acceptance

- [ ] `usage.CommandName` returns `bench worktree land` for `usage.WorktreeLand`,
      `bench worktree reset` for `usage.WorktreeReset`, and `bench worktree list`
      for `usage.WorktreeList`; a test asserts each case.
- [ ] No production file under `internal/worktree` holds a string literal
      `"bench worktree show"`, `"bench worktree build"`, `"bench worktree merge"`,
      `"bench worktree reset"`, `"bench worktree exec"`, or `"bench worktree land"`.
- [ ] The existing worktree refusal tests stay green with no edit.
