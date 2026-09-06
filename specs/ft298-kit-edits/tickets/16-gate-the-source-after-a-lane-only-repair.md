# Gate the source after a lane only repair

Blocked by: 15-record-each-dogfood-run-before-the-repair-re-review.md
Writes: .agents/commands/bench-final-check.md
Covers: none

## What to build

Add one rule to `.agents/commands/bench-final-check.md` under `Exit handoff`. The rule reads: After a lane-only repair commit and before the landing, run the whole-tree gate on the source. The lane skips the conformance checks the landing gate runs. Write it in ASD-STE100, in the shape of its neighbouring items, and add no other sentence. Source: roadmap row FT298.

## Acceptance

- [ ] the rule sits in the Exit handoff section beside the landing paragraph.
- [ ] the text names the whole-tree gate on the source.
- [ ] `bench gate-prose` passes on the edited file.
