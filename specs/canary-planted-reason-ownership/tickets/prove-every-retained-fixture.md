# Prove every retained fixture through its owner

Blocked by: none
Writes: `internal/conformance`, `internal/releaseevidence/release_index_test.go`, `internal/releaseevidence/package_artifact_test.go`, `tests/canary/package-core-guard/release-digest-binding-omitted/`, `tests/canary/package-core-guard/release-package-evidence-omitted/`, `tests/canary/package-core-guard/release-evidence-probe-base.txt`, `tests/canary/package-core-guard/unrouted-subcommand/MUTATE.json`

## What to build

Replace the partial hand-maintained fixture-bite coverage with one test-internal runner that derives cases from `canary.Fixtures`, resolves each fixture's exact registered non-meta owner and declared tier, and proves the fixture expectation appears after materialization and disappears after restoration. Repair the three current false proofs atomically: retire the two release fixtures and their orphaned shared base into omission-graded owning-package mutation tests over canonical index encoding and focused production-builder staging, and reanchor `unrouted-subcommand` to the current `commandRegistry` entry. Keep the runner and live wrapper in the architecture-censused fixture proof file; do not add a production fixture journey, authenticated release probe, full artifact build, nested sweep, or focused FT168 surface.

## Acceptance

- [ ] (covers PB1) All 182 retained fixtures receive completed attributed subtests through their exact owner, closing the 29 retained missing comparisons; only the two exact release fixtures retire with equivalent owning-package tests.
- [ ] (covers PB2) A synthetic canary directory containing a registered-family dev fixture and an explicit-`CHECK` ship fixture drives both through the same producer-derived runner; stored names and tier or binding-form filters are red.
- [ ] (covers PB6) The live test calls the runner once, records identity only after `runFixtureBite` returns, compares the observed set with a fresh `canary.Fixtures` key set, and an AST grader rejects fixture-name literals, manual sources, a second proof loop, or record-before-proof.
- [ ] (covers PB3) Each fixture's normalized non-empty `EXPECT` appears after mutation and is absent after restoration through the same registered owner.
- [ ] (covers PB4) Missing fixtures; absent, zero-byte, ASCII-whitespace, or non-ASCII-whitespace expectations; and unbound, unregistered, or meta checks fail with fixture attribution, while a resolution-only synthetic ship fixture returns `registry.Ship` without invoking its owner.
- [ ] (covers PB5) `fixture_bite_test.go` is an exact `directArchitectureTests` member and admits no generic repository/process constructor, gate, wrapper, `go test`, or `go run`; native ship-owner implementation files remain outside that harness census.
- [ ] (covers PB7) The release JSON-tag and omitted-wrapper-copy mutations make `TestReleaseIndexBindsComponentManifestDigest` and `TestBuildReleaseEvidenceIncludesRegisteredPackageEvidence` red without source-text matching, ship authentication, or a full artifact build; an independent census makes either replacement test's omission red and requires both retired fixture directories plus their shared base absent; `unrouted-subcommand` materializes its current command-registry mutation and bites through `subcommand-routing`; the retained producer therefore advances from 181/182 to 182/182.
