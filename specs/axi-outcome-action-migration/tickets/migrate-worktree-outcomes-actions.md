# Migrate worktree outcomes and actions

Blocked by: migrate-output-adapter.md
Ownership fence: `internal/worktree`, `cmd/bench/main.go`, `cmd/bench/main_test.go`
Integration surfaces: returned-output adapter and dispatch→migrate-output-adapter.md; shared carriers→implemented prerequisite `axi-carriers-and-registry`; final contraction→contract-outcome-action-routes.md
Contracts: mode enum, lifecycle kind, exact exit, ordered recovery/action tokens, fingerprint-bearing state, and empty class cross `internal/worktree`→`cmd/bench/main.go`, membership is query plus create/list/refresh/cleanup/release/recovery/resume, order is authority derivation then carry then render, and absence is exact empty/no-op, asserted by WT1 and WT2
Closure: WT1/query-route, WT2/seven-modes, WT2/action-inputs, WT2/fingerprint-authority, WT2/empty, WT2/bypass

## What to build

worktree query outcomes reach the shared result route. seven worktree lifecycle families route outcomes/actions without changing authority or empty forms.

## Acceptance

- [ ] [WT1] (covers OA3) worktree query outcomes reach the shared result route.
- [ ] [WT2] (covers OA7) seven worktree lifecycle families route outcomes/actions without changing authority or empty forms.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| WT1/query-route | restore local list result | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| WT2/seven-modes | bypass create | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| WT2/action-inputs | drop the recovery ref from a fixed action | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| WT2/fingerprint-authority | derive an action before fingerprint validation | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| WT2/empty | normalize empty list to completed cleanup | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |
| WT2/bypass | restore local recovery rendering | the independent route or compatibility test | apply the subject mutation, run the focused real producer plus compatibility case under its named bound, and require the specific red |

