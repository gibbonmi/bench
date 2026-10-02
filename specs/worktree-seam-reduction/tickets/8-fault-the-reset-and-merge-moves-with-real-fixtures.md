# 8. Fault the reset and merge moves with real fixtures

Blocked by: 7-fault-the-cleanup-reads-with-real-fixtures.md
Writes: internal/worktree/joins.go, internal/worktree/reset_apply.go, internal/worktree/merge.go, internal/worktree/reauthorize.go, internal/worktree/reset_apply_test.go, internal/worktree/reset_repair_test.go, internal/worktree/merge_test.go, internal/worktree/reauthorize_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: WS45, WS51, WS52, WS53, WS54, WS55, WS56, WS81

## What to build

Chunk: SR-C5.

Remove `reauthorizeUnlock`. The rollback test denies writes to the admin directory at
mode 0500, so the real unlock fails.

Remove `resetMove`. The move-fault tests plant a stale `HEAD.lock` in the admin
directory. `TestResetApplyExitsThreeWhenTheMoveDidNotLand` uses the ignore-rule drift: the
move keeps an ignored file that the checkpoint's rules do not ignore, so the post-move
check exits 3. `TestResetApplyKeepsIgnoredBytesAcrossAnIgnoreRuleChange` shows the shape.

The merge-reconcile tests depend on a probe. The target's declared lane script, from
ticket 3, plants a stale `index.lock` in the target's admin directory before the
publication. Then run `bench probe` on `merge.go` with the reconcile error branch omitted.
If the test turns red, remove `mergeReconcile`. A probe fails when it stays green, or when the fixture cannot make the converted test pass. After a failed
probe, keep the field and its two tests unchanged, and record one `bench learning` entry.

The permission fixture calls `capability.Capability` with `capability.Privilege` under
the root user. Each converted test keeps its name. `merge_test.go` is over its line
budget, so the change does not grow it.

For a field without a probe, a fixture that cannot make its test pass stops the build.
The build names that test, per decision 3.

This ticket removes the last field, so its review confirms the final field set against
the decision list and the recorded probe learnings.

## Acceptance

- [ ] A reauthorize whose admin directory denies writes rolls back and leaves the retained state unchanged.
- [ ] A stale `HEAD.lock` makes a reset apply exit 3, with the envelope when the checkout is attached and with `preserved=none` when it is detached.
- [ ] A move that keeps an ignore-rule drift file makes the reset apply exit 3 and names the preserved ref.
- [ ] A stale `index.lock` makes the merge exit 3 and leaves a merge that the reset apply reconciles, or the learning entry records the failed probe.
- [ ] The joins value declares the 15 named fields and each field whose probe failed, and no other field.
