# 11. Lift each read below a census entry

Blocked by: 10-fault-the-merge-reconcile-with-a-stale-index-lock.md
Writes: internal/worktree/snapshot.go, internal/worktree/subshell.go, internal/worktree/resume.go, internal/worktree/ownership.go, internal/worktree/worktree.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: WS68, WS83

## What to build

Chunk: SR-C6.

Lift each clock read that sits below a census entry. Each entry reads the clock once and
passes the instant down:

- `ClaimRecordedLease` reads the clock and passes the instant to `claimRecordedLease`;
- `CreateCommand` reads the instant and passes it to `createAttributed`;
- `ApplyAutomatic` reads the clock once, and its replan uses that instant;
- `Subshell` reads the clock once and calls `releaseCommandWith` in place of `ReleaseCommand`.

The change in `worktree.go` adds no net line, because the file is over its line budget.
The behavior stays the same, so the existing package tests stay green. This ticket lands
green alone, and ticket 12's live census then confirms that no read stays below an entry.

## Acceptance

- [ ] `claimRecordedLease`, `createAttributed`, `subshellAt`, and the replan in `applyAutomaticWithTerminal` hold no `currentTime` or `Home` reference.
- [ ] `subshellAt` calls no exported verb entry.
- [ ] `worktree.go` holds no more lines than at the ticket's base.
- [ ] The worktree package tests pass.
