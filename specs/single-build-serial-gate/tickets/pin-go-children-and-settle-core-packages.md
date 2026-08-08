# Pin Go children and settle core packages

Blocked by: serialize-gate-lineages-and-transfer-turns.md
Ownership fence: `internal/gate/`
Integration surfaces: transferable turn/progress owner→serialize-gate-lineages-and-transfer-turns.md; package settlements and width-one child policy→`internal/gate/`; canary-item and turn settlement/watchdog junction→close-all-resource-concurrency-routes.md; final resource enrollment→close-all-resource-concurrency-routes.md
Contracts: `PackageSettlement` (ordered package identity, terminal status) crosses real `CoreTestPackages` producer→progress watchdog in `internal/gate/`, membership is every enumerated package exactly once, ordering is enumeration order while continuing after reds, absence retains the current stall deadline and cannot be replaced by output, asserted by GW1-GW2 against the real package runner; gate-package child width is exactly `-p=1`, `-parallel=1`, `GOMAXPROCS=1` with ambient absence semantics stripped
Closure: GW1/build-package-slot-one, GW1/test-package-slot-one, GW1/vet-package-slot-one, GW1/run-package-slot-one, GW1/compile-package-slot-one, GW1/test-test-slot-one, GW1/compile-test-slot-one, GW1/build-runtime-thread-one, GW1/test-runtime-thread-one, GW1/vet-runtime-thread-one, GW1/run-runtime-thread-one, GW1/compile-runtime-thread-one, GW1/build-goflags-stripped, GW1/test-goflags-stripped, GW1/vet-goflags-stripped, GW1/run-goflags-stripped, GW1/compile-goflags-stripped, GW1/build-gomaxprocs-stripped, GW1/test-gomaxprocs-stripped, GW1/vet-gomaxprocs-stripped, GW1/run-gomaxprocs-stripped, GW1/compile-gomaxprocs-stripped, GW2/per-package-invocation, GW2/continue-after-red, GW2/package-settlement-order, GW2/cumulative-progress, GW2/package-stall, GW2/heartbeat-ignored

## What to build

Pin every gate-owned Go child to width one and change the core test phase from one multi-package child to ordered one-package invocations. Feed only terminal package/phase/item/turn settlements to the 45-minute no-progress watchdog.

## Acceptance

- [ ] [GW1] (covers ZC2) every Go build/test/vet/run/compile constructor owned by `internal/gate` has one package/test slot and one runtime thread regardless of ambient width flags; the canary-owned compile constructor closes in CS3 and the composed invariant closes in RC2.
- [ ] [GW2] (covers ZC8) core tests invoke every real enumerated package once in order, continue after reds, publish terminal settlements, allow cumulative wall beyond 45 minutes with progress, ignore output heartbeats, and time out one stalled package with code 124 and teardown.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| GW1/build-package-slot-one | omit or widen `-p=1` on build | argv matrix test | inspect the build constructor and expect its exact package-width failure |
| GW1/test-package-slot-one | omit or widen `-p=1` on test | argv matrix test | inspect the test constructor and expect its exact package-width failure |
| GW1/vet-package-slot-one | omit or widen `-p=1` on vet | argv matrix test | inspect the vet constructor and expect its exact package-width failure |
| GW1/run-package-slot-one | omit or widen `-p=1` on run | argv matrix test | inspect the run constructor and expect its exact package-width failure |
| GW1/compile-package-slot-one | omit or widen `-p=1` on compile | argv matrix test | inspect the compile constructor and expect its exact package-width failure |
| GW1/test-test-slot-one | omit or widen `-parallel=1` on test | argv matrix test | inspect the test constructor and expect its exact test-width failure |
| GW1/compile-test-slot-one | omit or widen `-parallel=1` on compile | argv matrix test | inspect the compile constructor and expect its exact test-width failure |
| GW1/build-runtime-thread-one | pass a non-one runtime width to build | closed environment test | invoke build and expect one authoritative `GOMAXPROCS=1` |
| GW1/test-runtime-thread-one | pass a non-one runtime width to test | closed environment test | invoke test and expect one authoritative `GOMAXPROCS=1` |
| GW1/vet-runtime-thread-one | pass a non-one runtime width to vet | closed environment test | invoke vet and expect one authoritative `GOMAXPROCS=1` |
| GW1/run-runtime-thread-one | pass a non-one runtime width to run | closed environment test | invoke run and expect one authoritative `GOMAXPROCS=1` |
| GW1/compile-runtime-thread-one | pass a non-one runtime width to compile | closed environment test | invoke compile and expect one authoritative `GOMAXPROCS=1` |
| GW1/build-goflags-stripped | preserve widening `GOFLAGS` for build | hostile environment test | export widening flags and expect build to remain width one |
| GW1/test-goflags-stripped | preserve widening `GOFLAGS` for test | hostile environment test | export widening flags and expect test to remain width one |
| GW1/vet-goflags-stripped | preserve widening `GOFLAGS` for vet | hostile environment test | export widening flags and expect vet to remain width one |
| GW1/run-goflags-stripped | preserve widening `GOFLAGS` for run | hostile environment test | export widening flags and expect run to remain width one |
| GW1/compile-goflags-stripped | preserve widening `GOFLAGS` for compile | hostile environment test | export widening flags and expect compile to remain width one |
| GW1/build-gomaxprocs-stripped | append ambient `GOMAXPROCS` for build | hostile environment test | export a wider value and expect one authoritative build value of one |
| GW1/test-gomaxprocs-stripped | append ambient `GOMAXPROCS` for test | hostile environment test | export a wider value and expect one authoritative test value of one |
| GW1/vet-gomaxprocs-stripped | append ambient `GOMAXPROCS` for vet | hostile environment test | export a wider value and expect one authoritative vet value of one |
| GW1/run-gomaxprocs-stripped | append ambient `GOMAXPROCS` for run | hostile environment test | export a wider value and expect one authoritative run value of one |
| GW1/compile-gomaxprocs-stripped | append ambient `GOMAXPROCS` for compile | hostile environment test | export a wider value and expect one authoritative compile value of one |
| GW2/per-package-invocation | restore one `go test` call over the whole package slice | real core-runner spy | enumerate three packages and expect three ordered invocations |
| GW2/continue-after-red | return after the first red package | red-middle package runner test | make package two red and expect package three to run plus aggregate red |
| GW2/package-settlement-order | publish settlement before completion or out of enumeration order | package settlement recorder | block each package and expect one terminal event in exact order |
| GW2/cumulative-progress | retain the one-shot whole-run 45-minute context | fake-clock cumulative test | settle packages below the window while total exceeds it and expect no timeout |
| GW2/package-stall | reset the timer without a terminal settlement | fake-clock stalled-package test | hold one package past 45 minutes and expect code 124 plus teardown |
| GW2/heartbeat-ignored | treat child output as progress | heartbeat-spam timeout test | emit output without settlement for 45 minutes and expect timeout |
