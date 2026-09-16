# Validate legacy prepared packs

Blocked by: none
Writes: internal/preflight, internal/chargeevidence (new), .agents/skills/bench-craft-delegate/references/charge-evidence-format.md (new), reviews/bounded-charge-evidence.md (new), internal/toon/toon_test.go, internal/conformance/data_handling_test.go
Covers: CE16, CE18, CE19, CE20, CE21, CE22, CE23, CE24, CE25, CE26, CE27, CE28, CE29, CE30, CE31, CE32, CE33, CE34, CE35, CE36, CE37, CE38, CE41, CE42, CE43, CE44, CE45, CE46, CE47, CE48, CE49, CE50, CE51, CE52, CE111, CE114, CE125, CE126, CE147, CE148, CE149

## What to build

Preserve legacy build output through a validated in-memory pack and its generated format reference.

Refactor the existing build charge producer through one validated in-memory prepared pack.
The existing legacy renderer consumes its exact decoded sources and metadata.
Preserve the enumerated compact and full output bytes, exits, refusal ordering, and movement behavior.
Review remains on its existing renderer at this checkpoint.

Implement the canonical format registry, strict pack reader, and generated shipped format reference together.
The real legacy renderer supplies the first consumer of the pack; no unused schema-only layer lands.
Keep source selection in preflight and serialization in chargeevidence.
Derive build metadata from the existing typed facts, then prove its registered schema.

This ticket creates no persistent artifact and adds no public command.
The root ticket starts from the pinned baseline and has no predecessor.
Record a differential fixture before refactoring, including the exact output and exit family.
Record omission mutations for the strict reader, legacy consumer, metadata schema, and reference projection.

Review chunk: CE-C1A

Use the shared contract in the spec; do not duplicate its field inventories here.
The successor starts only after this chunk review and its repair coverage close.

## Context price and checks

8 focused production files; at most 1,600 source lines. Read charge.go, preparation.go, charge_test.go, source adapters, and TOON before codec changes.
This slice owns 41 predicates, including separately named table cases.
Use one retained context and the existing shared fixture harness; do not copy private helpers across packages.
Do not add code to an over-budget file without moving its responsibility and headroom in this ticket.

- `bench test --package ./internal/preflight`
- `bench test --package ./internal/chargeevidence`

Record a biting omission or mutation for every independent expected schema or policy fact.
The completion-plan probe is the minimum named mutation, not a substitute for the acceptance cases.

## Acceptance

- [ ] CE16: Identical complete logical inputs produce identical evidence identities.
- [ ] CE18: A changed mode changes the evidence identity.
- [ ] CE19: A changed spec selector changes the evidence identity.
- [ ] CE20: A changed ticket selector changes the evidence identity.
- [ ] CE21: A changed base changes the evidence identity.
- [ ] CE22: A changed source tip changes the evidence identity.
- [ ] CE23: A changed role changes the evidence identity.
- [ ] CE24: A changed path changes the evidence identity.
- [ ] CE25: A changed requiredness changes the evidence identity.
- [ ] CE26: Changed source bytes change the evidence identity.
- [ ] CE27: Changed required metadata changes the evidence identity.
- [ ] CE28: The writer emits the documented canonical TOON profile.
- [ ] CE29: The reader refuses an unsupported manifest profile.
- [ ] CE30: The reader refuses noncanonical manifest bytes.
- [ ] CE31: The reader refuses unsupported container version.
- [ ] CE32: The reader refuses invalid header marker.
- [ ] CE33: The reader refuses nonzero reserved header bits.
- [ ] CE34: The reader refuses overflowing length arithmetic.
- [ ] CE35: The reader refuses overlapping page ranges.
- [ ] CE36: The reader refuses gapped page ranges.
- [ ] CE37: The reader refuses truncated container.
- [ ] CE38: The reader refuses trailing container data.
- [ ] CE41: Preparation refuses a required source in the absent state.
- [ ] CE42: Preparation refuses a required source in the empty state.
- [ ] CE43: Preparation refuses a required source in the symlink state.
- [ ] CE44: Preparation refuses a required source in the FIFO state.
- [ ] CE45: Preparation refuses a required source in the socket state.
- [ ] CE46: Preparation refuses a required source in the device state.
- [ ] CE47: Preparation refuses a required source in the directory state.
- [ ] CE48: Preparation refuses unsupported source bytes.
- [ ] CE49: Preparation refuses a dirty source checkout.
- [ ] CE50: Preparation refuses a missing selected ticket.
- [ ] CE51: Preparation refuses a source-tip mismatch.
- [ ] CE52: Preparation refuses an inactive assignment.
- [ ] CE111: Ordinary preflight preserves its enumerated baseline behavior.
- [ ] CE114: Required-source inventories derive from one executable policy.
- [ ] CE125: The pack reader refuses a page range beyond its source length.
- [ ] CE126: The manifest reader refuses a changed scalar type.
- [ ] CE147: The legacy build charge preserves its enumerated output through the validated in-memory pack.
- [ ] CE148: The shipped format reference equals the canonical format registry projection.
- [ ] CE149: Build metadata conforms to the exact registered metadata schema.
