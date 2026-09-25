# Align the drain's commit wording with the lane and landing contract

Blocked by: none
Writes: .agents/commands/bench-drain.md, internal/anchors/registry_data.go, internal/anchors/registry_retained_workflow.go, tests/canary/workflow-guidance-anchors/drain-anchor/files/dot-agents/commands/bench-drain.md
Covers: none

## What to build

This ticket is the drain-wording half of FT89. The drain says "commit on green" and
says that the gate is the cost of a commit. The operating guide requires a lane pass for
a worktree commit and the whole-project gate for the landing. So the drain states the
wrong price for its own batch commit.

The drain says that the pass commits once on a lane pass and lands through the gate. An
implement-now item lands green through its own landing. Each landing costs one gate.

The anchor rows that pin the old wording change with it. The ungrouped Require row
and the implement-now Require row take the new text. Two Forbid rows keep the retired
wording out of the drain. The drain canary fixture carries the new Require text, so its
only miss stays the one that its expectation names.

The section 8 heading keeps its text, because the recurrence contract check finds that
section by its exact heading. The retire-hint half of FT89 waits on the FT284 decision
and stays on the row.

## Acceptance

- [ ] The drain names a lane pass for its batch commit and the gate for its landing. It no longer says "commit on green" or "The gate is what a commit costs".
- [ ] The anchor registry requires the new wording and forbids the two retired phrases in the drain.
- [ ] A probe that restores "commit on green" in the drain turns the root conformance red.
- [ ] `bench test --changed` and the drain canary pass, and the gate is green at the landing.
