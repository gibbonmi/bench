# Move the landing faces into the shared registry

Blocked by: 01-create-the-shared-refusal-route-registry.md
Writes: internal/refusalroute/registry.go (new), internal/refusalroute/registry_test.go (new), internal/worktree/land_refusal.go, internal/worktree/land.go, internal/worktree/land_identity.go, internal/worktree/land_resume.go, internal/worktree/classifier.go, internal/worktree/identity_component.go, internal/worktree/build.go, internal/worktree/merge.go, internal/worktree/refusal_route_test.go (new), internal/worktree/identity_component_test.go, internal/worktree/land_surface_test.go, internal/worktree/land_tickets_only_test.go, internal/worktree/land_release_refusal_test.go, internal/worktree/land_folded_base_test.go, internal/worktree/land_identity_test.go, internal/worktree/land_resume_refusal_test.go, internal/worktree/land_journey_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: RR01, RR02, RR03, RR14, RR16, RR17, RR18, RR19, RR20

## What to build

Move every landing face from the worktree registry into the shared registry.
Remove the worktree copy of the inventory, so the shared registry is the one source.
Sweep every refusal site of the landing and map each site to one face.

Declare each land face of the spec's face inventory table with the authority and the route of that table.
Add the `land-handback` reviewer face.
The table changes the authority of the conflict faces and the destination faces, and the routes of `source-not-clean` and `resume-marker`.
Each other land face keeps its current route.
The conflict routes keep the hand merge text behind the `reviewer: ` prefix.

Each landing preflight route still ends with the caller's own re-run.
The skipped-proof sentence stays the first step of the route, after the authority marker.
The worktree `refusal` record printer composes its `next` label from the registry export.
Move the hand merge text into the registry, and change `merge.go` and `build.go` to read it from there.

Change the tests that compare a route with `landingRefusalFaceByName(...).route(...)` so that they read the shared registry.
`TestLandingRefusalRegistryHasAProducingFixture` walks the land faces of the shared registry.
Each new expectation derives from the fixture inputs and the declared faces, never from the renderer under test.

## Acceptance

- [ ] Every registered face declares a verb from the six write verbs, the authority `agent` or `reviewer`, and at least one route step.
- [ ] Each land face has exactly one producing fixture, and each fixture produces a registered land face.
- [ ] The landing composition-conflict refusal and the composition-conflict-pending refusal each print a `next=` value that starts with `reviewer: `.
- [ ] The landing destination-not-clean refusal prints a `next=` value that starts with `reviewer: `.
- [ ] The landing source-not-clean refusal prints a route that contains `bench commit --in ` and the source label.
- [ ] `TestLandCommandReportsEveryRefusalInOnePreflight` passes, and each preflight route ends with the caller's own re-run.
- [ ] No production file under `internal/worktree` declares a landing face inventory.
