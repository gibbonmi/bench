# Migrate the shift outcome

Blocked by: none
Ownership fence: `internal/shift`
Integration surfaces: shared outcome carrier→implemented prerequisite `axi-carriers-and-registry`; final contraction→contract-outcome-action-routes.md
Contracts: shift kind enum, exact exit integer, ordered result fields, bounded detail, and durable interrupted state cross `internal/shift` outcome→renderer, membership is complete/failed/usage/incomplete/no-op/interrupted, order is current resultFields, and absence is empty optional detail, asserted by SH1
Closure: SH1/six-kinds, SH1/six-exits, SH1/field-order, SH1/detail, SH1/interrupted, SH1/bypass

## What to build

all six shift outcomes retain exits 0/1/2/3/4/130 through the shared route and exact renderer.

## Acceptance

- [ ] [SH1] (covers OA8) all six shift outcomes retain exits 0/1/2/3/4/130 through the shared route and exact renderer.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| SH1/six-kinds | remove interrupted from the policy | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| SH1/six-exits | change no-op exit 4 | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| SH1/field-order | reorder one result field | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| SH1/detail | remove the existing detail bound | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| SH1/interrupted | lose interrupted state on fresh-process reload | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| SH1/bypass | emit directly from shift Outcome | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |

