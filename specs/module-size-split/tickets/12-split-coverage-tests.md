# Split the coverage tests

Blocked by: none
Writes: internal/coverage

## What to build

Split `internal/coverage/coverage_test.go` (783 lines) into parse, command, and schema test files. The shared fixtures (`stories`, the header constants, `mapShape`, `spec`) sit beside the parse tests; `mustWrite` moves with the command file. Pure moves only.

## Acceptance

- [ ] R07: `bench structure` no longer lists `internal/coverage/coverage_test.go`.
- [ ] R03: every created file counts at most 400 newlines.
- [ ] R08: `go test -list '.*' ./internal/coverage/` emits the same test-name set at base and tip.
- [ ] R11: each moved fixture has exactly one definition in the package.
- [ ] R12: `bench structure` prints no `DIR CROWDED` line for `internal/coverage/`.
- [ ] R18: `bench gate` exits zero before the commit.
