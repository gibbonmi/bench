# Route each checkpoint refusal by its cause

Blocked by: 09-route-each-commit-refusal-through-the-registry.md
Writes: internal/refusalroute/registry.go (new), internal/gate/checkpoint.go, internal/gate/run_transaction.go, internal/gate/gate.go, internal/gate/complete_checkpoint.go (new), internal/gate/complete_checkpoint_test.go (new), internal/gate/run_outcomes_test.go, internal/gate/review_checkpoint_test.go, internal/gate/refusal_route_test.go (new)
Covers: RR39, RR40, RR41, RR42, RR43, RR44

## What to build

This ticket starts only on a base that contains the FT392 landing.
FT392 creates `internal/gate/complete_checkpoint.go` and `internal/gate/complete_checkpoint_test.go`, and the dirty-checkout refusal lives there.
If the base does not contain FT392, stop and report.

Each cause in the pre-oracle checkpoint funnel maps to its own face.
Declare these gate faces in the registry with the routes of the spec's face inventory table:

- `checkpoint-completion-evidence`, which keeps the route that `routedRefusal` prints today
- `checkpoint-dirty-checkout`, for the FT392 dirty-checkout refusal
- `checkpoint-composition`
- `checkpoint-tip-moved`
- `checkpoint-subject-unavailable`, for a subject-capture fault that no other face claims
- `gate-handback`

The checkpoint prints `next=<route>` on stderr after its reason.
It no longer prints the fixed `help[1]{cmd,why}` row.
The gate resolves the `<label>` fact through `intent.AssignmentsOwning` at the refusal.

Add `TestCheckpointFacesFollowTheirRoutes` in `internal/gate/refusal_route_test.go`.
It walks the gate faces of the registry and drives one producing fixture for each face.
Each fixture follows its printed route and reruns the checkpoint.
Replace the pin of the write-access help row in `TestGateRunRetainsSubjectConstructionCause`.

## Acceptance

- [ ] A checkpoint refusal for a missing completion record prints `next=bench preflight review example` and no `help[1]{cmd,why}` row.
- [ ] A complete checkpoint on a dirty assignment checkout prints a `next=` route that contains `bench commit --in ` and the assignment label.
- [ ] A subject-capture fault prints `next=` with a route that contains `bench doctor` and `--fresh`, and it prints no `help[1]{cmd,why}` row.
- [ ] Each gate face has exactly one producing fixture that follows its route out of the face.
