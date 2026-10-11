# Hold the kit CLAUDE.md to the managed form

Blocked by: none
Writes: internal/adopt/claude_md_mirror_test.go
Covers: none

## What to build

The root `CLAUDE.md` of the kit must hold the exact bytes that `benchClaudeMD` in `internal/adopt/marker.go` returns. No check holds this rule now. When a line of the managed form goes away from the function or from the file, the package tests stay green.

Add one test in `internal/adopt`. The test reads the root `CLAUDE.md` of the checkout that it runs in. The test compares those bytes with the return value of `benchClaudeMD`. When the two differ, the test fails and shows both values.

## Acceptance

- [ ] `bench test --package ./internal/adopt` passes when the root `CLAUDE.md` holds the same bytes that `benchClaudeMD` returns.
- [ ] `bench test --package ./internal/adopt` fails when `benchClaudeMD` does not return the `@AGENTS.md` import line that the root `CLAUDE.md` holds.
- [ ] `bench test --package ./internal/adopt` fails when `benchClaudeMD` does not return the `@.bench/BENCH.md` import line that the root `CLAUDE.md` holds.
