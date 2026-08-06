# Register the current communication markers

Blocked by: pin-structured-phase-clause-set.md
Ownership fence: `internal/anchors`, `internal/conformance/validity_checks_test.go`
Integration surfaces: marker definitions and registry rows→`internal/anchors`; `bench anchors` query and shared evaluator→existing unchanged `cmd/bench` and `internal/anchors/registry.go` + AM1; canonical-source duplicate consumer→`internal/conformance/validity_checks_test.go`; relocation/comment canary ownership→existing unchanged `internal/conformance/registry/registry.go`, `internal/conformance/registry_test.go`, and `tests/canary/workflow-guidance-anchors` + AM3
Contracts: three exported marker strings with membership `NEVER assume, always verify`, `Clear beats dense`, and `Right-size the process`, registry order Roles→How to talk to me→Workflow, and absence invalid cross `internal/anchors` registry rows→`internal/conformance/validity_checks_test.go`, asserted by AM1 and AM2 against `Entries`, the real evaluator, and the real duplicate-source consumer

## What to build

Expand the FT156 registry over the three communication markers as they exist before the prose cut. Define each marker once in `internal/anchors`, consume those exports from the registry and `checkSharedRuleSingleSource`, and independently pin all three current tuples so omitting any row is observable. This ticket is the green guardrail the following migration updates atomically with the Roles sentence.

## Acceptance

- [ ] [AM1] Registry enumeration contains exactly the current Roles, How to talk to me, and Workflow `RequireInSection` tuples, and omitting any one of the three fails the independent expectation.
- [ ] [AM2] Each marker has one production definition in `internal/anchors` consumed by registry evaluation and `checkSharedRuleSingleSource`, with no residual literal copy in `internal/conformance`.
- [ ] [AM3] `bench anchors .bench/BENCH.md` reports all three rows and the existing relocation/comment canaries retain their targeted bite.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AM1 | omit the Workflow registry row while leaving its exported marker and shared-rule consumer intact | the independent registry-tuple test | remove the row, run `go test ./internal/anchors`, expect the missing-tuple failure |
| AM2 | replace the `Clear beats dense` exported-symbol use in `checkSharedRuleSingleSource` with the literal | the coordinator residual-definition sweep | make the replacement, enumerate the exact literal and exported-symbol consumers with `rg`, expect the duplicate production definition to be rejected |
| AM3 | swap the Clear marker's owning section from How to talk to me to Roles | the real graded-root conformance check | change only the section, run `BENCH_CONFORMANCE_ROOT=$PWD go test ./internal/conformance -run '^TestRootConformance$' -count=1`, expect the attributed missing-section marker diagnostic |
