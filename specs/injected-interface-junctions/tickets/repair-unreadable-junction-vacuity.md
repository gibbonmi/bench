# Repair the vacuous unreadable-metadata junction and its spec row

Blocked by: none
Ownership fence: `internal/specbuild/abandon_test.go`, `specs/injected-interface-junctions/spec.md`
Contracts: the spec's story-2 unreadable coverage row crosses spec→test suite; after this ticket it claims pre-composition refusal (fake-driven test + internal/worktree planner tests), not real-planner composition, asserted by RU2 against the tree
Assumptions: sol review round 1 findings SP1/C1/S1 are the authority for this repair; the Service refuses unreadable metadata at precondition classification before any owner call, so no owner composition exists for this shape; claims re-derived from the tree at pickup

## What to build

Delete `TestAbandonRefusesUnreadableCheckoutMetadataThroughRealPlanner` (it
asserts the planner is never reached, duplicating the fake-driven refusal
test); trim the `decayedOwner` comment's reference to it; amend the spec's
story-2 unreadable row to its honest disposition — refusal is pre-composition,
covered by the existing fake-driven test, with the planner's own unreadable
contract pinned in internal/worktree — and add a Won't-handle line naming why
the composition cannot be tested (it does not exist).

## Acceptance

- [ ] [RU1] the vacuous twin is gone, the `decayedOwner` comment names only the decayed/husk junction tests, and `go test ./internal/specbuild` is green.
- [ ] [RU2] the spec's unreadable row and Won't-handle line state the pre-composition refusal honestly and `bench coverage --check` passes.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RU1 | restore the deleted test name in the comment | review grep | `rg ThroughRealPlanner internal/specbuild/abandon_test.go` names only decayed/husk tests |
| RU2 | restore the old row text claiming real-planner composition | the coverage validator plus review | the row's claim would again contradict `internal/specbuild/precondition.go`'s pre-owner refusal; sol finding SP1 is the recorded red |
