# Route ordinary phase plumbing

Blocked by: introduce-run-scoped-bench-selection.md, own-gate-run-binary.md
Ownership fence: `internal/gate/`, `internal/conformance/gate_entry_test.go`, `internal/canary/gate_entry_test.go`
Integration surfaces: selected owner environment→own-gate-run-binary.md; selection API→introduce-run-scoped-bench-selection.md; selected helper migration→migrate-gate-helpers.md; selected suite environment→migrate-contract-preflight-helpers.md; nested phase propagation→propagate-selected-binary-to-nested-gates.md; serial scheduler consumer→serialize-phase-tables.md; closed constructor census→contract-ordinary-build-census.md
Contracts: selected executable argv crosses `internal/gate/phases.go`→gate entry, `GateGoCommand`, conformance, contract, shellcheck, and canary phases, membership is freshness, gate-phases, gofmt, test, race, conformance-suite, ordinary conformance, ordinary contract, shellcheck, and canary launch, ordering validates one selection before table construction and preserves phase dependencies after removing the build phase, asserted by PP1 against the real assembled phase table and shell entries
Closure: PP1/freshness-selected, PP1/gate-phases-selected, PP1/gofmt-selected, PP1/test-selected, PP1/race-selected, PP1/conformance-suite-selected, PP1/conformance-selected-env, PP1/contract-selected-env, PP1/shellcheck-selected-env, PP1/canary-selected-env, PP1/build-phase-removed, PP1/no-gate-go-run

## What to build

Route every ordinary phase command through the selected executable or its explicit environment. Remove the phase-owned Bench build and the `go run ./cmd/bench` gate-go route, updating the build-component declarations, phase expectations, and entry tests as one independently green landing; separating those declarations would strand the existing component-identity and manifest contracts red.

## Acceptance

- [ ] [PP1] (covers RS4) the real shell and phase table pass one identical selected executable through every enumerated ordinary consumer, contain no Bench build phase or gate-go `go run`, and retain the table's non-build dependencies.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PP1/freshness-selected | run freshness through a separate verifier | gate-entry argv test | assemble and run the shell entry and require the selected path's `freshness-check` marker |
| PP1/gate-phases-selected | route phase execution through wrapper or dist | gate-entry argv test | run the shell and require the same selected path's `gate-phases` marker |
| PP1/gofmt-selected | construct gofmt gate-go with `go run` | phase-table argv test | materialize gofmt and require selected executable argv |
| PP1/test-selected | construct test gate-go with `go run` | phase-table argv test | materialize test and require selected executable argv |
| PP1/race-selected | construct race gate-go with `go run` | phase-table argv test | materialize race and require selected executable argv |
| PP1/conformance-suite-selected | construct conformance-suite gate-go with `go run` | phase-table argv test | materialize the suite and require selected executable argv |
| PP1/conformance-selected-env | omit selection from ordinary conformance | phase child-env test | run the conformance marker and require exact selected path |
| PP1/contract-selected-env | omit selection from ordinary contract | phase child-env test | run the contract marker and require exact selected path |
| PP1/shellcheck-selected-env | omit selection from shellcheck | phase child-env test | run an installed shellcheck marker and require exact selected path |
| PP1/canary-selected-env | omit selection from canary launch | phase child-env test | run the canary marker and require exact selected path |
| PP1/build-phase-removed | restore the phase-owned subject builder | real phase-table test | enumerate the table and require no Bench build phase or builder argv |
| PP1/no-gate-go-run | retain `go run ./cmd/bench` in any ordinary gate-go constructor | assembled-argv census | enumerate all ordinary gate-go argv and require the selected executable as argv zero |
