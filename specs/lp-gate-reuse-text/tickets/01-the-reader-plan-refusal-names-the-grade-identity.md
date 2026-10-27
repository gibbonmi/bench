# The reader plan refusal names the grade identity

Blocked by: none
Writes: internal/gate/reader_plan.go, internal/gate/reader_plan_test.go, internal/gate/complete_checkpoint.go, CHANGELOG.md
Covers: none

## What to build

The reader plan compares the grade identity of the subject with the identity of the reader table.
When the two differ, the plan refuses the reader set, but the case in the refusal line names the oracle identity.
The oracle identity is a different value, so the line names an identity that the comparison does not read.

The refusal case now names the grade identity, so the line reads `gate: whole gate (reader set refused: new grade identity)`.
The reader-plan test of that case now pins the exact line, so a change to the label turns the test red.

The doc comment of `runUnderBroker` has a line break before the word `It` in its last sentence.
The comment now wraps without that break, and the code does not change.

A `CHANGELOG.md` entry under `[Unreleased]` records the changed refusal line.

## Acceptance

- [ ] A refused reader plan prints `new grade identity` as the case of a table with a different grade identity.
- [ ] A search of the tree, without `.git`, finds no copy of the old label.
- [ ] The reader-plan test of the identity case pins the exact refusal line, and it is red against the old label.
- [ ] The doc comment of `runUnderBroker` has no orphan line break, and the edit changes no code.
- [ ] `CHANGELOG.md` has one entry under `[Unreleased]` that states the changed refusal line.
