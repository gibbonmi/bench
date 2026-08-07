# Render doctor hook health

Blocked by: expose-hook-health-record.md
Ownership fence: `internal/adopt/doctor.go`, `internal/contract/surface/doctor_test.go`
Integration surfaces: hook-health producer→expose-hook-health-record.md; doctor row and exit posture→`internal/adopt/doctor.go` + D1/D2; doctor runtime contract→`internal/contract/surface/doctor_test.go` + D1/D2
Contracts: the complete hook-health record crosses `internal/adopt/link_hook.go`→`internal/adopt/doctor.go`, asserted by D1 against the real exported producer
Closure: D1/branch-provenance, D2/guess-green

## What to build

Render the effective branch and provenance in doctor’s managed-hook row while preserving a green exit for baked fallback provenance. Keep these two row facts together: separating provenance from the exit comparison strands the doctor runtime contract red when the baked fixture renders a new field but changes the accepted green posture.

## Acceptance

- [ ] [D1] Doctor’s managed pre-push row names the effective branch and provenance.
- [ ] [D2] A baked fallback branch leaves doctor’s exit code as green as a live-resolved branch.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| D1/branch-provenance | omit provenance from the doctor row | doctor runtime contract | run doctor in live and baked fixtures, expect both tokens in each row |
| D2/guess-green | return red for baked provenance | doctor runtime contract | compare live and baked fixture exits, expect equal green exits |
