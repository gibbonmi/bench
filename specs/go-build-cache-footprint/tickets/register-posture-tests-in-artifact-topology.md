# Register the grown posture test list in the artifact topology registry

Blocked by: none
Ownership fence: `internal/contract/surface/artifact/internal/fixture/topology_test.go`
Integration surfaces: posture top-level test family→existing `internal/contract/surface/artifact/posture` + RT1; other subject lists→unchanged entries in the fenced registry
Contracts: the posture test-name enumeration crosses `internal/contract/surface/artifact/internal/fixture/topology_test.go`→the posture subject package, asserted by RT1 against the real package files

## What to build

Candidate gate repair for promote's prospective red. `TestSubjectPackageTopology`
(`internal/contract/surface/artifact/internal/fixture/topology_test.go`) pins an
executable registry (`artifactTests`) of each artifact surface package's
top-level test names. The composition grew the posture package from four tests
to eighteen — fourteen from the original build's mode and grammar tables, three
from the stub-relocation and seal-signal repairs — without registering them, so
the contract phase reds with the exact got/want difference. Register the
eighteen names the posture package now holds; every other subject's list and
every other assertion in the contract stay untouched. This is a registry sync
for deliberate, reviewed additions — not a loosening: the contract's job is to
make additions conscious, and these were.

## Acceptance

- [ ] [RT1] `TestSubjectPackageTopology` passes against the composed tree: the posture entry enumerates exactly the eighteen top-level tests the package holds, and the prepared, offline, and distributable entries are byte-for-byte unchanged.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RT1 | drop one newly registered posture test name from the registry | the topology contract over the real package files | run `go test -count=1 ./internal/contract/surface/artifact/internal/fixture`, expect the got/want difference naming the dropped test |
