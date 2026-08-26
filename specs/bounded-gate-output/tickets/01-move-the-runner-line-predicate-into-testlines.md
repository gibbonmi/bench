# Move the runner-line predicate into testlines

Blocked by: none
Writes: internal/testlines/testlines.go (new), internal/testlines/testlines_test.go (new), internal/testreport/testreport.go, internal/testreport/testreport_test.go

## What to build

A new low package, `internal/testlines`, owns the classification of plain
`go test` output. `internal/testreport` composes that package and keeps no
second copy of the predicate.

This contract crosses into tickets 02, 03, and 04:

- `testlines.RunnerLine(line string) bool` reports a runner line.
- `testlines.FailureRows(lines []string) []string` returns the failure rows of one red stream.

`RunnerLine` is the predicate `testreport.runnerLine` holds today. Move it
whole and delete the original. `testreport.decode` then calls
`testlines.RunnerLine`.

`FailureRows` reads a plain stream, never `-json` output. A `--- FAIL:` line
at any indent starts a failure block and is the first row. A `# <package>`
line starts a build-error block and is the first row. Each later line is a
row until the next runner line ends the block.

A `FAIL <package>` line and a
`panic:` line are each a row on their own. An `ok` line and a `?` line are
never rows. A `--- PASS:` line and a `=== RUN` line end a block and add no
row. When no rule matches, `FailureRows` returns an empty slice; the caller
owns the fallback.

The package imports neither `internal/gate` nor `internal/testreport`. That
keeps it low enough for both callers.

## Acceptance

- [ ] A stream with two `--- FAIL:` lines yields two rows that hold those lines. (BG03)
- [ ] A FAIL block with two diagnostic lines yields two rows after the name row. (BG04)
- [ ] A stream with `FAIL\tgithub.com/x/y [build failed]` and `panic: boom` yields both lines as rows. (BG05)
- [ ] A stream with `# github.com/x/y` and `./x.go:12:3: undefined: y` yields both as rows. (BG31)
- [ ] An indented `--- FAIL:` line starts a block like a top-level line does.
- [ ] A `--- PASS:` line ends a block and adds no row.
- [x] A stream of only `WARNING: DATA RACE` lines yields rows (superseded by
      BG38, ticket 08: a `WARNING: DATA RACE` line now opens a block, so this
      no longer yields an empty slice).
- [ ] `testreport.decode` calls `testlines.RunnerLine`, and no copy of the predicate remains in `internal/testreport`.
