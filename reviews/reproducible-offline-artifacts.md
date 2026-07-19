# Reproducible offline artifacts — terminal review

Range: `390c419..7a1d454`

## Standards

- **High:** `internal/releaseevidence/release_plan.go:58` independently expands archive entries already owned by `scripts/release-plan.mjs:46`, so builder and authorization consumers can drift.
- **Medium:** `scripts/release-plan.json:3` stores `target_cardinality` alongside the target list, duplicating a derived count.
- **Medium:** `internal/conformance/workflow_checks_test.go:82` and `internal/preflight/native_proof_test.go:57` hardcode the target/runner matrix already defined by the release plan without a named mutation proving the duplicate expectation is necessary.
- **Medium:** `ROADMAP.md:40` and `ROADMAP.md:299` still direct cold sessions to write the already-implemented spec.
- **Medium:** `CHANGELOG.md:7` omits the user-visible reproducible offline archives and `commands --brief` surface.

## Spec

- **High:** `scripts/build-artifacts.sh:115` compares two builds made from the same source tree before final `release-index.json`, `SHA256SUMS`, and independently generated evidence exist, leaving acceptance row 250 unresolved.
- **High:** `scripts/smoke-offline.sh:97` executes native binaries outside the Node-only network sentinel, leaving acceptance row 245 open for native egress.
- **High:** `tests/canary/package-core-guard/offline-network-repair-allowed/MUTATE.json:1` changes a default that fails before any repair attempt, so acceptance row 259 lacks a biting repair-attempt mutation.
- **High:** `.github/workflows/native-runtime.yml:33` constructs and uploads artifacts before authoritative release identity/preflight authorization, contrary to acceptance row 262.
- **Medium:** `internal/contract/surface/artifact/artifact_fixture_test.go:206` interrupts only builder promotion with SIGTERM; acceptance row 261 requires SIGINT recovery proof at comparison, installation, registry service, and removal as well.

## Coverage

- **High:** `internal/contract/surface/artifact/offline_journey_test.go:60` disables repair control but never injects or observes a compiler/rebuild attempt, leaving acceptance row 247 without a rebuild-specific red signal.
- **Medium:** acceptance row 261 has no stage-specific interruption mutations or assertions for comparison, installation, registry service, and removal, nor proof that each interrupted stage reruns idempotently.
