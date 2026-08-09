# Capture the pinned baseline independently

Blocked by: authenticate-baseline-manifest.md
Ownership fence: `internal/axi/compatibility`, `cmd/bench/axi_compatibility_test.go`, `specs/axi-compatibility-oracle/testdata`
Integration surfaces: authenticated manifest→authenticate-baseline-manifest.md; canonical builder→`scripts/go-build.sh` exercised unchanged; candidate comparison→compare-four-observations.md
Contracts: absolute executable identity and raw observation bytes cross `cmd/bench/axi_compatibility_test.go`→`specs/axi-compatibility-oracle/testdata`, membership is pinned baseline versus distinct candidate, ordering is build then capture then seal, and absence refuses before comparison, asserted by BC1
Closure: BC1/baseline-only, BC1/distinct-candidate, BC1/immutable-expected

## What to build

expected observations come only from the pinned baseline executable and cannot be refreshed from the candidate.

## Acceptance

- [ ] [BC1] (covers CO2) expected observations come only from the pinned baseline executable and cannot be refreshed from the candidate.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| BC1/baseline-only | route expected capture through the candidate executable | paired capture provenance test | build both subjects once and require candidate-authored expected refusal |
| BC1/distinct-candidate | reuse the baseline path as the candidate path | paired capture provenance test | compare executable identities and require distinct-subject refusal |
| BC1/immutable-expected | rewrite expected bytes after candidate execution | immutable fixture test | run the candidate and require the fixture digest to remain unchanged |

