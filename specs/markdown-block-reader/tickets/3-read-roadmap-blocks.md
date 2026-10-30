# Read roadmap lines from the block reader

Blocked by: 1-add-block-reader.md
Writes: internal/roadmap/tree.go, internal/roadmap/sequence_projection.go, internal/roadmap/commitment_test.go, internal/roadmap/closure.go, internal/coverage/rows.go, internal/coverage/coverage_test.go, internal/spec/header_fence_test.go, internal/tickets/fence_test.go, internal/maps/fields.go, internal/preflight/ownership/premises_test.go
Covers: MB31, MB32, MB33, MB81, MB99, MB100, MB103

## What to build

Replace the fence toggles in `internal/roadmap/tree.go` and the section detector
in `internal/roadmap/sequence_projection.go` with the block reader's classification.
The toggles are `ParseDocument` at `tree.go:115`, `parseSequence` at `:276`, and
`rowNextDiagnostics` at `:344`. Each one matches its grammar on the raw text of
the unfenced lines. `ParseDocument` skips a fenced line for the row grammar and
for the section title, which it reads at `tree.go:118` and `:137`. Today it
matches a row on every line, fenced or not, so the skip is new behavior for a
backtick fence too.

`sequenceBounds` (`sequence_projection.go:10`) and `sectionBounds`
(`sequence_projection.go:17`, toggle at `:22`) use the reader's H2 headings, so a
fenced heading neither starts nor ends a section. `parseSequence` consumes those
bounds and reads sequence rows only from unfenced lines. `ProjectSequence`
consumes those bounds and renders its heading through the reader's H2 constant.

`closeDependencies` (`closure.go:47`) consumes `sectionBounds` at `closure.go:48`,
and `Close` calls it. The sequence
projection keeps its outcome grammar, and the dependency closure keeps its table
grammar. The exported signatures stay the same. `Close` reads a row heading only
on a body line, through the same predicate as `ParseDocument`, so a fenced row
example neither counts as a row nor leaves the index (MB103).
`rowNextDiagnostics` reads only unfenced lines for the `Next:` marker. It keeps
its own column-zero and separator grammar.

`internal/roadmap` is a crowded directory, so the ticket adds no file. The new
cases go in `commitment_test.go`, beside its sequence and dependency closure
tests. Do not recreate the deleted helper file or grow `tree_test.go`.

## Acceptance

- [ ] A tilde-fenced row example and a backtick-fenced row example each give no row and no malformed-row failure.
- [ ] A tilde-fenced `## Recommended sequence` line opens no sequence section.
- [ ] A tilde-fenced `Next:` line after a real one gives no duplicate diagnostic.
- [ ] `TestProjectSequenceSkipsFencedHeadings` proves that `ProjectSequence` replaces only the live sequence after a tilde-fenced sequence heading.
- [ ] `TestCommitmentDependencyClosureSkipsFencedHeadings` proves that `Close` removes FT1 from only the live dependency section after a tilde-fenced dependency heading.
- [ ] `TestCloseSkipsAFencedRowHeading` proves that `Close` removes only the live FT1 row beside a fenced example.
- [ ] `TestCommitmentSequenceClosure` and `TestCommitmentDependencyClosure` stay green without an edit.
- [ ] `bench roadmap` prints the same output on the live tree at the base and at the tip.
- [ ] `tree.go` and `sequence_projection.go` hold no string literal with a run of three backticks, and no literal `## `. The sites today are `tree.go:115`, `:118`, `:137`, `:276`, and `:344`, and `sequence_projection.go:22`, `:29`, and `:36`.
- [ ] `internal/roadmap/tree.go` does not grow.
