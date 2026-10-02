# 8. Fault the reset and merge moves with real fixtures

Blocked by: 7-fault-the-cleanup-reads-with-real-fixtures.md
Writes: internal/worktree/joins.go, internal/worktree/reset_apply.go, internal/worktree/merge.go, internal/worktree/reauthorize.go, internal/worktree/reset_apply_test.go, internal/worktree/reset_repair_test.go, internal/worktree/merge_test.go, internal/worktree/reauthorize_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: WS45, WS51, WS52, WS53, WS54, WS55, WS56, WS81

## What to build

Chunk: SR-C5.

Remove `reauthorizeUnlock`. The rollback test denies writes to the admin directory at
mode 0500, so the real unlock fails.

Remove `resetMove`. The move-fault tests plant a stale `HEAD.lock` in the admin
directory. `TestResetApplyExitsThreeWhenTheMoveDidNotLand` needs a real move that exits
0 without a move. If no real fixture gives that move, stop the build and name the test
for a reviewer decision, per decision 3.

The merge-reconcile tests depend on a probe. Plant a stale `index.lock` in the target
checkout after the branch moves, then run `bench probe` on `merge.go` with the reconcile
error branch omitted. If the test turns red, remove `mergeReconcile`. If the probe stays
green, keep the field and its two tests unchanged, and record one `bench learning` entry.

The permission fixture calls `capability.Capability` with `capability.Privilege` under
the root user. Each converted test keeps its name. `merge_test.go` is over its line
budget, so the change does not grow it.

This ticket removes the last field, so its review confirms the final field set against
the decision list and the recorded probe learnings.

## Acceptance

- [ ] A reauthorize whose admin directory denies writes rolls back and leaves the retained state unchanged.
- [ ] A stale `HEAD.lock` makes a reset apply exit 3, with the envelope when the checkout is attached and with `preserved=none` when it is detached.
- [ ] A move that exits 0 without a move makes the reset apply exit 3, or the build stops and names the test.
- [ ] A stale `index.lock` makes the merge exit 3 and leaves a merge that the reset apply reconciles, or the learning entry records the failed probe.
- [ ] The joins value declares the 15 named fields and each field whose probe stayed green, and no other field.
