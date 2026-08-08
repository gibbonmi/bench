# Migrate contract, preflight, and release proof consumers

Blocked by: register-real-build-proof-identities.md, migrate-contract-and-preflight-artifact-consumers.md
Ownership fence: `internal/contract/`, `internal/preflight/`, `scripts/build-artifacts.sh`, `scripts/release-preflight.sh`, `scripts/native-proof.sh`, `scripts/build-offline-archives.sh`
Integration surfaces: closed proof request→register-real-build-proof-identities.md; ordinary consumer adapters→migrate-contract-and-preflight-artifact-consumers.md; migrated contract/preflight/release inventory→route-real-build-proofs-through-registered-identities.md; legacy proof helper retirement→enforce-the-executable-artifact-contract.md
Contracts: registered proof identities cross contract, preflight, and release proof owners→artifact-store proof backend, membership is every real compiler-observing proof in the owned paths, ordering preserves each existing assertion while replacing authorship, absence refuses without raw script or Go-builder fallback, asserted by MC1 against the real inventory
Closure: MC1/contract-proofs, MC1/preflight-proofs, MC1/release-script-proofs, MC1/reproducibility-slots

## What to build

Migrate contract, preflight, and release-script real-build proofs to registered identities, including two independent reproducibility authorship slots.

## Acceptance

- [ ] [MC1] (covers local) every owned proof consumer resolves a registered identity and retains its existing compiler-observing assertion without a private builder.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| MC1/contract-proofs | retain one direct contract builder | structural proof-consumer audit | enumerate contract proof sites and expect its unregistered-consumer diagnostic |
| MC1/preflight-proofs | retain one direct preflight builder | structural proof-consumer audit | enumerate preflight proof sites and expect its unregistered-consumer diagnostic |
| MC1/release-script-proofs | retain one direct release-script builder | structural proof-consumer audit | enumerate owned scripts and expect its unregistered-consumer diagnostic |
| MC1/reproducibility-slots | collapse first and second authorship into one slot | existing reproducibility comparison | request both slots and expect two backend events before byte comparison |
