# Realign runtime commit fixtures with prospective authorization

Blocked by: adopt-exact-landing-in-commit.md
Ownership fence: `internal/contract/runtime/runtime_commit_reduced_test.go`, `internal/contract/runtime/runtime_gate_action_proof_test.go`, `internal/contract/runtime/runtime_gate_proof_helpers_test.go`, `internal/contract/runtime/runtime_gate_test.go`
Integration surfaces: retired-reduction gate tallies→`internal/contract/runtime/runtime_commit_reduced_test.go` plus PRF3; Story 5 result assertions and common-Git coordination→`internal/contract/runtime/runtime_gate_action_proof_test.go` plus PRF1; R2 and oracle-run marker setup→`internal/contract/runtime/runtime_gate_proof_helpers_test.go` plus PRF2; runtime gate consumer assertions→`internal/contract/runtime/runtime_gate_test.go` plus PRF1-PRF2
Contracts: common Git-directory marker paths and prospective refusal diagnostics cross fixture producers in `internal/contract/runtime/runtime_gate_proof_helpers_test.go` and `internal/contract/runtime/runtime_gate_action_proof_test.go`→their runtime proof consumers, asserted by PRF1-PRF4 against the real compiled command

Prerequisite repair for the exact prospective landing adapter. The adapter's
required scoped package floor exposed runtime fixtures that still exchange gate
markers through checkout-local `.git` paths or assert the former ambient-gate
diagnostic stage. A prospective checkout represents `.git` as a gitfile, and its
authorization deliberately excludes foreign ambient bytes.

## What to build

Move the affected runtime-only gate markers to the repository's common Git
directory and update the two Story 5 result assertions to the exact refusal stage
now owned by prospective authorization. Preserve the existing semantic claims:
the resolved full gate runs, oracle changes cannot reuse green, failures refuse
the commit, and the destination and checkout remain unchanged.

## Acceptance

- [ ] [PRF1] The R14 subject-build and start-failure rows assert the stage-specific prospective refusal actually returned by `bench commit`, retain their no-landing state checks, and complete green without weakening the standalone gate-action proofs.
- [ ] [PRF2] The R2 ignored-input and oracle-bound verdict fixtures count gate executions through the common Git directory, so their deliberate oracle/input changes prove a fresh prospective authorization run rather than observing only the ambient seed run.
- [ ] [PRF3] Both retired-reduction commit fixtures record their resolved full-gate executions through the common Git directory; confined and mixed named-path commits each show the seed run plus one prospective run, no reduced phase run, and the intended commit lands.
- [ ] [PRF4] The six-row repro and the full `internal/contract/runtime` package complete green with no production change and no generated verification artifact left behind.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PRF1 | restore either R14 result assertion to its pre-prospective gate-result diagnostic | existing R14 result row | run the six-row repro below and expect the changed R14 row to fail on the stage-specific refusal text |
| PRF2 | address one ignored-input or oracle-run marker through checkout-local `.git` again | existing R2 or oracle-bound fixture | run the six-row repro and expect its run count to remain at the ambient seed value |
| PRF3 | address `full-runs` through checkout-local `.git` again | existing retired-reduction fixtures | run the six-row repro and expect both commit rows to report only the seed `full` marker |
| PRF4 | restore any checkout-local marker path while running the complete runtime package | the six-row repro and full `internal/contract/runtime` package | run `go test -count=1 ./internal/contract/runtime`; expect the owning proof row to fail without leaving a generated verification artifact |

Focused repro:

```bash
go test -count=1 ./internal/contract/runtime -run '^(TestFT78Story5ProofLedger|TestFT78Story2ProofLedger|TestCommitMixedStagedSetRunsFullGate|TestCommitPaysFullGateForConfinedStagedSet|TestRuntimeGateContracts)$'
```
