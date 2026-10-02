# 10. Fault the merge reconcile with a stale index lock

Blocked by: 9-fault-the-reset-move-with-real-fixtures.md
Writes: internal/worktree/joins.go, internal/worktree/merge.go, internal/worktree/merge_test.go, internal/worktree/reset_repair_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: WS45, WS52, WS53, WS81

## What to build

Chunk: SR-C5.

The two merge-reconcile tests depend on a probe. The target's declared lane script, from
ticket 3, plants a stale `index.lock` in the target's admin directory before the
publication. The reconcile then fails after the branch moves. Run `bench probe` on
`merge.go` with the reconcile error branch omitted, scoped to
`TestMergeExitsThreeWhenTheReconcileFails`.

If the test turns red, remove `mergeReconcile`, and the merge calls the real reconcile
directly. A probe fails when it stays green, or when the fixture cannot make the converted
test pass. After a failed probe, keep the field and its two tests unchanged, and record one
`bench learning` entry that names the field and the probe.

Each converted test keeps its name. `merge_test.go` is over its line budget, so the change
does not grow it.

This ticket removes the last field, so its review confirms the final field set. The review
compares `joins.go` with the decision list and the recorded probe learnings of tickets 4,
5, 7, and this ticket.

## Acceptance

- [ ] A stale `index.lock` makes the merge exit 3 and leaves a merge that the reset apply reconciles, or the learning entry records the failed probe.
- [ ] The joins value declares the 15 named fields and each field whose probe failed, and no other field.
- [ ] Each failed probe has one learning entry, and no converted test was skipped.
