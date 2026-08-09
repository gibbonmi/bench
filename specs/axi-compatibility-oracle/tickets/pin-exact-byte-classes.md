# Pin exact byte and empty classes

Blocked by: compare-four-observations.md
Ownership fence: `internal/axi/compatibility`, `cmd/bench/axi_compatibility_test.go`, `specs/axi-compatibility-oracle/testdata`
Integration surfaces: four-observation comparator→compare-four-observations.md; helper classes→`decisions/byte-preserving-axi-foundation/assets/ft173-helper-compatibility-census.md`; public producers→existing renderers exercised unchanged
Contracts: case-pair IDs and raw bytes cross helper census→`specs/axi-compatibility-oracle/testdata`, membership is every default/full, empty, cap, control, and newline class, order follows stable IDs, and absence refuses, asserted by EB1
Closure: EB1/default-full, EB1/empty-classes, EB1/bound-edges, EB1/control-bytes, EB1/final-newline

## What to build

every helper-census default/full pair, empty class, bounds edge, control byte, and final newline remains exact.

## Acceptance

- [ ] [EB1] (covers CO5) every helper-census default/full pair, empty class, bounds edge, control byte, and final newline remains exact.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| EB1/default-full | change one full-mode producer while default remains exact | exact compatibility test | run the paired case IDs and require the full delta |
| EB1/empty-classes | normalize spec-build one-row empty to zero rows | exact compatibility test | run the empty case and require raw output mismatch |
| EB1/bound-edges | change one owner cap input | exact compatibility test | run below/at/above case IDs and require the boundary delta |
| EB1/control-bytes | drop one producer control escape | exact compatibility test | run the hostile byte case and require raw mismatch |
| EB1/final-newline | trim one producer final newline | exact compatibility test | run the case and require byte-length mismatch |

