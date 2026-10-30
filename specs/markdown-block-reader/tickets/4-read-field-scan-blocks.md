# Read field-scan and ticket lines from the block reader

Blocked by: 1-add-block-reader.md
Writes: internal/maps/fields.go, internal/maps/tickets_test.go, internal/tickets/tickets.go, internal/tickets/fence_test.go (new), internal/preflight/ownership/metadata.go, internal/preflight/ownership/premises_test.go
Covers: MB34, MB35, MB36, MB37, MB95, MB102

## What to build

Make `FieldScan.Scan` (`internal/maps/fields.go:74`) set `Fenced` from the block
reader. A fence marker line and each fenced line stay `Fenced`, and they match no
field. `Scan` keeps its return shape, the lines and the duplicate diagnostics, so
the decision-map output does not change.

`FieldLine` carries `FenceRole` and `FenceLabel` (`fields.go:45-49`), and `Scan`
sets them at `fields.go:83-84`. The reader-derived `Scan` keeps both fields. An
opener line takes the open role, and its label is the trimmed info string. A
closer takes the close role, and its label is the trimmed text after its marker
run. The planning-facts reader in `internal/preflight/ownership/metadata.go` keys
on these fields, so it keeps its grammar and refuses a nested marker.

`ParseTicket` in `internal/tickets` calls the block reader on its content for the
fence fault, instead of its odd count of backtick lines at `tickets.go:109-116`.
The diagnostic keeps the exact text `unterminated fence`. Delete the
`fenceMarker` constant from `internal/tickets`. `ParseTicket` ignores the
reader's frontmatter fault and reads each line as body, as it does today.

## Acceptance

- [ ] A `Writes:` line inside a tilde fence is marked fenced and matches no field.
- [ ] A ticket with an unterminated `~~~` block gives `unterminated fence`.
- [ ] A four-backtick block that holds a three-backtick line gives no `unterminated fence` diagnostic.
- [ ] A `Covers:` line inside that block sets no `Covers` value.
- [ ] A ticket whose first line is an unclosed `---` gives the same diagnostics as it does today.
- [ ] `TestMetadataGrammarFenceRoles` in `internal/preflight/ownership/premises_test.go` stays green at the ticket tip.
- [ ] `TestParseTicketHostileInput` (`tickets_test.go:307`) and `TestLiftedFieldScanSkipsFencedLines` (`maps_parse_test.go:438`) stay green without an edit.
- [ ] `ParseTicket` gives the same diagnostics for each ticket under `specs/` at the base and at the tip.
- [ ] `internal/maps` and `internal/tickets` hold no string literal with a run of three backticks. The sites today are `fields.go:83` and `tickets.go:39`.
- [ ] No file in `Writes:` grows past 400 lines.
