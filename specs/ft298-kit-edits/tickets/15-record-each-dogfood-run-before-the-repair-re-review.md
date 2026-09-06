# Record each dogfood run before the repair re review

Blocked by: 14-state-the-landing-base-apart-from-the-review-base.md
Writes: .agents/commands/bench-review-implementation.md
Covers: none

## What to build

Add one rule to `.agents/commands/bench-review-implementation.md` under `Review modes`. The rule reads: The coordinator records every dogfood run in the spec before the repair-scoped re-review starts. An unrecorded run is a blocking finding. Write it in ASD-STE100, in the shape of its neighbouring items, and add no other sentence. Source: roadmap row FT298.

## Acceptance

- [ ] the Review modes section states the rule after the repair-ticket paragraph.
- [ ] the text names the unrecorded run as a blocking finding.
- [ ] `bench gate-prose` passes on the edited file.
