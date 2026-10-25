# Record each repair red before the confirming round

Blocked by: none
Writes: .agents/skills/bench-craft-delegate/references/delegation-discipline.md, .agents/commands/bench-implement-spec.md
Covers: none

## What to build

The author charge asks for a witness: one demonstrated red of one independent test expectation.
The repair-charge template asks for no witness return.
In two repair cycles, a repair session demonstrated the reds of its new test expectations.
Only the planned probe reached the review record, so a Standards finding stayed open and the chunk took a record-only round.

The repair-charge template gets a witness return field.
That field asks for the demonstrated red of each independent test expectation that the repair adds.
Each red names the mutation, the killing test, and the restore.

The Land section of the implementation command states that the review record holds each red of that return before the confirming round's charge.
`bench record verification --probe-result <file>` retains the red of a planned probe requirement.
The orchestrator writes each other red as prose in the review record outside the fenced record block.
The implementation command is at its line budget, so the diff joins two wrapped lines of that file and changes none of their words.

The project instructions count an independent test expectation only when its red is recorded and demonstrated.
That rule stays in the project instructions, because each linked repo owns its own instructions.

## Acceptance

- [ ] The repair-charge template has a `Witness return:` field that asks for each demonstrated red of an independent test expectation that the repair adds.
- [ ] The Land section states that the review record holds each demonstrated repair red before the orchestrator prepares the confirming round's charge.
- [ ] The Land section names only the existing `bench record verification --probe-result <file>` form for a planned probe red.
- [ ] The implementation command stays inside its prose line budget, and the budget table does not change.
- [ ] Every anchor in `internal/anchors` stays green, and the diff changes no Go file.
