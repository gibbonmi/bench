# Give the index-file git form one owner

Blocked by: none
Writes: internal/git/tree.go, internal/landing/gitexec.go
Covers: none

## What to build

The quality survey of 2026-09-29, in its "Small certain cuts" table, found that `internal/landing/gitexec.go` restates the index-file git form that `internal/git/tree.go` owns. The form is `git -C <root> <args>` with `GIT_INDEX_FILE=<idx>` appended to the process environment. The landing package spells it three times, in `indexRun`, `indexOutputRaw`, and `indexOutput`. The git package spells it once, in `idxCommand`.

Export `idxCommand` from `internal/git` as `IndexCommand`. Make the three landing helpers build their command through `IndexCommand`. Each helper keeps its own result shape:

- `indexRun` returns the error of `Run`.
- `indexOutputRaw` returns the bytes and the error of `Output`.
- `indexOutput` returns the output of `Output` with `strings.TrimSpace`.

Thus the landing path keeps the same argv, environment, captured standard error, and error values.

This change touches the `bench worktree land` seam. It must not change behavior.

## Acceptance

- [ ] `internal/landing` contains no `GIT_INDEX_FILE` text and no `os.Environ` call.
- [ ] `internal/git` has the exported `IndexCommand` and no `idxCommand`.
- [ ] The three landing helpers keep their names, signatures, and result shapes.
- [ ] `go vet ./...` passes.
- [ ] `bench test --changed` passes for the changed packages.
- [ ] `bench test --check git-plumbing-owner` passes.
- [ ] A probe that removes the `GIT_INDEX_FILE` environment line from `IndexCommand` makes a landing test fail.
