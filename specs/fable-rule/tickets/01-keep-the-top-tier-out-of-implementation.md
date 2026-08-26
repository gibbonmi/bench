# Keep the top tier out of implementation

Blocked by: none
Writes: .agents/skills/bench-craft-line/SKILL.md, projects/benchkit.md, capture/agent-performance/claude-models.md, specs/exec-census/spec.md

## What to build

The reviewer directed on 2026-08-26 that the top tier implements nothing,
code or guidance prose, unless the reviewer names it for the run. Today the
leverage override in `craft-line` routes guidance prose to top + high, the
decision table routes an uncertain seam to top + high, and the profile's
`Lines` section repeats the override. Move the leverage override to mid +
high in both places, and route the uncertain-seam row to mid + high with an
offer of the top tier to the reviewer. The escalation ladder's pause on any
top-tier bump stays as the enforcement. The scorecard records the decision
in its current decisions. The staged exec-census spec's third story group
moves from `fable / high` to `opus / high`, like any other implementation.

## Acceptance

- [ ] `craft-line` routes the leverage override to mid + high and states that the top tier implements nothing unless the reviewer names it.
- [ ] `craft-line`'s decision table routes a genuinely uncertain seam to mid + high with the top tier offered to the reviewer.
- [ ] The profile's `Lines` section routes skill, command, and doc authoring to the mid model at high effort.
- [ ] The scorecard's current decisions record the rule with its date.
- [ ] `specs/exec-census/spec.md` group three declares `Line: opus / high.`
- [ ] The anchors registry and the prose budgets stay green.
