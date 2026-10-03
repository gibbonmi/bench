# Read the review record fence from the block reader

Blocked by: 1-add-block-reader.md
Writes: internal/reviewrecord/parse.go, internal/reviewrecord/write.go, internal/reviewrecord/record_test.go, internal/reviewrecord/recordtest/fixture.go
Covers: MB48, MB49, MB50, MB51

## What to build

Make `locate` in `internal/reviewrecord` find the one fence opener whose info
string equals the record name. The span runs from the byte after the opener line
to the first byte of its closer line, by the reader's byte offsets. A named
opener inside another fenced block is a fenced line, so it is not a record.

The duplicate, unterminated, and missing refusals keep their exact message
bytes. `Render` and `fenced` keep one shared locate. Render the fence through the
reader's fence constant in `write.go` and in `recordtest/fixture.go`, and delete
the package's own `fenceMarker` constant.

`internal/reviewrecord` is a crowded directory, so the ticket adds no file.

## Acceptance

- [ ] A four-backtick block that quotes the record fence before the real record reads the real payload with no duplicate error.
- [ ] `TestRenderRefusesAnUnterminatedFence`, `TestRenderRefusesADuplicateFence`, and `TestRenderKeepsProseAroundTheFence` stay green without an edit.
- [ ] `internal/reviewrecord` holds no literal that starts with three backticks.
