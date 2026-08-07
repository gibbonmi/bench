# Render hook health through guards

Blocked by: expose-hook-health-record.md
Ownership fence: `internal/guards/guards.go`, `internal/contract/axi/axi_guards_test.go`
Integration surfaces: hook-health producer→expose-hook-health-record.md; guards TOON and brief routes→`internal/guards/guards.go` + G1/G2/G4; SessionStart consumer→`.bench/hooks/session-start.sh` + G4; guards AXI contract→`internal/contract/axi/axi_guards_test.go` + G1/G2/G4
Contracts: the complete hook-health record crosses `internal/adopt/link_hook.go`→`internal/guards/guards.go`, asserted by G1 and G4 against the real exported producer
Closure: G1/guards-fields, G2/single-derivation, G4/brief-fields

## What to build

Render branch, provenance, and currency from the shared hook-health record in `bench guards` and `bench guards --brief`, without preserving any second derivation. Keep the full and brief renderers together: a full-only cut strands the guards AXI contract red because the SessionStart consumer receives a brief record missing the same fields.

## Acceptance

- [ ] [G1] `bench guards` renders branch, provenance, and currency from the shared record.
- [ ] [G2] No package outside `internal/adopt` derives hook health, and `internal/adopt` has one derivation site.
- [ ] [G4] `bench guards --brief` carries branch, provenance, and currency to the SessionStart banner.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| G1/guards-fields | derive a marker-only row without currency | guards AXI contract | run guards in a baked stale fixture, expect all shared-record cells |
| G2/single-derivation | retain exported `ClassifyPrePush` and a guards-local parser | compile and call-site checks | run package tests and the call-site assertion, expect the duplicate derivation failure |
| G4/brief-fields | truncate the added fields from brief output | guards AXI contract | run `guards --brief`, expect branch, provenance, and currency |
