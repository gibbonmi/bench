# Route real-build proofs through registered identities

Blocked by: register-real-build-proof-identities.md, migrate-gate-real-build-proof-consumers.md, migrate-contract-preflight-release-proof-consumers.md, migrate-conformance-release-target-proof-consumers.md, migrate-preprelease-install-proof-consumers.md
Ownership fence: `internal/artifactstore/`, `internal/conformance/`
Integration surfaces: closed proof registry→register-real-build-proof-identities.md; migrated gate proofs→migrate-gate-real-build-proof-consumers.md; migrated contract/preflight/release proofs→migrate-contract-preflight-release-proof-consumers.md; migrated release-target proofs→migrate-conformance-release-target-proof-consumers.md; migrated preprelease/install routes→migrate-preprelease-install-proof-consumers.md; closed proof inventory→enforce-the-executable-artifact-contract.md
Contracts: the finite `ProofRequest` registry crosses every migrated proof owner→private backend in `internal/artifactstore/`, membership is the union of the four migration inventories, ordering preserves independent authorship per registered slot, absence of an inventory row or presence of a raw proof builder is structural red, asserted by RP1 against real compiler-observing assertions and the closed census
Closure: RP1/changed-source, RP1/alternate-artifact, RP1/planted-bytes, RP1/prospective-execution, RP1/build-authorship, RP1/release-target, RP1/reproducibility-slots, RP1/compiler-failure, RP1/finite-slot-registry, RP1/no-caller-nonce

## What to build

Contract the completed proof migrations into one closed inventory. Prove every retained real build uses a typed identity, every existing compiler-observing assertion still bites, and no unregistered proof route remains.

## Acceptance

- [ ] [RP1] (covers SB5) changed-source, alternate-artifact, planted-byte, prospective, authorship, release-target, reproducibility, and compiler-failure proofs build exactly their finite registered identities and no caller can invent an extra slot.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RP1/changed-source | classify a changed-source proof as the unchanged canonical identity | existing freshness changed-source assertion | mutate the source, run its focused proof, and expect the stale-digest failure |
| RP1/alternate-artifact | route the alternate package through the canonical package identity | existing alternate-artifact assertion | request the alternate package and expect its distinct executable/digest assertion to fail |
| RP1/planted-bytes | materialize canonical bytes for a planted-byte proof | existing planted-binary refusal | plant the registered hostile bytes and expect the refusal owner to observe them |
| RP1/prospective-execution | reuse the base artifact for unpublished prospective source | existing prospective exact-tree test | change unpublished `cmd/bench`, execute prospective preparation, and expect the old-behavior failure |
| RP1/build-authorship | reuse a prior authorship slot | existing attestation/authored-at assertion | request the later authoring event and expect its digest/time identity to differ |
| RP1/release-target | compute a release request without target tuple | existing cross-target artifact assertion | build two registered targets and expect distinct target records/bytes |
| RP1/reproducibility-slots | collapse `first` and `second` into one slot | existing reproducibility comparison | request both registered slots and expect two independent backend events before byte comparison |
| RP1/compiler-failure | satisfy the broken-source canary from a cached success | existing compiler/linker failure canary | inject invalid source and expect the real compiler diagnostic rather than an artifact |
| RP1/finite-slot-registry | omit a required authorship slot from the composed registry | proof registry membership test | enumerate the migrated proof assertions and expect the omitted slot diagnostic |
| RP1/no-caller-nonce | reintroduce a caller-defined nonce at one migrated site | closed-consumer audit | enumerate proof request construction and expect the invented-field diagnostic |
