# Close gate-entry fixture dependency closure

Blocked by: none
Ownership fence: `internal/conformance/gate_entry_test.go`
Integration surfaces: freshness source-copy fixture→`internal/conformance/gate_entry_test.go`; `internal/artifactstore/digest` import closure→`internal/conformance/gate_entry_test.go`
Contracts: none crosses
Closure: GF1/digest-package-source

## What to build

The real gate-entry conformance fixture copies the complete production source closure needed by `internal/freshness/freshness.go`, so it reaches the gate-entry refusal assertions instead of failing while compiling its synthetic module.

## Acceptance

- [ ] [GF1] (covers local) both real gate-entry tests compile the synthetic kit after freshness imports the shared digest leaf and retain their existing refusal and replacement assertions.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| GF1/digest-package-source | omit `internal/artifactstore/digest/digest.go` from the gate-entry fixture source closure | the existing real gate-entry conformance tests | run `go test -count=1 ./internal/conformance -run '^TestGateEntry(RefusesUnverifiedBinaryBeforeGatePhases|RejectsLegacyBeforeRunningOldTableAndRunsReplacementOnce)$'` and expect the missing-package compile diagnostic before the gate-entry assertions |
