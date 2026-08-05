# Exercise pre-lock faults after reuse

Blocked by: reuse-exact-green-before-gate-lock.md
Ownership fence: `internal/gate/runner_test.go`
Integration surfaces: pre-lock reusable-green behavior→reuse-exact-green-before-gate-lock.md; private R17 tuple→existing `internal/contract/runtime/runtime_gate_partial_proof_test.go` plus PF1; Story 5 proof registry→existing `internal/contract/runtime/runtime_gate_action_registry_test.go` plus PF1
Contracts: the private fault tuple with operation hits, reuse result, durable state, attempts, and temporary-file count crosses the compiled `internal/gate/runner_test.go` bridge→the runtime R17 proof consumer, asserted by PF1 against the real compiled bridge binary

## What to build

Drive the lock-open and lock-acquisition fault injections with non-reusable retained evidence now that exact green reuse occurs before lock access, while preserving their one-attempt failure posture, repeated-call proof, interrupted-pending durability, and zero temporary evidence.

## Acceptance

- [ ] [PF1] The compiled R17 private bridge reaches both lock-open and lock-acquisition injections on each call after pre-lock reuse is enabled, reports one operation attempt per call, returns no reusable green, and leaves interrupted-pending state with no temporary evidence; the two public runtime proof rows and the full `internal/gate` package are green.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PF1 | run the pre-lock fault pair at the reusable seed clock again | existing private-bridge tuple and public R17 rows | run `go test -count=1 ./internal/contract/runtime -run '^TestFT78Story5ProofLedger$/^R17$/^fault-lock-(open|acquisition)$' -v`, expect both rows to fail with `reusable:true failed:false hits:0` |
