# Give the red-source fold an exit

Blocked by: 05-route-each-merge-refusal-through-the-registry.md
Writes: internal/refusalroute/registry.go (new), internal/worktree/merge.go, internal/landing/merge.go, internal/worktree/merge_route_test.go (new), internal/worktree/merge_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/worktree/parallel_census_test.go, internal/worktree/merge_refusal.go, internal/worktree/land.go, internal/worktree/refusal_route_follow_test.go, internal/refusalroute/faces_land.go (new), internal/refusalroute/faces_merge.go (new), internal/worktree/merge_target_grade_test.go (new), internal/worktree/merge_caller_root_test.go
Covers: RR21, RR22, RR23, RR24, RR25, RR26, RR55, RR56, RR57, RR59

## What to build

A fold red picks its face by its cause.
Declare the `merge-target-red` agent face and the `merge-fold-red` reviewer face with the routes of the spec's face inventory table.

A `candidate` gate kind carries the attribution, and it takes `merge-fold-red`.
On an `inherited` or a `lane fail` fold, the merge grades the target tip alone before it picks the face.
That run uses the same lane, or the whole gate when the project declares no lane.
It runs on the refusal path only, and it changes no publish rule.

- When the target tip alone is red, the fold refuses with `merge-target-red`. The route repairs each failing check in the target, commits through `bench commit --in <label>`, and reruns the fold.
- When the target tip alone is green, the fold refuses with `merge-fold-red`. The route hands back the open FT342 decision.

Promote the collision 5a repro to `TestRedSourceFoldNamesAnExit`, with no build tag.
The test drives a red source that folds `main`, then a fold into a dirty target.
It follows the printed route: it removes the red file, runs the printed commit step, and runs the printed fold step.
Extend `TestMergeFacesFollowTheirRoutes` with a fixture for each new face.

## Acceptance

- [ ] A fold of `main` into a target whose committed tip fails its lane prints a `refused{` record that contains `next=`.
- [ ] The target-red route contains `bench commit --in ` and the target label, and it ends with `bench worktree merge --from `, the incoming spelling, and the target.
- [ ] After the fixture removes the red file and runs the printed commit step, the printed fold step exits 0.
- [ ] A fold into a dirty target prints a `refused{` record whose `next=` contains `bench commit --in `.
- [ ] A `lane fail` fold whose incoming commit adds the red to a lane-green target prints a `next=` value that starts with `reviewer: `.
- [ ] A `lane fail` fold whose target tip alone fails the lane prints a `next=` value that does not start with `reviewer: `.
- [ ] A `candidate` fold red prints a `next=` value that starts with `reviewer: `.
- [ ] An `inherited` fold whose target tip alone grades green prints a `next=` value that starts with `reviewer: `.
- [ ] An `inherited` fold whose target tip alone grades red prints a `next=` value with `bench commit --in ` and the target label.
