# Retire the step-scoped anchor kind

Blocked by: 5-prune-hand-pins.md
Writes: internal/anchors, cmd/bench/anchors_command.go, cmd/bench/anchor_help_test.go, cmd/bench/anchors_dir_test.go, internal/gate/lane_select_test.go, internal/gate/lane_test.go, CHANGELOG.md, tests/canary/workflow-guidance-anchors, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: AP33, AP34, AP35, AP36, AP37, AP38

## What to build

Remove the nine step-scoped rows from the chunk-chain and calibration registry files, with their test rules. Remove the fixture directories that `assets/fixture-census.md` lists for this ticket. Then retire the step-scoped kind, the step field, and the step narrowing from the anchors package, with each caller and each test.

The kind tests walk a hard-coded kind list, so remove the step kind from that list in the same commit. The case-folded emphasis row and its help test stay together. The `bench anchors` table loses its `step` column, and `CHANGELOG.md` records that change. Update each test that pins the old header or the empty-state marker.

`bench anchors` must still print a definitive empty state for a file with no rows. A non-Markdown file with no rows must still take the unknown lane class. No touched file may grow, and no new file may go into `internal/anchors`.

## Acceptance

- [ ] The anchors package declares no step-scoped kind, no step field, and no step narrowing, and the ticket record shows each step symbol gone.
- [ ] Each surviving match kind has at least one registry row.
- [ ] `bench anchors` on a file with no rows prints a zero-row table at exit 0, with no `step` column.
- [ ] `bench test --package ./internal/anchors`, `bench test --package ./cmd/bench`, `bench test --package ./internal/gate`, and `bench test --check docs-currency-workflow` pass.
