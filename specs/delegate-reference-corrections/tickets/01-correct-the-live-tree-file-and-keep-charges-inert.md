# Correct the live-tree inventory path and keep reviewer charges inert

Blocked by: none
Writes: specs/delegate-reference-corrections/tickets/01-correct-the-live-tree-file-and-keep-charges-inert.md (new), .agents/skills/bench-craft-delegate/references/delegation-discipline.md, .agents/skills/bench-craft-delegate/references/cross-harness-reviewers.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_calibration.go, internal/anchors/registry_debug_loop.go, tests/canary/workflow-guidance-anchors/delegate-live-tree-inventory-fence/, tests/canary/workflow-guidance-anchors/cross-harness-reviewer-exec-heredoc/
Covers: none

## What to build

The `craft-delegate` references give two incorrect instructions. This ticket
corrects both.

The live-tree test rule in `delegation-discipline.md` tells a charge to fence
`internal/conformance/tier_test.go`. The live-tree inventory is in
`internal/conformance/tier_live_tree_test.go`. The rule names that file. Its
anchor row, its test expectation, and its canary move to the new path.

The two direct reviewer recipes in `cross-harness-reviewers.md` put the charge
in a double-quoted argument. In that argument, the coordinator's shell runs each
backtick and each `$()` as a command substitution, and a Bench charge cites
commands in backticks. Each direct recipe reads its charge from a file on stdin.
The shell does not expand that text, and the file closes stdin at its end. The
paragraph above the recipes states this rule one time. The anchor rows move to
the new recipes.

The two `bench worktree exec` reviewer forms also put the charge in a
double-quoted argument. Each exec form carries the charge in the body of a quoted
heredoc. The quoted delimiter keeps the body inert, and the heredoc closes the
child's stdin at its end. The anchor row, its test expectation, its diagnostic,
and the `cross-harness-reviewer-exec-heredoc` canary move to the new form.

The anchor registry file and its test file are over their line budget, so this
ticket adds no row to them. Forbid rows for the retired text need a new anchor
registry file, which is outside this ticket.

## Acceptance

- [ ] The live-tree test rule in `delegation-discipline.md` names `internal/conformance/tier_live_tree_test.go`. The anchor check fails when the rule names `internal/conformance/tier_test.go` in its place.
- [ ] The two direct reviewer recipes in `cross-harness-reviewers.md` read the charge from a file on stdin. The anchor check fails when a recipe puts the charge in a double-quoted argument in its place.
- [ ] The two `bench worktree exec` reviewer forms carry the charge in the body of a quoted heredoc. The anchor check fails when the Claude exec form puts the charge in a double-quoted argument again.
- [ ] The `delegate-live-tree-inventory-fence` canary mutates the current path text and bites through its owning check.
- [ ] The `cross-harness-reviewer-exec-heredoc` canary mutates the current exec form and bites through its owning check.
