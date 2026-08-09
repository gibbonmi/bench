# Migrate roadmap aggregates

Blocked by: none
Ownership fence: `internal/roadmap`
Integration surfaces: shared aggregate/empty carriers→implemented prerequisite `axi-carriers-and-registry`; final contraction→contract-aggregate-empty-routes.md
Contracts: roadmap typed context/drain facts cross `internal/roadmap`→shared aggregate; domain is owner source/drain state; order is current block/field order; absence distinguishes zero/unknown/degraded, asserted by RA1
Closure: RA1/order, RA1/totals, RA1/unknown, RA1/zero, RA1/degraded, RA1/emitted-less-total, RA1/route

## What to build

Migrate roadmap aggregates through the shared carriers without changing owner semantics or public bytes.

## Acceptance

- [ ] [RA1] (covers AE3) migrate roadmap aggregates preserve order, totals, unknown, zero, degraded, emitted-less-total, route.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RA1/order | reorder context facts | the independent roadmap renderer test | render the context fixture and require current order |
| RA1/totals | derive totals from visible rows | the independent roadmap context test | supply more source facts than visible and require owner totals |
| RA1/unknown | coerce unknown to zero | the independent roadmap context test | make one source unreadable and require unknown |
| RA1/zero | omit explicit zero | the independent roadmap drain test | run an empty disposition and require zero remains present |
| RA1/degraded | drop degraded evidence | the independent roadmap context test | make one source unreadable and require evidence facts |
| RA1/emitted-less-total | equate emitted and total | the independent roadmap context test | truncate the visible set and require both owner values |
| RA1/route | leave one context aggregate local | the independent roadmap route test | invoke the context command and require the missing route marker |
