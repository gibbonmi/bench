# Route the commitment policy refusals through faces

Blocked by: 09-route-each-commit-refusal-through-the-registry.md, 10-route-each-checkpoint-refusal-by-its-cause.md
Writes: internal/refusalroute/registry.go (new), internal/commitment/delivery.go, internal/commitment/repository/admission.go, internal/commitment/repository/candidate.go, internal/commitment/repository/continuation.go, internal/commitment/repository/publication.go, internal/commitment/repository/light_path.go, internal/commitment/repository/verification.go, internal/commitment/repository/readiness.go, internal/commitment/repository/publication_test.go, internal/commit/commit.go, internal/commit/refusal_route_test.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/refusalroute/faces_commitment.go (new), internal/commitment/repository/light_path_test.go, internal/commit/commitment_test.go, internal/worktree/land_refusal.go, internal/worktree/commitment_light_landing_test.go, internal/worktree/land_identity_test.go, internal/commit/assessment_span_test.go (new), internal/commit/commit_test.go, internal/commit/commitment_route_test.go (new), internal/worktree/land.go, internal/refusalroute/routetest/routetest.go, internal/refusalroute/routetest/routetest_test.go, internal/worktree/refusal_route_follow_test.go, internal/worktree/land_refusal_fixture_test.go, internal/worktree/commitment_landing_fixture_test.go
Covers: RR61, RR62, RR63, RR64

## What to build

Each commitment policy error that carries an inline route tail becomes a typed error that names its face.
Its sentence drops the tail.
The spec section "The commitment route tails" lists each site and its face.
Declare each of those faces in the registry with the authority and the route of the spec's face inventory table.

`bench commit` prints the route of the commitment face when the commitment policy refuses its candidate.
`bench preflight` prints the sentence beside its own recovery column, and that column does not change.
No reader parses the tail text.

This ticket runs after ticket 09, because both write `internal/commit/commit.go` and `internal/commit/refusal_route_test.go`.
Add the commit rows of this ticket to `TestCommitFacesFollowTheirRoutes`.

The walk test uses the shared route walk in `internal/refusalroute/routetest` and does not copy it.

## Acceptance

- [ ] A commit from an assignment with no delivery binding prints a `next=` route that contains `bench commitment start --outcome`.
- [ ] A commit whose candidate changes the protected commitment prints a `next=` value that starts with `reviewer: `.
- [ ] A light-path commit with a path outside its ticket's `Writes:` line prints an agent route that names the ticket and not `bench commitment start`.
- [ ] The `assignment has no current delivery binding` refusal sentence does not contain `run bench`.
