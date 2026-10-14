# Name the production paths in the unbound commit refusal

Blocked by: none
Writes: internal/commitment/repository/candidate.go, internal/commitment/repository/unbound_refusal_test.go (new), internal/refusalroute/registry.go, internal/refusalroute/route.go, internal/refusalroute/faces_commitment.go, internal/refusalroute/route_test.go, internal/refusalroute/routetest/routetest.go, internal/commit/commitment_route_test.go, internal/worktree/land_identity_test.go
Covers: none

## What to build

A commit in an assignment that has no current delivery binding can include production paths. The refusal now names each production path that needs a binding. The refusal names the first five paths and the total count when more than five paths need a binding.
The refusal takes the path list from the one place that computes the production paths.

The face `commitment-unbound` keeps its name. Its route now states three ways out ahead of the bind command.
The ways out are to bind the outcome, to remove the production paths from the change, or to cover them with a one-ticket light-path folder.

## Acceptance

- [ ] An unbound commit of one production path refuses with a sentence that names that path.
- [ ] An unbound commit of more than five production paths names the first five and the total count.
- [ ] A planning path in the commit is not named in the refusal.
- [ ] The printed route of the face names the remove alternative and the light-path folder alternative, and still ends with the bind command.
