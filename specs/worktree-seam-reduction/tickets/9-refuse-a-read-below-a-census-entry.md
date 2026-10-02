# 9. Refuse a read below a census entry

Blocked by: 8-fault-the-reset-and-merge-moves-with-real-fixtures.md
Writes: internal/worktree/single_read_census_test.go (new), internal/worktree/parallel_census_test.go, internal/worktree/snapshot.go, internal/worktree/subshell.go, internal/worktree/resume.go, internal/worktree/ownership.go, internal/worktree/effects.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: WS57, WS58, WS59, WS60, WS61, WS62, WS63, WS64, WS65, WS66, WS67, WS68, WS69, WS70

## What to build

Chunk: SR-C6.

Add the single-read census in a new test file. A census entry is a function or a method
declaration whose name is exported. The census derives its read set from three sources:

- each function that `effects.go` declares;
- each package-qualified function that an `effects.go` body calls;
- each exported gate function `X` whose name plus `AtKit` names another exported gate function.

A read's kind is the package-qualified call that the census reaches from it. The census
accepts a read only in an entry's own body, outside each function literal and loop body.
It accepts only the first read of each kind there. It reports each other read with one of the
three messages in the spec's Implementation decisions.

Lift each read that the live census reports. `ClaimRecordedLease` reads the clock and
passes the instant down. `ApplyAutomatic` reads the clock once. `Subshell` reads the clock
once and calls `releaseCommandWith` in place of `ReleaseCommand`.

Lower `worktreeSerialCeiling` to the live serial count, and make the ceiling check refuse
a set below the ceiling with
`the package holds <n> serial tests, below the ceiling of <c>: lower worktreeSerialCeiling to <n> in this change`.
Move `serialSet` and `serialCeilingBreach` into the new file, so that
`parallel_census_test.go` does not grow. Keep `TestSerialSetStaysBelowTheCeiling` in
`parallel_census_test.go`.

## Acceptance

- [ ] Each synthetic census case in the spec's rows WS57 to WS67 draws its exact report or no report.
- [ ] The census over the live package reports nothing.
- [ ] The live serial set equals `worktreeSerialCeiling`, and the ceiling is below 46.
- [ ] A synthetic serial set below its ceiling draws the below-ceiling refusal.
- [ ] `parallel_census_test.go` holds fewer lines than at the ticket's base.
