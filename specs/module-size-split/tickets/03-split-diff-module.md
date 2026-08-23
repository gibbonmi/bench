# Split the diff module

Blocked by: none
Writes: internal/diff

## What to build

Split `internal/diff/diff.go` (900 lines) into exactly three files. The retained `diff.go` keeps the command grammar and the render functions. `snapshot.go` takes the git-shell parsing, the snapshot identity, and the drift detection. `range.go` takes the movement types and the range resolution.

The directory holds 10 files today and lands at exactly 12. Pure moves only. `SetSnapshotAfterReadForTest` and every exported symbol stay in package `diff` unchanged.

## Acceptance

- [ ] R01: `bench structure` no longer lists `internal/diff/diff.go`.
- [ ] R03: every created file counts at most 400 newlines.
- [ ] R04: `go build ./...` exits zero with no edit outside `internal/diff/`.
- [ ] R08: `go test -list '.*' ./internal/diff/` emits the same test-name set at base and tip.
- [ ] R12: `bench structure` prints no `DIR CROWDED` line for `internal/diff/`.
- [ ] R18: `bench gate` exits zero before the commit.
