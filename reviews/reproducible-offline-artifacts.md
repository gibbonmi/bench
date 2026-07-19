# Reproducible offline artifacts — final terminal review

Range: `390c419..a725a29`

## Standards

- **High:** `scripts/assemble-offline-archive.mjs:14`, `scripts/release-plan.mjs:57`, `scripts/release-plan.mjs:64`, and `scripts/build-release-evidence.mjs:16` derive packaged-evidence membership with different predicates instead of one production-policy owner.
- **Medium:** `scripts/release-plan.mjs:15` and `scripts/offline-registry.mjs:8` hardcode the OS/architecture domain already declared by `scripts/release-plan.json:3`.

## Spec

- **High:** `.github/workflows/native-runtime.yml:52` uploads only `dist/artifacts`, so the evidence job lacks the required `dist/reproducibility.json` and cannot finalize evidence.
- **High:** `scripts/build-artifacts.sh:112` overlays dirty and untracked working-tree content into the second generation, while the focused workflow prerequisite does not validate checkout state; artifacts can claim clean `HEAD` provenance for uncommitted bytes.

## Coverage

- **High:** `tests/canary/package-core-guard/offline-archive-digest-binding-omitted/MUTATE.json:1` only renames a diagnostic and does not weaken archive-digest validation.
- **Medium:** `tests/canary/package-core-guard/native-proof-aggregation-bypassed/MUTATE.json:1` removes both aggregation calls rather than omitting exactly one target proof.
- **Medium:** `scripts/compare-artifacts.sh:8` and `scripts/aggregate-native-proofs.sh:6` resolve the repository root from an unresolved symlink path; row 18 lacks a real symlink-invocation proof for these seams and both fail when invoked through an external symlink.
