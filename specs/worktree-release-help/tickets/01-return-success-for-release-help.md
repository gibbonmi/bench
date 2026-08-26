# Return success for release help

Blocked by: none
Writes: internal/worktree/worktree.go, internal/worktree/worktree_test.go

## What to build

The release command prints its usage and exits successfully when the user requests help.
Invalid release arguments still print the usage and exit with a grammar error.

## Acceptance

- [ ] `bench worktree release --help` prints the release usage on standard output and exits 0.
- [ ] An invalid release argument prints the release usage on standard error and exits 2.
