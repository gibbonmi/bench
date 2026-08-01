# Record a component ancestor slot as its own class

Blocked by: Identify a component by its inputs and execution closure
Ownership fence: new `internal/gate/component_slots.go`, new
`internal/gate/component_slots_test.go`
Assumptions: the retained-evidence store is `<git-common-dir>/bench-gate-evidence`,
written through `durableReplaceAt` and read through `loadVerdict`; `strictJSON`,
`requireObjectFields`, and `strictRecordTime` are package-internal validators
this ticket calls without editing. Re-derive from the tree at pickup.

## What to build

N slots in the existing retained-evidence store, one per evidence-skipped
component, each addressed by that component's identity. A slot record is its own
validated class: it names its component, carries no inherited fields, and is
refused outright if it fails validation, names a different component than the
slot resolves for, or carries any field of the reduced or full verdict classes.
Every refusal is a run-the-component answer, never a repair.

Content addressing with no freshness bound: a slot stands until its component's
identity moves. Authorship is at execution — a run that executed a component
green authors that component's slot; a run that skipped it leaves the slot
byte-identical; a red component invalidates its own slot and touches no other's.
This ticket lands the class, its store operations, and its refusals; the gate
wiring that calls them arrives with the partition ticket.

**Evidence authorship.** `bench gate` is the canonical producer of these slots.
The `gate-phases` plumbing entry records no verdict cache and authors no slot, so
a phase table run through it never publishes evidence a later run consumes.

The slot class is a new member of an enumerated family. Trace an existing
sibling before adding it: the verdict record's field-set validation, the two
ready classes in `verdict.go`, and the evidence-store naming in `engine.go` are
what a record class already appears in — a new class that skips any of them is
readable by nothing or validated by nothing.

## Acceptance

- [ ] PC9 — a component skipped by a run leaves its slot byte-identical, including its recorded time.
- [ ] PC10a — a slot record carrying any verdict-class inherited field is refused and its component runs.
- [ ] PC10b — a slot record naming a different component than the slot resolves for is refused and its component runs.
- [ ] PC10c — a slot record failing field-set, schema, or time validation is refused and its component runs.
- [ ] PC11 — authoring a slot for one component leaves every other component's slot bytes unchanged, and invalidating one leaves the others present.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PC9 | re-stamp `recorded_at` when a slot is read | `TestSkippedComponentSlotIsByteIdentical` | author a slot, read its bytes, resolve it as a skip, re-read and compare bytes |
| PC10a | validate the slot with the full ready field set | `TestSlotWithInheritedFieldsRefuses` | write a slot carrying `ancestor`, resolve it, expect a refusal naming the class |
| PC10b | drop the component-name comparison from resolution | `TestSlotNamingAnotherComponentRefuses` | author `vet`'s slot, copy its bytes to `test`'s slot path, resolve `test`, expect a refusal |
| PC10c | accept a record whose schema is unset | `TestMalformedSlotRefuses` | write three malformed slots (bad field set, bad schema, future time), resolve each, expect refusals |
| PC11 | author every component's slot on any green | `TestSlotAuthorshipIsPerComponent` | author `vet`'s slot only, assert `test`'s slot absent; invalidate `vet`, assert `test`'s slot bytes unchanged |
