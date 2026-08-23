# Split the skillsindex tests

Blocked by: none
Writes: internal/skillsindex

## What to build

Split `internal/skillsindex/skillsindex_test.go` (832 lines) into render, allowlist, and reference test files. `writeFile` and `reference` stay with the render file and `siblingTemps` with the reference file, each with a single definition, because `command_test.go` calls them. Pure moves only.

## Acceptance

- [ ] R07: `bench structure` no longer lists `internal/skillsindex/skillsindex_test.go`.
- [ ] R03: every created file counts at most 400 newlines.
- [ ] R08: `go test -list '.*' ./internal/skillsindex/` emits the same test-name set at base and tip.
- [ ] R11: each moved fixture has exactly one definition in the package.
- [ ] R12: `bench structure` prints no `DIR CROWDED` line for `internal/skillsindex/`.
- [ ] R18: `bench gate` exits zero before the commit.
