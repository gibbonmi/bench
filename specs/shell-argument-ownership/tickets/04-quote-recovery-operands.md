# Render the recovery line through the help-row renderer

Blocked by: 01-quote-tree-build-arguments.md, 02-quote-selected-history.md, 03-quote-citation-commands.md, 05-quote-evidence-commands.md, 06-share-rebuild-and-fixture-quotes.md
Writes: internal/worktree/list.go, internal/worktree/list_actions_test.go, internal/worktree/path_identifier_test.go, internal/worktree/shell_arguments_test.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/worktree/unlanded_route_test.go, internal/worktree/identifier_operand_test.go, internal/worktree/parallel_census_test.go, internal/axi/action.go, internal/axi/action_test.go, internal/worktree/missing_tree_recovery_test.go
Covers: S16

## What to build

Render recoveryRoute.line from the same argument list as recoveryRoute.action, through the AXI executable renderer.
The line and its help row then have one spelling, which ticket 07 switches to the surviving quote owner.
This ticket changes no output byte for an accepted path.
Preserve raw refusal structure, line policy, identity, and command authority.
A path that the AXI renderer refuses keeps the current refusal escape policy of the line.

Review chunk: S-E.
Ticket 07 follows this ticket in the same chunk, and the chunk review covers both.
Run this checkpoint with ticket 07 unbuilt.
Preserve all existing assertions and successful behavior outside the approved delta.

## Acceptance

Read a real recovery line with a hostile accepted path.
The shell receives the exact path, the line equals the command of its help row, and the raw line retains its existing reason and control policy.

- [ ] RecoveryRoute.line renders the same command as its help row, so its words and path quote through the surviving owner. (S16).

## Checkpoint verification

- `bench test --package ./internal/worktree --run 'RecoveryRouteArgument|MissingTreeRoute|UnlandedRelease|IdentifierOperand'`
- `bench test --check axi-query-registry`
- `bench test --check subcommand-routing`

## Required omission evidence

Drop the path from the shared argument list, or drop a route word.
The real producer fixture and the line-and-help-row agreement assertion must fail.
A bypass that restores the separate word-join spelling reads the same bytes until the ticket 07 switch, so its red belongs to the S-E omission probe.

Record the exact failing assertion, restored source, and passing result in the review pickup.
A shared-owner test alone cannot prove its caller or transaction outcome.

The co-owned command registries satisfy the existing package binding.
Keep command membership, routes, help grammar, and declared inventories unchanged.
Modify a registry only when the approved caller change requires its exact disposition.
