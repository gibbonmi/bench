# Migrate spec-build outcomes and actions

Blocked by: none
Ownership fence: `internal/specbuild`, `cmd/bench/specbuild.go`, `cmd/bench/specbuild_test.go`
Integration surfaces: shared carriers→implemented prerequisite `axi-carriers-and-registry`; final contraction→contract-outcome-action-routes.md
Contracts: operation name, lifecycle kind, exact exit, ordered action tokens or prose disposition, and empty class cross `internal/specbuild`→`cmd/bench/specbuild.go`, membership is all nine operations, operation order is grammar order, and absence is explicit no-run state, asserted by SB1
Closure: SB1/nine-operations, SB1/exits, SB1/actions, SB1/prose, SB1/empty, SB1/bypass

## What to build

all nine spec-build operations route exact outcomes and actions while preserving exits, empty status, and orchestration prose.

## Acceptance

- [ ] [SB1] (covers OA5) all nine spec-build operations route exact outcomes and actions while preserving exits, empty status, and orchestration prose.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| SB1/nine-operations | bypass the start operation | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| SB1/exits | normalize a refusal exit to usage | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| SB1/actions | drop the slug from one fixed action input | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| SB1/prose | mark release-assignment prose executable | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| SB1/empty | normalize one-row no-run status to zero rows | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| SB1/bypass | restore a local result for promote | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |

