# Refuse a partial verdict at prep-release

Blocked by: Carry the partition in the verdict and refuse its reuse
Ownership fence: `internal/preprelease/preprelease.go`, `internal/contract/surface/preprelease/preprelease_test.go`
Assumptions: `Refusal(inspection)` today branches on
`inspection.Reduced && inspection.Status == "green"` and points at
`bench gate --fresh`; the ship tier already refuses without a current dev-green
verdict. Re-derive from the tree at pickup.

## What to build

A release answers for the whole tree, so a verdict that skipped components
cannot authorize one. The refusal generalizes from the reduced marker to the
partition and says which components were skipped — a refusal that named only
"partial" would leave the maintainer guessing what went ungraded — and points at
the same single escape, `bench gate --fresh`.

## Acceptance

- [ ] [PC16a] a partial green is refused, and the refusal names every skipped component.
- [ ] [PC16b] the refusal points at `bench gate --fresh` and at re-running prep-release.
- [ ] [PS27] a reduced green is still refused with its existing wording.
- [ ] [PS28] a full green still passes the refusal check.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PC16a | emit a fixed "verdict is partial" string with no component list | `TestPrepReleaseRefusalNamesSkippedComponents` | build an inspection carrying a two-component partition, call `Refusal`, assert both names appear |
| PC16b | drop the `--fresh` clause | `TestPrepReleaseRefusalPointsAtFresh` | call `Refusal` on a partial green, assert the string contains `bench gate --fresh` |
| PS27 | route the reduced case through the partial branch | existing reduced-refusal surface test | call `Refusal` on a reduced green, compare the string |
| PS28 | refuse whenever the partition field is non-nil, including empty | `TestPrepReleaseAcceptsAFullGreen` | call `Refusal` on a full green inspection, assert no refusal |
