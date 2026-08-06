# Remove the dead attested publisher and restore the one-entry enumeration

Blocked by: single-source-fake-go-builder-stub.md
Ownership fence: `internal/gate/build_attestation.go`, `internal/gate/build_attestation_test.go`, `internal/freshness/publication_topology_test.go`, `specs/go-build-cache-footprint/spec.md`
Integration surfaces: attestation authoring for the fixture pair→existing `internal/gate/component_decision.go` path + AP2; caller enumeration→`internal/freshness/publication_topology_test.go`; spec invariant advertisement→`specs/go-build-cache-footprint/spec.md`
Contracts: the one-entry production caller enumeration crosses `internal/freshness/publication_topology_test.go`→`specs/go-build-cache-footprint/spec.md`, asserted by AP1 against the real tree

## What to build

Repair for review finding `spec-07-attested-publisher-unreachable-in-production`.
`internal/gate/build_attestation.go:publishAttestedBuild` has no production call
site: the gate's build phase runs the builder as a child and authors its
attestation through `attestExecutedBuild`, so the install-and-attest
alternative is dead code, and the topology contract's two-publisher enumeration
asserts source presence, not reachability. Delete `publishAttestedBuild`; its
test authors the attested published fixture through the surviving real parts
(`freshness.Publish` composition stays out of production gate code — the test
file is walk-exempt). Shrink the topology contract's expected caller set to
exactly the subject-operation entry, dropping the `gateAttestationPublisher`
constant, its package-private assertion, and the bites fixture's
missing-gate-publisher leg; the all-package collection stays — it is what makes
the one-entry claim enforceable. Amend the spec's invariant paragraph and
acceptance row 2 to the one-entry wording: the production call graph enumerates
every `freshness.Publish` caller in any package as exactly the subject-operation
entry.

## Acceptance

- [ ] [AP1] The topology contract enumerates every production `freshness.Publish` caller in any package as exactly the subject-operation entry; any other caller — including a reintroduced gate-side publisher — or the entry missing is red, and the spec's invariant language states the same one-entry scope.
- [ ] [AP2] `publishAttestedBuild` is absent from the tree, and the build-attestation contracts still prove attestation-seal agreement and refusal against a fixture pair authored through the surviving production paths.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AP1 | reintroduce a production `freshness.Publish` caller inside `internal/gate` | the topology contract over the real tree | add the caller, run `go test -count=1 ./internal/freshness -run TestFreshnessPublicationTopology`, expect the exact unexpected-call diagnostic |
| AP2 | make the attestation fixture record a digest that disagrees with the seal it sits beside | the build-attestation refusal contract | author the disagreeing pair, run the attestation inspection tests, expect the refusal-with-rebuild verdict |
