# Move canary and stripped proofs

Blocked by: 02-assemble-branch-native-gate-drivers.md, 03-own-bounded-system-journeys.md
Ownership fence: `internal/canary/`, `internal/conformance/`, `internal/gate/`, `internal/systemtest/`, `tests/canary/`
Integration surfaces: canary fixture registry and owning check registry→`internal/canary/`; top-level canary dispatch journey→`internal/systemtest/`; stripped source inventory and selected executable→`internal/systemtest/`; phase argv→`internal/gate/`
Contracts: every canary fixture identity crosses `tests/canary/`→one production owning check with complete membership and no duplicate owner, asserted by CM1 against the real registries; stripped package inventory and selected executable identity cross the stripped producer→one system journey with absent excluded paths, asserted by ST1 and ST2
Closure: CM1/complete-ownership, CM1/one-owner, CM1/specific-red, CM1/restored-green, CM1/no-inner-gate, CM1/no-nested-go, CT1/top-level-dispatch, CT1/aggregation, ST1/one-materialization, ST1/no-ordinary-rerun, ST2/package-shape, ST2/wrapper-resolution, ST2/selected-binary, ST2/excluded-path-refusal

## What to build

Make canary fixtures direct mutation inputs to their registered production checks, retaining one top-level real-binary dispatch journey. Replace the stripped phase schedule with one selected-binary distribution smoke journey over one stripped subject.

## Acceptance

- [x] [CM1] (covers NC1) every canary mutation has exactly one registered production owner, produces its mutation-specific red, restores to green, and starts no inner gate or nested Go test.
- [x] [CT1] (covers NC2) one bounded system journey proves top-level `bench canary` dispatch and aggregation while ordinary canary checks remain in-process.
- [x] [ST1] (covers ST1) one stripped subject is materialized for one system smoke journey and no contract, conformance, or ordinary package suite reruns against it.
- [x] [ST2] (covers ST2) the stripped journey proves installed package shape, wrapper resolution, selected-binary identity, and excluded-path dependency refusal.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CM1/complete-ownership | add a fixture without an owner | canary registry totality test | mutate the fixture inventory, run the direct test, expect missing-owner diagnostic |
| CM1/one-owner | assign a second owner to a fixture | canary registry totality test | mutate the owner registry, run the direct test, expect duplicate-owner diagnostic |
| CM1/specific-red | make an owning check always green | direct mutation test | mutate the check, run the fixture mutation, expect missing targeted red |
| CM1/restored-green | leave mutated input in place | direct mutation test | omit restoration, rerun the owner, expect restored-green failure |
| CM1/no-inner-gate | add a gate invocation to a canary owner | architecture census | add the constructor, run the census, expect inner-gate diagnostic |
| CM1/no-nested-go | add a nested Go test to a canary owner | architecture census | add the constructor, run the census, expect nested-Go diagnostic |
| CT1/top-level-dispatch | bypass canary command dispatch | system canary journey | mutate dispatch, run the tagged journey, expect missing selected command observation |
| CT1/aggregation | discard one owner failure | system canary journey | mutate aggregation, run the tagged journey, expect missing failure diagnostic |
| ST1/one-materialization | request a second stripped repository | system repository budget | mutate the journey plan, run the tagged suite, expect stripped count diagnostic |
| ST1/no-ordinary-rerun | add a package-universe command to stripped journey | phase argv and architecture census | mutate the journey, run focused checks, expect forbidden rerun diagnostic |
| ST2/package-shape | retain a forbidden package file | stripped journey | mutate the inventory, run the tagged journey, expect package-shape failure |
| ST2/wrapper-resolution | route the wrapper to the source checkout | stripped journey | mutate the wrapper, run the tagged journey, expect route mismatch |
| ST2/selected-binary | copy or build a second executable | system identity ledger | mutate the journey, run the tagged suite, expect selected-binary mismatch |
| ST2/excluded-path-refusal | let the stripped subject read an excluded path | stripped journey | mutate the command input, run the tagged journey, expect excluded dependency accepted |
