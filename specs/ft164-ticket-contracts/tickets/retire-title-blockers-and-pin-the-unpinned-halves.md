# Retire title blockers and pin the unpinned halves

Blocked by: point-craft-spec-at-the-ticket-contracts.md
Ownership fence: `.agents/skills/bench-craft-tickets/SKILL.md`, `internal/conformance/docs_workflow_helpers_test.go`, `internal/conformance/fixture_bite_test.go`
Assumptions: step 3 of the breakdown method still reads "by sibling ticket title" while the template and field bullets key on basenames; the AB1 placeholder needle and mutation row leave the AB2 row unpinned; the junction needle pins the relocation half only; `forbidInSection` exists beside `requireInSection`. Re-derive from the tree at pickup.

## What to build

FT164 repair round: the review found the breakdown procedure contradicting the
basename contract, and the falsification pass found two taught sentences whose
deletion or partial revert leaves the gate green. Step 3 of the numbered method
names blockers by sibling ticket file basename, with the title form forbidden
in that section; the second template placeholder row and the junction-creation
half of the junction rule each gain their own section-scoped needle and
byte-exact mutation row.

## Acceptance

- [ ] [RT1] step 3 of the breakdown method names blockers by sibling ticket file basename, and the phrase "by sibling ticket title" is forbidden inside the Draft-the-breakdown section.
- [ ] [RT2] the AB2 placeholder row has a section-scoped needle and a mutation row, so reverting either template row alone reds the gate.
- [ ] [RT3] the junction-creation sentence has its own needle and mutation row, so deleting it while keeping the relocation half reds the gate.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RT1 | restore "by sibling ticket title" in step 3 | the `title blocker forbidden` mutation subtest | swap the basename phrase back to the title form, run the anchor check, expect the forbid diagnostic |
| RT2 | revert the AB2 row to the old unlabeled placeholder | the `template row two` mutation subtest | swap AB2's line for the unlabeled form, run the anchor check, expect the second-row diagnostic |
| RT3 | delete the junction-creation sentence, keep relocation | the `junction creation` mutation subtest | delete the sentence, run the anchor check, expect the junction-creation diagnostic |
