# Classify workflow-guidance fixtures

Blocked by: none
Ownership fence: `internal/conformance/registry_test.go`
Contracts: the `workflow-guidance-anchors` fixture-family directory crosses the canary inventory→`internal/conformance/registry_test.go`, asserted by FC1 against the complete real fixture enumeration

## What to build

Register the existing `workflow-guidance-anchors` family once under its real Go conformance owners, so every current and newly added fixture is classified without each fixture ticket deriving the same family fact.

## Acceptance

- [ ] [FC1] The complete canary fixture registry classifies `workflow-guidance-anchors` under its actual conformance source owners.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| FC1 | omit the family registration while its directory remains | the complete fixture registry test | run `TestCanaryFixtureRegistryClassifiesEveryFixture` and expect the named unclassified-family failure |
