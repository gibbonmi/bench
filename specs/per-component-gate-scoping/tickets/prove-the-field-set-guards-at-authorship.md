# Prove the field-set guards at authorship

Blocked by: Record a component ancestor slot as its own class; Attest the build
the gate itself ran
Ownership fence: `internal/gate/component_slots_test.go`, `internal/gate/build_attestation_test.go`
Assumptions: `requireObjectFields` is reached on the author path at
`component_slots.go` and `build_attestation.go` before any record is written;
`componentSlotFields` and `buildAttestationFields` are the declared exact field
sets; `strictJSON` already refuses unknown and duplicate names on the read path,
which is why no read-path mutation can red these guards. Re-derive from the
tree at pickup.

## What to build

Review finding C1: deleting both `requireObjectFields` calls leaves the suite
green — the guards' only distinct reach is the author path, where a field added
to `componentSlotRecord` or `buildAttestationRecord` without updating its
`*Fields` slice would refuse at authorship, and no test constructs that state.
A guard with no recorded red is outside the gate's bite standard.

Bind each record struct to its declared field set where the disagreement is
observable: a test per record class that derives the field names from the Go
struct's json tags by reflection and requires exact set equality with the
declared slice. Struct-and-slice drift then reds the suite at the binding,
which is the state the authorship guard exists to refuse.

## Acceptance

- [ ] [RF1] a test derives `componentSlotRecord`'s json field names by reflection and asserts exact set equality with `componentSlotFields`.
- [ ] [RF2] a test derives `buildAttestationRecord`'s json field names by reflection and asserts exact set equality with `buildAttestationFields`.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RF1 | remove one name from `componentSlotFields` | the RF1 test | run it, expect red naming the missing field |
| RF2 | remove one name from `buildAttestationFields` | the RF2 test | run it, expect red naming the missing field |
