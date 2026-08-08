# Enforce the executable artifact contract

Blocked by: share-artifacts-across-local-processes.md, resolve-target-aware-artifact-identities.md, prepare-gate-artifacts-before-scheduling.md, migrate-gate-fixture-artifact-consumers.md, migrate-contract-and-preflight-artifact-consumers.md, register-real-build-proof-identities.md, migrate-gate-real-build-proof-consumers.md, migrate-contract-preflight-release-proof-consumers.md, migrate-conformance-release-target-proof-consumers.md, migrate-preprelease-install-proof-consumers.md, route-real-build-proofs-through-registered-identities.md
Ownership fence: `internal/artifactstore/`, `internal/conformance/`, `internal/gate/`, `internal/contract/`, `internal/preflight/`, `internal/preprelease/`, `internal/releaseevidence/registry.json`, `.bench/gate.sh`, `.bench/gate-prospective.sh`, `.github/workflows/native-runtime.yml`, `package.json`, `scripts/go-build.sh`, `scripts/build-artifacts.sh`, `scripts/release-preflight.sh`, `scripts/native-proof.sh`, `scripts/build-offline-archives.sh`, `scripts/gen-platform-packages.sh`
Integration surfaces: temporary legacy-builder coexistence→share-artifacts-across-local-processes.md; target request closure→resolve-target-aware-artifact-identities.md; prepared entry→prepare-gate-artifacts-before-scheduling.md; ordinary gate migration→migrate-gate-fixture-artifact-consumers.md; ordinary contract/preflight migration→migrate-contract-and-preflight-artifact-consumers.md; proof registry→register-real-build-proof-identities.md; gate proof migration→migrate-gate-real-build-proof-consumers.md; contract/preflight/release proof migration→migrate-contract-preflight-release-proof-consumers.md; conformance release-target migration→migrate-conformance-release-target-proof-consumers.md; preprelease/install migration→migrate-preprelease-install-proof-consumers.md; composed proof closure→route-real-build-proofs-through-registered-identities.md; closed executable contract→admit-one-gate-per-common-repository.md, serialize-gate-lineages-and-transfer-turns.md, and close-all-resource-concurrency-routes.md; final executable registry and typed backend→internal/artifactstore and structural conformance
Contracts: executable registry rows (consumer ID, artifact/proof class, target domain, materialization policy, adapter owner) cross `internal/artifactstore/`→every audited caller, membership is closed and code-derived, ordering is select record before execute/materialize, absence or unregistered argv refuses, asserted by EC1-EC2 against real adapters and producer artifacts
Closure: EC1/gate-entry-consumers, EC1/gate-plumbing-consumers, EC1/gate-fixture-consumers, EC1/contract-consumers, EC1/preflight-consumers, EC1/prospective-consumers, EC1/release-consumers, EC1/verifier-consumers, EC1/proof-consumers, EC1/assembled-argv-audit, EC2/direct-exec-policy, EC2/immutable-link-policy, EC2/mutable-copy-policy, EC2/detach-before-mutation, EC2/writable-link-refusal

## What to build

Contract the wide artifact migration: delete legacy ordinary constructors, close the executable/proof registry, and turn on the structural typed-call and parsed-argv audit. The final registry is the only source for consumer membership and materialization policy.

## Acceptance

- [ ] [EC1] (covers SB4) every gate-entry, gate-plumbing, gate-fixture, contract, preflight, prospective, release, verifier, and proof consumer crosses the typed backend/adapter registry, and any assembled direct builder argv is red.
- [ ] [EC2] (covers SB8) every consumer obeys its registered direct-exec, immutable-link, or mutable-copy policy and no mutation can change the store artifact or another root.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| EC1/gate-entry-consumers | bypass the selected verifier/CLI record in `.bench/gate.sh` | structural consumer audit plus gate-entry refusal test | restore a direct construction/execution route and expect its unregistered-consumer diagnostic |
| EC1/gate-plumbing-consumers | rebuild `GateGoArgv` from separate `go`, `run`, and package literals | Go AST/parsed-argv audit | apply the assembled argv mutation and expect the forbidden authoring command |
| EC1/gate-fixture-consumers | restore one ordinary `internal/gate` builder helper | registry completeness audit | add the helper call and expect consumer-without-adapter failure |
| EC1/contract-consumers | restore one ordinary contract current-root build | registry completeness audit | add the helper call and expect consumer-without-adapter failure |
| EC1/preflight-consumers | restore one ordinary preflight current-root build | registry completeness audit | add the helper call and expect consumer-without-adapter failure |
| EC1/prospective-consumers | invoke `scripts/go-build.sh` outside the prospective adapter | parsed shell/argv audit | route around the adapter and expect the unregistered proof diagnostic |
| EC1/release-consumers | invoke a release builder outside its registered target/proof adapter | release consumer audit | add the raw script call and expect missing registry ownership |
| EC1/verifier-consumers | add a second verifier builder outside the backend | artifact-class author audit | introduce the command and expect the duplicate class-author diagnostic |
| EC1/proof-consumers | create a raw compiler-observing builder not in the proof registry | proof membership audit | add the direct command and expect unknown proof owner |
| EC1/assembled-argv-audit | split a forbidden command into variables/constants that contain no full substring | structural audit mutation test | construct and resolve the argv, then expect the same forbidden-command diagnostic |
| EC2/direct-exec-policy | copy a direct-exec artifact to an unverified path | materialization policy matrix | resolve the consumer and expect a policy mismatch before execution |
| EC2/immutable-link-policy | give an arbitrary mutator an immutable link | materialization policy matrix | resolve the consumer and expect mutable-consumer/link refusal |
| EC2/mutable-copy-policy | hardlink a mutable-copy consumer | existing inode/digest isolation test | mutate one root and expect store/sibling digest to remain fixed |
| EC2/detach-before-mutation | remove the atomic detach from an immutable-link mutation API | existing detach witness | backdate or rewrite one root and expect the shared-inode failure |
| EC2/writable-link-refusal | accept a writable multiply-linked completed artifact | artifact validator test | plant the mode/link-count combination and expect refusal before materialization |
