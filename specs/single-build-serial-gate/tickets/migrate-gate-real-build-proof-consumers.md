# Migrate gate real-build proof consumers

Blocked by: register-real-build-proof-identities.md, migrate-gate-fixture-artifact-consumers.md
Ownership fence: `internal/gate/`
Integration surfaces: closed proof request→register-real-build-proof-identities.md; migrated gate fixture seams→migrate-gate-fixture-artifact-consumers.md; migrated gate proof inventory→route-real-build-proofs-through-registered-identities.md; legacy proof helper retirement→enforce-the-executable-artifact-contract.md
Contracts: registered gate proof identities cross gate proof owners→artifact-store proof backend, membership is every real compiler-observing proof under `internal/gate/`, ordering preserves each existing assertion while replacing only authorship, absence of a typed record refuses without direct-builder fallback, asserted by MG1 against the real gate proof inventory
Closure: MG1/changed-source, MG1/alternate-artifact, MG1/planted-bytes, MG1/prospective-execution, MG1/build-authorship, MG1/compiler-failure

## What to build

Migrate the real-build proof consumers under `internal/gate/` to the registered proof API while preserving their independently observed source, artifact, bytes, prospective, authorship, and compiler-failure behavior.

## Acceptance

- [ ] [MG1] (covers local) every gate-owned real-build proof resolves its registered identity and no gate proof retains an unregistered direct builder.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| MG1/changed-source | classify a changed-source proof as unchanged | existing freshness assertion | mutate source and expect its stale-digest failure |
| MG1/alternate-artifact | route the alternate package through the canonical identity | existing alternate-artifact assertion | request the alternate package and expect its distinct digest assertion to fail |
| MG1/planted-bytes | materialize canonical bytes for a planted-byte proof | existing planted-binary refusal | plant the hostile bytes and expect the refusal owner to observe them |
| MG1/prospective-execution | reuse the base artifact for unpublished source | existing prospective exact-tree test | change unpublished `cmd/bench` and expect the old-behavior failure |
| MG1/build-authorship | reuse a prior authorship slot | existing authored-at assertion | request the later event and expect its identity to differ |
| MG1/compiler-failure | satisfy broken source from cached success | existing compiler/linker failure canary | inject invalid source and expect a real compiler diagnostic |
