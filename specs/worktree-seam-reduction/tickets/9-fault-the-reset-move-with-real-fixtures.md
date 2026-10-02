# 9. Fault the reset move with real fixtures

Blocked by: 8-fault-the-reauthorize-unlock-with-a-denied-admin-directory.md
Writes: internal/worktree/joins.go, internal/worktree/reset_apply.go, internal/worktree/reset_apply_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: WS54, WS55, WS56

## What to build

Chunk: SR-C5.

Remove `resetMove` from the joins value. The reset apply calls the real move directly.

The two move-fault tests plant a stale `HEAD.lock` in the admin directory. Their red
evidence is the fixture: without the `HEAD.lock` plant, the apply exits 0.

`TestResetApplyExitsThreeWhenTheMoveDidNotLand` uses the ignore-rule drift. The move keeps
an ignored file that the checkpoint's rules do not ignore, so the post-move check exits 3.
`TestResetApplyKeepsIgnoredBytesAcrossAnIgnoreRuleChange` shows the shape. Run `bench probe`
on `reset_apply.go` with the post-move check omitted, and record the red.

Each converted test keeps its name. The field has no probe fallback. If a fixture cannot
make its test pass, stop the build and name the test, per decision 3.

## Acceptance

- [ ] A stale `HEAD.lock` makes a reset apply exit 3, with the envelope when the checkout is attached and with `preserved=none` when it is detached.
- [ ] Without the `HEAD.lock` plant, each move-fault apply exits 0, and the verification note records it.
- [ ] A move that keeps an ignore-rule drift file makes the reset apply exit 3 and names the preserved ref.
- [ ] The joins value declares no `resetMove` field.
