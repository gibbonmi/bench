# Keep AXI conformance on its in-process seam

Blocked by: none
Ownership fence: `internal/conformance/axi_disclosure_test.go`, `internal/conformance/ordinary_build_census_test.go`, `internal/specbuild/disclosure_test.go`
Integration surfaces: production operation/refusal axes→existing `internal/spec` and `internal/specbuild` accessors consumed by `internal/conformance/axi_disclosure_test.go`; checked-in matrix outputs→existing `internal/specbuild/testdata/axi-cases.jsonl`; real-service observation and exit status→`internal/specbuild/disclosure_test.go`; ordinary architecture census→`internal/conformance/ordinary_build_census_test.go`
Contracts: every applicable production cell crosses `DisclosureCells`→its checked-in output in the process-free conformance owner, asserted by PB2; real-service bytes and class-derived exit status cross `ObserveDisclosureCell`→the checked-in output only in `internal/specbuild/disclosure_test.go`, asserted by PB3; the ordinary-test process rule crosses `projects/benchkit.md`→the architecture census, asserted by PB1
Closure: PB1/no-conformance-observation-call, PB2/applicable-cell-fixture-totality, PB2/help-envelope-shape, PB3/exact-real-service-fixture-bytes, PB3/class-exact-runtime-exit, PB3/usage-runtime-exit, PB3/single-real-service-observation-owner

## What to build

Close the accepted Terra Standards finding
P1-ordinary-conformance-process-repository. The ordinary AXI conformance check
derives the production axes and validates every checked-in applicable output,
help envelope, and help spelling without calling
`ObserveDisclosureCell`, constructing a repository, or starting Git. The
`internal/specbuild` owner becomes the one test that drives every applicable
real-service cell, compares its exact bytes to those same fixtures, and checks
the class-derived 0/1 exit plus a real usage/2 control. Removing the
conformance-layer process therefore does not remove the runtime oracle or
pretend that the byte-only fixture independently encodes exit status.

## Acceptance

- [ ] [PB1] (covers local) (P1-ordinary-conformance-process-repository) ordinary AXI conformance contains no call path to the process-backed `ObserveDisclosureCell` seam.
- [ ] [PB2] (covers SB7) (P1-ordinary-conformance-process-repository) process-free conformance still derives every applicable operation/refusal cell and requires a checked-in structured output with a terminal help envelope for each cell.
- [ ] [PB3] (covers SB2) (P1-ordinary-conformance-process-repository) `internal/specbuild/disclosure_test.go` is the single real-service observation owner; it compares every applicable cell's exact public bytes with `matrix/<operation>/<class>`, requires success cells to exit 0 and refusal cells to exit 1, and exercises a real usage control that exits 2.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PB1/no-conformance-observation-call | restore the `ObserveDisclosureCell` call in ordinary conformance | the ordinary architecture census | scan ordinary conformance production and test functions, require no reference to the process-backed observation entry, and require the candidate with the restored call to fail naming the file and symbol |
| PB2/applicable-cell-fixture-totality | delete one applicable cell's checked-in matrix output | the AXI disclosure conformance check | derive `DisclosureCells`, load the JSONL outputs, and require the omitted operation/class key to be named |
| PB2/help-envelope-shape | remove the terminal `help[]` block from one applicable checked-in output | the AXI disclosure conformance check | mutate one loaded output and require the exact operation/class to fail structured-output validation |
| PB3/exact-real-service-fixture-bytes | change one applicable checked-in matrix output while leaving the real service unchanged | the specbuild real-service observation test | load `matrix/<operation>/<class>`, observe that cell through the real service, compare exact public bytes, and require the changed fixture to fail naming the cell |
| PB3/class-exact-runtime-exit | make one observed refusal return 0 or one success return 1 | the specbuild real-service observation test | derive expected 0/1 from the production refusal class, observe every applicable cell, and require the mutated cell's actual exit to fail |
| PB3/usage-runtime-exit | make the real usage control return any exit other than 2 | the specbuild real-service observation test | invoke the public service with invalid usage and require actual exit 2 |
| PB3/single-real-service-observation-owner | add a second `ObserveDisclosureCell` consumer outside `internal/specbuild/disclosure_test.go` | the ordinary architecture census | enumerate every production call site of the observation entry and require the sole consumer to be the specbuild real-service owner test |
