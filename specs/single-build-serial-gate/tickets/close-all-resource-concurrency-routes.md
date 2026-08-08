# Close all resource concurrency routes

Blocked by: enforce-the-executable-artifact-contract.md, serialize-gate-lineages-and-transfer-turns.md, reclaim-interrupted-resource-turns.md, admit-one-gate-per-common-repository.md, pin-go-children-and-settle-core-packages.md, serialize-canary-stages.md
Ownership fence: `internal/artifactstore/`, `internal/gate/`, `internal/canary/`, `internal/conformance/`, `internal/contract/surface/publication/publication_canary_test.go`, `internal/bounds/bounds.go`, `projects/benchkit.md`, `decisions/gate-concurrency.md`, `decisions/gate-budget.md`, `decisions/gate-pipeline.md`, `decisions/gate-critical-path.md`, `docs/adr/0009-canary-concurrency-budget.md`
Integration surfaces: artifact/proof resource classes→enforce-the-executable-artifact-contract.md; transferable turn→serialize-gate-lineages-and-transfer-turns.md; exact-once interruption settlement→reclaim-interrupted-resource-turns.md; repository admission→admit-one-gate-per-common-repository.md; Go/package progress→pin-go-children-and-settle-core-packages.md; canary stages→serialize-canary-stages.md; nested publication Go child→internal/contract/surface/publication/publication_canary_test.go; final resource and coordination-test registries→owned paths; concurrency-policy advertisements→projects/benchkit.md + decisions/gate-concurrency.md + decisions/gate-budget.md + decisions/gate-pipeline.md + decisions/gate-critical-path.md + docs/adr/0009-canary-concurrency-budget.md
Contracts: resource registry rows (class, typed launch adapter, turn-transfer posture, progress owner) cross producers in `internal/artifactstore/`, `internal/gate/`, `internal/canary/`, and `internal/contract/surface/publication/publication_canary_test.go`→process-tree recorder in `internal/gate/`, membership is phase, stripped materialization/cleanup, Go child, package/fixture materialization, publication/sealing, canary item, and proof classes, ordering is one transferred turn lineage, absence is conformance red; coordination-test registry in `internal/conformance/` is closed and every absent `internal/gate` `t.Parallel` is forbidden; `ProgressSettlement` crosses package, phase, canary-item, and turn owners→watchdog in `internal/gate/`, and only terminal membership resets the deadline; asserted by RC1-RC4 against real producers
Closure: RC1/phase-class, RC1/stripped-materialization-class, RC1/stripped-cleanup-class, RC1/go-child-class, RC1/publication-canary-go-child, RC1/package-fixture-class, RC1/publication-class, RC1/canary-class, RC1/proof-class, RC1/max-runnable-one, RC2/no-cli-width-knob, RC2/no-env-width-knob, RC2/no-manifest-width-knob, RC2/no-unregistered-tparallel, RC2/coordination-caller-serialization, RC2/high-core-invariance, RC2/ambient-gomaxprocs-invariance, RC2/ambient-goflags-invariance, RC2/ambient-test-width-invariance, RC3/canary-item-progress, RC3/turn-progress, RC4/profile-outer-serial, RC4/profile-inner-width-absent, RC4/profile-worker-arithmetic-absent, RC4/profile-operator-lever-absent, RC4/gate-concurrency-product-budget-retired, RC4/gate-concurrency-operator-lever-retired, RC4/gate-budget-pool-retired, RC4/gate-budget-positive-width-retired, RC4/gate-budget-work-conserving-retired, RC4/canary-adr-worker-pool-retired, RC4/canary-adr-inner-width-retired, RC4/canary-adr-operator-budget-retired, RC4/gate-pipeline-outer-concurrency-retired, RC4/gate-pipeline-gomaxprocs-two-retired, RC4/gate-pipeline-worker-derivation-retired, RC4/gate-critical-path-outer-concurrency-retired

## What to build

Contract the serialization migration: close the resource registry over every producer, remove all 196 current gate `t.Parallel` calls, allow only named lightweight coordination callers, and enforce that no CLI/env/manifest route can widen an ordinary gate.

## Acceptance

- [ ] [RC1] (covers ZC5) every registered resource class transfers the one turn and the process-tree recorder never observes more than one runnable resource actor.
- [ ] [RC2] (covers ZC7) no public width knob or unregistered gate `t.Parallel` exists, coordination callers serialize resource children, and host/ambient width cannot widen the gate.
- [ ] [RC3] (covers local) the progress watchdog accepts terminal canary-item and turn settlements through the same closed reset-source registry as package and phase settlements, and rejects output as a reset source.
- [ ] [RC4] (covers local) the benchkit profile, gate concurrency, gate budget, gate pipeline, gate critical-path decisions, and canary-budget ADR describe the current capacity-one state with serial canary work and no `GOMAXPROCS` operator lever or positive-width policy.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RC1/phase-class | bypass the turn in phase execution | process-tree active recorder | run two ready class probes and expect maximum runnable one |
| RC1/stripped-materialization-class | omit stripped-subject materialization enrollment | resource registry audit plus blocking materializer | start materialization beside another class and expect audit or double-actor failure |
| RC1/stripped-cleanup-class | omit stripped-subject cleanup enrollment | resource registry audit plus blocking cleanup | start cleanup beside another class and expect audit or double-actor failure |
| RC1/go-child-class | launch one Go child outside the typed turn adapter | resource registry audit | add the raw launch and expect missing class ownership |
| RC1/publication-canary-go-child | launch the nested publication `go test` directly | publication canary resource test | materialize its marker beside another registered actor and expect one turn plus width-one argv before execution |
| RC1/package-fixture-class | materialize a package/fixture outside the turn | blocking materialization test | overlap its request with another class and expect maximum one |
| RC1/publication-class | publish/seal outside the turn | blocking publisher test | overlap publication with another class and expect maximum one |
| RC1/canary-class | bypass the turn in one canary item kind | canary/resource integration test | run the mutated item beside another class and expect maximum one |
| RC1/proof-class | author a proof artifact outside the turn | proof/resource integration test | request the proof beside another identity and expect maximum one |
| RC1/max-runnable-one | widen the turn capacity to two | all-class process-tree recorder | drive one member of every class and expect the recorded maximum to stay one |
| RC2/no-cli-width-knob | add a CLI capacity option | usage structural test | advertise or parse the option and expect the forbidden-knob diagnostic |
| RC2/no-env-width-knob | add an environment capacity option | environment registry test | advertise or consume the variable and expect the forbidden-knob diagnostic |
| RC2/no-manifest-width-knob | add a manifest capacity option | manifest schema test | advertise or parse the field and expect the forbidden-knob diagnostic |
| RC2/no-unregistered-tparallel | add `t.Parallel` outside the coordination registry | Go AST gate test audit | add the call and expect its exact test-name diagnostic |
| RC2/coordination-caller-serialization | let two registered coordination callers delegate resource children simultaneously | coordination registry integration test | release both callers and expect one resource child marker at a time |
| RC2/high-core-invariance | derive capacity from detected cores | high-core fake resolver test | report a large host width and expect capacity exactly one |
| RC2/ambient-gomaxprocs-invariance | honor ambient `GOMAXPROCS` | hostile environment end-to-end scheduler test | export a wider runtime value and expect one runnable actor plus authoritative child values of one |
| RC2/ambient-goflags-invariance | honor ambient widening `GOFLAGS` | hostile environment end-to-end scheduler test | export widening Go flags and expect one runnable actor plus width-one argv |
| RC2/ambient-test-width-invariance | honor ambient test parallelism flags | hostile environment end-to-end scheduler test | supply a wider test flag and expect one runnable actor plus `-parallel=1` |
| RC3/canary-item-progress | omit terminal canary-item settlements from the watchdog reset registry | fake-clock junction test | settle canary items inside each window while cumulative time exceeds it and expect no timeout |
| RC3/turn-progress | omit terminal turn settlements from the watchdog reset registry | fake-clock junction test | settle transferred turns inside each window while cumulative time exceeds it and expect no timeout |
| RC4/profile-outer-serial | retain the profile's four-concurrent-phases statement | profile currency assertion | inspect the gate section and expect exactly one outer lineage with no concurrent-phase advertisement |
| RC4/profile-inner-width-absent | retain `CanaryInnerWidth` in the profile | profile currency assertion | inspect the canary policy paragraph and expect no retired inner-width name |
| RC4/profile-worker-arithmetic-absent | retain derived canary worker arithmetic in the profile | profile currency assertion | inspect the canary policy paragraph and expect one serial iterator with no worker formula |
| RC4/profile-operator-lever-absent | retain `GOMAXPROCS=8 bench gate` as an operator lever | profile currency assertion | inspect the gate policy and expect no supported widening invocation |
| RC4/gate-concurrency-product-budget-retired | retain the worker-times-width product budget as current | decision currency assertion | inspect the decision and expect capacity one without product arithmetic |
| RC4/gate-concurrency-operator-lever-retired | retain outer width as an operator lever | decision currency assertion | inspect the decision and expect no supported widening route |
| RC4/gate-budget-pool-retired | retain a runnable resource pool | decision-map currency assertion | inspect the active map and expect the pool answer to point at capacity one |
| RC4/gate-budget-positive-width-retired | retain a positive configurable width | decision-map currency assertion | inspect the active map and expect no positive-width build target |
| RC4/gate-budget-work-conserving-retired | retain work-conserving scheduling | decision-map currency assertion | inspect the active map and expect strict serial admission |
| RC4/canary-adr-worker-pool-retired | retain a canary worker pool as accepted architecture | ADR currency assertion | inspect the current-state record and expect one serial iterator |
| RC4/canary-adr-inner-width-retired | retain fixed positive canary inner width | ADR currency assertion | inspect the current-state record and expect width-one children without an inner-width policy |
| RC4/canary-adr-operator-budget-retired | retain an operator-controlled canary budget | ADR currency assertion | inspect the current-state record and expect no operator-controlled budget |
| RC4/gate-pipeline-outer-concurrency-retired | retain concurrent outer phases in the pipeline decision | decision currency assertion | inspect the current state and expect one ordered phase lineage |
| RC4/gate-pipeline-gomaxprocs-two-retired | retain `GOMAXPROCS=2` child policy in the pipeline decision | decision currency assertion | inspect the current state and expect runtime width one |
| RC4/gate-pipeline-worker-derivation-retired | retain worker derivation arithmetic in the pipeline decision | decision currency assertion | inspect the current state and expect no derived worker formula |
| RC4/gate-critical-path-outer-concurrency-retired | retain concurrent gate branches in the critical-path decision | decision currency assertion | inspect the current state and expect a singular resource lineage |
