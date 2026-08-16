# Realign the consumers, glossary, and docs

Blocked by: project-one-row-shape-across-schemas.md, remake-craft-spec-and-craft-tickets-on-their-sources.md
Writes: .agents/skills/bench-craft-delegate/SKILL.md, .agents/skills/bench-craft-review/SKILL.md, .agents/commands/bench-implement-spec.md, CONTEXT.md, CHANGELOG.md, docs/field-guide.html, docs/reporesident-distillation.md, .bench/BENCH-reference.md

## What to build

Every consumer describes the shape it actually receives. `craft-delegate`'s
write-delegation charge carries behavior and seam rather than a red signal;
`craft-review`'s Coverage axis and `/bench-implement-spec`'s task-seeding prose
name the reduced projection. `CONTEXT.md` gains glossary entries for **coverage
row** and **acceptance row** with Avoid lists — landing in this change, not
ahead of it, so the glossary never describes a tree that does not exist. The
field guide and the distillation doc stop teaching the retired column, and
CHANGELOG carries one entry for the reduced schema and the single review round.
`bench skills-index --check` is green after the skill descriptions change.

## Acceptance

- [ ] No guidance file describes the coverage map as carrying a red signal.
- [ ] `craft-delegate`'s charge names behavior and seam.
- [ ] `CONTEXT.md` defines **coverage row** and **acceptance row**, each with an
      Avoid list, and stays glossary-only.
- [ ] The field guide and distillation doc describe the reduced schema.
- [ ] CHANGELOG carries the entry; `bench skills-index --check` is green.
- [ ] `bench gate` is green over the whole change.
