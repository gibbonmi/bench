# Read the tip tree for the deleted-declaration test of bench consumers --changed

Blocked by: none
Writes: internal/consumers/blast.go, internal/consumers/hunks.go, internal/consumers/command.go, internal/consumers/resolve.go, internal/consumers/blast_edges_test.go, CHANGELOG.md
Covers: none

## What to build

Roadmap row FT329 records a false `blast_deleted` row for
`systemtest.systemLandingRaceFixture`. The pair made a body-only edit, and the
function stayed at the same file and line.

The repro is a body-only edit to that function in its `//go:build system` file,
committed through `bench commit`. Then `bench consumers --changed --base <parent>`
reports `blast_deleted[1]` with `systemtest.systemLandingRaceFixture`.

The cause is the set of names that the tip still declares. The deletion test
built that set from the loaded packages only. The loader uses the default build
context, so it omits every file that only a build tag selects. A kept
declaration in such a file therefore reads as deleted.

The deletion test now reads every Go file that sits directly in each affected
directory of the tip tree. It parses these files the same way it parses the base
side. The loaded packages still supply the consumer rows.

## Acceptance

- [x] A body edit in a build-tagged file gives no `blast_deleted` row for the kept declaration.
- [x] A declaration that the same removed run deletes from that file still gives its `blast_deleted` row.
- [x] The existing blast tests stay green.

## Verification

`TestBodyEditInBuildTaggedFileIsNotADeletion` was red before the fix. Its output
held a false `target.Fixture,target/tagged.go,5` row. The test is green after
the fix, and the `internal/consumers` package passes.
