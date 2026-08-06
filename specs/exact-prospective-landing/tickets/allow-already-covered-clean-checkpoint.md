# Allow an already-covered clean checkpoint

Blocked by: allow-clean-assignment-refresh.md
Ownership fence: `internal/specbuild/checkpoint.go`, `internal/specbuild/checkpoint_receipt_test.go`
Integration surfaces: coordinator checkpoint receipt and changed-path ownership validation→`internal/specbuild/checkpoint.go`; clean base-tree checkpoint proof→`internal/specbuild/checkpoint_receipt_test.go` plus NC1; non-empty ownership and fence refusals→`internal/specbuild/checkpoint_receipt_test.go` plus NC2
Contracts: receipt ownership and live changed paths cross coordinator evidence→`internal/specbuild/checkpoint.go` validation, asserted by NC1-NC2 against the real checkpoint lifecycle in `internal/specbuild/checkpoint_receipt_test.go`

## What to build

Permit checkpoint evidence for an assignment whose worktree tree exactly equals its
base and whose ownership list is empty. This is the lifecycle closure for a ticket
whose acceptance rows are independently verified as already covered after a
prerequisite repair. Do not weaken validation for any non-empty payload: changed
paths must still equal receipt ownership and remain inside the assignment fence.

## Acceptance

- [ ] [NC1] A clean assignment can checkpoint with empty ownership and already-covered row outcomes, producing a provisional checkpoint whose tree equals the assignment base tree.
- [ ] [NC2] Any non-empty changed path still must be declared by the receipt and lie inside the ticket ownership fence.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| NC1 | require `insideFence` for an empty changed-path set | `TestCheckpointAdmitsCleanAlreadyCoveredAssignment` | run `go test ./internal/specbuild -run '^TestCheckpointAdmitsCleanAlreadyCoveredAssignment$' -count=1`; expect the honest no-op checkpoint to be refused |
| NC2 | make ownership optional for non-empty changed paths | `TestCheckpointRereadsEveryLiveFact` | run `go test ./internal/specbuild -run '^TestCheckpointRereadsEveryLiveFact$' -count=1`; expect the unexplained-path or outside-fence receipt case to stop refusing |
