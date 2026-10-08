# Route each commit refusal through the registry

Blocked by: 08-route-the-commit-exit-3-to-the-reset-plan.md
Writes: internal/refusalroute/registry.go (new), internal/commit/commit.go, internal/commit/refusal_route_test.go (new), internal/commit/dry_run_test.go, internal/commit/landing_test.go, internal/refusalroute/faces_commit.go (new), internal/landing/landing.go, internal/landing/landing_test.go
Covers: RR36, RR37, RR38

## What to build

Each commit refusal prints `next=<route>` on its stderr refusal line.
Declare these commit faces in the registry with the routes of the spec's face inventory table:

- `commit-primary-checkout`
- `commit-red`
- `commit-infrastructure`
- `commit-handback`

A raising site derives the `<label>` fact through `intent.AssignmentsOwning` over its own root.
When no assignment owns the root, the slot prints `<label>`.
A red commit routes to the repair of each failure that the run reports, then the caller's own commit command.
It no longer prints only `run bench gate --fresh`.
A commit in the primary checkout routes to `bench worktree create --request <opaque-id> --label <work-item>`.

The commit help line `exit 3: published; the checkout did not reconcile — paste next= to repair` composes its label from the registry export.
The help line text does not change.

Add `TestCommitFacesFollowTheirRoutes`.
It walks the commit faces of the registry and drives one producing fixture for each face.
Each fixture follows its printed route and reruns the commit.

## Acceptance

- [ ] Each commit face has exactly one producing fixture that follows its route out of the face.
- [ ] A red commit refusal prints a `next=` line on stderr whose route ends with the caller's commit command.
- [ ] A commit in the primary checkout prints a `next=` route that contains `bench worktree create --request`.
