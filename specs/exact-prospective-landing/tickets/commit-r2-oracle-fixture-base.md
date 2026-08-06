# Commit the R2 oracle fixture base

Blocked by: realign-runtime-commit-fixtures-with-prospective-authorization.md
Ownership fence: `internal/contract/runtime/runtime_gate_proof_helpers_test.go`
Integration surfaces: ignored-input fixture subject construction→`internal/contract/runtime/runtime_gate_proof_helpers_test.go`; R2 proof-ledger consumer→existing `internal/contract/runtime/runtime_gate_proof_test.go` plus R2F1-R2F3
Contracts: committed gate and manifest bytes plus the ignored declared input cross `internal/contract/runtime/runtime_gate_proof_helpers_test.go` fixture setup→the real compiled `bench commit` invocation, asserted by R2F1-R2F2 against the R2 proof-ledger consumer

Prerequisite repair for the exact prospective landing adapter. The R2 ignored-input
fixture rewrites an ignored gate script and manifest but commits only `work.txt`.
An ambient `bench gate` can see those uncommitted oracle files; the exact prospective
checkout correctly starts from the expected commit and therefore has no gate to run.

## What to build

Make the R2 expected base self-contained by force-adding only its ignored gate script
and gate-input manifest before the fixture commit. Keep `ft78-ignored` ignored and
uncommitted: the standalone seed gate must see its green ambient bytes, while the
prospective commit gate must exclude those foreign ignored bytes and refuse.

## Acceptance

- [ ] [R2F1] The ignored-input fixture's expected commit contains its executable gate and manifest but not `ft78-ignored`; the standalone seed gate runs green against ambient `green` bytes.
- [ ] [R2F2] After the ignored input changes to `red` and named `work.txt` changes, `bench commit` runs the prospective gate, refuses without moving `HEAD`, and the common-directory tally records exactly the seed plus prospective runs.
- [ ] [R2F3] The focused R2 row, the six-row prospective fixture repro, and full `internal/contract/runtime` complete green with no production or adapter changes.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| R2F1 | omit the forced gate and manifest addition from the fixture base | existing R2 ignored-input refusal proof | run `go test -count=1 ./internal/contract/runtime -run '^TestFT78Story2ProofLedger$/^R2$/^ignored-input-refusal$' -v`; expect prospective authorization to find no gate |
| R2F2 | address the gate-run tally through checkout-local `.git` instead of the common Git directory | existing R2 ignored-input refusal proof | run `go test -count=1 ./internal/contract/runtime -run '^TestFT78Story2ProofLedger$/^R2$/^ignored-input-refusal$' -v`; expect the tally to remain at the ambient seed value |
| R2F3 | omit the forced gate and manifest addition while running the complete runtime package | `TestFT78Story2ProofLedger` and the runtime package suite | run `go test -count=1 ./internal/contract/runtime`; expect the R2 proof ledger to fail without leaving a generated verification artifact |

Focused repro:

```bash
go test -count=1 ./internal/contract/runtime -run '^TestFT78Story2ProofLedger$/^R2$/^ignored-input-refusal$' -v
```
