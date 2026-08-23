# Split the worktree land files

Blocked by: none
Writes: internal/worktree/land.go, internal/worktree/land_test.go

## What to build

Split `internal/worktree/land.go` (717 lines) into `land.go` (first-run flow), `land_resume.go`, `land_identity.go`, and `land_refusal.go`. Split `land_test.go` (2210 lines) into journey-family files. The families: real-git journeys, resume, flags, reauthorization, release refusal, freshness, spec-less, spec-amendment, and tickets-only. One `land_fixtures_test.go` holds the cross-family fixtures: `buildLandingBinary`, the landing fixtures, `landArgs`, and `specLessLandArgs`.

Pure moves only. The directory's pre-existing crowding violation stays a single violation. Its file count worsens, and the spec flags this call for reviewer veto. Ticket 13 shares the `internal/worktree/` fence and must not run in parallel with this ticket.

## Acceptance

- [ ] R02: `bench structure` no longer lists `internal/worktree/land.go`.
- [ ] R06: `bench structure` no longer lists `internal/worktree/land_test.go`.
- [ ] R03: every created file counts at most 400 newlines.
- [ ] R04: `go build ./...` exits zero.
- [ ] R08: `go test -list '.*' ./internal/worktree/` emits the same test-name set at base and tip.
- [ ] R11: each moved fixture has exactly one definition in the package.
- [ ] R18: `bench gate` exits zero before the commit.
