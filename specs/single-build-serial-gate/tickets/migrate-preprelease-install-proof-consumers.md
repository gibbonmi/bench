# Migrate preprelease and install proof consumers

Blocked by: register-real-build-proof-identities.md, prepare-gate-artifacts-before-scheduling.md
Ownership fence: `internal/preprelease/`, `package.json`, `scripts/gen-platform-packages.sh`, `internal/releaseevidence/registry.json`, `.github/workflows/native-runtime.yml`
Integration surfaces: closed proof request→register-real-build-proof-identities.md; prepared gate/build entry→prepare-gate-artifacts-before-scheduling.md; live preprelease/install callers and advertisements→owned paths; migrated inventory→route-real-build-proofs-through-registered-identities.md; legacy route retirement→enforce-the-executable-artifact-contract.md
Contracts: selected artifact or registered proof records cross preprelease, package-install, platform-package, evidence-registry, and native-runtime callers→their existing commands, membership is every live caller or advertisement in the owned census, ordering resolves or delegates to the typed producer before execution, absence of typed enrollment is structural red, asserted by MI1 against the exact live inventory
Closure: MI1/preprelease-steps, MI1/package-prepare, MI1/platform-packages, MI1/release-evidence, MI1/native-runtime-workflow

## What to build

Migrate or explicitly preserve as external drivers every live release/install caller and advertisement found by the executable census. No owned site may construct the Bench CLI through an unregistered local fallback.

## Acceptance

- [ ] [MI1] (covers local) preprelease steps, package prepare, platform-package generation, release evidence, and native-runtime workflow all cross a selected artifact or registered proof producer with their live semantics intact.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| MI1/preprelease-steps | retain one raw build step in the preprelease registry | exact step-registry test | enumerate live steps and expect the untyped route diagnostic |
| MI1/package-prepare | restore direct local construction in package prepare | package-script structural test | inspect the prepare command and expect the raw-builder diagnostic |
| MI1/platform-packages | retain direct construction in platform package generation | script structural test | inspect the generator and expect the raw-builder diagnostic |
| MI1/release-evidence | advertise a retired raw builder | release-evidence registry test | enumerate executable claims and expect the stale-route diagnostic |
| MI1/native-runtime-workflow | restore direct construction in the workflow | workflow structural test | inspect both native jobs and expect the stale-route diagnostic |
