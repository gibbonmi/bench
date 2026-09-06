# Enumerate the exact match tests on a changed message

Blocked by: 11-enumerate-the-importers-a-widened-pattern-reds.md
Writes: .agents/skills/bench-craft-spec/references/map-discipline.md
Covers: none

## What to build

Add one rule to `.agents/skills/bench-craft-spec/references/map-discipline.md` under `Before the map locks`. The rule reads: A spec that changes a rendered message enumerates the exact-match tests on that text. Write it in ASD-STE100, in the shape of its neighbouring items, and add no other sentence. Source: roadmap row FT298.

## Acceptance

- [ ] the rule sits as one list item under `## Before the map locks`.
- [ ] the item names exact-match tests as the enumeration.
- [ ] `bench gate-prose` passes on the edited file.
