# Add the batch-approval amendment carve-out to craft-spec

Blocked by: none
Writes: .agents/skills/bench-craft-spec/SKILL.md, projects/benchkit.md
Covers: none

## What to build

`roadmap/FT298.md` records one reviewer decision, now resolved. A batch
approval licenses a narrow exception to the rule that a build may not edit
its own acceptance rows. Add that exception to
`.agents/skills/bench-craft-spec/SKILL.md`, right after the sentence "A
build may not edit its own spec's acceptance rows, budget targets, or
ownership fences." (currently line 62).

Add a new paragraph, separate from the existing one, with these four
sentences, in this order and close to this wording:

1. A batch approval licenses one narrow exception to that rule.
2. A build may amend its own acceptance row when the code contradicts the
   row's literal premise.
3. The amendment keeps the row's verdict unchanged and cites the
   contradicting evidence under Build decisions.
4. It never licenses a build to loosen what the row counts as passing.

Do not name a specific past spec or roadmap row inside this rule. A named
example rots once that spec retires. State the rule generically.

Every sentence obeys `references/ste-prose.md`: 25 words per sentence, six
sentences per paragraph. Check `.agents/skills/bench-craft-spec/SKILL.md`'s
guidance-prose-budget row in `projects/benchkit.md` (currently 152 lines;
the file is presently 151 lines, so it has almost no headroom). If your
edit pushes the file over budget, raise that row's number in
`projects/benchkit.md` by the exact overage. Make that change in the same
commit. Use the same headroom route `craft-tickets` uses for a code file.
State whether you needed the headroom route.

## Acceptance

- [ ] `craft-spec/SKILL.md` states the four-sentence exception as its own paragraph, directly after the "may not edit" sentence.
- [ ] The new paragraph names no specific spec or roadmap row.
- [ ] `bench gate-prose` passes on every touched file.
- [ ] `bench test --check guidance-prose-budgets` passes.
