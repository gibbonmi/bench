# Verify the headroom route before a charge dispatches

Blocked by: none
Writes: .agents/skills/bench-craft-delegate/references/delegation-discipline.md
Covers: none

## What to build

A drained learning names a gap. The FT310 ticket named an over-budget write
target with no headroom route, and the commit lane refused on growth. That
gap is a missing verification step, not a missing rule: `craft-tickets`
already states that a ticket names its headroom route. Nothing runs
`bench structure` to confirm it before the charge dispatches.

Add one rule to `delegation-discipline.md`, under `## In the charge`. Read
the whole file first. A prior fold already removed one over-budget-route
bullet from this section. Do not restate that rule. This ticket adds
verification, not naming.

Add this rule, close to this wording:

A charge whose ticket names a write target runs `bench structure` against
that target before dispatch. The run confirms the ticket's stated headroom
route, when the path is over budget.

Every sentence obeys `references/ste-prose.md`: 25 words per sentence, six
sentences per paragraph. This file is not in the guidance-prose-budget
table (only `.agents/skills/*/SKILL.md` is budgeted), so no headroom route
applies here.

## Acceptance

- [ ] `delegation-discipline.md`'s "In the charge" section states the `bench structure` verification rule.
- [ ] The new rule does not restate the headroom-route-naming rule `craft-tickets` already owns.
- [ ] `bench gate-prose` passes on the touched file.
