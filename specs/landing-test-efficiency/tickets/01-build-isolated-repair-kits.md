# 01 Build isolated repair kits from selected source assets

Blocked by: none
Writes: internal/adopt/adopttest/ (new), internal/adopt/repairtest/session_test.go, internal/adopt/repairtest/kit_test.go
Covers: LTE1, LTE2, LTE3, LTE4, LTE5, LTE6, LTE7, LTE8, LTE9, LTE10, LTE11, LTE12, LTE13, LTE14

## What to build

Review chunk: LTE-C1.

Introduce the selected-kit fixture owner and migrate consumer repair sessions to it. Copy required assets from the canonical kit into private storage.
The helper must not import `adopt`. Each scenario retains a private repository, home, compatibility records, and writable kit.
Keep real Link, Doctor, and undo calls. Keep the kit-repair branch's repository-as-kit behavior without first constructing an unused consumer kit.

Before changing fixtures, record each affected test's predicate and source-asset requirements. Include the privacy tests that call `linkedSession`.
Measure fixture setup and focused test cost during required verification. Stop for a scope decision if the measurements refute the fixture-cost premise.
Keep the complete assertions in repair_test.go and privacy_test.go. Do not migrate them to shared writable linked repositories.

This chunk introduces the helper that ticket 02 consumes. Its accepted interface returns a private kit containing only explicitly selected source assets.
Missing requested assets fail construction. Bytes and executable modes survive copying, and no destination shares a writable inode with another fixture.

## Acceptance

- [ ] The selected-kit copy test detects an omitted requested asset and an executable-mode change.
- [ ] The isolation test mutates one destination and proves another destination and the canonical source remain unchanged.
- [ ] The helper excludes an unrelated source asset from the materialized kit.
- [ ] Managed repair, modified conflicts, foreign conflicts, privacy, undo, idempotence, and authority-boundary assertions retain their original predicates.
- [ ] The equivalent-replacement undo case still uses file identity, as well as bytes and mode.
- [ ] Kit repair retains broker update, undo, and aliased-shim cases.
- [ ] The verification record contains source-asset inventory, observed fixture costs, and named omission results.

## Verification

Run focused tests for adopttest and repairtest. Record a red for missing managed repair and for lost undo identity, then restore exactly. Run the repair package once after restoration. Do not add parallel environment or working-directory mutation.

