# 1. Inspect one pinned harness record

Blocked by: none
Writes: internal/harnesstranscript (new), internal/harnesses/command.go, internal/harnesses/command_test.go, internal/harnesses/observed_test.go (new), cmd/bench/main.go, .agents/skills/bench-craft-cli/SKILL.md, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: ME1, ME2, ME3, ME4, ME5, ME6, ME7, ME8, ME9, ME10, ME11, ME12, ME17, ME18, ME19

## What to build

Add the opt-in record view from input through typed observations to TOON output and command tests.
Keep the existing compiled views unchanged.
The new reader owns its metric inventory and source mapping.
Fixtures preserve the relevant source shapes without copying private transcript contents.
The command exposes all reliable dimensions and names every unavailable dimension.

## Acceptance

- [ ] A pinned multibyte fixture yields the known text-byte total through the public command.
- [ ] Outer calls, nested calls, duplicate completions, and unmatched calls retain distinct results.
- [ ] Absent, empty, malformed, and unsupported inputs produce their specified availability and exit code.
- [ ] Native cumulative snapshots do not inflate the token total.
- [ ] A differential run preserves both compiled views over the registered harness inventory.
- [ ] Hostile record contents cannot create a sentinel file or execute a command.
