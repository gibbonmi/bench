# Compare all four public observations

Blocked by: capture-pinned-baseline.md, derive-complete-case-membership.md
Ownership fence: `internal/axi/compatibility`, `cmd/bench/axi_compatibility_test.go`
Integration surfaces: baseline capture→capture-pinned-baseline.md; complete cases→derive-complete-case-membership.md; exact byte classes→pin-exact-byte-classes.md; hostile processes→exercise-hostile-process-cases.md
Contracts: raw stdout bytes, raw stderr bytes, integer exit, boolean argv acceptance, and fresh-run identity cross `cmd/bench/axi_compatibility_test.go`→`internal/axi/compatibility`, ordered by stable case ID then run number, with no missing observation allowed, asserted by PC1
Closure: PC1/stdout, PC1/stderr, PC1/exit, PC1/acceptance, PC1/fresh-rerun

## What to build

every case compares raw stdout, raw stderr, exit, and accepted/rejected status twice from fresh state.

## Acceptance

- [ ] [PC1] (covers CO4) every case compares raw stdout, raw stderr, exit, and accepted/rejected status twice from fresh state.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PC1/stdout | change one candidate stdout byte | paired comparator test | run the case twice and require the stdout delta |
| PC1/stderr | move one candidate refusal from stderr to stdout | paired comparator test | run the refusal and require both stream deltas |
| PC1/exit | change one candidate exit code | paired comparator test | run the case and require the exit delta |
| PC1/acceptance | accept one formerly rejected argv spelling with unchanged bytes | paired comparator test | run the spelling and require the acceptance delta |
| PC1/fresh-rerun | reuse mutable state from the first run | fresh-state test | run in two fresh roots and require identical independent observations |

