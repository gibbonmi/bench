# Quote direct recovery-route operands

Blocked by: 01-quote-tree-build-arguments.md, 03-quote-citation-commands.md
Writes: internal/worktree/list.go, internal/worktree/list_actions_test.go, internal/worktree/path_identifier_test.go, internal/worktree/shell_arguments_test.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: S16

## What to build

Migrate recoveryRoute.line to the surviving quote owner.
Change only its optional path operand and preserve raw refusal structure, line policy, identity, and command authority.
AXI action spelling remains at its current policy until the separate global switch.

Review chunk: S-B3.
The predecessor supplies the accepted contract named above.
Run this checkpoint with successor tickets unbuilt.
Preserve all existing assertions and successful behavior outside the approved delta.

## Acceptance

Read a real recovery line with a hostile accepted path.
The shell receives the exact path and the raw line retains its existing reason and control policy.

- [ ] RecoveryRoute.line quotes its path through the surviving owner. (S16).

## Checkpoint verification

- `bench test --package ./internal/worktree --run 'RecoveryRouteArgument|List|PathIdentifier'`
- `bench test --check axi-query-registry`
- `bench test --check subcommand-routing`

## Required omission evidence

Bypass the owner only at recoveryRoute.line.
The raw command fixture must fail independently of AXI help tests.

Record the exact failing assertion, restored source, and passing result in the review pickup.
A shared-owner test alone cannot prove its caller or transaction outcome.

The co-owned command registries satisfy the existing package binding.
Keep command membership, routes, help grammar, and declared inventories unchanged.
Modify a registry only when the approved caller change requires its exact disposition.
