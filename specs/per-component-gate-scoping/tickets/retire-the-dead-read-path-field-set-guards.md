# Retire the dead read-path field-set guards

Blocked by: Prove the field-set guards at authorship
Ownership fence: `internal/gate/component_slots.go`, `internal/gate/build_attestation.go`, `internal/gate/component_slots_test.go`, `internal/gate/build_attestation_test.go`
Assumptions: `strictJSON` refuses unknown and duplicate names before either
read-path `requireObjectFields` call is reached, and every missing field of the
two flat records is caught by the next semantic check (schema mismatch, empty
component, non-content-address digest, strict record time); struct-to-field-set
drift already reds the existing authorship assertions, because both authors
marshal the struct and validate the bytes before publishing. Re-derive from the
tree at pickup.

## What to build

Composed-review round 2 (S4, C1): the read-path `requireObjectFields` calls in
the slot and attestation validators are dead — no mutation reds them — and the
round-1 binding tests duplicate coverage the authorship assertions already
carry, under comments whose warrant is false. Remove the two dead calls, the
two binding tests, and their `jsonFieldNames` helper. The field-set slices stay:
the record-class family registry consumes them.

## Acceptance

- [ ] [RG1] with both read-path `requireObjectFields` calls removed, every existing slot and attestation refusal test still passes (full `internal/gate` suite green).
- [ ] [RG2] struct-to-field-set drift still reds at authorship: renaming one json tag in either record struct reds an existing authorship test (demonstrated and restored).
- [ ] [RG3] the two round-1 binding tests and the `jsonFieldNames` helper are removed and nothing else references the helper.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RG1 | none — the criterion is the suite's green over the removal | full `internal/gate` suite | run it, expect green |
| RG2 | rename `json:"identity"` on `componentSlotRecord` | the existing authorship assertion in `component_slots_test.go` | run it, expect red; restore |
| RG3 | `rg jsonFieldNames internal/` | textual absence check | expect zero matches after removal |
