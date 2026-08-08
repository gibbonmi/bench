# Remove the prospective constructor carrier

Blocked by: repair-unused-prospective-kit-seam.md, repair-parallelism-evidence.md
Ownership fence: `internal/gate/evaluation.go`, `specs/gate-test-concurrency/implementation-evidence.md`
Integration surfaces: `internal/gate/evaluation.go` prospective constructor→existing `internal/gate/evaluation_test.go` + RPC1; RPS1 mutation ledger→`specs/gate-test-concurrency/implementation-evidence.md` + RPC1
Contracts: prospective evaluation in `internal/gate/evaluation.go` carries runtime root, identity root, and tree identity but no kit identity because prospective scoping returns before kit-dependent evaluation, asserted by RPC1
Closure: RPC1/prospective-constructor, RPC1/evidence-accuracy

## What to build

Delete `newProspectiveTreeEvaluationAtKit` and construct the prospective
evaluation directly in `newProspectiveTreeEvaluation` without assigning a kit.
Widen the zero-unused-carrier audit across the complete prospective evaluation
path, and correct the implementation ledger so RPS1 claims only the engine seam
it actually proved while this repair records the constructor proof.

## Acceptance

- [ ] [RPC1] (covers local) the prospective path has no private kit-bearing constructor or stored kit identity, and the exact-candidate evidence attributes both unused-carrier repairs to the audit each one actually proved, closing S4 and C4.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RPC1/prospective-constructor | restore a four-argument prospective constructor that stores the forwarded kit without any reachable consumer | the complete prospective-carrier source audit | apply, run the audit across `internal/gate`, expect it to name the restored declaration, forwarding call, and stored field while compile-only remains green |
| RPC1/evidence-accuracy | leave RPS1 claiming zero carrier hits across the prospective path after the restored constructor is detected | the exact mutation ledger audit | apply, compare the carrier audit result with the RPS1 and RPC1 ledger rows, expect the overclaim or missing repair attribution to be named |
