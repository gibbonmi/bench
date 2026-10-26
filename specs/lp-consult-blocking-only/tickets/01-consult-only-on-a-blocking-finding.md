# Consult only on a blocking finding or an open behavior

Blocked by: none
Writes: .agents/skills/bench-craft-line/SKILL.md
Covers: none

## What to build

`craft-line`'s "Delegated author lines" tells the orchestrator to decide an advisory finding itself. Only a blocking finding or a behavior that the spec leaves open takes a delegated diagnostic consultation. The consultation pair and its top-binding tiebreak stay as they are.

In the gate-reuse GR-C1 review, all ten findings were advisory, and the consultation pair and the tiebreak ran one after the other for about ten minutes.

## Acceptance

- [ ] The section states the rule, and the skill stays inside its line budget.
- [ ] `bench gate` passes at the landing.
