# Project fixture-bite diagnostics to pinned paths

Blocked by: none
Writes: internal/canary/inventory.go, internal/canary/inventory_test.go, internal/conformance/fixture_bite_test.go
Covers: none

## What to build

A fixture-bite failure shows the diagnostics that name a path the fixture pins.
A catch-all owner returns hundreds of diagnostics over the materialized tree,
and the one relevant line hides among them. The failure message keeps the
`want %q` clause, prints the projected lines, and adds one trailing line with
the count of the omitted diagnostics.

The canary package exports `PinnedPaths(root, fixtureDir)`. This accessor wraps
the unexported reader that `FixturePins` already drives, so the map builder and
the new caller share one source. The live pin map does not change.

If no diagnostic names a pinned path, the message prints the first 20
diagnostics and the omitted count. The reader then still sees evidence.

## Acceptance

- [ ] A fixture-bite failure message lists only the diagnostics that name a pinned path, and one trailing line reports the omitted count.
- [ ] An empty projection falls back to the first 20 diagnostics, and the trailing line reports the omitted count.
- [ ] `FixturePins` returns the same map for the live tree as before the change.
