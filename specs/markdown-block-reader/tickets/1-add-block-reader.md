# Add the block reader and compose the prose check with it

Blocked by: none
Writes: internal/markdown/ (new), internal/prose/prepare.go, internal/prose/prose_test.go
Covers: MB1, MB2, MB3, MB4, MB5, MB6, MB7, MB8, MB9, MB10, MB11, MB12, MB13, MB14, MB15, MB16, MB17, MB18, MB19, MB20, MB21, MB22, MB23

## What to build

Create `internal/markdown`, a leaf package that imports only the standard
library. Its one entry point reads document bytes and returns the classified
lines, the removed comment spans, and at most one fault. The spec's
`The block reader` decision states each line fact and the four rules in their
fixed order. The comment rule is the rejoin rule of `stripCommentsMapped` in
`internal/anchors/locate.go`. The fence rule is the rule of `stripFences` in
`internal/prose/prepare.go`.

The package exports the fence marker and the H2 prefix as constants. Later
tickets render fences and headings through them.

Then make `prose.prepare` read the block reader. It blanks each line that is
not a body line, and it keeps the comment-free text of each body line. It maps
the reader's fault to `KindFrontmatter`, `KindComment`, or `KindFence` at the
fault line. Delete `stripFrontmatter`, `stripComments`, `stripFences`, and
`fenceMarker` from `internal/prose`.

Put the fence edges in one table test, `TestReadFenceEdges`, and the other
block rules in `TestReadBlockRules`. Run `go list` on the new package to confirm
its import edges. Record the old and new `bench gate-prose` output over every
tracked Markdown file in the commit message. The two outputs must be equal.

## Acceptance

- [ ] `TestReadFenceEdges` holds MB1 to MB7, and each row is red against a reader that drops its rule.
- [ ] `TestReadBlockRules` holds MB8 to MB19.
- [ ] `TestFindings` and `TestParagraphsRefusesAnUnterminatedDelimiter` stay green without an edit.
- [ ] `TestFindingsSkipsARejoinedComment` gives no sentence finding for a 30-word sentence inside a rejoined comment.
- [ ] `TestProseMechanicsHoldsOnTheLiveTree` stays green.
- [ ] `internal/prose` holds no fence, comment, or frontmatter literal.
