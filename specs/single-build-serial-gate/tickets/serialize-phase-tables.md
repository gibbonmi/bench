# Serialize every phase table

Blocked by: route-ordinary-phase-plumbing.md, migrate-gate-helpers.md
Ownership fence: `internal/gate/`
Integration surfaces: final phase declarations→route-ordinary-phase-plumbing.md; selected gate fixtures→migrate-gate-helpers.md; serial scheduler→serialize-primary-stripped-schedule.md; scheduler constructor census→contract-ordinary-build-census.md; cancellation lifecycle→contract-run-directory-lifecycle.md
Contracts: ordered `Phase` declarations cross `internal/gate/phases.go`→the scheduler in `internal/gate/runner.go`, membership is every outer and inner phase-table invocation, ordering chooses the first ready declaration and settles it before launching another, while results preserve needs, optional skips, unrelated execution, aggregate red, and process-group cancellation, asserted by SR1-SR3 against the real scheduler
Closure: SR1/max-active-one, SR1/stable-ready-order, SR1/outer-serial, SR1/inner-serial, SR2/needs-order, SR2/red-dependent-skip, SR2/optional-skip, SR2/unrelated-after-red, SR2/aggregate-red, SR3/interrupt-code, SR3/timeout-code, SR3/descendant-reap

## What to build

Contract the scheduler to one active phase process without replacing its DAG semantics with a simple fail-fast loop. Replace the positive overlap test with a max-active and stable-ready-order proof, then retain and adapt the existing dependency, skip, red, and process-group tests.

## Acceptance

- [ ] [SR1] (covers SG1) every outer and inner phase-table invocation has maximum active phase count one and launches simultaneously ready phases in stable declaration order.
- [ ] [SR2] (covers SG3) serial scheduling preserves needs order, red-dependent and optional skips, unrelated execution after red, and aggregate-red outcome.
- [ ] [SR3] (covers SG4) interrupt and timeout preserve their established codes, signal the active process group, and reap descendants before scheduler return.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| SR1/max-active-one | launch two ready phases before either settles | scheduler barrier test | provide four edge-free phases and require bounded max active equals one |
| SR1/stable-ready-order | choose ready work from map iteration | scheduler order test | run the same table repeatedly and require declaration-order records |
| SR1/outer-serial | retain `runPhasesConcurrent` for outer mode | outer command test | invoke outer phases and require max active one |
| SR1/inner-serial | add a separate concurrent inner path | inner command test | invoke an inner manifest and require max active one |
| SR2/needs-order | ignore a declared need | scheduler needs test | block the producer and require the consumer never starts early |
| SR2/red-dependent-skip | run a dependent after its need reds | scheduler red test | red the producer and require a dependent skip record |
| SR2/optional-skip | treat a missing optional executable as red | scheduler optional test | omit shellcheck and require the established skip outcome |
| SR2/unrelated-after-red | fail fast after the first red | scheduler aggregation test | red one phase and require a later unrelated marker still runs |
| SR2/aggregate-red | return green after recording a red | scheduler aggregation test | produce one red plus greens and require red result |
| SR3/interrupt-code | translate context interrupt to ordinary red | process-group test | interrupt the active phase and require code 130 |
| SR3/timeout-code | translate deadline to interrupt or red | process-group test | expire the gate deadline and require code 124 |
| SR3/descendant-reap | return after the leader exits while its descendant lives | orphan-descendant test | cancel a leader that leaves a child and require the group is gone before return |
