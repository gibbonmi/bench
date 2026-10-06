# 03 Consolidate conformance dispatch and diagnostic fixtures

Blocked by: none
Writes: internal/conformance/tier_test.go, internal/conformance/tier_live_tree_test.go, internal/conformance/fixture_bite_test.go, internal/conformance/validity_checks_test.go
Covers: LTE18, LTE19, LTE20, LTE21, LTE22, LTE23, LTE24

## What to build

Review chunk: LTE-C3.

Fold the dev-membership and timing-cardinality assertions into TestTimingOrderStable. The resulting fixture invokes the dispatcher exactly twice.
Retain exact registry membership, ship exclusion, row format, row count, order, and reset checks. Do not compare real elapsed values between runs.

Scope the hostile-path, tracked-mode, and absent-versus-empty fixtures to their registered diagnostic owners. Use the existing dispatcher scope when it preserves routing proof.
Preserve every diagnostic string currently asserted by those fixtures. Check owner subject selection before changing each call.
Retain the complete fixture universe, registered-owner red, restoration, and architecture omission checks.

Update live-tree classification only when actual detection changes. Reduce existing bodies to provide headroom in the oversized fixture file.

## Acceptance

- [ ] Two complete dev executions retain the three former tests' distinct assertions.
- [ ] Omitting a registered execution makes the consolidated fixture fail.
- [ ] Removing timing reset or changing dispatch order makes the consolidated fixture fail.
- [ ] The hostile root fixture still reports invalid package JSON through the registered owner.
- [ ] The executable-mode fixture still uses tracked Git mode and detects its mismatch.
- [ ] Missing and empty inputs retain each currently asserted diagnostic.
- [ ] Every retained fixture still reaches its registered owner and passes its restoration check.

## Verification

Run the three revised diagnostic tests, the consolidated timing test, and the fixture universe. Demonstrate the execution-omission, reset-omission, and order-change reds before accepting consolidation. Run the conformance package after restoring those mutations.

