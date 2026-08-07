# Repair recompose bootstrap marker

Blocked by: none
Ownership fence: `internal/specbuild/recompose.go`, `internal/specbuild/promotion_recompose_test.go`
Integration surfaces: promotion recomposition→`internal/specbuild/recompose.go` + RBM1; recomposition contract→`internal/specbuild/promotion_recompose_test.go` + RBM1
Contracts: a working-branch advance crosses `Promote`→`recomposePromotion`→`gate.Bootstrap`, asserted by RBM1 against the live green marker rather than the run's recorded base
Closure: RBM1/bootstrap-recorded-base

## What to build

When promotion recomposes a run onto a moved working tip, bootstrap the gate
against the live green marker `greenMarker(s.root, subject.branch)` — not the
run's recorded base. The recorded base stays unchanged until
`finishRecomposition` completes, so an interrupted recomposition never leaves
the run pointing past evidence it does not have. Carry the strong regression:
when a sibling commit has advanced the green marker past the recorded base,
bootstrap must receive the live marker as its expected commit while the record
still holds the old base at call time.

## Acceptance

- [ ] [RBM1] (covers local) With the green marker advanced past the recorded base by a sibling commit, `Promote` bootstraps against the live marker, retains the recorded base until recomposition completes, and recomposes onto the working tip.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RBM1/bootstrap-recorded-base | pass `run.Base` instead of the live green marker to `gate.Bootstrap` | recomposition contract | advance the green marker to a sibling commit, advance the working tip, run `Promote`, expect the bootstrap-subject assertion red because expected != live marker |
