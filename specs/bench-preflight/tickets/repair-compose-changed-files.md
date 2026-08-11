# Repair: compose the changed-file set from the diff package

Blocked by: none
Ownership fence: `internal/diff/`, `internal/preflight/`
Integration surfaces: NUL-pair parser→existing `parseNameStatusZ`/`changedFiles` in `internal/diff/diff.go`; consumer→`internal/preflight/gather.go` `changedFilePaths` (replaced)
Contracts: the changed-file set (ordered repo-relative paths, `bench diff`'s exact committed+index+tracked-worktree semantics) crosses `internal/diff`→`internal/preflight/`, asserted by RS1 against the real exported function
Closure: RS1/exported-changed-set-consumed, RS1/no-inline-name-status-parse

## What to build

Review finding S1: `internal/preflight/gather.go:262-277` re-implements `git
diff --name-status --no-renames -z` plus the NUL-pair parse that
`internal/diff` already owns — the duplication the spec's compose-never-
re-derive decision forbids by name. Export a changed-path resolution from
`internal/diff` (composing the existing `changedFiles`/`parseNameStatusZ`,
mirroring the `ResolveReviewBase` export precedent), make `bench diff`'s own
path and `internal/preflight`'s gatherer both consume it, and delete
`changedFilePaths`' inline reimplementation and its stated-duplication
comment. Bare `bench diff` output stays byte-identical; every existing
preflight CLI test stays green unchanged.

## Acceptance

- [ ] [RS1] (covers local) `internal/preflight` derives its changed set through the exported diff-package function; no `--name-status` invocation or NUL-pair parse remains in `internal/preflight`, and all existing preflight and diff tests pass unchanged.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RS1/exported-changed-set-consumed | make the exported function drop tracked-worktree changes | the existing preflight out-of-fence CLI tests (worktree-seeded) | `go test ./internal/preflight -run 'OutOfFence\|Fence'`, expect the missed-path failure |
| RS1/no-inline-name-status-parse | (structural) leave the inline parse in gather.go | the verifying coordinator's grep for `name-status` under `internal/preflight/` | run the grep, expect the surviving hit as the red |
