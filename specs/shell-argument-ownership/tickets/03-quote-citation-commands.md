# Quote citation command arguments

Blocked by: 01-quote-tree-build-arguments.md
Writes: internal/consumers/citation.go, internal/consumers/citation_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: S14

## What to build

Use sanitize.ShellQuote for each citation command argv value.
Keep supplied operands, citation identity, hashing, line permission, and replay behavior.
Preserve consumer query output and its existing authority assertions.

Review chunk: S-B2.
The predecessor supplies the accepted contract named above.
Run this checkpoint with successor tickets unbuilt.
Preserve all existing assertions and successful behavior outside the approved delta.

## Acceptance

A citation with spaces and embedded quotes recovers each original argument through the real printed command.
Its independent fixture pins the approved safe-token quote delta without changing citation identity.

- [ ] Citation commands preserve the supplied argv through the surviving owner. (S14).

## Checkpoint verification

- `bench test --package ./internal/consumers`
- `bench test --check axi-query-registry`
- `bench test --check subcommand-routing`

## Required omission evidence

Return the former bare spelling at citation.cmd or omit a supplied flag.
The exact fixture or recovered argv must turn red.

Record the exact failing assertion, restored source, and passing result in the review pickup.
A shared-owner test alone cannot prove its caller or transaction outcome.

The co-owned command registries satisfy the existing package binding.
Keep command membership, routes, help grammar, and declared inventories unchanged.
Modify a registry only when the approved caller change requires its exact disposition.
