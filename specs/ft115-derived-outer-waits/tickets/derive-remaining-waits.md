# Derive the remaining outer test waits

Blocked by: none
Writes: internal/guards, internal/git/admin_readers_test.go, internal/git/worktree_admin_enum_test.go, internal/status/status_signals_test.go, internal/status/status_spec_count_test.go (new), internal/census, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: none

## What to build

Replace the six selected one-second outer waits with bounds.TestDeadline.
Use the actual inner Git bound for the two Git waits.
Use zero where the operation has no inner timed bound.
Report an exhausted window through bounds.TestTimeoutVerdict.
Preserve every result, cleanup, and refusal assertion.
Relocate affected tests when an oversized file needs headroom.

The coordinator reconciles the roadmap after this fix lands.

## Acceptance

- [ ] All six selected waits derive their windows from the bounds owner.
- [ ] Each Git outer window exceeds its active inner bound.
- [ ] Timeout failures name the exhausted window.
- [ ] Guard cancellation still waits for worker cleanup.
- [ ] Status and census still refuse FIFO inputs without opening them.
