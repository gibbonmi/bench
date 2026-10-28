# Add the trace times and the drain budget rule

Blocked by: 03-state-the-probe-only-rule.md, 04-derive-the-retro-cost-section.md
Writes: internal/roadmap/retro_cost.go (new), internal/roadmap/retro_cost_test.go (new), .agents/commands/bench-drain.md, internal/anchors/registry_build_cost.go (new), internal/anchors/registry_build_cost_test.go (new), internal/roadmap/retro_scaffold.go, internal/spec/landing.go, internal/spec/landing_test.go (new), internal/otelrecord/reader.go, internal/otelrecord/registry.go
Covers: BC14, BC19, BC20, BC21, BC22, BC71

## What to build

The cost section gains one trace line. It reads the record of the Bench home through `otelrecord.ReadSpans`. It keeps each span whose trace holds a span with an assignment ID of the landed version 2 plan.
The filter line in `retro_cost.go` reads `if ids[span.Attributes[otelrecord.AttrAssignmentID]] {`, because the probe-05 command swaps that exact line.
For the kept traces, the line states the count and the summed milliseconds of the commit, lane, and landing spans. An absent record, a version 1 plan, and a record without a span of the plan's assignments each print `unknown`.

Section 4 of `.agents/commands/bench-drain.md` gains two sentences:

- "When a retro's `## Cost` section marks a line `over budget`, propose one removal or one merge of a build step for that line."
- "A budget breach never makes the gate red."

The drain guidance repeats no budget number. Add the Forbid anchors of BC14 to `registry_build_cost.go`. BC22 and BC71 are review-owned, because ADR 0012 keeps rule wording out of the anchor registry.

## Acceptance

- [ ] The trace line states the counts and the summed milliseconds of the plan's assignment (BC19).
- [ ] A commit span of another assignment adds nothing (BC20).
- [ ] An absent record prints `unknown` (BC21).
- [ ] The drain section holds both sentences (BC22, BC71).
- [ ] `bench-drain.md` holds neither budget number (BC14).
