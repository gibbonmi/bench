# Signal stale hook status

Blocked by: expose-hook-health-record.md
Ownership fence: `internal/status/status.go`, `internal/contract/runtime/runtime_status_test.go`
Integration surfaces: hook-health producer→expose-hook-health-record.md; status five-row budget→`internal/status/status.go` + G3/G5; status runtime contract→`internal/contract/runtime/runtime_status_test.go` + G3/G5
Contracts: the complete hook-health record crosses `internal/adopt/link_hook.go`→`internal/status/status.go`, asserted by G3 against the real exported producer
Closure: G3/status-budget, G5/linked-worktree-single-row

## What to build

Spend one status signal only for a stale managed hook, preserve silence for a baked managed branch, and keep linked worktrees from duplicating the row. Keep these signal-budget cases together: a stale-only cut strands the status runtime contract red when the linked-worktree fixture emits a second pre-push row.

## Acceptance

- [ ] [G3] `bench status` stays silent for a managed baked branch and emits one signal for a stale hook.
- [ ] [G5] A linked worktree renders no duplicate pre-push status row.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| G3/status-budget | emit a status signal for every baked branch | status runtime contract | compare baked-current and stale fixtures, expect silence then exactly one signal |
| G5/linked-worktree-single-row | remove the primary-checkout guard | status runtime contract | run status in a linked worktree, expect no duplicate pre-push row |
