# Record the FT156 Contracts anchor gap

Blocked by: none
Ownership fence: `ROADMAP.md`
Contracts: none crosses

## What to build

FT156 records the verified false green: deleting only the ticket template's `Contracts:` line passes a real graded-root `TestRootConformance` because no section-scoped template requirement anchors that line. The roadmap rider already carries this outcome and assigns the oracle gap to FT156 without changing the oracle; this ticket is its separate ownership receipt.

## Acceptance

- [ ] [CA1] the FT156 roadmap row records the exact `Contracts:`-line deletion, the real graded-root false green, and the missing section-scoped template requirement.
- [ ] [CA2] the outcome assigns the oracle gap to FT156 without claiming that the conformance check was fixed.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CA1 | omit the concrete `Contracts:` deletion from the rider | the roadmap review | inspect FT156 against the recorded graded-root probe, expect the deleted line and the missing section-scoped requirement to remain explicit |
| CA2 | claim the false green was fixed | the semantic reviewer | compare the rider with the unchanged conformance oracle, expect it to record the gap under FT156 rather than claim a fix |
