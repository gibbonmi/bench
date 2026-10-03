# Read field-scan and ticket lines from the block reader

Blocked by: 1-add-block-reader.md
Writes: internal/maps/fields.go, internal/maps/tickets_test.go, internal/tickets/tickets.go, internal/tickets/fence_test.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: MB34, MB35, MB36, MB37

## What to build

Make `FieldScan.Scan` in `internal/maps` set `Fenced` from the block reader. A
fence marker line and each fenced line stay `Fenced`, and they match no field.

Make `ParseTicket` in `internal/tickets` read the reader's fence fault instead of
its odd count of backtick lines. The diagnostic keeps the exact text
`unterminated fence`. Delete the `fenceMarker` constant from `internal/tickets`.

## Acceptance

- [ ] A `Writes:` line inside a tilde fence is marked fenced and matches no field.
- [ ] A ticket with an unterminated `~~~` block gives `unterminated fence`.
- [ ] A four-backtick block that holds a three-backtick line gives no `unterminated fence` diagnostic.
- [ ] A `Covers:` line inside that block sets no `Covers` value.
- [ ] `TestParseTicketHostileInput` stays green without an edit.
