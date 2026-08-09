# Migrate guard aggregates

Blocked by: none
Ownership fence: `internal/guards`
Integration surfaces: shared aggregate/empty carriers→implemented prerequisite `axi-carriers-and-registry`; final contraction→contract-aggregate-empty-routes.md
Contracts: guard ScanResult strings cross `internal/guards/guards.go`→shared aggregate; domain is complete/incomplete with unknown/none literals; order is status-inspected-total-omitted-reason; absence is explicit unknown/none, asserted by GA1
Closure: GA1/status-complete, GA1/status-incomplete, GA1/inspected, GA1/total-unknown, GA1/omitted, GA1/reason-timeout, GA1/reason-none, GA1/order, GA1/route

## What to build

Migrate guard aggregates through the shared carriers without changing owner semantics or public bytes.

## Acceptance

- [ ] [GA1] (covers AE1) migrate guard aggregates preserve complete/incomplete status, inspected, total/unknown, omitted, reason/none, order, route.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| GA1/status-complete | render a complete scan as incomplete | the independent guard scan test | run the complete scan fixture and require exact status failure |
| GA1/status-incomplete | normalize incomplete to complete | the independent guard timeout test | run the bounded timeout fixture and require incomplete status |
| GA1/inspected | derive inspected from rendered rows | the independent guard scan test | inspect an omitted entry and require the owner count |
| GA1/total-unknown | coerce unknown total to zero | the independent guard timeout test | run before enumeration closes and require unknown |
| GA1/omitted | drop the omitted value | the independent guard scan test | limit visible entries and require owner omitted count |
| GA1/reason-timeout | replace timeout with none | the independent guard timeout test | run the bounded timeout and require timeout reason |
| GA1/reason-none | emit timeout on a complete scan | the independent guard scan test | run the complete scan and require none reason |
| GA1/order | reorder fields | the independent exact renderer test | render metadata and require current order |
| GA1/route | bypass shared aggregate | the independent route test | invoke the real guard command and require the missing route marker |
