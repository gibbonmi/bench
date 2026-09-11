# 1. Select worktree path facts

Blocked by: none
Writes: internal/worktree/list.go, internal/worktree/path.go, internal/worktree/list_actions_test.go, internal/worktree/list_selected_test.go (new), internal/usage/worktree.go, cmd/bench/main.go, cmd/bench/worktree_leaves.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/package-core-guard/unrouted-subcommand
Covers: QU1, QU2, QU3, QU9, QU10, QU16, QU17, QU18

## What to build

Add the fixed selected view through the existing worktree list owner.
Reuse the assignment selector without taking active-worktree authority.
Keep bare output unchanged and report every distinct failed operand.
The selected view omits the old unselected inventory.

## Acceptance

- [ ] Resolved rows include identity, path, and state, including non-active state.
- [ ] Ambiguous and absent targets remain beside successful results at exit 1.
- [ ] Aliases of one identity produce one row in first-request order.
- [ ] Bare invocations match the baseline input matrix.
- [ ] Hostile operands cannot erase other results.
- [ ] The selected result names `bench worktree list` as its complete-detail action.
- [ ] The baseline matrix covers list actions, path identifiers, request tokens, landed state, and hostile landed-cleanup callers.
