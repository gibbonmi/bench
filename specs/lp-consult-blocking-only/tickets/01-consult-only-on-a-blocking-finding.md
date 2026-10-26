# Consult only on a blocking finding that would stop for the reviewer

Blocked by: none
Writes: .agents/skills/bench-craft-line/SKILL.md
Covers: none

## What to build

`craft-line`'s "Delegated author lines" limits a delegated diagnostic consultation to one case: a blocking finding that would otherwise stop for the reviewer. The orchestrator decides every other finding itself, which includes each advisory finding and each repair. The consultation pair and its top-binding tiebreak stay as they are.

In the gate-reuse GR-C1 review, all ten findings were advisory, and the consultation pair and the tiebreak ran one after the other for about ten minutes. The reviewer stated that the consultants replace a reviewer stop only.

## Acceptance

- [ ] The section states the rule, and the skill stays inside its line budget.
- [ ] `bench gate` passes at the landing.
