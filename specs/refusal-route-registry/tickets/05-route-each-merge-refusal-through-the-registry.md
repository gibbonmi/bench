# Route each merge refusal through the registry

Blocked by: 03-prove-each-agent-route-passes-the-wired-guards.md, 04-render-the-recovery-matrix-from-the-registry.md
Writes: internal/refusalroute/registry.go (new), internal/worktree/merge.go, internal/worktree/land.go, internal/worktree/land_identity.go, internal/landing/landing.go, internal/landing/merge.go, internal/landing/landing_reviewed_test.go, internal/worktree/merge_route_test.go (new), internal/worktree/refusal_route_test.go (new), internal/worktree/merge_test.go, internal/worktree/merge_from_sha_test.go, internal/worktree/land_bench_home_test.go, internal/commit/dry_run_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/worktree/parallel_census_test.go, internal/worktree/refusal_route_follow_test.go, internal/worktree/merge_refusal.go (new), internal/worktree/land_refusal_fixture_test.go (new), internal/worktree/land_refusal.go
Covers: RR27, RR28, RR30, RR31, RR67

## What to build

Each merge refusal prints the route of its face in the `next=` field of its `refused{...}` record.
Declare these merge faces in the registry with the authority and the route of the spec's face inventory table:

- `merge-target-not-clean`
- `merge-sibling-not-clean`
- `merge-conflict`
- `merge-infrastructure`
- `merge-published-unreconciled`
- `merge-handback`

Declare these land faces for the landing's own authorization red, which `LandReviewed` raises:

- `land-red`, for a red kind
- `land-infrastructure`, for an infrastructure kind

The merge conflict prints the same reviewer route as the landing conflict.
A dirty sibling routes to `bench commit --in <sibling-label> -m <msg> -- <path>...`, then the re-run.
It no longer prints the `bench worktree exec` form.
The merge exit 3 record prints its reset route through the `merge-published-unreconciled` face.

`landing.refusalMessage` keeps its prefix, its kind, and its explanation, and it drops its inline action.
The merge's single infrastructure retry matches the typed kind, not the old sentence.
Change each test that reads the old inline action, so that it reads the new sentence.
The landing maps the typed kind of its authorization red to `land-red` or `land-infrastructure`, so no landing red loses its route.
Each other `landReviewed` error that no face claims hands back through `land-handback`, so no landing refusal prints without a route.
Add `TestLandingRedRouteNamesTheRepair` in `internal/worktree/refusal_route_test.go`.

Add `TestMergeFacesFollowTheirRoutes` in `internal/worktree/merge_route_test.go`.
It walks the merge faces of the registry and drives one producing fixture for each face.
Each fixture follows its printed route and reruns the merge.
The `merge-published-unreconciled` fixture is `TestMergeExitsThreeWhenTheReconcileFails`.

## Acceptance

- [ ] A merge composition conflict prints a `next=` value that starts with `reviewer: `.
- [ ] A fold with a dirty sibling prints a route that contains `bench commit --in ` and no `bench worktree exec`.
- [ ] The `inherited` authorization refusal sentence equals `prospective authorization refused: inherited (the gate ran red on the composed tree and no green baseline attributes the red to this diff)`.
- [ ] The merge retries an empty-reason infrastructure refusal exactly once.
- [ ] A landing whose composed tree grades red prints a `refused{` record whose `next=` contains `bench commit --in ` and ends with the caller's re-run.
- [ ] A `landReviewed` error that no face claims prints a `next=` value that starts with `reviewer: `.
- [ ] Each merge face that this ticket declares has one producing fixture that follows its route out of the face.
