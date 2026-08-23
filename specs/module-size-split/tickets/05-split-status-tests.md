# Split the status test file

Blocked by: none
Writes: internal/status

## What to build

Split `internal/status/status_test.go` (1574 lines) into six family files plus one fixture file. The families: signals, the producible-actions table, counters, gate-cache, render, and command routing. `status_fixtures_test.go` holds `initRepo` and `gitRun`, which every family calls.

Group-local fixtures move with their family. The directory lands at exactly 12 files. Pure moves only. `internal/status/status.go` stays untouched under its standing accept.

## Acceptance

- [ ] R05: `bench structure` no longer lists `internal/status/status_test.go`.
- [ ] R03: every created file counts at most 400 newlines.
- [ ] R08: `go test -list '.*' ./internal/status/` emits the same test-name set at base and tip.
- [ ] R11: each moved fixture has exactly one definition in the package.
- [ ] R12: `bench structure` prints no `DIR CROWDED` line for `internal/status/`.
- [ ] R18: `bench gate` exits zero before the commit.
