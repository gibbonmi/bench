# 13. Make the serial ceiling exact

Blocked by: 12-refuse-a-read-below-a-census-entry.md
Writes: internal/worktree/serial_ceiling_test.go (new), internal/worktree/parallel_census_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: WS69, WS70

## What to build

Chunk: SR-C6.

Make the serial ceiling check refuse a set below the ceiling as well as above it. Put the
refusal text in a pure renderer, `belowCeilingRefusal(n, c)`:
`the package holds <n> serial tests, below the ceiling of <c>: lower worktreeSerialCeiling to <n> in this change`.
`serialCeilingBreach` calls the renderer. The test asserts a non-empty breach that equals
the renderer's output for a count of 1 and a ceiling of 2.

Lower `worktreeSerialCeiling` to the live serial count, which is below 46.

Move `serialSet`, `serialCeilingBreach`, and the new renderer into a new file,
`serial_ceiling_test.go`, with the new test. The single-read census file of ticket 12 does
not hold them, so this ticket stays apart from the census logic. Keep
`TestSerialSetStaysBelowTheCeiling` in `parallel_census_test.go`, so its row cell stays
valid. Raise `worktreeTestCount` by the number of new top-level tests in this commit.

## Acceptance

- [ ] The live serial set equals `worktreeSerialCeiling`, and the ceiling is below 46.
- [ ] A synthetic serial set of 1 under a ceiling of 2 draws a non-empty breach equal to `belowCeilingRefusal(1, 2)`.
- [ ] `parallel_census_test.go` holds fewer lines than at the ticket's base.
- [ ] The package declares exactly `worktreeTestCount` top-level tests at this commit.
