# 3. Retire response spills with their assignment

Blocked by: 1-bound-exec-output.md
Writes: internal/responsebound/ (new), internal/worktree/lifecycle.go, internal/worktree/response_spill_test.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: BO16, BO17

## What to build

Add a spill drop to the response owner package. The retirement path in `internal/worktree/lifecycle.go` calls it beside `census.Drop`, so the retirement of an assignment removes that assignment's spill directory. The drop refuses an identifier that is not an assignment id, as `census.Drop` does.

After each new spill in the `primary` or `none` scope, the owner removes the oldest files until 64 remain.

## Acceptance

- [ ] The retirement of an assignment removes its spill directory.
- [ ] The 65th spill in the `primary` scope leaves the newest 64 files.
