# Split the maps tests

Blocked by: none
Writes: internal/maps

## What to build

Split `internal/maps/maps_test.go` (937 lines) into command, parse, and graph test files. `hasTreeDiagnostic` moves with the parse group, and `hasDiagnostic` moves with the graph group. The active-rows projection tests join the command file, because both exercise the public `Command` entry point. Pure moves only.

## Acceptance

- [ ] R07: `bench structure` no longer lists `internal/maps/maps_test.go`.
- [ ] R03: every created file counts at most 400 newlines.
- [ ] R08: `go test -list '.*' ./internal/maps/` emits the same test-name set at base and tip.
- [ ] R12: `bench structure` prints no `DIR CROWDED` line for `internal/maps/`.
- [ ] R18: `bench gate` exits zero before the commit.
