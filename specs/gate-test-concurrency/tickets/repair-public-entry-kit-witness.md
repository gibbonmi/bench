# Repair the public-entry kit witness

Blocked by: none
Ownership fence: `internal/gate/kit_root_test.go`
Integration surfaces: `internal/gate/phases.go` public entry→`internal/gate/kit_root_test.go` + RKW1; `internal/gate/manifest.go` absent-manifest selection→`internal/gate/kit_root_test.go` + RKW1
Contracts: non-empty `BENCH_KIT` crosses the public Go entry into `internal/gate/kit_root_test.go` as one distinct path value; an absent manifest selects the built-in table; the capture observes the table's kit argument without replacing entry resolution
Closure: RKW1/distinct-kit-entry, RKW1/ignore-environment-red

## What to build

Close review finding P1 with one serial `PhasesCommand` test. Give it a graded
root with no phase manifest and a distinct non-empty `BENCH_KIT`, stub only the
built-in table provider, invoke the public Go entry, and assert the provider
receives the distinct kit. The test stays serial because it owns both process
environment and the package-level provider for its lifetime.

The wrapper test remains useful for shell routing, but it is not this witness:
its fake executable prints the environment without entering Go. The new test
must red when `kitRoot` is temporarily changed to ignore non-empty
`BENCH_KIT`, then return green after exact restoration.

## Acceptance

- [ ] [RKW1] (covers local) the real public Go entry with no manifest passes a distinct non-empty `BENCH_KIT` to built-in phase-table resolution, and an ignore-environment mutation makes the focused test red before exact restoration, repairing KC2/P1.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RKW1/distinct-kit-entry | make `kitRoot` always return the graded root | the new serial public-entry test | apply, run the focused test, expect the captured kit mismatch red; restore exact bytes and rerun green |
| RKW1/ignore-environment-red | bypass `kitRoot` at `PhasesCommand` and pass root directly | the same test | apply, run the focused test, expect the captured kit mismatch red |
