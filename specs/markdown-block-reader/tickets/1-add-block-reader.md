# Add the block reader and compose the prose check with it

Blocked by: none
Writes: internal/markdown/ (new), internal/prose/prepare.go, internal/prose/prose_test.go
Covers: MB1, MB2, MB3, MB4, MB5, MB6, MB7, MB8, MB9, MB10, MB11, MB12, MB13, MB14, MB15, MB16, MB17, MB18, MB19, MB20, MB21, MB22, MB23, MB71, MB72, MB73, MB74, MB75, MB76, MB77, MB78, MB79, MB88, MB89, MB90, MB92

## What to build

Create `internal/markdown`, a leaf package that imports only the standard
library. It exports two entry points over one pipeline: `Read` for a whole
document, and `ReadFragment` for a section or State body with no frontmatter
pass. Each one returns the classified lines, the removed comment spans, and at
most one fault. The spec's
`The block reader` decision states each line fact, the four rules in their fixed
order, and the fault order. Obey these exact predicates:

- Frontmatter opens and closes only on a line that is exactly `---` after the carriage-return strip.
- A fence indent is space (U+0020) and tab (U+0009) only.
- A closer's trailing text does not stop the close.
- An H2 title is trimmed by `strings.TrimSpace`.

The comment rule gives the result of `stripCommentsMapped` in
`internal/anchors/locate.go`, rejoin included. The earliest open opener wins,
and the search for its closer starts after the opener's four bytes. One linear
pass gives that result.

The package exports the fence marker and the H2 prefix as constants. Later
tickets render fences and headings through them.

Then make `prose.prepare` read the block reader. It blanks each line that is
not a body line, and it keeps the comment-free text of each body line. It maps
the reader's fault to `KindFrontmatter`, `KindComment`, or `KindFence` at the
fault line. Delete `stripFrontmatter`, `stripComments`, `stripFences`, and
`fenceMarker` from `internal/prose`, the rune test at `prepare.go:104` included.

Put the fence edges in one table test, `TestReadFenceEdges`, and the other
block rules in `TestReadBlockRules`, both in `internal/markdown/markdown_test.go`.
Run `go list` on the new package to confirm its import edges. Record the old and
new `bench gate-prose` output over every tracked Markdown file in the commit
message. The two outputs must be equal.

## Acceptance

- [ ] `TestReadFenceEdges` holds MB1 to MB7, MB73, and MB78, and each row is red against a reader that drops its rule.
- [ ] `TestReadBlockRules` holds MB8 to MB19, MB71, MB72, MB74 to MB77, MB79, MB88 to MB90, and MB92.
- [ ] `TestFindings` and `TestParagraphsRefusesAnUnterminatedDelimiter` stay green without an edit.
- [ ] `TestFindingsSkipsARejoinedComment` gives no sentence finding for a 30-word sentence inside a rejoined comment.
- [ ] `TestProseMechanicsHoldsOnTheLiveTree` stays green, and `bench gate-prose` gives the same output on every tracked Markdown file at the base and at the tip.
- [ ] `internal/prose` holds no string literal with a run of three backticks or tildes, and no literal `---`, comment delimiter, or `## `.
- [ ] Each new file in `internal/markdown` stays under 400 lines, and `internal/prose/prepare.go` does not grow.
