# Resolve the implement-spec guidance contradictions

Blocked by: none
Writes: .agents/commands/bench-implement-spec.md, .agents/skills/bench-craft-line/SKILL.md, .agents/skills/bench-craft-tdd/SKILL.md, .agents/skills/bench-craft-delegate/SKILL.md
Covers: none

## What to build

A reader of `/bench-implement-spec` gets one instruction for each decision the
phase makes. The ladder tells a retained author what to do after a second
diff-owned red, and it never changes the implementation model in silence. The
phase names `/bench-review-implementation` as the owner of the three axes, and it defines
the diff condition that stops a `--full` run. Gate anchors pin much of this prose, so
each clarification arrives as a new sentence beside the pinned words.

## Acceptance

- [ ] `craft-line` tells a retained author to raise effort, and to ask the reviewer before a tier move.
- [ ] `craft-line` sets low effort as the repair default, and it permits a higher effort for a repair at risk.
- [ ] The Build section sends the author from the focused checks to the commit, and it does not end the ticket at the checks.
- [ ] The charge paragraph names the ticket graph as the approval unit, and it names the author of the supplement.
- [ ] The Exit handoff states the landing route once, and the "Land" section keeps the detail.
- [ ] Every gate-anchored sentence keeps its exact words; a new sentence carries each clarification.
- [ ] The phase names `/bench-review-implementation` where the chunk review starts.
- [ ] The `--full` pause names an observable diff condition.
- [ ] `craft-tdd` contains no first-person voice, and it points to `craft-line` for the declaration.
- [ ] `craft-delegate` separates the pre-approved diagnostic consultation from the ladder's top-tier pause.
