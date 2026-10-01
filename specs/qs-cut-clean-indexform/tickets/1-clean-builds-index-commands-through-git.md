# Build the clean index-file git commands through git.IndexCommand

Blocked by: none
Writes: internal/worktree/clean.go
Covers: none

## What to build

The `internal/git` function `IndexCommand` builds the index-file git form: `git -C <root> <args>` with `GIT_INDEX_FILE=<idx>` appended to the process environment. Its comment says that it is the one owner of this form. A review found that `internal/worktree/clean.go` still spells the form seven times. These calls are in `worktreeTree`, `realIndexTree`, and `conflictTree`. Each call gives `GIT_INDEX_FILE=` to the local helper `gitInput`, and that helper builds the same argv and environment.

Make each index-file call in `clean.go` build its command through `git.IndexCommand`. Keep the stdin feed and the error text in one local runner that takes a prebuilt `*exec.Cmd`:

- The runner writes the input bytes to stdin.
- The runner returns stdout without trailing newlines.
- On failure, the runner returns `git <verb>: <stderr>`.

The helper `gitInput` keeps its signature for its plain callers, which include `commitTree` with its author environment. It builds its command and then calls the same runner.

This change touches the `bench worktree clean` seam. It must not change behavior. The argv, the environment, the stdin, the error text, and the exit codes stay the same.

## Acceptance

- [ ] `internal/worktree/clean.go` contains no `GIT_INDEX_FILE` text.
- [ ] Each index-file call in `clean.go` builds its command through `git.IndexCommand`.
- [ ] `gitInput` keeps its signature, and its callers outside `clean.go` do not change.
- [ ] `go vet ./...` passes.
- [ ] `bench test --changed` passes for the changed packages.
- [ ] `bench test --check git-plumbing-owner` passes.
- [ ] A probe that removes the `GIT_INDEX_FILE` environment line from `IndexCommand` makes a clean test in `internal/worktree` fail.
