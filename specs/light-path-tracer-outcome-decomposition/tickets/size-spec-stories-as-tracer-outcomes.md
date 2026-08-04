# Size spec stories as tracer outcomes

Blocked by: none
Ownership fence: `.agents/commands/bench-write-spec.md`, `specs/light-path-tracer-outcome-decomposition/tickets/size-spec-stories-as-tracer-outcomes.md`
Contracts: the canonical story-sizing rule crosses `.agents/skills/bench-craft-spec/SKILL.md`→`.agents/commands/bench-write-spec.md` by named reference; STO1 asserts canonical ownership and post-seam ordering, STO2 asserts implementation-skill reading before story lock, and absence means stories are not locked

## What to build

`craft-spec` owns the canonical story-sizing rule. After seams are explicit, `/bench-write-spec` locks stories by applying that named rule instead of restating it. Before stories lock, the phase reads `craft-tickets`, `craft-delegate`, `craft-tdd`, and `craft-seams` so the planned outcomes reflect how implementation will be cut, charged, and verified. Seams retain test attachment and ownership fences, while `craft-tickets` retains later ticket slicing.

## Acceptance

- [ ] [STO1] `craft-spec` is the single canonical owner of story sizing, and `/bench-write-spec` applies its named rule only after seams are explicit without restating it.
- [ ] [STO2] `/bench-write-spec` reads `craft-tickets`, `craft-delegate`, `craft-tdd`, and `craft-seams` before locking stories.
- [ ] [STO3] seams remain responsible for test attachment and ownership fences, while `craft-tickets` remains the single owner of ticket slicing.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| STO1 | replace the post-seam named reference with a pre-seam copy of the canonical rule | the semantic reviewer | inspect canonical ownership and command order, expect the duplication and early lock to be rejected |
| STO2 | remove one implementation skill from the pre-lock read set | the semantic reviewer | enumerate the four required skills in the command step, expect the missing skill to make the phase incomplete |
| STO3 | tell the spec to split one ticket per seam | the consistency review | read both surfaces with `craft-tickets`, expect duplicate and conflicting slicing ownership |
