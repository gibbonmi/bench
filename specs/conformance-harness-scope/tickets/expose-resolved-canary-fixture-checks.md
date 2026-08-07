# Expose resolved canary fixture checks

Blocked by: none
Ownership fence: `internal/canary/canary.go`, `internal/canary/scope_test.go`
Integration surfaces: fixture-level `CHECK` and family fallback resolution→`internal/canary/canary.go` + CR1/CR2; exported resolved check→scope-direct-conformance-fixture-bites.md + FB1/FB2 after lifecycle refresh; real `default-branch-refabricated` fixture→existing `tests/canary/package-core-guard/default-branch-refabricated/CHECK` + CR1
Contracts: resolved check name crosses `internal/canary/canary.go`→the refreshed conformance fixture-bite assignment, asserted by CR1/CR2 against the real canary inventory and later by FB1/FB2 against the real scoped runner
Closure: CR1/check-override, CR1/real-default-branch-fixture, CR2/family-fallback, CR2/shared-precedence

## What to build

`canary.Fixtures` exposes the check each discovered conformance fixture actually
grades. The value reuses canary's existing resolution rule: a fixture-level
`CHECK` wins over its family's registry binding, while a fixture without `CHECK`
inherits the family binding. Keep that precedence as one canary-owned fact so the
sweep and external inventory cannot disagree. Pin the public inventory with a
canary-package regression using the real `default-branch-refabricated` fixture,
whose family is `package-core-guard` but whose check is
`default-branch-single-source`.

## Acceptance

- [ ] [CR1] (covers local) `canary.Fixtures` reports `default-branch-single-source` for the real `default-branch-refabricated` fixture instead of its `package-core-guard` family default.
- [ ] [CR2] (covers local) a fixture without `CHECK` inherits its registered family check, and exported inventory resolution shares the same `CHECK`-over-family precedence used by the sweep.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CR1/check-override | ignore a present fixture `CHECK` while exporting the resolved check | the canary public-inventory regression | apply the omission, run `go test ./internal/canary -run '^TestFixturesExposeResolvedChecks$' -count=1`, require the reported check to fall back wrongly to `package-core-guard` |
| CR1/real-default-branch-fixture | replace the regression's real-tree lookup with a synthetic fixture lacking the checked-in override | the coordinator-owned real-fixture probe | run the regression against `tests/canary/package-core-guard/default-branch-refabricated`, require the literal `default-branch-single-source`, then verify the real `CHECK` remains the subject |
| CR2/family-fallback | export only fixture marker values and omit the family lookup | the canary public-inventory regression | create a fixture with no `CHECK`, run `go test ./internal/canary -run '^TestFixturesExposeResolvedChecks$' -count=1`, require its resolved check to equal the registered family check |
| CR2/shared-precedence | add a second precedence implementation for exported inventory instead of reusing canary's existing resolver | the canary scope controls plus review | invert `CHECK`-over-family precedence at the shared resolver, run `go test ./internal/canary -run '^(TestFixturesExposeResolvedChecks|TestSweepScopesFixtureRunsToTheirCheck)$' -count=1`, require both public inventory and sweep assertions to red |
