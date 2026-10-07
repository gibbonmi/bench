# 05 Share source observations with the architecture census

Blocked by: 04-share-source-observations-for-git-policy.md
Writes: internal/conformance/ordinary_build_census_test.go, internal/conformance/tier_live_tree_test.go
Covers: LTE28, LTE33

## What to build

Review chunk: LTE-C5.

Move architecture traversal and parsing onto the accepted snapshot interface. Retain the ordinary phase assertions and architecture classification logic.
Keep whole-tree coverage and the existing excluded directories. Keep the current ship-tag bypass and strict traversal and parse errors.
Preserve the direct byte-source helper used by classification tests. It may adapt through the same parser owner without duplicating parsing policy.

Run the old and new visitor against matching synthetic trees before removing the old walker. Compare exact ordered diagnostics, including paths and line numbers.
Keep the independent phase-name and argv expectations. This migration does not alter the phase schedule or race membership.

## Acceptance

- [ ] The census still detects retired entries and forbidden nested process or repository constructors.
- [ ] Root and nested Go sources remain checked outside the declared exclusions.
- [ ] Ship sources and excluded directories retain their previous disposition.
- [ ] Malformed and unreadable included sources retain their error diagnostics.
- [ ] Removing an ordinary phase still fails the independent phase expectation.
- [ ] The old walker is absent from the resulting visitor path.

## Verification

Run architecture equivalence cases and TestBranchNativeArchitectureCensus. Demonstrate a retained forbidden constructor red and the phase-omission red. Restore the source exactly before the chunk checkpoint.

