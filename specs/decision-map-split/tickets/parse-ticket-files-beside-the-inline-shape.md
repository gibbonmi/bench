# Parse ticket files and index sections beside the inline shape

Blocked by: none
Writes: internal/maps/tickets.go (new), internal/maps/tickets_test.go (new), internal/maps/maps.go, internal/maps/validation.go, internal/maps/schema.go, tests/canary/workflow-guidance-anchors/decision-map-asset-path, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/decision-map-integrity/, internal/conformance/decision_map_integrity_test.go
Covers: DS1, DS2, DS4, DS5, DS6, DS7, DS8, DS9, DS10, DS11, DS12, DS13, DS14, DS16, DS21, DS33, DS42

## What to build

This is the expand step. The parser learns the split shape and keeps the inline shape. For each map candidate `<dir>/<topic>.md`, discovery lists the direct `.md` children of `<dir>/<topic>/tickets/`. It skips dotfiles and README, and it parses each file as one ticket. The basename without `.md` is the id and must match `[1-9][0-9]*`. The file's `# ` line is the title.

`Blocked by:`, `Type:`, `### Question`, and `### Answer` keep today's field rules and today's `ticket #n:` diagnostics. The ticket parse reuses the shared field scan the way `internal/tickets` does.

The index parse gains `## Notes` and `## Decisions so far`. Both are required for a map that has a tickets folder. Both stay optional for an inline map until the contract ticket lands. A gist line is `- [<name>](<topic>/tickets/<n>.md): <gist>`, resolved by the `<n>` segment.

For a map with a tickets folder, three drift rules run. A resolved ticket with no gist reds. A gist that links an unresolved ticket reds. A gist that links a missing ticket reds. A line in Decisions so far without the link grammar reds. A tickets folder whose index file is absent reds, and an empty tickets folder reds with `missing decision ticket`.

Each new diagnostic gets one fixture under `tests/canary/decision-map-integrity/`: `ticket-basename`, `tickets-absent`, `tickets-empty`, `orphan-tickets-folder`, `notes-missing`, `decisions-so-far-missing`, `gist-missing`, `gist-unresolved`, `gist-missing-file`, and `gist-malformed`. Each new fixture carries its own split map under its `files/` folder, because the family's base fixture stays inline until the rebuild ticket lands. The inventory in `decision_map_integrity_test.go` names each one.

`schema.go` is over its line budget. Any line this ticket adds there moves the same headroom out of that file in the same commit. New code lands in `tickets.go`. The status seam gets one test that writes a split ready map and a split shaping map by hand and reads `appendMaps`.

## Acceptance

- [ ] [T1] A split map with an index and `tickets/1.md` validates with no diagnostic, and the same map with `tickets/07.md` reds naming the file.
- [ ] [T2] A ticket file missing `Type:` reds with `ticket #1: missing Type`, and an unsupported Type reds with today's message.
- [ ] [T3] An index with a tickets folder and no Notes or no Decisions so far reds with the missing-section diagnostic.
- [ ] [T4] A resolved ticket with no gist reds, and a gist to an unresolved ticket reds, each with its coverage-map diagnostic.
- [ ] [T10] A gist to a missing file reds, and a malformed gist line reds, each with its coverage-map diagnostic.
- [ ] [T5] `decisions/beta/tickets/1.md` with no `decisions/beta.md` reds with `beta/tickets: no map index at beta.md`.
- [ ] [T6] A ready split map under `specs/x/decisions/` validates, and `Status: shaping` there reds with `compiled map must be ready`.
- [ ] [T7] `decisions/my map.md` with `decisions/my map/tickets/1.md` validates.
- [ ] [T8] Every existing inline map in the repository still validates, and `bench maps` prints the same rows as before this ticket.
- [ ] [T9] `appendMaps` returns `1 ready map(s)` for one ready split map and `1 unresolved map(s)` for one shaping split map.
- [ ] [T11] Each new fixture bites under the fixture-bite test, and deleting any one reds the inventory test.
