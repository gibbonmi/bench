# Pass one go test timeout from the bounds owner

Blocked by: none
Writes: internal/bounds/bounds.go, internal/bounds/bounds_test.go, internal/gate/phases.go, internal/gate/gate_go.go, internal/gate/lane_test.go, internal/gate/branch_native_phases_test.go, internal/conformance/ordinary_build_census_test.go, internal/conformance/bounds_policy_test.go, internal/preprelease/preprelease_test.go, internal/testreport/testreport_test.go, internal/testreport/check_test.go, internal/testreport/root_conformance_test.go, projects/benchkit.md, .bench/commitment.json
Covers: TC1, TC2, TC3, TC4, TC5, TC6, TC7, TC8, TC9, TC10, TC11, TC38

## What to build

Every Bench-owned go test command passes one per-binary timeout from one constant. Declare `GoTestTimeout` in `internal/bounds` beside `GateTimeout`, with the value of ten minutes. `gate.BaseTestArgv` appends the token `-timeout=` plus `bounds.GoTestTimeout.String()` right after `-count=1` and before the caller operands. The six callers that `bench consumers gate.BaseTestArgv` names inherit the token, so no caller adds its own copy.

This ticket supplies one value to its successors: the exported constant `bounds.GoTestTimeout`. TC-C3 reads it for the warning threshold and the `timeout_ms` cell.

Each exact-argv test builds the token from `bounds.GoTestTimeout` by its own expression. `testreport_test.go` declares one helper that returns the token. The argv lines of `check_test.go` and `root_conformance_test.go` call that helper, so `check_test.go`, which is over its line limit, does not grow.

`TestProductionPolicyValues` pins `GoTestTimeout` as ten minutes, beside its sibling values. That pin is the only literal timeout in the tests. Record the mutation `GoTestTimeout = 5 * time.Minute` and the observed red of that test in the commit or the review evidence.

The bounds policy owners map gains the entry `internal/gate/phases.go` with the token `bounds.GoTestTimeout`. The profile phase table in `projects/benchkit.md` shows the token in the `test`, `race`, and `system` rows.

The ticket names `.bench/commitment.json` only because the spec-stage close commits the deliverable binding inside the graded range. The author does not edit that file.

## Acceptance

- [ ] The release core, gate `test`, gate `system`, gate `race`, ship conformance, and focused `bench test` argv each carry the timeout token right after `-count=1`.
- [ ] Each named check child, the system check child, and the root conformance check child receive the timeout token before `-json`.
- [ ] `bench test --check bounds-policy` turns red when `internal/gate/phases.go` does not contain `bounds.GoTestTimeout`.
- [ ] `TestProductionPolicyValues` turns red under the mutation `GoTestTimeout = 5 * time.Minute`, and the ticket evidence records that red.
- [ ] No test of this ticket other than `TestProductionPolicyValues` spells the timeout as a literal.
- [ ] `bench structure --growth` stays green, and `check_test.go` keeps its line count.
