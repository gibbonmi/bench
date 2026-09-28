# 1. Select worktree path facts

Blocked by: none
Writes: internal/worktree/list.go, internal/worktree/path.go, internal/worktree/list_actions_test.go, internal/worktree/list_selected.go (new), internal/worktree/list_selected_test.go (new), internal/worktree/testdata/selected-defaults.json (new), internal/worktree/unlanded_route_test.go, internal/usage/worktree.go, cmd/bench/main.go, cmd/bench/worktree_leaves.go, cmd/bench/selected_queries_test.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/package-core-guard/unrouted-subcommand, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, CHANGELOG.md, tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns, tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary
Covers: QU1, QU2, QU3, QU9, QU10, QU16, QU17, QU18, QU26

## What to build

Add the fixed selected view through the existing worktree list owner.
Reuse the assignment selector without taking active-worktree authority.
Keep bare output unchanged and report every distinct failed operand.
The selected view omits the old unselected inventory.
Start the selected route only when `--view` or `--target` is present, before the bare list grammar parses.

Port from the stream reference source, and verify each ported line against the current tree:

- `194b7dba:internal/worktree/list_selected.go`
- `194b7dba:internal/worktree/list_selected_test.go`
- `194b7dba:internal/worktree/testdata/selected-defaults.json`
- `194b7dba:internal/worktree/list.go`
- `194b7dba:internal/usage/worktree.go`
- `194b7dba:cmd/bench/selected_queries_test.go`
- `194b7dba:CHANGELOG.md`

`main.go` now routes help through `Kind: commandHelp` (`cmd/bench/main.go:111`), so do not port the stream's `helpCommand` move.

The stream tip includes the accepted QU-C1 review repair of stream ticket 5.
Keep that repair: derive the baseline expectation through the TOON owner, and deduplicate resolved identities before the stored-path refusal.

## Acceptance

- [ ] Resolved rows include identity, path, and state, including recovered state with valid recovery metadata.
- [ ] Ambiguous and absent targets remain beside successful results at exit 1.
- [ ] Aliases of one identity produce one row in first-request order, also when the stored path is unrepresentable.
- [ ] Bare invocations match the baseline input matrix, and the pre-disclosure argument fixture stays unchanged.
- [ ] Hostile operands cannot erase other results, and each unsafe operand keeps its first request ordinal.
- [ ] Selected stored paths keep tab, newline, and return without the loss of other results.
- [ ] The selected result names `bench worktree list` as its complete-detail action.
- [ ] The baseline matrix covers list actions, path identifiers, request tokens, landed state, unlanded routes, and hostile landed-cleanup callers.
- [ ] Three resolved targets print exactly 6 lines with no spill line through the bounded dispatcher.

## Headroom

The new selected-view file owns the projection and its grammar.
The existing list entry calls that owner and retains the bare path.
The new command test file owns selected-route integration checks.
Existing over-budget files receive only required routing or expectation edits.
