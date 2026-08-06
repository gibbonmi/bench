# Allow an already-covered clean integration

Blocked by: allow-already-covered-clean-checkpoint.md
Ownership fence: `internal/specbuild/integrate.go`, `internal/specbuild/checkpoint_receipt_test.go`
Integration surfaces: clean checkpoint replay and attributed candidate advance→`internal/specbuild/integrate.go`; clean integration and release proof→`internal/specbuild/checkpoint_receipt_test.go` plus NI1; non-empty ticket-fence refusal→`internal/specbuild/checkpoint_receipt_test.go` plus NI2
Contracts: checkpoint patch bytes and changed paths cross checkpoint evidence→`internal/specbuild/integrate.go` replay and release, asserted by NI1-NI2 against the real integration lifecycle in `internal/specbuild/checkpoint_receipt_test.go`

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

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| NI1 | require `insideFence` for an empty changed-path set | `TestIntegrateAdmitsCleanAlreadyCoveredCheckpointOverCurrentCandidate` | run `go test ./internal/specbuild -run '^TestIntegrateAdmitsCleanAlreadyCoveredCheckpointOverCurrentCandidate$' -count=1`; expect the honest no-op integration to be refused |
| NI2 | skip fence validation for a non-empty changed-path set | `TestIntegrateRefusesNonEmptyCheckpointPathsOutsideTheFence` | run `go test ./internal/specbuild -run '^TestIntegrateRefusesNonEmptyCheckpointPathsOutsideTheFence$' -count=1`; expect the outside-fence integration to stop refusing |
