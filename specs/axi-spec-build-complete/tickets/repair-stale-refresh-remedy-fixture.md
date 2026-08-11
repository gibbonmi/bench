# Repair the stale refresh-remedy matrix fixture

Blocked by: none
Ownership fence: `internal/specbuild/testdata/axi-cases.jsonl`, `internal/specbuild/disclosure_test.go`
Integration surfaces: none crosses
Contracts: none crosses
Closure: SF1/refresh-remedy-fixture-matches-production

## What to build

Close the accepted Spec finding (P1) from the Terra/xhigh review of candidate
`45f32da76c90f513dc381850d0b673d5c74da120`: the `repair-stale-refresh-remedy.md`
round changed `RefusalForClass`'s `RefusalStaleRefresh`/`RefusalSpentRefresh`
case to emit `refreshAction(slug, "", "")`, which now renders
`bench spec build assign 'build demo' --ticket <ticket> --request <request>
--refresh <receipt>` for the no-known-values case (`internal/specbuild/disclosure.go`
around line 281). But the two matching entries in
`internal/specbuild/testdata/axi-cases.jsonl` — `matrix/assign/stale-refresh-receipt`
and `matrix/assign/spent-refresh-receipt` — still pin the pre-repair remedy with
no `--refresh <receipt>`, in both the inline error line and the appended
`help[1]{command}` block. Nothing currently catches this drift: fixtures whose
`id` starts with `matrix/` are exempt from `TestHistoricalFixtureOraclePinsAll122CasesAnd244Payloads`'s
live-reproduction equality check (`internal/specbuild/disclosure_test.go` around
line 363, `if !strings.HasPrefix(fixture.ID, "matrix/")`), so this specific pair
is unguarded by any test today. This violates SB6's "every operation and
lifecycle state receives a reviewed old-to-new fixture pair" requirement
(`specs/axi-spec-build-complete/spec.md` around line 61).

Fix, in two parts:

1. Update the `new` field of both `matrix/assign/stale-refresh-receipt` and
   `matrix/assign/spent-refresh-receipt` in `internal/specbuild/testdata/axi-cases.jsonl`
   to the exact bytes the real code emits today. Do not hand-derive these bytes:
   drive the real disclosure observation path (`ObserveDisclosureCell` /
   `newDisclosureObservation` in `internal/specbuild/disclosure_observation.go`,
   the same fixture machinery `TestEveryApplicableDisclosureCellUsesOneRealPublicObservation`
   already drives) for the `assign/stale-refresh-receipt` and
   `assign/spent-refresh-receipt` cells, capture the exact returned output, and
   copy those exact bytes into the fixture's `new` field (JSON-escaped as the
   file's existing entries already are — reuse the file's own escaping, don't
   hand-encode). Leave `old` untouched; it is the pre-AXI-migration baseline and
   is unaffected by this repair round.
2. Change each entry's `deltas` field from `["help"]` to `["help","remedy"]`
   (matching the file's existing precedent for a fixture whose remedy text
   changed, e.g. `matrix/assign/missing-run`), then extend
   `requireOnlyNamedDeltas`'s `"help,remedy"` case in
   `internal/specbuild/disclosure_test.go` (around line 505) with a third
   `else if` arm — alongside its existing `matrix/checkpoint/invalid-evidence-receipt`
   and default-start-remedy arms — that requires
   `strings.Contains(fixture.New, "--refresh <receipt>")` for these two fixture
   IDs specifically. Do not weaken or remove either existing arm.

## Acceptance

- [ ] [SF1] (covers local) (P1) `matrix/assign/stale-refresh-receipt` and
  `matrix/assign/spent-refresh-receipt` in `internal/specbuild/testdata/axi-cases.jsonl`
  carry the exact current production bytes (captured from the real disclosure
  observation path, not hand-derived) including `--refresh <receipt>` in both
  the inline error line and the help block, tagged `["help","remedy"]`, and
  `requireOnlyNamedDeltas` positively asserts the `--refresh <receipt>` remedy
  for both — a regression that drops `--refresh <receipt>` from either fixture
  or from the production remedy itself turns a focused test red.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| SF1/refresh-remedy-fixture-matches-production | revert either fixture's `new` field to omit `--refresh <receipt>`, or drop the new `else if` arm from `requireOnlyNamedDeltas` | `TestHistoricalFixtureOraclePinsAll122CasesAnd244Payloads` and a focused `requireOnlyNamedDeltas` assertion | run `go test ./internal/specbuild/... -run TestHistoricalFixtureOraclePinsAll122CasesAnd244Payloads` and require the mutation caught, either by the deltas-shape assertion or by the new positive `--refresh <receipt>` requirement |
