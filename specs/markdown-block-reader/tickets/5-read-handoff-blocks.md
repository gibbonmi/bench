# Read handoff lines from the block reader

Blocked by: 1-add-block-reader.md
Writes: internal/handoffdoc/fence.go, internal/handoffdoc/document.go, internal/handoffdoc/document_test.go, internal/handoff/state_file.go, internal/handoff/state_file_test.go
Covers: MB38, MB39, MB40, MB82, MB93

## What to build

Replace `isFence` in `internal/handoffdoc` with the block reader. `OpenFence`
and `UnfencedLines` read a State body through `ReadFragment`, so a leading `---`
opens no frontmatter. `Parse` reads the whole file through `Read`. `OpenFence`
reports the opener line of the reader's fence fault. `UnfencedLines` yields the
raw text of the unfenced lines with trailing whitespace removed, as it does
today. `splitSections` and `parseSection` read the reader's H2 headings and fence
classes.

`readStateFile` in `internal/handoff` refuses the reader's comment fault with a
new fault beside `faultStateFileFence`. The new fault names the opener line. The
handoff file is untracked, so no prose grade refuses an open comment in it.

The `Parse` refusal for an open fence keeps its exact message bytes, and
`OpenFenceRepair` stays. Render each heading through the reader's H2 constant, so
`internal/handoffdoc` keeps no `## ` literal. Each exported signature stays the
same.

## Acceptance

- [ ] A State body with a four-backtick block that holds a three-backtick line and a `## X` line parses as one section.
- [ ] `OpenFence` on `~~~` followed by a three-backtick line reports line 1.
- [ ] `OpenFence` on a State body of `---`, three backticks, `x`, and `---` reports line 2.
- [ ] `UnfencedLines` yields no line from inside a block that a longer run closes.
- [ ] `readStateFile` refuses a State with an unterminated comment and names the opener line.
- [ ] `TestParseKeepsAFencedHeadingInsideState` and `TestOpenFenceNamesTheOpeningLine` stay green without an edit.
- [ ] `handoffdoc.Parse` of the primary checkout's `capture/session-handoff.md` gives the same sections at the base and at the tip.
- [ ] `internal/handoffdoc` holds no string literal with a run of three backticks or tildes, and no literal `## `. The sites today are `fence.go:17`, `document.go:29`, `document.go:263`, `document.go:281`, and `document.go:324`.
- [ ] `internal/handoffdoc/document.go` (400 lines) does not grow, and no other file in `Writes:` grows past 400 lines.
