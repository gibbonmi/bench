# Enumerate every production sealed-publication caller

Blocked by: none
Ownership fence: `internal/freshness/publication_topology_test.go`, `specs/go-build-cache-footprint/spec.md`
Integration surfaces: spec row-2 invariant advertisement→`specs/go-build-cache-footprint/spec.md`; gate attested-build publisher symbol→existing path `internal/gate/build_attestation.go` + TP1; production caller walk→`internal/freshness/publication_topology_test.go`
Contracts: the two-entry production caller enumeration (subject-operation entry plus gate attested-build publisher) crosses `internal/freshness/publication_topology_test.go`→`specs/go-build-cache-footprint/spec.md`, asserted by TP1 against the real tree

## What to build

Repair for review finding `spec-03-topology-ignores-non-main-callers`. The
topology contract collects `freshness.Publish` callers only from `package main`,
so production library callers are invisible. The resolved design is proven and
pinned at `/tmp/bench-ft195-topology-resolution.patch`, with evidence in
`/tmp/bench-ft195-review-debug-receipt.md`: collection drops the package filter
and enumerates every production caller in any package; the expected set is
exactly the builder-visible subject-operation entry plus the gate's
package-private attested-build publisher, the latter asserted package-private so
it cannot become a second command surface. The spec's "exactly one entry"
language (line 75 and acceptance row 2) is amended to the decided scope using
the wording drafted in the receipt.

## Acceptance

- [ ] [TP1] The topology contract enumerates every production `freshness.Publish` caller in any package as exactly the subject-operation entry plus the gate's package-private attested-build publisher; any other caller, either expected caller missing, or the gate publisher turning exported is red, and the spec's invariant language states this same two-publisher scope.
- [ ] [TP2] The bites fixture proves the former blind spot in both directions: an injected `package main` caller and an injected non-main library caller each produce their exact unexpected-call diagnostic, and a fixture missing the gate publisher reports it missing.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| TP1 | add a production library caller of `freshness.Publish` outside the enumerated set | the topology contract over the real tree | add the caller file, run `go test -count=1 ./internal/freshness -run TestFreshnessPublicationTopology`, expect the exact unexpected-call diagnostic |
| TP2 | restore the `package main`-only collection filter | the extended bites fixture | run `go test -count=1 ./internal/freshness -run TestFreshnessPublicationTopologyBites`, expect the missing library-caller diagnostic to fail the test |
