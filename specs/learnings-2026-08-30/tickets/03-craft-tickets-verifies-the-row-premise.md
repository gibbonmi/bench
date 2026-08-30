# Make craft-tickets verify a roadmap row's premise before the ticket

Blocked by: none
Writes: .agents/skills/bench-craft-tickets/SKILL.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go

## What to build

On 2026-08-30 a ticket author implemented the FT223 row's decided fix as
written. The row called an inherited verdict a non-failure. The code showed
the opposite: the gate ran red and no green baseline attributed the red. The
mismatch surfaced only because the author read the kind's definition before
the acceptance rows were written. `craft-tickets` states no such step.

Add one rule to the `Draft the breakdown` section of
`.agents/skills/bench-craft-tickets/SKILL.md`, in ASD-STE100 prose. A ticket
that implements a roadmap row's decided fix first verifies the row's premise
against the code. It reads the definition of every named kind, state, or
error the row builds on. When the premise and the code disagree, the author
surfaces the mismatch as a reviewer decision before the ticket is written.
The author does not implement the row's fix as written.

Pin the rule's key sentence in `internal/anchors/registry_data.go` beside the
existing `bench-craft-tickets` needles, with a diagnostic in the same shape.
Show the pin red on removal in `internal/anchors/registry_data_test.go` the
way the existing pin tests do. The file stays inside its
`projects/benchkit.md` budget of 100 lines. Every other anchored sentence in
the file keeps its bytes and line breaks.

## Acceptance

- [ ] `.agents/skills/bench-craft-tickets/SKILL.md` states that a ticket for a roadmap row's decided fix verifies the row's premise against the code before the ticket is written.
- [ ] The same section states that a premise mismatch is a reviewer decision, not a fix to implement as written.
- [ ] `internal/anchors/registry_data.go` pins the key sentence, and `internal/anchors/registry_data_test.go` shows the pin red on removal.
- [ ] `TestGuidanceProseBudgetsHoldOnTheLiveTree` stays green with no change to the budget table.
- [ ] `TestProseMechanicsHoldsOnTheLiveTree` and the fixture-bite tests stay green.
- [ ] `bench gate-prose . -- .agents/skills/bench-craft-tickets/SKILL.md` passes.
