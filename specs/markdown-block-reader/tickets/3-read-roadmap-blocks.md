# Read roadmap lines from the block reader

Blocked by: 1-add-block-reader.md
Writes: internal/roadmap/tree.go, internal/roadmap/sequence_projection.go, internal/roadmap/commitment_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: MB31, MB32, MB33, MB81, MB99, MB100

## What to build

Replace the fence toggles in `internal/roadmap/tree.go` and the section detector
in `internal/roadmap/sequence_projection.go` with the block reader's classification.
Each one matches its grammar on the raw text of the unfenced lines. `ParseDocument` skips a fenced line for the row grammar and for
the section title. Today it matches a row on every line, fenced or not, so the
skip is new behavior for a backtick fence too.

`sequenceBounds` and `sectionBounds` use the reader's H2 headings, so a fenced
heading neither starts nor ends a section. `parseSequence` consumes those bounds
and reads sequence rows only from unfenced lines. `ProjectSequence` consumes
those bounds and renders its heading through the reader's H2 constant.

`closeDependencies` consumes `sectionBounds` through `Close`. The sequence
projection keeps its outcome grammar, and the dependency closure keeps its table
grammar. The exported signatures stay the same.
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
- [ ] `TestCommitmentSequenceClosure` and `TestCommitmentDependencyClosure` stay green without an edit.
- [ ] `bench roadmap` prints the same output on the live tree at the base and at the tip.
- [ ] `tree.go` and `sequence_projection.go` hold no string literal with a run of three backticks, and no literal `## `.
- [ ] `internal/roadmap/tree.go` does not grow.
