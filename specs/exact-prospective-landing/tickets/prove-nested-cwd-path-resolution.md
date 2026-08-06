# Prove nested-CWD path resolution

Blocked by: adopt-exact-landing-in-commit.md
Ownership fence: `internal/contract/runtime/runtime_commit_test.go`
Integration surfaces: invocation-CWD path resolution→existing `internal/commit/commit.go` plus NCWD1; real compiled command proof→`internal/contract/runtime/runtime_commit_test.go` plus NCWD1
Contracts: a valid relative path supplied below the repository root crosses the real `bench commit` adapter→the exact landing owner, asserted by NCWD1 against distinct root and nested files in `internal/contract/runtime/runtime_commit_test.go`

## What to build

Add the missing real-binary proof for a valid `bench commit` invocation below the
repository root. From `sub/`, name `a.txt`; the command must land `sub/a.txt` and
must not attribute or land a distinct root-level `a.txt`. Preserve the existing
usage-error nested-CWD case and do not change production code unless this public
proof exposes an actual resolution defect.

## Acceptance

- [ ] [NCWD1] Running the compiled command from `sub/` with argument `a.txt` lands the nested file, leaves the same-named root file outside the commit, and returns with the named nested path clean.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| NCWD1 | resolve a valid relative argument from repository root instead of invocation CWD | real-binary nested-CWD runtime case | run `go test ./internal/contract/runtime -run '^TestCommitResolvesValidPathsFromNestedCWD$' -count=1`; expect the landed tree to contain the wrong `a.txt` path |
