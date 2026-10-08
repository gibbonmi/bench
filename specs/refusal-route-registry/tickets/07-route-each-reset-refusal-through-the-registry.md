# Route each reset refusal through the registry

Blocked by: 06-give-the-red-source-fold-an-exit.md
Writes: internal/refusalroute/registry.go (new), internal/worktree/reset.go, internal/worktree/reset_apply.go, internal/worktree/reset_restore.go, internal/worktree/path.go, internal/worktree/merge_route_test.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/worktree/parallel_census_test.go, internal/worktree/merge_refusal.go, internal/worktree/refusal_route_follow_test.go, internal/refusalroute/faces_reset.go (new), internal/worktree/reset_apply_test.go, internal/worktree/reset_restore_refusal_test.go, internal/worktree/clean_landed.go, internal/worktree/clean_landed_test.go, internal/worktree/worktree.go, internal/worktree/lifecycle.go, internal/worktree/worktree_test.go, internal/worktree/lifecycle_test.go, internal/worktree/missing_tree_recovery_test.go (new), internal/worktree/ownership.go
Covers: RR29

## What to build

Each reset refusal prints the route of its face in the `next=` field of its `refused{...}` record.
Declare these reset faces in the registry:

- `reset-checkout-conflicted`
- `reset-plan-stale`
- `reset-tree-missing`
- `reset-handback`

Each face keeps the current route of the reset verb unchanged.
The `reset-tree-missing` face serves the missing-tree refusal.
Its route clears its cause: `bench worktree clean --landed` retires a landed assignment whose tree is missing, and `bench worktree release` releases an unlanded one.
A reset refusal that has no face of its own hands back through `reset-handback`.
The reset plan's `--apply` pointer renders through the registry's field function, and its text does not change.

Extend `TestMergeFacesFollowTheirRoutes` to walk the reset faces too.
Each reset fixture produces its face, follows the printed route, and reruns the reset.
After this ticket, the walk covers every merge face and every reset face.

## Acceptance

- [ ] Each merge face and each reset face has exactly one producing fixture that follows its route out of the face.
- [ ] Each reset refusal prints the same route text that the reset verb printed before this ticket.
- [ ] A reset plan still prints `next=bench worktree reset --to ` with its `--apply` fingerprint.
- [ ] For a landed assignment whose tree is missing, the printed `bench worktree clean --landed` route and its apply retire the assignment, and the refusal no longer prints.
- [ ] For an unlanded assignment whose tree is missing, the printed `bench worktree release` route releases the assignment, and the refusal no longer prints.
