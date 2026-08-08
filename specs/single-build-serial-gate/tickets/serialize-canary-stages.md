# Serialize canary stages

Blocked by: serialize-gate-lineages-and-transfer-turns.md
Ownership fence: `internal/canary/canary.go`, `internal/canary/canary_test.go`, `internal/canary/canary_concurrency_test.go`, `internal/canary/compiled_bite_test.go`, `internal/canary/runner_junction_test.go`, `internal/canary/abort_test.go`, `internal/bounds/bounds.go`, `internal/conformance/bounds_policy_test.go`, `tests/canary/package-core-guard/bounds-canary-width-unconsumed/`, `tests/canary/package-core-guard/bounds-duplicate-canary-width/`
Integration surfaces: transferable turn descriptor→serialize-gate-lineages-and-transfer-turns.md; canary compile/baseline/bite settlement→iterator and tests in `internal/canary/`; canary-item settlement/watchdog junction→close-all-resource-concurrency-routes.md; retired worker/width advertisements→`internal/conformance/bounds_policy_test.go` and the owned canary fixtures; final resource enrollment→close-all-resource-concurrency-routes.md
Contracts: `CanaryItemSettlement` (stage, package/group/fixture identity, terminal status) crosses iterator→run progress owner in `internal/canary/`, membership is every selected compile then baseline then bite, ordering is deterministic stage order with each prerequisite settled first, absence refuses or records red and no worker fallback exists, asserted by CS1-CS3 against the real selected inventory
Closure: CS1/compile-serial, CS1/baseline-serial, CS1/bite-serial, CS1/stage-order, CS1/error-order, CS1/eachindex-removed, CS1/fixtureworkers-removed, CS1/inner-width-code-removed, CS1/inner-width-conformance-removed, CS2/width-unconsumed-fixture-retired, CS2/duplicate-width-fixture-retired, CS3/compile-package-slot-one, CS3/compile-test-slot-one, CS3/compile-runtime-thread-one, CS3/compile-goflags-stripped, CS3/compile-duplicate-gomaxprocs-stripped

## What to build

Replace canary's worker pool with one deterministic item iterator, preserve compile→baseline→bite ordering and sorted error reporting, and retire `fixtureWorkers` plus `CanaryInnerWidth` from code, tests, and bounds policy.

## Acceptance

- [ ] [CS1] (covers ZC3) canary compiles, baselines, and bites one item at a time in deterministic prerequisite/error order, with worker and inner-width policy absent.
- [ ] [CS2] (covers local) the two canaries that advertise or mutate `CanaryInnerWidth` are retired or replaced by serial-policy mutations that still bite.
- [ ] [CS3] (covers local) the canary `go test -c` child carries `-p=1`, `-parallel=1`, and exactly one `GOMAXPROCS=1` regardless of ambient width.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CS1/compile-serial | launch two package compiles from workers | blocking fake runner | select two packages and expect maximum active compiles one |
| CS1/baseline-serial | launch two baselines concurrently | blocking fake runner | select two groups and expect maximum active baselines one |
| CS1/bite-serial | launch two fixture bites concurrently | blocking fake runner | select two fixtures and expect maximum active bites one |
| CS1/stage-order | start a bite before its compile/baseline settles | stage recorder test | block the prerequisite and expect no bite marker |
| CS1/error-order | report errors in completion order | existing sorted-error test | return reversed failures and expect stable fixture order |
| CS1/eachindex-removed | retain the `eachIndex` worker helper | structural canary audit | enumerate helper symbols and expect no `eachIndex` owner or call |
| CS1/fixtureworkers-removed | retain the `fixtureWorkers` policy | structural canary audit | enumerate width symbols and expect no `fixtureWorkers` owner or call |
| CS1/inner-width-code-removed | leave `CanaryInnerWidth` in bounds code | structural bounds audit | inspect exported bounds and expect the retired-symbol failure |
| CS1/inner-width-conformance-removed | leave `CanaryInnerWidth` in conformance policy | bounds policy test | inspect the policy registry and expect the retired-bound failure |
| CS2/width-unconsumed-fixture-retired | retain the fixture mutation that requires `bounds.CanaryInnerWidth` consumption | canary fixture inventory test | apply the fixture family to the serialized source and expect no stale old-literal or retired-policy expectation |
| CS2/duplicate-width-fixture-retired | retain the duplicate-width fixture after the bound is removed | canary fixture inventory test | apply the fixture family and expect its replacement serial-policy mutation to bite |
| CS3/compile-package-slot-one | omit or widen `-p=1` on `go test -c` | canary compile argv test | inspect the real compile constructor and expect the exact package-width failure |
| CS3/compile-test-slot-one | omit or widen `-parallel=1` on `go test -c` | canary compile argv test | inspect the real compile constructor and expect the exact test-width failure |
| CS3/compile-runtime-thread-one | pass a value other than one in the compile child environment | canary compile environment test | invoke the constructor and expect one authoritative `GOMAXPROCS=1` entry |
| CS3/compile-goflags-stripped | preserve ambient widening `GOFLAGS` | hostile canary compile environment test | export widening flags and expect the compile argv/environment to remain width one |
| CS3/compile-duplicate-gomaxprocs-stripped | append rather than replace ambient `GOMAXPROCS` | hostile canary compile environment test | export a wider value and expect exactly one authoritative value of one |
