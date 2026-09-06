# Probe the fixture source on the coverage axis

Blocked by: 03-probe-the-production-adapter-through-its-grading-test.md
Writes: .agents/skills/bench-craft-review/references/finding-discipline.md
Covers: none

## What to build

Add one rule to `.agents/skills/bench-craft-review/references/finding-discipline.md` under `Where an axis under-reads`. The rule reads: The Coverage axis probes a test's fixture source, not only its assertion. A fixture that names a symbol the production file declares can stay green while the assertion never runs. Write it in ASD-STE100, in the shape of its neighbouring items, and add no other sentence. Source: roadmap row FT298.

## Acceptance

- [ ] the rule sits as one list item under `## Where an axis under-reads`.
- [ ] the item names the fixture source as the probe target.
- [ ] `bench gate-prose` passes on the edited file.
