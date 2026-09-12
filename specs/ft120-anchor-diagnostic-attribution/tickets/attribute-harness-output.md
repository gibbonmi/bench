# Attribute anchor harness diagnostics to their file

Blocked by: none
Writes: internal/anchors, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: none

## What to build

Limit a failed anchorHarness assertion's displayed diagnostics to its subject file.
Preserve the complete diagnostic slice for membership and cross-talk assertions.
Use registry ownership when a diagnostic does not contain its file path.
Move the shared harness out of the oversized registry test file.

The coordinator reconciles the roadmap after this fix lands.

## Acceptance

- [ ] A failed assertion displays its own file's diagnostics.
- [ ] The assertion omits diagnostics owned only by other files.
- [ ] A custom diagnostic prefix retains correct file attribution.
- [ ] Membership and cross-talk failures still fail the test.
