## Standards

1 finding. Worst issue: the required packaged-evidence set has multiple production
owners.

1. **High — hard violation: packaged evidence is independently selected for
   wrapper and platform packages.** `AGENTS.md`'s one-source-per-fact rule
   requires production policy and executable registries to be single-sourced.
   Membership is expressed by `package_mode` at
   `internal/releaseevidence/requirements.json:17-25`, but the wrapper separately
   selects `LICENSE` and the entire governance tree at
   `scripts/wrapper-assets.json:2,19`, while platform packages hardcode them at
   `scripts/build-release-evidence.mjs:179-180`. The builder uses the registry
   only for validation (`scripts/build-release-evidence.mjs:19,79-87`). Derive
   wrapper and platform evidence copies from `packageEvidence`; reserve
   `wrapper-assets.json` for non-evidence assets.

## Spec

5 findings. Worst issue: the special-file acceptance test passes without
exercising its claimed failure.

1. **High — the special-file acceptance test is a false positive.** Row 212
   requires special-file sources to fail before artifact promotion
   (`specs/governed-release-evidence.md:212`). The FIFO fixture omits
   `internal/releaseevidence/requirements.json`
   (`internal/contract/surface/artifact/artifact_test.go:81-99`), which the builder
   reads before the FIFO (`scripts/build-artifacts.sh:31-32`). Supply a complete
   source fixture and assert the special-file diagnostic plus prior-set
   preservation.
2. **Medium — release-index binding is only partially acceptance-proven.** Row
   217 requires independent recomputation of artifact, component, SBOM,
   inventory, policy, target, identity, phase, and requirement relationships
   (`specs/governed-release-evidence.md:217`).
   `internal/preflight/release_index_test.go:71-90` recomputes only final artifact
   SHA-256, while `internal/conformance/native_workflow_test.go:288-299` checks
   component digests only for presence. Independently recompute every named
   relationship.
3. **Medium — the profile command matrix is incomplete.** Row 219 requires
   missing/unknown profile exit 2 plus green full public and bank fixtures
   (`specs/governed-release-evidence.md:219`). Tests cover green public, red bank,
   and focused publish at `internal/preflight/preflight_test.go:144-169,255-288`,
   but not missing profile, unknown profile, or green bank. Row 215's conditional
   record reason path also lacks a fixture. Add those command cases.
4. **Medium — two producer-record failure classes remain unproved.** Row 220
   requires missing, empty, malformed, unknown-version, duplicate-key,
   identity-mismatch, and digest-mismatch fixtures
   (`specs/governed-release-evidence.md:220`).
   `internal/preflight/review_fixes_test.go:55-90` adds only the latter four;
   profile tests cover missing, but no empty or syntactically malformed producer
   record exists. Add distinct built-command cases.
5. **Medium — encoding-failure preservation remains unproved.** Row 223 requires
   encode-failed runs to preserve the prior generation
   (`specs/governed-release-evidence.md:223`). Existing fault tests cover
   cancellation, unsafe inputs, and promotion failure, but `canonicalJSON` has no
   fault seam and `internal/releaseevidence/release_evidence.go:104-115` has no
   encoding-failure preservation test. Add a deterministic fault seam and test.

## Coverage

1 finding. Worst issue: artifact inspection has no total archive budget.

1. **High — aggregate archive exhaustion is unbounded.** A tgz containing
   millions of individually sub-256 MiB members, or multi-gigabytes of compressed
   input, is read wholesale and accumulated in a files map at
   `internal/releaseevidence/artifact_evidence.go:151,195,238`. The sole limit and
   test cover one oversized member (`internal/releaseevidence/artifact_evidence.go:228-236`;
   `internal/preflight/archive_hostile_test.go:11-31`). Add compressed-size,
   aggregate-expanded-byte, and member-count limits plus a built-command test
   proving prompt red and prior-generation preservation.
