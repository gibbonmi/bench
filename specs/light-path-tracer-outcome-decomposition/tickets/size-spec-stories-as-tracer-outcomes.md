# Size spec stories as tracer outcomes

Blocked by: none
Ownership fence: `.agents/commands/bench-write-spec.md`, `.agents/skills/bench-craft-spec/SKILL.md`
Contracts: the story-sizing guidance as ordered prose crosses `.agents/skills/bench-craft-spec/SKILL.md`→`.agents/commands/bench-write-spec.md`; STO1 asserts the same tracer-outcome domain rule, STO2 asserts implementation-skill reading before story lock, and absence means stories are not locked

## What to build

Spec authors size stories as independently deliverable and demonstrable tracer outcomes. Seams remain the places where tests attach and ownership fences fall; `craft-tickets` owns the later ticket slicing and the spec does not pre-decide ticket boundaries. Before stories lock, the phase reads `craft-tickets`, `craft-delegate`, `craft-tdd`, and `craft-seams` so the planned outcomes reflect how implementation will be cut, charged, and verified.

## Acceptance

- [ ] [STO1] both spec surfaces reject a horizontal engineering layer wearing a story name and require each story to deliver and demonstrate a complete outcome on its own.
- [ ] [STO2] `/bench-write-spec` reads `craft-tickets`, `craft-delegate`, `craft-tdd`, and `craft-seams` before locking stories.
- [ ] [STO3] seams remain responsible for test attachment and ownership fences, while `craft-tickets` remains the single owner of ticket slicing.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| STO1 | size one story as a parser layer with no deliverable behavior | the semantic reviewer | apply the story-sizing rule, expect the horizontal layer to be rejected |
| STO2 | remove one implementation skill from the pre-lock read set | the semantic reviewer | enumerate the four required skills in the command step, expect the missing skill to make the phase incomplete |
| STO3 | tell the spec to split one ticket per seam | the consistency review | read both surfaces with `craft-tickets`, expect duplicate and conflicting slicing ownership |
