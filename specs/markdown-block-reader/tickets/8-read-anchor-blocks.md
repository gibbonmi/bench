# Read anchor sections and comments from the block reader

Blocked by: 1-add-block-reader.md
Writes: internal/anchors/locate.go, internal/anchors/match.go, internal/anchors/locate_test.go, internal/anchors/registry.go, internal/reviewrecord/parse.go, internal/reviewrecord/record_test.go
Covers: MB52, MB53, MB55, MB56, MB85, MB94, MB108, MB109

## What to build

Make `stripCommentsMapped` (`locate.go:58`) in `internal/anchors` a projection of
the reader's removed comment spans. It keeps its rune-origin mapping, so `Locate`
names the same physical lines. Make `scopeRunesMapped` (`locate.go:81`) read the
reader's fence classes instead of its own backtick toggle. The file read and the
section walk `sectionRunesMapped` both use `Read`. The walk takes its frontmatter
state from the original document, so a comment cannot turn a line into a
frontmatter marker (MB109).

The section opener compares the reader's H2 title with the anchor title, so an
indented `## ` line no longer opens a section. Today the opener builds its
heading from the literal at `locate.go:134` and compares the trimmed line at
`:136`. The section closer at `:137` reads the reader's H2 heading instead.

Delete `commentOpen`, `commentClose`, `fenceMark`, and `headingMark`, the literals
at `locate.go:46-48`. Update the rule comment at `match.go:51` that names the
backtick fence rule.

`internal/anchors` is a crowded directory, so the ticket adds no file. The new
cases go in `locate_test.go`. The test callers keep their calls:
`internal/conformance/docs_workflow_helpers_test.go:236` and `:96`,
`locate_test.go:253-254`, and `match_test.go:69` and `:154`.

The MB-C4 review repair also lands here. The section walk reads the whole
comment-free document through `Read`, as the comment strip does (MB94). A record
closer whose label is the record name refuses as a duplicate (MB108).

## Acceptance

- [ ] `TestLocateReadsTheBlockReader` proves that leading frontmatter with a comment opener or a tilde marker keeps section `S` (MB94).
- [ ] `TestRenderRefusesANestedRecordOpener` proves the duplicate refusal for a second opener right after a record (MB108).

- [ ] A needle after a tilde-fenced `## Other` line inside section `S` locates inside `S`.
- [ ] A tilde-fenced `## S` line before the real one gives a section count of 1.
- [ ] An indented `  ## S` line gives no section `S`.
- [ ] `TestLocateReadsTheBlockReader` proves that a first line of a comment and `---` hides no section (MB109).
- [ ] `TestLocateStripsRejoinedComments` (`locate_test.go:247`) stays green without an edit.
- [ ] `bench test --check docs-currency-workflow` is green, so each live anchor keeps its verdict.
- [ ] `internal/anchors/locate.go` holds no string literal with a run of three backticks, and no literal comment delimiter or `## `. The sites today are `locate.go:46`, `:47`, `:48`, and `:134`.
- [ ] `internal/anchors/locate.go` (363 lines) does not grow, and no other file in `Writes:` grows past 400 lines.
