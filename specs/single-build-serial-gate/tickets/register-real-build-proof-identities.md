# Register real-build proof identities

Blocked by: share-artifacts-across-local-processes.md, resolve-target-aware-artifact-identities.md
Ownership fence: `internal/artifactstore/`
Integration surfaces: proof backend→share-artifacts-across-local-processes.md; target-aware fields→resolve-target-aware-artifact-identities.md; closed `ProofRequest`→migrate-gate-real-build-proof-consumers.md; closed `ProofRequest`→migrate-contract-preflight-release-proof-consumers.md; closed `ProofRequest`→migrate-conformance-release-target-proof-consumers.md; closed `ProofRequest`→migrate-preprelease-install-proof-consumers.md; migrated inventory→route-real-build-proofs-through-registered-identities.md; proof registry contraction→enforce-the-executable-artifact-contract.md
Contracts: `ProofRequest` (proof class, exact source/package/bytes/authorship/target/failure fields, closed slot) crosses registered proof owners→private artifact backend, membership is a finite proof-class and authorship-slot registry, ordering is validate membership then resolve exact identity, absence or caller-invented fields refuse before backend execution, asserted by PR1 against the real registry
Closure: PR1/finite-class-registry, PR1/finite-slot-registry, PR1/no-caller-nonce, PR1/unknown-class-refusal

## What to build

Expand the artifact store with the closed typed request and registry needed by retained real-build proofs. Do not migrate any proof consumer in this ticket.

## Acceptance

- [ ] [PR1] (covers local) every retained proof kind and authorship event has a finite typed identity, while unknown classes, slots, and caller nonces refuse before backend execution.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PR1/finite-class-registry | omit one required proof class | proof registry membership test | enumerate the staged proof inventory and expect the omitted-class diagnostic |
| PR1/finite-slot-registry | omit one required authorship slot | proof registry membership test | enumerate the staged proof inventory and expect the omitted-slot diagnostic |
| PR1/no-caller-nonce | accept an arbitrary proof nonce | closed-request test | request an unknown nonce and expect refusal before the backend marker |
| PR1/unknown-class-refusal | accept an unknown proof class | closed-request test | request the unknown class and expect refusal before the backend marker |
