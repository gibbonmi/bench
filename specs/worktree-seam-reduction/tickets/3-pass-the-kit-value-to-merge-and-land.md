# 3. Pass the kit value to merge and land

Blocked by: 2-carry-an-ambient-value-below-each-verb-entry.md
Writes: internal/worktree/joins.go, internal/worktree/merge.go, internal/worktree/land.go, internal/worktree/merge_test.go, internal/worktree/merge_caller_root_test.go, internal/worktree/delegated_integration_test.go, internal/worktree/land_effects_test.go, internal/worktree/verb_fixture_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: WS17, WS18, WS19, WS20, WS21, WS22, WS23, WS24, WS25

## What to build

Chunk: SR-C3.

Remove the `mergeLane` and `kitSourceCheckout` fields from the joins value. The merge
verb resolves its lane with `gate.LaneForCommitAtKit(target, kit)`, and the land verb
decides its install notice with `gate.KitSourceCheckoutAtKit(root, kit)`. Both take the
kit from the ambient value.

Convert each merge test that replaced `mergeLane`. `mergeFixture` and the delegated
journey commit a phase manifest whose lane appends to the tally. Each merge call passes a
kit value apart from the target. A test that needs another lane check commits
its own manifest lane in the target before the merge. Convert the two land tests that
used `kitCheckoutJoins` to a kit value. The value names the destination for the install
notice, and it names another directory for the repair route.

Each converted test keeps its name. `merge_test.go` is over its line budget, so the
conversion does not grow it.

## Acceptance

- [ ] A merge with a manifest tally lane and a kit apart from the target appends one byte to the tally.
- [ ] A declared lane check that exits 1 refuses the merge and leaves the branch tip unchanged.
- [ ] The fast-forward, edited-checkout, prose-placeholder, caller-root, and delegated-journey tests pass through the manifest lane.
- [ ] A broker-changing landing with the kit at the destination names the install step.
- [ ] A landing with a kit apart from the destination names the installed repair route.
- [ ] The joins value declares no `mergeLane` and no `kitSourceCheckout` field.
