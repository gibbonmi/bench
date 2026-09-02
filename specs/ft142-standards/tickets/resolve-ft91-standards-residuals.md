# Resolve the FT91 standards residuals

Blocked by: none
Writes: ROADMAP.md, roadmap/FT142.md, internal/canary/inventory.go, internal/conformance/fixture_bite_test.go, internal/conformance/tier_test.go, internal/preprelease/preprelease.go
Covers: none

## What to build

Remove the two FT91 standards defects that remain in the current tree. Make the
release package registry the only package-list owner. State the dev-tier test
failure as a current fact instead of a change record.

The audit must also verify the resolved residuals. The canary inventory owns the
`CHECK` filename, and the tree has no second contract import prefix.

## Acceptance

- [ ] The ship-step comment refers to the release package registry without listing its members or count.
- [ ] The dev-tier test failure describes the current tier contract without FT91 change provenance.
- [ ] The tree has one owner for the `CHECK` marker filename.
- [ ] The canary and gate packages do not derive a second contract import prefix.
- [ ] FT142 retains only its runtime tracks, and FT294 becomes the next recommended item.
