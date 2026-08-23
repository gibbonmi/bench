# Split the roadmap context tests

Blocked by: none
Writes: internal/roadmap

## What to build

Split `internal/roadmap/context_test.go` (975 lines) into command, row-selector, and occurrence test files. Fold `context_types.go` into `context_parse.go` in the same change, so the directory lands at exactly 12 files. The trailing parse-failure test joins the row-selector file. Pure moves only.

## Acceptance

- [ ] R06: `bench structure` no longer lists `internal/roadmap/context_test.go`.
- [ ] R03: every created file counts at most 400 newlines.
- [ ] R04: `go build ./...` exits zero.
- [ ] R08: `go test -list '.*' ./internal/roadmap/` emits the same test-name set at base and tip.
- [ ] R12: `bench structure` prints no `DIR CROWDED` line for `internal/roadmap/`.
- [ ] R18: `bench gate` exits zero before the commit.
