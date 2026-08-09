# Migrate outline aggregates

Blocked by: none
Ownership fence: `internal/outline`
Integration surfaces: shared aggregate/empty carriers→implemented prerequisite `axi-carriers-and-registry`; final contraction→contract-aggregate-empty-routes.md
Contracts: outline integer/boolean facts cross `internal/outline/outline.go`→shared aggregate; domain is owner-derived scan/project facts; order is current metadata order; absence retains explicit zeros, asserted by OA1
Closure: OA1/tracked, OA1/scanned, OA1/skipped, OA1/total, OA1/emitted, OA1/omitted, OA1/truncated, OA1/order, OA1/route

## What to build

Migrate outline aggregates through the shared carriers without changing owner semantics or public bytes.

## Acceptance

- [ ] [OA1] (covers AE2) migrate outline aggregates preserve tracked, scanned, skipped, total, emitted, omitted, truncated, order, route.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| OA1/tracked | derive tracked from result rows | the independent outline metadata test | include a tracked file with no symbols and require owner count |
| OA1/scanned | derive scanned from result rows | the independent outline metadata test | include a scanned empty file and require owner count |
| OA1/skipped | omit one skipped input | the independent outline skip test | include an oversized input and require skip facts |
| OA1/total | derive total from visible rows | the independent outline boundary test | cross the cap and require owner total |
| OA1/emitted | report total as emitted | the independent outline boundary test | cross the cap and require emitted count |
| OA1/omitted | derive omitted before owner filtering | the independent outline metadata test | mix skipped and bounded rows and require exact omission |
| OA1/truncated | infer truncated from skipped state | the independent outline metadata test | skip without row bound and require current boolean |
| OA1/order | reorder fields | the independent exact renderer test | render metadata and require current order |
| OA1/route | bypass shared aggregate | the independent route test | invoke outline and require the missing route marker |
