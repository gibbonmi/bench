# Remove the dead kit wrappers

Blocked by: none
Ownership fence: `internal/gate/component_decision.go`, `internal/gate/evaluation.go`
Integration surfaces: `internal/gate/evaluation.go` explicit constructor→`internal/gate/gate.go` + RDW1; `internal/gate/component_decision.go` explicit scoping function→`internal/gate/evaluation.go` + RDW1
Contracts: the explicit functions in `internal/gate/evaluation.go` and `internal/gate/component_decision.go` have present callers and no zero-value forwarding aliases
Closure: RDW1/dead-wrapper-removal

## What to build

Delete the zero-caller forwarding wrappers `newEngineEvaluation` and
`scopeComponentsForIdentityGenerations`. Their callers already use the
explicit forms that name both root and kit, so the wrappers advertise no
supported boundary and retain no behavior.

## Acceptance

- [ ] [RDW1] (covers local) the two zero-caller forwarding wrappers are absent and every remaining explicit function has a present caller, repairing S1.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RDW1/dead-wrapper-removal | reintroduce either zero-caller forwarding wrapper | the zero-caller declaration search | apply, run the search, expect it to name the wrapper declaration |
