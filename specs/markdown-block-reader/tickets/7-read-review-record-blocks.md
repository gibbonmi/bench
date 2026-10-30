# Read the review record fence from the block reader

Blocked by: 1-add-block-reader.md
Writes: internal/reviewrecord/parse.go, internal/reviewrecord/write.go, internal/reviewrecord/assignment.go, internal/reviewrecord/record_test.go, internal/reviewrecord/recordtest/fixture.go
Covers: MB48, MB49, MB50, MB51, MB96, MB107

## What to build

Make `locate` in `internal/reviewrecord` find the one fence opener whose info
string equals the record name. The span runs from the byte after the opener line
to the first byte of its closer line, by the reader's byte offsets. A named
opener inside another fenced block is a fenced line, so it is not a record. The
payload is the raw source bytes of the span. `locate` ignores the reader's
frontmatter fault and reads each line as body, as it does today.

The duplicate, unterminated, and missing refusals keep their exact message
bytes. `Render` and `fenced` keep one shared locate. `assignment.go:96` calls
`locate(data, planFence)`, and it keeps that call.

Delete the package's own `fenceMarker` constant at `parse.go:290`. Its readers at
`write.go:41` and `parse.go:304` and `:309` move to the reader's fence constant.
Render the fence in `recordtest/fixture.go:87` through the same constant.
`recordtest/fixture.go` has 396 lines, so its edit must not grow it.

`internal/reviewrecord` is a crowded directory, so the ticket adds no file.

The record writer encodes the less-than sign of each excerpt in the payload, so
an excerpt that holds an HTML comment opener opens no comment for the block
reader (MB107). Each existing payload field keeps its decoded value.

## Acceptance

- [ ] `TestRenderEncodesACommentOpenerInAnExcerpt` proves that a record whose excerpt holds an unterminated comment opener reads whole after a render (MB107).

- [ ] A four-backtick block that quotes the record fence before the real record reads the real payload with no duplicate error.
- [ ] A record whose first line is an unclosed `---` reads its payload as it does today.
- [ ] `TestRenderRefusesAnUnterminatedFence`, `TestRenderRefusesADuplicateFence`, and `TestRenderKeepsProseAroundTheFence` stay green without an edit.
- [ ] `ReadPlan` of each staged spec and `Read` of each file under `reviews/` give the same payload at the base and at the tip.
- [ ] The non-test files of `internal/reviewrecord` and `internal/reviewrecord/recordtest` hold no string literal that contains a run of three backticks. The sites today are `parse.go:290` and `recordtest/fixture.go:87`.
- [ ] `internal/reviewrecord/recordtest/fixture.go` does not grow, and no other file in `Writes:` grows past 400 lines.
