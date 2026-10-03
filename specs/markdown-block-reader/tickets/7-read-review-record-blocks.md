# Read the review record fence from the block reader

Blocked by: 1-add-block-reader.md
Writes: internal/reviewrecord/parse.go, internal/reviewrecord/write.go, internal/reviewrecord/record_test.go, internal/reviewrecord/recordtest/fixture.go
Covers: MB48, MB49, MB50, MB51, MB96

## What to build

Make `locate` in `internal/reviewrecord` find the one fence opener whose info
string equals the record name. The span runs from the byte after the opener line
to the first byte of its closer line, by the reader's byte offsets. A named
opener inside another fenced block is a fenced line, so it is not a record. The
payload is the raw source bytes of the span. `locate` ignores the reader's
frontmatter fault and reads each line as body, as it does today.

The duplicate, unterminated, and missing refusals keep their exact message
bytes. `Render` and `fenced` keep one shared locate. Render the fence through the
reader's fence constant in `write.go` and in `recordtest/fixture.go`, and delete
the package's own `fenceMarker` constant.

`internal/reviewrecord` is a crowded directory, so the ticket adds no file.

## Acceptance

- [ ] A four-backtick block that quotes the record fence before the real record reads the real payload with no duplicate error.
- [ ] A record whose first line is an unclosed `---` reads its payload as it does today.
- [ ] `TestRenderRefusesAnUnterminatedFence`, `TestRenderRefusesADuplicateFence`, and `TestRenderKeepsProseAroundTheFence` stay green without an edit.
- [ ] `ReadPlan` of each staged spec and `Read` of each file under `reviews/` give the same payload at the base and at the tip.
- [ ] `internal/reviewrecord` holds no string literal that contains a run of three backticks. The sites today are `parse.go:244` and `recordtest/fixture.go:81`.
- [ ] No file in `Writes:` grows past 400 lines.
