# Classify the example-agreement canary family

Blocked by: land-the-example-agreement-check.md
Ownership fence: `internal/conformance/registry_test.go`
Assumptions: `TestCanaryFixtureRegistryClassifiesEveryFixture` refuses any fixture directory absent from `canaryFixtureFamilyRegistry` in `registry_test.go`; the `conformanceGoFixture` constructor is the classification shape the sibling conformance-owned families use; the land ticket's fence missed this registry because its charge never traced a sibling family — the gap this spec's registry-tracing rule exists to catch. Re-derive from the tree at pickup.

## What to build

FT164 repair: the `tests/canary/example-agreement` family, landed by the
example-agreement ticket, is classified in the canary fixture registry so the
classification test recognizes it. Four lines: one `canaryFixtureFamilyRegistry`
entry constructed with `conformanceGoFixture`, naming the check file and its
binding site, mirroring the sibling conformance-owned entries.

## Acceptance

- [ ] [CF1] `TestCanaryFixtureRegistryClassifiesEveryFixture` passes with the `example-agreement` family present on disk.
- [ ] [CF2] the whole `internal/conformance` package suite passes with no unclassified-fixture failure.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CF1 | delete the new registry entry | `TestCanaryFixtureRegistryClassifiesEveryFixture` | remove the entry, run the classification test, expect the unclassified `wrapped-example-fence` failure |
| CF2 | misname the family key in the entry | `TestCanaryFixtureRegistryClassifiesEveryFixture` | rename the key, run the classification test, expect the same unclassified failure |
