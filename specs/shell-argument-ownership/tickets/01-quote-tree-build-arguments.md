# Quote tree-build arguments through the existing owner

Blocked by: none
Writes: internal/sanitize/shellquote_test.go (new), internal/worktree/tree_target.go, internal/worktree/shell_arguments_test.go (new), internal/treetarget/build_test.go, internal/worktree/target_refusal_test.go, internal/systemtest/tree_target_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: S01, S02, S03, S04, S15

## What to build

Use existing sanitize.ShellQuote in PrintTreeBuildRefusal and add its accepted owner oracles.
Keep the owner algorithm unchanged and always quoted.
Update all three tree-build readers from the accepted eleven-file closure in this same checkpoint.
Keep fixed command tokens, whole-line control escaping, child non-start, and destination authority.

Review chunk: S-A.
This first caller checkpoint consumes the existing sanitize.ShellQuote owner.
Run this checkpoint with successor tickets unbuilt.
Preserve all existing assertions and successful behavior outside the approved delta.

## Acceptance

A safe alpha label prints bench worktree build with a quoted alpha operand.
An embedded quote or metacharacter label recovers the exact original shell argument.
The owner preserves an empty argument, while tree-target validation keeps its existing refusal policy.
Retain every control-byte, argv, environment, cwd, and child-marker assertion in the three existing fixtures.

- [ ] The shared owner renders an empty string as one empty POSIX-shell argument. (S01).
- [ ] An embedded single quote survives a POSIX-shell round trip. (S02).
- [ ] Spaces, glob characters, dollar signs, and shell operators round-trip literally. (S03).
- [ ] A safe nonempty value still renders with surrounding single quotes. (S04).
- [ ] Tree-build refusal commands preserve the exact assignment label. (S15).

## Checkpoint verification

- `bench test --package ./internal/sanitize`
- `bench test --package ./internal/worktree --run 'TreeBuildArgument|TargetRefusal'`
- `bench test --package ./internal/treetarget`
- `bench test --check system`
- `bench test --check axi-query-registry`
- `bench test --check subcommand-routing`

## Required omission evidence

Return bare alpha, split an embedded quote, drop the empty owner argument, or bypass the owner at PrintTreeBuildRefusal.
Exact fixtures and NUL-framed shell argv must detect the named mutation.
BENCH_KIT is supplied by bench test --check system; never run the tagged suite directly.

Record the exact failing assertion, restored source, and passing result in the review pickup.
A shared-owner test alone cannot prove its caller or transaction outcome.

The co-owned command registries satisfy the existing package binding.
Keep command membership, routes, help grammar, and declared inventories unchanged.
Modify a registry only when the approved caller change requires its exact disposition.

## Separate owner and caller evidence

Run owner omission and preservation witnesses with the sanitize package command.
Run PrintTreeBuildRefusal bypass and refusal witnesses with the named worktree package command.
The treetarget and system checks preserve their three exact tree-build output fixtures.
BENCH_KIT is supplied by bench test --check system.
A sanitize-only result cannot satisfy the caller obligations.
