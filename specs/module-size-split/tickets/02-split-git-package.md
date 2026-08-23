# Split the git package files

Blocked by: none
Writes: internal/git

## What to build

Split `internal/git/git.go` (761 lines) along responsibility. The retained `git.go` keeps the process primitives, the ref checks, and the folded branch-lifecycle group. `worktree_admin.go` takes admin-dir integrity and enumeration. `status.go` takes the porcelain parsing and the repo-fact types.

Split `git_test.go` (928 lines) into family files. The set: a shared `testhelpers_test.go`, two worktree-admin files (enumeration versus hostile shapes), a fact-family file, and a refs file. Fold files as needed so the directory lands at 12 or fewer. Pure moves only.

## Acceptance

- [ ] R01: `bench structure` no longer lists `internal/git/git.go`.
- [ ] R05: `bench structure` no longer lists `internal/git/git_test.go`.
- [ ] R03: every created file counts at most 400 newlines.
- [ ] R04: `go build ./...` exits zero with no edit outside `internal/git/`.
- [ ] R08: `go test -list '.*' ./internal/git/` emits the same test-name set at base and tip.
- [ ] R11: each moved fixture has exactly one definition in the package.
- [ ] R12: `bench structure` prints no `DIR CROWDED` line for `internal/git/`.
- [ ] R18: `bench gate` exits zero before the commit.
