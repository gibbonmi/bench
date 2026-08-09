# Contract the ordinary-build census

Blocked by: route-ordinary-phase-plumbing.md, migrate-gate-helpers.md, migrate-contract-preflight-helpers.md, propagate-selected-binary-to-nested-gates.md, serialize-phase-tables.md, serialize-primary-stripped-schedule.md
Ownership fence: `internal/conformance/`
Integration surfaces: migrated phase constructors→route-ordinary-phase-plumbing.md; migrated gate helpers→migrate-gate-helpers.md; migrated contract/preflight helpers→migrate-contract-preflight-helpers.md; migrated nested gates→propagate-selected-binary-to-nested-gates.md; final scheduler constructors→serialize-phase-tables.md and serialize-primary-stripped-schedule.md; closed census→contract-run-directory-lifecycle.md and align-profile-and-changelog.md
Contracts: exact build-constructor inventory crosses repository source→the structural test in `internal/conformance/`, membership is the single gate owner, single focused-test owner, forbidden ordinary constructors, and the enumerated cross-target, release, artifact-posture/reproducibility, changed-source, alternate-package, linker/publication, and test-executable exceptions, ordering is classify the assembled call site before allowing it, and an unknown member reds, asserted by BC1 against the post-migration tree
Closure: BC1/single-gate-owner, BC1/single-test-owner, BC1/no-ordinary-go-build, BC1/no-ordinary-go-run-bench, BC1/no-ordinary-go-run-verifier, BC1/no-ordinary-subject-builder, BC1/cross-target-exceptions, BC1/release-exceptions, BC1/artifact-posture-exceptions, BC1/changed-source-exceptions, BC1/alternate-package-exception, BC1/linker-publication-exception, BC1/canary-test-compile-exception, BC1/gate-test-compile-exception, BC1/unknown-member-red

## What to build

Add one structural census that reads real assembled argv and source constructors. Its allowlist is the spec's exact exception set, not a substring category. The census is the contraction seam: after it lands, a new ordinary build or a new claimed exception cannot enter unnoticed.

## Acceptance

- [ ] [BC1] (covers RS9) the default conformance suite proves exactly one gate owner and one focused-test owner, rejects every forbidden ordinary constructor, accepts only each enumerated independent proof, and reds on any unclassified member.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| BC1/single-gate-owner | add a second gate-owner builder constructor | structural census | run default conformance and require the duplicate owner is named red |
| BC1/single-test-owner | add a second focused-test builder constructor | structural census | run default conformance and require the duplicate owner is named red |
| BC1/no-ordinary-go-build | add direct `go build` to an ordinary helper | structural census | run default conformance and require the exact path/constructor is named red |
| BC1/no-ordinary-go-run-bench | add `go run ./cmd/bench` to ordinary phase argv | assembled-argv census | run default conformance and require the phase/member is named red |
| BC1/no-ordinary-go-run-verifier | add `go run ./internal/freshness/check` to gate entry | shell-entry census | run default conformance and require the verifier constructor is named red |
| BC1/no-ordinary-subject-builder | add subject-mode `scripts/go-build.sh` to a consumer | structural census | run default conformance and require the consumer is named red |
| BC1/cross-target-exceptions | remove one of build-artifacts, native-proof, or stress cross-compile from the declared set | exception-set test | run default conformance and require the missing exact member is red |
| BC1/release-exceptions | remove release-preflight or release-only preprelease gate-go from the set | exception-set test | run default conformance and require the missing release member is red |
| BC1/artifact-posture-exceptions | collapse artifact-mode or reproducibility proof members into one wildcard | exception-set test | add a new posture builder and require it is unknown and red |
| BC1/changed-source-exceptions | route either named changed-source proof through selected bytes | changed-source plus census test | mutate each proof source and require distinct output and declared membership |
| BC1/alternate-package-exception | remove build-attestation's alternate package from the set | exception-set test | run the attestation proof and require its exact member remains declared |
| BC1/linker-publication-exception | replace the reduced fixture's compiler/linker build with selected bytes | reduced fixture proof | mutate the fixture package and require its published bytes follow the mutation |
| BC1/canary-test-compile-exception | prohibit canary `go test -c` | compiled canary proof | compile a bite package and require the distinct test executable remains allowed |
| BC1/gate-test-compile-exception | prohibit runtime gate partial-proof `go test -c` | partial-proof test | compile and execute the gate test binary and require the exact member remains allowed |
| BC1/unknown-member-red | accept a constructor by broad directory or category prefix | structural census mutation | add an unlisted subject builder beside an allowed file and require default conformance red |
