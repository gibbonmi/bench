# Prove the repair needles and true up the inventory

Blocked by: retire-title-blockers-and-pin-the-unpinned-halves.md
Ownership fence: `internal/conformance/fixture_bite_test.go`, `specs/ft164-ticket-contracts/spec.md`
Assumptions: the basename require needle's diagnostic appears nowhere in the mutation table so deleting the require stays green; the forbid's additive power is unexercised because the only proof is a swap the require also catches; the additive-contradiction row shape exists at the `additive direct working branch permission` subtest; the spec's anchor-inventory paragraph still says 24 needles while the repairs added three. Re-derive from the tree at pickup.

## What to build

FT164 repair round 3: the round-2 review found the RT1 basename needle
undefended and the forbid unexercised additively, and the spec's needle
enumeration — its own designated audit source — under-counting by three. Two
mutation rows land in the bite harness, and the spec's inventory paragraph
states the true per-section counts and total.

## Acceptance

- [ ] [PN1] a deletion mutation row proves the basename require needle's own diagnostic fires when the basename sentence is deleted.
- [ ] [PN2] an additive mutation row proves the forbid fires when a title-keyed sentence is added beside the intact basename sentence.
- [ ] [PN3] the spec's anchor-inventory paragraph counts the repair needles per section and totals 27.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PN1 | delete the basename require needle's new table row | the `basename blocker` mutation subtest | remove the row, run the mutation harness, expect the basename diagnostic unproven (row absent means the subtest disappears — review confirms presence; the row itself proves the needle) |
| PN2 | drop the additive row and keep only the swap | the `title blocker additive` mutation subtest | remove the additive row, run the harness, same review-confirmed presence; the row itself proves the forbid bites on addition |
| PN3 | leave the inventory at 24 | review | read the inventory paragraph against the registered needle count |
