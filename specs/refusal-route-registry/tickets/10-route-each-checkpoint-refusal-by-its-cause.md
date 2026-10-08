# Route each checkpoint refusal by its cause

Blocked by: 09-route-each-commit-refusal-through-the-registry.md
Writes: internal/refusalroute/registry.go (new), internal/gate/checkpoint.go, internal/gate/run_transaction.go, internal/gate/gate.go, internal/gate/complete_checkpoint.go, internal/gate/complete_checkpoint_test.go, internal/gate/run_outcomes_test.go, internal/gate/review_checkpoint_test.go, internal/gate/refusal_route_test.go (new), internal/refusalroute/faces_gate.go (new)
Covers: RR39, RR40, RR41, RR42, RR43, RR44, RR66, RR68

## What to build

The landed FT392 complete checkpoint raises three refusals of its own in `executeCompleteCheckpoint`, through `operational`.
The constant `cleanCheckoutRefusal` in `internal/gate/complete_checkpoint.go` is the one source of the dirty-checkout sentence, and it stays so.

Each cause in the pre-oracle checkpoint funnel maps to its own face.
Declare these gate faces in the registry with the routes of the spec's face inventory table:

- `checkpoint-completion-evidence`, which keeps the route that `routedRefusal` prints today
- `checkpoint-dirty-checkout`, for the `cleanCheckoutRefusal` refusal; the face declares no sentence and adds the route
- `checkpoint-composition`, for the `complete checkpoint unavailable` compose refusal
- `checkpoint-tip-moved`
- `checkpoint-subject-unavailable`, for a subject-capture fault that no other face claims, and for the `complete checkpoint source unavailable` refusals
- `gate-handback`

The checkpoint prints `next=<route>` on stderr after its reason.
It no longer prints the fixed `help[1]{cmd,why}` row.
The gate resolves the `<label>` fact through `intent.AssignmentsOwning` at the refusal.

Add `TestCheckpointFacesFollowTheirRoutes` in `internal/gate/refusal_route_test.go`.
It walks the gate faces of the registry and drives one producing fixture for each face.
Each fixture follows its printed route and reruns the checkpoint.
Replace the pin of the write-access help row in `TestGateRunRetainsSubjectConstructionCause`.
Extend `TestCompleteCheckpointRefusesADirtyCheckout` and `TestCompleteCheckpointRefusesAnUntransformableSpec` with the route assertions of RR41 and RR66.

The walk test uses the shared route walk in `internal/refusalroute/routetest` and does not copy it.

## Acceptance

- [ ] A checkpoint refusal for a missing completion record prints `next=bench preflight review example` and no `help[1]{cmd,why}` row.
- [ ] A complete checkpoint with an uncommitted tracked edit prints `cleanCheckoutRefusal` and a `next=` route that contains `bench commit --in `.
- [ ] A complete checkpoint on a dirty assignment checkout prints a `next=` route that contains the assignment label after `--in `.
- [ ] A complete checkpoint on a spec with no `Status: staged` line prints a `next=` value that starts with `reviewer: `.
- [ ] A subject-capture fault prints `next=` with a route that contains `bench doctor` and `--fresh`, and it prints no `help[1]{cmd,why}` row.
- [ ] Each gate face has exactly one producing fixture that follows its route out of the face.
