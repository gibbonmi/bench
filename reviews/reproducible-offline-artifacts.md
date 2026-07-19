# Reproducible offline artifacts — closure review

Range: `390c419..53f9e54`

## Standards

- **Medium:** `internal/releaseevidence/offline_archive.go:98`, `internal/releaseevidence/artifact_evidence.go:175`, `internal/releaseevidence/artifact_evidence.go:206`, and `scripts/native-proof.sh:26` hardcode archive-manifest paths already defined by `scripts/release-plan.json`.

## Spec

- **High:** `scripts/verify-release-artifact.mjs:16` accepts a minimal index containing one matching artifact without validating approved status, schema, completeness, reproducibility, or native proofs.
- **Medium:** publish phases in `internal/releaseevidence/registry.json:2` construct artifacts before identity and ancestry checks in `internal/preflight/command.go:161`, allowing mismatched release state to materialize candidates before rejection.

## Coverage

- **High:** `internal/preflight/native_proof_test.go:72` lacks negative mutations for digest, musl, and strip bindings, leaving authoritative native-proof authorization without a red signal or behavior-owned canary.
- **Medium:** `internal/contract/surface/artifact/artifact_fixture_test.go:214` sends SIGTERM at the build promotion seam but does not prove the required SIGINT cleanup and idempotent rerun.
