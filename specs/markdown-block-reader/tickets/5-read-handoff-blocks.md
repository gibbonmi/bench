# Read handoff lines from the block reader

Blocked by: 1-add-block-reader.md
Writes: internal/handoffdoc/fence.go, internal/handoffdoc/document.go, internal/handoffdoc/document_test.go
Covers: MB38, MB39, MB40

## What to build

Replace `isFence` in `internal/handoffdoc` with the block reader. `OpenFence`
reports the opener line of the reader's fence fault. `UnfencedLines` yields the
unfenced lines with trailing whitespace removed, as it does today.
`splitSections` and `parseSection` read the reader's H2 headings and fence
classes.

The `Parse` refusal for an open fence keeps its exact message bytes, and
`OpenFenceRepair` stays. Render each heading through the reader's H2 constant, so
`internal/handoffdoc` keeps no `## ` literal. Each exported signature stays the
same.

## Acceptance

- [ ] A State body with a four-backtick block that holds a three-backtick line and a `## X` line parses as one section.
- [ ] `OpenFence` on `~~~` followed by a three-backtick line reports line 1.
- [ ] `UnfencedLines` yields no line from inside a block that a longer run closes.
- [ ] `TestParseKeepsAFencedHeadingInsideState` and `TestOpenFenceNamesTheOpeningLine` stay green without an edit.
