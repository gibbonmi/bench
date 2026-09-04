# Restate the Shape age and derive the ignore path

Blocked by: read-a-legacy-document-as-main.md
Writes: internal/handoff/text.go, internal/adopt/init.go, internal/conformance/handoff_single_source_test.go
Covers: HS5

## What to build

Repair for review findings S1, S2, and Sp2. Verify the premise first. The
last paragraph of `ShapeSection` in `internal/handoff/text.go` states that
`bench status` dates the file by its last write. `internal/adopt/init.go`
spells the document path as a literal beside `handoffdoc.DocumentPath`.
Then rewrite that paragraph to the per-section rule the working agreement
states, and read the path from the leaf.

## Acceptance

- [ ] The Shape text states that `bench status` dates each section by the commits past its recorded tip, and `main` by the file write time.
- [ ] `TestHandoffShapeSingleSourcedBites` and `TestHandoffShapeNamesTheSectionGrammar` pass.
- [ ] `internal/adopt/init.go` holds no literal spelling of the document path.
- [ ] Self-probe: restore the file-age sentence, and report which check reds it or that none does.
