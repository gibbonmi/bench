# Migrate conformance release-target proof consumers

Blocked by: register-real-build-proof-identities.md
Ownership fence: `internal/conformance/cross_compile_default_test.go`, `internal/conformance/cross_compile_stress_test.go`, `internal/conformance/package_core_checks_test.go`
Integration surfaces: closed proof request→register-real-build-proof-identities.md; default/stress release-target producers and shared caller→owned paths; migrated release-target inventory→route-real-build-proofs-through-registered-identities.md; legacy proof helper retirement→enforce-the-executable-artifact-contract.md
Contracts: target-aware proof records cross default-tag and stress-tag release-target producers→their shared conformance caller, membership and signature parity cover both build-tag siblings, ordering resolves a target record before the caller observes bytes, absence or target collapse refuses, asserted by MR1 against both build-tag variants
Closure: MR1/default-target-producer, MR1/stress-target-producer, MR1/shared-caller, MR1/build-tag-parity, MR1/release-target-identity

## What to build

Migrate both build-tag siblings of the release-target producer and their shared caller together so neither compilation mode retains a raw builder or a stale signature.

## Acceptance

- [ ] [MR1] (covers local) default and stress conformance variants compile against one target-aware proof contract, and two registered targets retain distinct records and bytes.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| MR1/default-target-producer | retain the default-tag direct builder | default-tag focused conformance test | run the default variant and expect the raw-builder diagnostic |
| MR1/stress-target-producer | retain the stress-tag direct builder | stress-tag focused conformance test | run the stress variant and expect the raw-builder diagnostic |
| MR1/shared-caller | bypass the returned proof record in the shared caller | package-core focused test | execute the shared caller and expect its record-observation failure |
| MR1/build-tag-parity | change only one sibling's signature | dual-tag compile test | compile both tag variants and expect exact contract parity |
| MR1/release-target-identity | omit target tuple from the request | existing cross-target assertion | build two registered targets and expect distinct records and bytes |
