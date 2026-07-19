# Reproducible offline artifacts — delta review

Range: `53f9e54..3a09150`

- **Medium (Standards):** `scripts/verify-release-artifact.mjs:28`, `scripts/verify-release-artifact.mjs:36`, and `internal/contract/surface/artifact/reproducibility_test.go:247` reconstruct artifact cardinality, target matrix, and filenames instead of consuming canonical `artifactRecords` from the release plan.
- **Medium (Coverage):** `internal/preflight/focused_publish_test.go:13` proves identity failure stops construction but does not prove ancestry failure does; the existing ancestry test only asserts a red exit.
- **Medium (Coverage):** `scripts/release-plan.mjs:73` strips a hard-coded `evidence/` prefix to recover packaged-evidence paths, so a valid relocated `{package_evidence}` plan entry fails despite correct archive projection.
