# Delete the canary fixture-owner table

Blocked by: none
Writes: internal/conformance/registry_test.go, internal/conformance/registry_validation_test.go, internal/conformance/canary_fixtures_test.go (new), internal/conformance/tier_live_tree_test.go
Covers: none

## What to build

The quality survey of 2026-09-29, card 03, found that the canary fixture-owner table grades subjects that no longer exist. The table in `internal/conformance/registry_test.go` gives every fixture the one owner `conformance` and a list of source files. Most shell sources in the list are retired gate fragments. The shell-twin test reads a missing fragment as empty text, so its check passes and tests nothing. The live `.bench/gate.sh` is a 33-line wrapper, and the gate entry contract already refuses each retired fragment name in it.

Three other owners hold each fact that the table restates:

- The family binding in `internal/conformance/registry` and the `CHECK` file bind each fixture to its check.
- The universal fixture bite requires a non-empty `EXPECT` file, a bound family, and a registered check that is not meta for each fixture.
- The conformance family check reports each unbound canary family.

Delete the table and the three tests that read it. Move the `canaryFixturePaths` helper, which the fixture bite uses, into a new file named for it. Do not put the helper into `fixture_bite_test.go`, because that file is above its structure budget grant. Remove the two deleted test names from the live-tree classification.

## Acceptance

- [ ] `internal/conformance/registry_test.go` and `internal/conformance/registry_validation_test.go` are gone.
- [ ] No Go source names `fixtureRegistration`, `canaryFixtureRegistry`, `canaryFixtureFamilyRegistry`, or `fixtureRegistrationFor`.
- [ ] `canaryFixturePaths` is unchanged in `internal/conformance/canary_fixtures_test.go`.
- [ ] The live-tree classification names no deleted test, and `TestClassifiedLiveTreeInventoryNamesDetectedTests` passes.
- [ ] The conformance package tests and the fixture bite pass.
