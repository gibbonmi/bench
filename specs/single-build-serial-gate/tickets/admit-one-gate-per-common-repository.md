# Admit one gate per common repository

Blocked by: prepare-gate-artifacts-before-scheduling.md, enforce-the-executable-artifact-contract.md
Ownership fence: `.bench/gate.sh`, `.bench/gate-prospective.sh`, `internal/gate/gate.go`, `internal/gate/admission.go`, `internal/gate/admission_test.go`
Integration surfaces: prepared gate entry→prepare-gate-artifacts-before-scheduling.md; admitted gate lineage→serialize-gate-lineages-and-transfer-turns.md; repository-common admission→close-all-resource-concurrency-routes.md
Contracts: `GateAdmission` (common Git directory, owner identity, acquired state) crosses public gate entry→`internal/gate/admission.go`, membership is one active gate per common Git directory, ordering is acquire before verifier/artifact preparation and release after descendant reap, absence or contention refuses without waiting or work, asserted by GA1 against the two-worktree process harness
Closure: GA1/common-gitdir-lock, GA1/prework-refusal

## What to build

Add the repository-common admission lock as its own green slice after executable preparation exists. It touches no turn-transfer state and strands no scheduler red, so it does not belong in the lineage ticket.

## Acceptance

- [ ] [GA1] (covers ZC9) two worktrees sharing one common Git directory admit one gate and refuse the second before verifier, artifact, store, or phase work.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| GA1/common-gitdir-lock | key admission by worktree Git directory | two-worktree gate harness | start one held gate and expect the second to refuse without artifact or phase markers |
| GA1/prework-refusal | acquire admission after verifier or artifact preparation | two-worktree prework marker test | hold the first gate, start the second, and expect no verifier, store, backend, or phase marker |
