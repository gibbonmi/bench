# Allow an already-covered clean integration

Blocked by: allow-already-covered-clean-checkpoint.md
Ownership fence: `internal/specbuild/integrate.go`, `internal/specbuild/checkpoint_receipt_test.go`
Integration surfaces: clean provisional checkpoint -> candidate replay; empty changed-path set -> ticket fence validation; attributed integration commit -> assignment release
Contracts: `internal/specbuild/integrate.go` treats an empty checkpoint payload as the identity transformation while every non-empty checkpoint path remains inside the ticket fence; `internal/specbuild/checkpoint_receipt_test.go` proves a clean already-covered checkpoint integrates, advances by one attributed commit without changing the candidate tree, and releases the assignment

## What to build

Permit integration of a verified provisional checkpoint whose changed-path set is
empty. This is the lifecycle closure paired with clean checkpoint acceptance: the
empty payload is an identity transformation over the current candidate, but still
produces the existing attributed integration commit and durable assignment release.
Do not weaken fence validation for any non-empty payload or any other checkpoint
provenance check.

## Acceptance

- [ ] [NI1] A clean already-covered checkpoint integrates over the current candidate, preserves its tree, records one attributed candidate advance, and releases the assignment.
- [ ] [NI2] Every non-empty checkpoint path still must lie inside the ticket ownership fence and retain existing conflict and provenance refusals.

## Red mutations

- [ ] [MNI1] Requiring `insideFence` for an empty changed-path set refuses the honest no-op integration.
- [ ] [MNI2] Skipping fence validation for a non-empty changed-path set makes an existing outside-fence integration case green.
