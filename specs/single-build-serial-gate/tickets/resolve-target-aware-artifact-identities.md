# Resolve target-aware artifact identities

Blocked by: share-artifacts-across-local-processes.md, rendezvous-artifact-consumers-across-processes.md
Ownership fence: `internal/artifactstore/`
Integration surfaces: opaque store identity→share-artifacts-across-local-processes.md; target-aware `ArtifactRequest`→prepare-gate-artifacts-before-scheduling.md; target-aware `ArtifactRequest`→register-real-build-proof-identities.md; closed request fields→enforce-the-executable-artifact-contract.md
Contracts: `ArtifactRequest` (artifact class, target-selected source digest, toolchain, GOOS, GOARCH, CGO, flags, mode, closed proof slot) crosses resolver→store in `internal/artifactstore/`, membership is the closed request-field registry, ordering is select target files then frame identity fields before store lookup, absence or an unknown class/slot refuses, asserted by TA1 against target fixtures and the real store
Closure: TA1/windows-source-selection, TA1/unix-source-selection, TA1/architecture-selection, TA1/cgo-selection, TA1/toolchain-identity, TA1/goos-identity, TA1/goarch-identity, TA1/build-mode-identity, TA1/build-tags-identity, TA1/artifact-class-identity, TA1/proof-slot-identity, TA1/output-path-exclusion, TA1/root-spelling-exclusion, TA1/mtime-exclusion, TA1/invocation-exclusion, TA1/process-exclusion

## What to build

Expand the opaque store with the exact target-aware request resolver before any gate consumer enrolls. This is separate from store publication because today's host-derived file enumeration can land green while the store remains unconsumed, whereas each target-selection class needs its own red mutation.

## Acceptance

- [ ] [TA1] (covers SB2) target-selected source, toolchain, target, build mode, artifact class, and registered proof slot determine identity, while output path, root spelling, mtime, invocation, and process do not.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| TA1/windows-source-selection | omit Windows-selected files from a Windows request | Windows build-tag identity fixture | mutate a Windows-only source and expect a distinct identity |
| TA1/unix-source-selection | omit Unix-selected files from a Unix request | Unix build-tag identity fixture | mutate a Unix-only source and expect a distinct identity |
| TA1/architecture-selection | enumerate files under the host architecture | architecture-tag identity fixture | mutate a non-host architecture source, request that architecture, and expect a distinct identity |
| TA1/cgo-selection | ignore CGO file selection | CGO identity fixture | toggle CGO selection with distinct sources and expect distinct identities |
| TA1/toolchain-identity | omit the selected toolchain from the preimage | request-field registry test | resolve equal sources under two toolchains and expect distinct identities |
| TA1/goos-identity | omit GOOS from the preimage | request-field registry test | resolve equal selected bytes for two operating systems and expect distinct identities |
| TA1/goarch-identity | omit GOARCH from the preimage | request-field registry test | resolve equal selected bytes for two architectures and expect distinct identities |
| TA1/build-mode-identity | omit build mode from the preimage | request-field registry test | resolve two registered modes and expect distinct identities |
| TA1/build-tags-identity | omit build tags from the preimage | request-field registry test | resolve equal sources under two registered tag sets and expect distinct identities |
| TA1/artifact-class-identity | omit artifact class from the preimage | request-field registry test | request CLI and verifier classes over equal bytes and expect distinct identities |
| TA1/proof-slot-identity | omit the registered proof slot from the preimage | request-field registry test | request two finite authorship slots and expect distinct identities |
| TA1/output-path-exclusion | include output path in the preimage | incidental-field exclusion test | resolve two destinations and expect one identity and one store record |
| TA1/root-spelling-exclusion | include absolute or symlinked root spelling | incidental-field exclusion test | resolve the same Git subject through two root spellings and expect one identity |
| TA1/mtime-exclusion | include source mtime | incidental-field exclusion test | touch selected files without byte changes and expect the same identity |
| TA1/invocation-exclusion | include caller invocation | incidental-field exclusion test | resolve through two registered callers and expect the same identity |
| TA1/process-exclusion | include process identity | fresh-process exclusion test | resolve from two processes and expect the same identity and record |
