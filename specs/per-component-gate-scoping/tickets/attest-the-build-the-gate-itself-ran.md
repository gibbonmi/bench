# Attest the build the gate itself ran

Blocked by: Expose the build inputs and seal digests the gate reads; Record a
component ancestor slot as its own class
Ownership fence: `internal/gate/build_attestation.go`, `internal/gate/build_attestation_test.go`, `internal/freshness/freshness.go`, `internal/gate/component_slots.go`, `internal/gate/kitshaped_fixture_test.go`, `internal/gate/record_classes.go`, `internal/gate/record_classes_test.go`, `internal/gate/verdict.go`
Assumptions: `freshness.Check` already refuses a missing binary, a missing seal,
a source-digest mismatch, an executable-digest mismatch, and a symlinked or
irregular sidecar; this ticket adds a check on top of it and weakens none of
them. Re-derive from the tree at pickup.

## What to build

A seal proves a binary matches its sources. It does not prove the gate built it:
an attacker or a mistaken script can plant a binary and recompute a
self-consistent seal, and every digest then agrees. The attestation closes that
gap — a gate-evidence record, in the same store as the ancestor slots, authored
only when the build phase runs green inside a gate, recording the executable
digest of the binary that build produced.

The build skip decision therefore needs both: `freshness.Check` passing *and*
the seal's executable digest matching the attestation. A planted binary plus a
recomputed seal fails attestation, so the build runs, overwrites both, and
re-authors the attestation. A green build republishes the seal through the
existing `Publish` and re-authors the attestation in the same step, so the two
can never disagree after a successful build.

This ticket lands the attestation record class, its authorship, and its
verification. Wiring it into the build skip is the next ticket's work.

**Evidence authorship.** `bench gate` authors the attestation when its build
phase runs green; `bench gate --fresh` re-authors it; `gate-phases` authors
none. No other command writes it.

The attestation is a new record class. Trace a sibling before adding it: the
slot class's validation, the verdict record's exact-field-set discipline, and
the evidence-store naming are the registries it must join.

## Acceptance

- [ ] [PC5] a planted binary with a recomputed, self-consistent seal fails attestation.
- [ ] [PS20] an attestation authored by a green build matches that binary's executable digest, and verification passes.
- [ ] [PS21] a missing attestation, an attestation naming a different executable digest, and a malformed attestation each fail verification.
- [ ] [PS22] authoring an attestation replaces the previous one atomically and leaves every component slot's bytes unchanged.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PC5 | verify only that the seal parses, without comparing the attested digest | `TestPlantedSealFailsAttestation` | publish a valid binary, replace it with a different one and recompute its seal by hand, verify, expect failure |
| PS20 | author the attestation from the seal rather than from the built binary | `TestGreenBuildAttestsItsOwnBinary` | run a build in the kit-shaped fixture, author, verify, expect success |
| PS21 | treat an absent attestation as satisfied | `TestAttestationRefusals` | remove, mismatch, and corrupt the attestation in three subtests, verify each, expect failure |
| PS22 | write the attestation with a plain `os.WriteFile` into the slot directory | `TestAttestationDoesNotDisturbSlots` | author a component slot, author an attestation twice, assert the slot bytes unchanged and one attestation present |
