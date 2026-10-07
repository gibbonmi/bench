# Route each reset refusal through the registry

Blocked by: 06-give-the-red-source-fold-an-exit.md
Writes: internal/refusalroute/registry.go (new), internal/worktree/reset.go, internal/worktree/reset_apply.go, internal/worktree/reset_restore.go, internal/worktree/path.go, internal/worktree/merge_route_test.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
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
The reset plan's `--apply` pointer renders through the registry's field function, and its text does not change.

Extend `TestMergeFacesFollowTheirRoutes` to walk the reset faces too.
Each reset fixture produces its face, follows the printed route, and reruns the reset.
After this ticket, the walk covers every merge face and every reset face.

## Acceptance

- [ ] Each merge face and each reset face has exactly one producing fixture that follows its route out of the face.
- [ ] Each reset refusal prints the same route text that the reset verb printed before this ticket.
- [ ] A reset plan still prints `next=bench worktree reset --to ` with its `--apply` fingerprint.
