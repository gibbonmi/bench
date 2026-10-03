# Start each repair cycle with debug

Blocked by: none
Writes: .agents/commands/bench-debug.md, .agents/skills/bench-craft-line/SKILL.md, .agents/skills/bench-craft-line/references/bounded-repair-policy.md, .agents/skills/bench-craft-delegate/references/delegation-discipline.md, internal/anchors/registry_debug_loop.go, internal/anchors/registry_retained_workflow.go, internal/conformance/implementation_continuation_test.go
Covers: none

## What to build

The reviewer decided that each repair starts with the bug path. The author
diagnoses the failure and records its cause before the repair edit. The rule
applies to a fix-loop repair, a post-review repair cycle, and a fix-and-gate
repair.

The "How it meets the rest of Bench" section of the bug path owns the rule.
A diff-owned red names the fix-loop repair. A review finding names both
review repairs. Two anchors in the debug-loop family hold the rule there.
The file stays inside its prose budget, because the exit handoff loses its
restatement of the first invariant.

The `craft-line` continuation section lets the ticket author invoke debug only
after a reassessment. That sentence becomes a pointer to the owner. Its anchor
and its conformance expectation change with it. The bounded repair policy also
points to the owner.

The repair-charge template keeps its `Debug route:` field. The field tells the
repair author to start with the bug path and to return the recorded cause.

## Acceptance

- [ ] The bug path's integration section states that each repair of a diff-owned red or a review finding starts with the bug path.
- [ ] The `craft-line` continuation section and the bounded repair policy link to the integration section and do not restate the rule.
- [ ] The repair-charge template requires the bug path at the start and the recorded cause in the return.
- [ ] A deletion of the owner sentence turns the anchor check red.
- [ ] The prose gate on each edited Markdown file and the root conformance pass stay green.
