# Pin the restored constraints with anchors, fixtures, and a budget row

Blocked by: 01-restore-the-debug-phase-constraints.md
Writes: internal/anchors/registry_data.go, tests/canary/workflow-guidance-anchors/, projects/benchkit.md

## What to build

Make the restoration durable at the two conformance seams that already grade
guidance prose.

- Add `Require` anchor rows in `internal/anchors/registry_data.go` for the
  restored needles in `.agents/commands/bench-debug.md`: the two stop-gate
  sentences, the `Tighten the loop` step name, the reference-file pointer, the
  `A green proxy only narrows a hypothesis` sentence, and the `- [ ]` spelling
  of the first Phase 1 criterion and the first Phase 2 confirmation. Add two
  `Require` rows over
  `.agents/skills/bench-debug/references/loop-constructions.md` for the menu
  needles `Bisection harness` and `Property / fuzz loop`. Existing debug rows
  stay untouched.
- Add two `workflow-guidance-anchors` canary fixtures, each with `BASE` copying
  the phase file, a single-occurrence `MUTATE.json` anchor, and an `EXPECT`
  naming the new anchor's exact diagnostic: one softens the Phase 1 stop-gate
  sentence, one deletes the reproduction-economics sentence.
- Add an exact `.agents/commands/bench-debug.md` row to the
  `Guidance prose budgets` table in `projects/benchkit.md`, deriving the limit
  from the file ticket 01 landed at the table's 0-5-line headroom convention
  (the spec estimates 170 against ~165 restored lines).
- No `registry_test.go` edit: the `workflow-guidance-anchors` family already
  has its `canaryFixtureFamilyRegistry` registration.

The contract with ticket 01: every needle registered here is a sentence
ticket 01 landed, occurring exactly once in the file.

## Acceptance

- [ ] The kit gate is green with the new anchor rows over the real tree (DR1–DR9 needles all match).
- [ ] Both new fixtures red through the registered owner with their exact diagnostics, and each restore re-runs green (`TestEveryRetainedFixtureBitesThroughRegisteredOwner` passes) (DP1, DP2).
- [ ] The budget table carries the exact bench-debug row and the restored file passes under it (DP3).
