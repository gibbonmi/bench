# Split the landing package files

Blocked by: none
Writes: internal/landing

## What to build

Split `internal/landing/landing.go` (953 lines) into responsibility files. The retained `landing.go` keeps the `Owner` orchestration and the fingerprint group. `composition.go` takes the merge-tree engine. `spectree.go` takes the spec-transition tree edits. `attribution.go` takes path validation and composition. `gitexec.go` takes the raw git process helpers.

Split `landing_test.go` (973 lines) into `landing_test.go`, `landing_reviewed_test.go`, and `landing_helpers_test.go`. The shared fixtures move once into the helpers file. The directory lands at 12 files. Pure moves only; no symbol changes name, visibility, or signature.

## Acceptance

- [ ] R01: `bench structure` no longer lists `internal/landing/landing.go`.
- [ ] R05: `bench structure` no longer lists `internal/landing/landing_test.go`.
- [ ] R03: every created file counts at most 400 newlines.
- [ ] R04: `go build ./...` exits zero with no edit outside `internal/landing/`.
- [ ] R08: `go test -list '.*' ./internal/landing/` emits the same test-name set at base and tip.
- [ ] R11: each moved fixture has exactly one definition in the package.
- [ ] R12: `bench structure` prints no `DIR CROWDED` line for `internal/landing/`.
- [ ] R18: `bench gate` exits zero before the commit.
