# Account for every delegated participant

Blocked by: none
Writes: internal/assessment, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: DI13, DI14, DI15, DI16, DI17, DI18, DI19

## What to build

Add the orchestration role and explicit multi-assignment collection to the existing assessment command.
Keep one normalized selector implementation for the old and new input forms.
Validate event ownership across all batches before the stored record changes.
Preserve the attempts array as the sole participant inventory.

The phase consumer receives one ordinary-work record with complete attempt history and independently partial cost categories.
The existing show and comparison projections include orchestration without a new report command.
Old records remain readable and old singular imports retain their outputs.
The registry files are required co-named consumers; CLI grammar remains unchanged.

## Acceptance

- [ ] Show and comparison include orchestration and every supplied attempt state.
- [ ] One import collects two explicitly named assignments.
- [ ] Both input forms together, foreign assignments, and conflicting event mappings refuse without changing the stored record.
- [ ] Identical repeated native events contribute once.
- [ ] Updates preserve failed, cancelled, replaced, and incomplete attempts.
- [ ] Estimates, actual charges, currencies, measured zero, and unknowns remain distinct.
- [ ] Fixed overlapping intervals retain wall-time union and separate effort time.
- [ ] Differential fixtures preserve singular imports.
