# Allow an already-covered clean checkpoint

Blocked by: allow-clean-assignment-refresh.md
Ownership fence: `internal/specbuild/checkpoint.go`, `internal/specbuild/checkpoint_receipt_test.go`
Integration surfaces: coordinator checkpoint receipt -> changed-path ownership validation; already-covered ticket rows -> clean provisional checkpoint
Contracts: `internal/specbuild/checkpoint.go` accepts an exact base tree with empty ownership only when the receipt otherwise validates, while every non-empty changed path still must match receipt ownership and remain inside the ticket fence; `internal/specbuild/checkpoint_receipt_test.go` proves clean closure and preserves hostile ownership refusals

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

- [ ] [MNC1] Requiring `insideFence` for an empty changed-path set refuses the honest no-op checkpoint.
- [ ] [MNC2] Treating all ownership as optional makes the existing unexplained-path or outside-fence receipt cases green.
