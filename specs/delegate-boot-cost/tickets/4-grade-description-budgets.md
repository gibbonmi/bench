# 4. Grade the description budgets from the profile

Blocked by: 2-pin-bench-agent-definitions.md, 3-trim-skill-descriptions.md
Writes: .bench/structure.budgets, internal/conformance/tier_test.go, internal/conformance/skill_description_budgets_test.go (new), internal/conformance/checks_test.go, internal/conformance/registry/registry.go, internal/conformance/registry_test.go, internal/gate/lane_select_test.go, projects/benchkit.md, tests/canary/skill-description-budgets (new), .agents/skills/bench-craft-skills/SKILL.md, tests/canary/guidance-prose-budgets/over-budget-skill, tests/canary/line-routing/line-binding-prose-drift, tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading, tests/canary/workflow-guidance-anchors/benchkit-review-round-owner, tests/canary/workflow-guidance-anchors/benchkit-review-round-routing, tests/canary/workflow-guidance-anchors/benchkit-spec-ownership, tests/canary/workflow-guidance-anchors/benchkit-system-suite-route
Covers: DB14, DB15, DB16, DB17, DB18, DB19, DB20, DB21, DB23, DB24, DB25, DB26, DB29, DB30

## What to build

Add the `Skill description budgets` section to the project profile with two glob rows at 250 characters.
Add the `skill-description-budgets` check that parses that section and grades every skill and command description by collapsed rune count.
It reports a missing heading, a malformed row, a missing description, an over-budget description with its count and limit, and a symlinked or special subject.
Register the check with the benchkit-profile input, so the lane runs it on a profile edit.
Add one sentence to the craft-skills skill that names the profile table as the budget's owner.

## Acceptance

- [ ] The live tree passes the check through `bench test --check skill-description-budgets`.
- [ ] Each canary fixture turns the check red with its expected diagnostic text.
- [ ] A description that continues on a second line reports a one-line diagnostic.
- [ ] A multibyte description counts its collapsed runes.
- [ ] A lowered profile cell lowers the limit in the diagnostic.
- [ ] A symlinked skill directory and a special file each return a refusal diagnostic.
- [ ] The benchkit-profile lane expectation lists the new check.
- [ ] The invocation policy rows are unchanged and the existing adapter canary stays red on a flip.
