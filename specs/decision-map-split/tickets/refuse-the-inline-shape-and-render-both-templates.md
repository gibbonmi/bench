# Refuse the inline shape and render both templates

Blocked by: split-every-map-with-the-migration-program.md, rebuild-the-integrity-fixture-family-on-the-split-shape.md
Writes: internal/maps/, scripts/, internal/status/status_signals_test.go, internal/status/status_producible_test.go, internal/status/status_command_test.go, internal/conformance/decision_map_integrity_test.go, tests/canary/decision-map-integrity/, tests/canary/workflow-guidance-anchors/decision-map-asset-path, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: DS3, DS17, DS18, DS19, DS20, DS40

## What to build

This is the contract step. An `## #n:` heading in a map index reds with `inline ticket #n: move it to <topic>/tickets/<n>.md`. Notes and Decisions so far become required for every map, and the three drift rules run for every map. The inline ticket parse and the migration program are deleted together.

`DecisionMapTemplate` renders the index skeleton with Notes and Decisions so far and no inline ticket. A new `DecisionTicketTemplate` renders one ticket skeleton, and `bench maps --ticket-template` prints it. The two template flags and `--count` are mutually exclusive. The asset-rule sentence becomes `A map-owned asset stays in the map's assets folder, decisions/<topic>/assets/.` The anchor needle that pins `decisions/assets/` in `schema.go`, its mutation-table row, and the `decision-map-asset-path` fixture change with it.

Every test that writes a template map to one file writes the index and one ticket from the two templates instead. Add the `inline-ticket-heading` fixture and its inventory entry.

## Acceptance

- [ ] [C1] A map index carrying `## #2: Old` reds with `inline ticket #2: move it to alpha/tickets/2.md`, and the fixture bites.
- [ ] [C2] `bench maps --template` prints `## Notes` and `## Decisions so far` and no `## #1:` line.
- [ ] [C3] `bench maps --ticket-template` prints a skeleton with `# `, `Blocked by: none`, `Type: Research`, `### Question`, and `### Answer`.
- [ ] [C4] `bench maps --template --ticket-template` exits 2 with the mutual-exclusion help line.
- [ ] [C5] The index skeleton holds the new asset-rule sentence, and the `decision-map-asset-path` fixture bites on it.
- [ ] [C6] `scripts/split-decision-maps/` does not exist, and no function in `internal/maps` parses an `## #` heading.
- [ ] [C7] The status, conformance, and command-registry tests pass with two-file maps.
