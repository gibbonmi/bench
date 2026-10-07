# Route the commit exit 3 to the reset plan

Blocked by: 07-route-each-reset-refusal-through-the-registry.md
Writes: internal/refusalroute/registry.go (new), internal/commit/commit.go, internal/commit/landing_test.go, internal/commit/refusal_route_test.go (new), internal/worktree/commit_route_test.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/worktree/parallel_census_test.go
Covers: RR32, RR33, RR34, RR35

## What to build

The commit publishes, and then the checkout fails to reconcile.
The exit 3 record keeps the shape `committed{published_commit=…,path=…,next=…}`.
Its `next=` value is the route of the `commit-published-unreconciled` face: `bench worktree reset --to <published-commit> <checkout>`.
The route no longer names `git restore`.

`<checkout>` is the commit's root path, shell-quoted.
It is the placeholder `<checkout>` when the path is not line-safe.
The reset verb keeps the dirty layer in a recoverable envelope, and then it prints its own `--apply` command.

Create `internal/commit/refusal_route_test.go`.
Later commit tickets add their rows to this file.
Promote the collision 8 restore repro to these rows, with no build tag.
Replace the two `TestPublicationRemainder...` pins of the `git restore` route.

`TestCommitExitThreeRouteReconcilesTheCheckout` in `internal/worktree/commit_route_test.go` drives an exit 3 in an assignment worktree.
It runs the printed reset plan and then its `--apply`.

## Acceptance

- [ ] A commit exit 3 prints `next=` with the value `bench worktree reset --to <published-commit> <checkout>` for its published commit and its root.
- [ ] A commit exit 3 `next=` value does not contain `git restore`.
- [ ] A commit exit 3 on a root path that is not line-safe prints the placeholder `<checkout>`.
- [ ] After a commit exit 3 in an assignment worktree, the printed reset plan and its `--apply` leave `git status --porcelain` empty at the published commit.
