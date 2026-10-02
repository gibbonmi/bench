# 8. Fault the reauthorize unlock with a denied admin directory

Blocked by: 7-fault-the-cleanup-reads-with-real-fixtures.md
Writes: internal/worktree/joins.go, internal/worktree/reauthorize.go, internal/worktree/reauthorize_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: WS51

## What to build

Chunk: SR-C5.

Remove `reauthorizeUnlock` from the joins value. The reauthorize verb calls the real unlock
directly. The unlock-failure case of the rollback test denies writes to the admin
directory at mode 0500, so the real unlock fails. Then run `bench probe` on
`reauthorize.go` with the unlock error branch omitted, and record the red in the
verification note.

The permission fixture calls `capability.Capability` with `capability.Privilege` under the
root user. The converted test keeps its name. The field has no probe fallback. If the
fixture cannot make the test pass, stop the build and name the test, per decision 3.

## Acceptance

- [ ] A reauthorize whose admin directory denies writes rolls back and leaves the retained state unchanged.
- [ ] The unlock probe turns the test red, and the verification note records it.
- [ ] The joins value declares no `reauthorizeUnlock` field.
