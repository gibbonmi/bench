# Share rebuild and fixture shell quoting

Blocked by: 01-quote-tree-build-arguments.md, 04-quote-recovery-operands.md
Writes: internal/freshness/freshness_verify.go, internal/freshness/freshness_verify_test.go, internal/diff/identity_test.go, internal/preprelease/preprelease_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: S18, S37

## What to build

Replace freshness.shellQuote with sanitize.ShellQuote and remove that independent algorithm.
Remove the freshness quoteForShell fixture helper and independent diff and preprelease fixture quote helpers.
Use the surviving owner for fixture script arguments while keeping their observed argv expectations independent.
Keep systemtest.shellQuoteJSON because it encodes a separate JSON tool envelope.

Review chunk: S-D1.
The predecessor supplies the accepted contract named above.
Run this checkpoint with successor tickets unbuilt.
Preserve all existing assertions and successful behavior outside the approved delta.

## Acceptance

Execute the real rebuild action against hostile paths and retain its exact cwd and arguments with no output-byte delta.
Migrated fixture scripts preserve their observed argument bytes and count.

- [ ] Freshness rebuild commands use sanitize.ShellQuote for each path. (S18).
- [ ] Fixture shell arguments use sanitize.ShellQuote instead of independent helpers. (S37).

## Checkpoint verification

- `bench test --package ./internal/freshness`
- `bench test --package ./internal/diff`
- `bench test --package ./internal/preprelease`
- `bench test --check axi-query-registry`
- `bench test --check subcommand-routing`

## Required omission evidence

Restore the independent freshness algorithm or a fixture helper.
The review-owned definition census must reject the copy.
A malformed quote in a hostile fixture input must fail the independent observed argv assertion.

Record the exact failing assertion, restored source, and passing result in the review pickup.
A shared-owner test alone cannot prove its caller or transaction outcome.

The co-owned command registries satisfy the existing package binding.
Keep command membership, routes, help grammar, and declared inventories unchanged.
Modify a registry only when the approved caller change requires its exact disposition.

## Verification ownership

Run malformed rebuild-operand evidence in the freshness package.
Run hostile diff fixture argv evidence in the diff package.
Run hostile release fixture argv evidence in the preprelease package.
Independent review grades removal of surviving algorithms through the definition census.
A freshness-only result cannot prove either other fixture package.
