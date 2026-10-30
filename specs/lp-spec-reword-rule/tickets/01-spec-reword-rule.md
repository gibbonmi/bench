# A reworded spec sentence keeps each predicate of the old sentence

Blocked by: none
Writes: .agents/skills/bench-craft-spec/SKILL.md, ROADMAP.md, roadmap/FT417.md
Covers: none

## What to build

Roadmap row `FT417` decides one spec-authoring rule.
A spec author who rewords an anchored lead can drop a rule clause, and the spec review catches it only in a second iteration.

Step 3 of the `craft-spec` authoring process gains the reword rule.
Unless the decision source deletes a clause, a reworded sentence keeps each predicate of the old sentence.
New content joins as new sentences.

The `craft-spec` review rubric gains one sentence that points to that rule, so the spec review round grades each reword against it.
The rule has one source; the rubric sentence does not restate it.
The rubric keeps its six questions, because the `dg-15-mandatory-red-capitalized` canary pins the text `The round asks six questions:`.

Retire `FT417` in the same commit: delete `roadmap/FT417.md` and its heading line in `ROADMAP.md`.

## Acceptance

- [ ] Step 3 of the authoring process in `.agents/skills/bench-craft-spec/SKILL.md` states the reword rule.
- [ ] The review rubric in `.agents/skills/bench-craft-spec/SKILL.md` keeps its six questions, and a new sentence after them points to the reword rule.
- [ ] `.agents/skills/bench-craft-spec/SKILL.md` stays inside its guidance prose budget.
- [ ] `roadmap/FT417.md` does not exist, and `ROADMAP.md` holds no `FT417` heading.
- [ ] The `bench commit` lane passes.
