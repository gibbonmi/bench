# Register the current communication markers

Blocked by: pin-structured-phase-clause-set.md
Ownership fence: `internal/anchors`, `internal/conformance/validity_checks_test.go`
Integration surfaces: marker definitions and registry rows→`internal/anchors`; `bench anchors` query and shared evaluator→existing unchanged `cmd/bench` and `internal/anchors/registry.go` + AM1; canonical-source duplicate consumer→`internal/conformance/validity_checks_test.go`; relocation/comment canary ownership→existing unchanged `internal/conformance/registry/registry.go`, `internal/conformance/registry_test.go`, and `tests/canary/workflow-guidance-anchors` + AM3
Contracts: three exported marker strings with membership `NEVER assume, always verify`, `Clear beats dense`, and `Right-size the process`, registry order Roles→How to talk to me→Workflow, and absence invalid cross `internal/anchors` registry rows→`internal/conformance/validity_checks_test.go`, asserted by AM1 and AM2 against `Entries`, the real evaluator, and the real duplicate-source consumer
Closure: AM1/roles-tuple, AM1/clear-tuple, AM1/workflow-tuple, AM1/registry-order, AM2/roles-single-source, AM2/clear-single-source, AM2/workflow-single-source, AM3/query-enumeration, AM3/relocation-canary, AM3/comment-canary

## What to build

Expand the FT156 registry over the three communication markers as they exist before the prose cut. Define each marker once in `internal/anchors`, consume those exports from the registry and `checkSharedRuleSingleSource`, and independently pin all three current tuples so omitting any row is observable. This ticket is the green guardrail the following migration updates atomically with the Roles sentence.

## Acceptance

- [ ] [AM1] Registry enumeration contains exactly the current Roles, How to talk to me, and Workflow `RequireInSection` tuples, and omitting any one of the three fails the independent expectation.
- [ ] [AM2] Each marker has one production definition in `internal/anchors` consumed by registry evaluation and `checkSharedRuleSingleSource`, with no residual literal copy in `internal/conformance`.
- [ ] [AM3] `bench anchors .bench/BENCH.md` reports all three rows and the existing relocation/comment canaries retain their targeted bite.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AM1/roles-tuple | omit the Roles registry row | the independent registry-tuple test | remove the row, run `go test ./internal/anchors`, expect the missing-Roles-tuple failure |
| AM1/clear-tuple | omit the Clear registry row | the independent registry-tuple test | remove the row, run `go test ./internal/anchors`, expect the missing-Clear-tuple failure |
| AM1/workflow-tuple | omit the Workflow registry row | the independent registry-tuple test | remove the row, run `go test ./internal/anchors`, expect the missing-Workflow-tuple failure |
| AM1/registry-order | swap two FT187 rows | the independent registry-tuple test | reorder the rows, run `go test ./internal/anchors`, expect the tuple-order failure |
| AM2/roles-single-source | replace the Roles exported-symbol use with its literal | the residual-definition sweep | make the replacement, enumerate literal and symbol consumers with `rg`, expect the duplicate production definition |
| AM2/clear-single-source | replace the Clear exported-symbol use with its literal | the residual-definition sweep | make the replacement, enumerate literal and symbol consumers with `rg`, expect the duplicate production definition |
| AM2/workflow-single-source | replace the Workflow exported-symbol use with its literal | the residual-definition sweep | make the replacement, enumerate literal and symbol consumers with `rg`, expect the duplicate production definition |
| AM3/query-enumeration | drop one FT187 row from `Entries` | the real `bench anchors` query | run `bench anchors .bench/BENCH.md`, expect the row count and named tuple to be absent from the required enumeration |
| AM3/relocation-canary | swap the Clear marker's owning section from How to talk to me to Roles | the real graded-root conformance check | change only the section, run `BENCH_CONFORMANCE_ROOT=$PWD go test ./internal/conformance -run '^TestRootConformance$' -count=1`, expect the attributed missing-section marker diagnostic |
| AM3/comment-canary | wrap one registered marker in an HTML comment | the existing commented-anchor canary | mutate the subject, run the owning fixture test, expect the required-anchor diagnostic |
