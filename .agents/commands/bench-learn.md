---
description: Learnings integration — drain .bench/learnings.md, the journal the agent appended to during real use, and propose promotions. Runs each candidate through the bench-craft-synthesis discipline; lands promotions in the rule/skill/command they fix, adds a CHANGELOG line, and marks the source entries resolved. Proposes, never auto-merges. Maintenance, not a workflow phase.
disable-model-invocation: true
---

# /bench-learn — drain the learnings journal into the kit

Fold *your own* changes — the friction the agent logged during real use — back into
the kit. This is the internal, usage-driven input — the counterpart to `/bench-update`,
which syncs against upstream. It has nothing to do with the source repos. It
**proposes; you merge.**

## 1. Drain the journal (the input)

- Read `.bench/learnings.md`. Each **open** entry is a candidate sourced from how the
  kit actually behaved in real work — a deviation, a should-have-asked, a recurring
  friction. Tag every candidate `learnings`.
- Group entries that point at the same fix so one change resolves several, and note
  any that are one-off context rather than a general rule (those get dismissed, not
  promoted).

## 2. Run the synthesis discipline

Hand the candidates to `bench-craft-synthesis`; it owns the discipline end to end
(respect closed decisions, assess, the three quality loops, propose-don't-merge).
Don't restate it here.

## 3. Record (the output)

Only after the loops pass and I've signed off, apply and record: land each promotion in
the rule, skill, or command it fixes, add a `CHANGELOG.md` line, and mark the source
entries in `.bench/learnings.md` resolved (promoted or dismissed, one line of why) so
they're never re-reviewed.

The merge is mine; never auto-apply.
