# Migrate guard aggregates

Blocked by: none
Ownership fence: `internal/guards`
Integration surfaces: ordered aggregate carrier→`internal/axi/aggregate.go` exercised by GA1; guard scan producer→`internal/guards/guards.go` exercised by GA1; guard scan renderer→`internal/guards/guards.go` exercised by GA1; legacy carrier contraction→contract-aggregate-empty-routes.md
Contracts: the `ScanResult` status/inspected/total/omitted/reason strings cross `internal/guards/guards.go`→`internal/axi/aggregate.go`; the status domain is exactly `complete`/`incomplete`, total and omitted are decimal digits or the literal `unknown`, reason is `timeout` or `none`; order is status, inspected, total, omitted, reason; absence is carried as the literal `unknown`/`none`, never as `0` or an empty cell, asserted by GA1 against the real `Scan` producer
Closure: GA1/status-complete, GA1/status-incomplete, GA1/inspected, GA1/total-unknown, GA1/omitted-count, GA1/omitted-unknown, GA1/reason-timeout, GA1/reason-none, GA1/order, GA1/route

## What to build

`bench guards` supplies its already-derived scan facts to the shared ordered aggregate
carrier and renders the identical `guard_scan[1]{status,inspected,total,omitted,reason}`
block. `Scan` keeps every derivation: the status string, the inspected count, the total
(`len(candidates)` or `unknown` when enumeration itself timed out), the omitted count, and
the reason. The carrier only transports the ordered typed facts.

AE1 is an already-covered row: `TestScanEnumerationTimeoutUsesUnknownCounts`,
`TestScanTimeoutPreservesPartialRowsAndHonestCounts`, and
`TestCommandAlwaysEmitsCompleteGuardScanMetadata` stay exactly as they are and remain the
named existing controls. This ticket adds the subject mutation the row lacks — a new
`TestGuardScanAggregateRouteCarriesOwnerFacts` in `internal/guards` that drives
`Command(nil)` in a real temporary repository and asserts the rendered block was produced
through the shared carrier from `Scan`'s own values, so a byte-equal local re-render is red.

Tree condition that must hold when this ticket is refreshed: `internal/axi/aggregate.go`
exists and declares the exported ordered-aggregate type `Aggregate` with its typed fact
entry `Fact`. If that path or either symbol is absent, stop and report rather than build —
the prerequisite `axi-carriers-and-registry` build has not landed.

## Acceptance

- [ ] [GA1] (covers AE1) `guard_scan` renders `Scan`'s own complete/incomplete status, inspected count, total, omitted, and reason in that order through the shared aggregate carrier.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| GA1/status-complete | in `Scan`, return `Status: "incomplete"` on the path where every candidate was inspected | `TestCommandAlwaysEmitsCompleteGuardScanMetadata` (`internal/guards`) | run `go test ./internal/guards -run TestCommandAlwaysEmitsCompleteGuardScanMetadata -count=1 -timeout 60s`; the `strings.Contains(out, "complete")` assertion fails because the block renders `status=incomplete`; the whole scan is bounded by `bounds.GuardScanTimeout` through `bounds.Context`, so a wedged inspection ends the scan rather than hanging the test |
| GA1/status-incomplete | in `Scan`'s cancelled-inspection branch (`guards.go:195`), return `Status: "complete"` | `TestScanTimeoutPreservesPartialRowsAndHonestCounts` (`internal/guards`) | run `go test ./internal/guards -run TestScanTimeoutPreservesPartialRowsAndHonestCounts -count=1 -timeout 60s`; the `got.Status != "incomplete"` check fails; the test's own `inspectGuard` stub blocks only on `ctx.Done()` and the test cancels that context itself, so the blocked worker always returns |
| GA1/inspected | in `Scan`, set `Inspected` from `len(rows)` instead of the inspected counter | `TestScanTimeoutPreservesPartialRowsAndHonestCounts` (`internal/guards`) | run `go test ./internal/guards -run TestScanTimeoutPreservesPartialRowsAndHonestCounts -count=1 -timeout 60s`; the `got.Inspected != "1"` check fails when a guard that emits no row is inspected; bounded by the test's explicit `cancel()` after the `blocked` channel closes |
| GA1/total-unknown | in `Scan`'s enumeration-timeout branch (`guards.go:177,182`), return `Total: "0"` instead of `Total: "unknown"` | `TestScanEnumerationTimeoutUsesUnknownCounts` (`internal/guards`) | run `go test ./internal/guards -run TestScanEnumerationTimeoutUsesUnknownCounts -count=1 -timeout 60s`; the `got.Total != "unknown"` check fails against the rendered `0`; the test's `context.WithTimeout(..., 20*time.Millisecond)` bounds the stubbed enumeration that waits on `ctx.Done()` |
| GA1/omitted-count | in `Scan`, compute `Omitted` as `len(candidates) - len(rows)` instead of `len(candidates) - inspected` | `TestScanTimeoutPreservesPartialRowsAndHonestCounts` (`internal/guards`) | run `go test ./internal/guards -run TestScanTimeoutPreservesPartialRowsAndHonestCounts -count=1 -timeout 60s`; the `got.Omitted != "1"` check fails; bounded by the test's `cancel()` on the `blocked` handshake |
| GA1/omitted-unknown | in `Scan`'s enumeration-timeout branch, return `Omitted: "0"` instead of `Omitted: "unknown"` | `TestScanEnumerationTimeoutUsesUnknownCounts` (`internal/guards`) | run `go test ./internal/guards -run TestScanEnumerationTimeoutUsesUnknownCounts -count=1 -timeout 60s`; the `got.Omitted != "unknown"` check fails; bounded by the test's 20ms enumeration context |
| GA1/reason-timeout | in `Scan`'s enumeration-timeout branch, return `Reason: "none"` | `TestScanEnumerationTimeoutUsesUnknownCounts` (`internal/guards`) | run `go test ./internal/guards -run TestScanEnumerationTimeoutUsesUnknownCounts -count=1 -timeout 60s`; the `got.Reason != "timeout"` check fails; bounded by the test's 20ms enumeration context |
| GA1/reason-none | in `Scan`'s all-inspected return, set `Reason: "timeout"` | `TestGuardScanAggregateRouteCarriesOwnerFacts` (`internal/guards`, new) | run `go test ./internal/guards -run TestGuardScanAggregateRouteCarriesOwnerFacts -count=1 -timeout 60s`; the assertion that the completed scan's fifth fact equals `none` fails against `timeout`; the fixture repository is a `t.TempDir()` `git init` and the scan is bounded by `bounds.GuardScanTimeout` |
| GA1/order | supply the aggregate facts to the carrier as status, total, inspected, omitted, reason | `TestCommandAlwaysEmitsCompleteGuardScanMetadata` (`internal/guards`) | run `go test ./internal/guards -run TestCommandAlwaysEmitsCompleteGuardScanMetadata -count=1 -timeout 60s`; the `strings.Contains(out, "guard_scan[1]{status,inspected,total,omitted,reason}:")` assertion fails against the reordered header; bounded by `bounds.GuardScanTimeout` |
| GA1/route | keep the pre-migration local `toon.Table("guard_scan", ...)` call and never construct the shared aggregate | `TestGuardScanAggregateRouteCarriesOwnerFacts` (`internal/guards`, new) | run `go test ./internal/guards -run TestGuardScanAggregateRouteCarriesOwnerFacts -count=1 -timeout 60s`; the assertion that the guard scan aggregate was constructed through `axi.Aggregate` with `Scan`'s five owner facts fails with no aggregate observed, even though the rendered bytes are unchanged; bounded by `bounds.GuardScanTimeout` |
