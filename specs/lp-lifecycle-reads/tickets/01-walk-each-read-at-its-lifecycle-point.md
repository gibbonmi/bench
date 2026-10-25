# Walk each read at its lifecycle point

Blocked by: none
Writes: .agents/skills/bench-craft-spec/references/map-discipline.md, CHANGELOG.md
Covers: none

## What to build

The FT445 spec review rejected iteration 1 for two reads that the spec did not walk.
The prepared run kept the tip-moved guard, and the record commit during the run voids the HEAD read of that guard.
The delta eligibility read the chunk tip in the record, and that tip still names the earlier freeze tip at a repair commit.
The `Workflow-step writes` class of the pre-review proof checklist traces only the writes of a new workflow step.
The class now also lists each read in the callee chain of the step.
For each read, the class names the HEAD, the record state, and the capture source at the lifecycle point of the step.

## Acceptance

- [ ] The `Workflow-step writes` class lists each read in the callee chain of each new workflow step.
- [ ] For each read, the class names the HEAD, the record state, and the capture source at the lifecycle point of the step.
- [ ] The class label and the other checklist classes stay the same.
- [ ] `CHANGELOG.md` has one entry for the changed class.
