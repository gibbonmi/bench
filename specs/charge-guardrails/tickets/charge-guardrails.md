# Guard the charge returns, the fixture reuse, and the red route

Blocked by: none
Writes: .agents/skills/bench-craft-delegate/SKILL.md, .agents/skills/bench-craft-delegate/references/delegation-discipline.md, internal/anchors/registry_data.go, internal/anchors/registry_charge_binding.go, internal/anchors/registry_charge_binding_test.go, tests/canary/workflow-guidance-anchors/delegate-coverage-row-charge/MUTATE.json, tests/canary/workflow-guidance-anchors/delegate-coverage-row-red-green/MUTATE.json, tests/canary/workflow-guidance-anchors/delegate-coverage-row-red-green/EXPECT
Covers: none

## What to build

One write charge went to two author models, and the charge guidance had three gaps. First,
the charge asked a blocked author for a commit sha and a preflight result, and one author
made a raw commit. Second, the charge quoted only part of the
one-source rule and named no fixture helper, so one author copied two fixture harnesses.
Third, the red-first rule required a red before the edit, but the build accepted a red from
a probe.

The delegation discipline states that a charge asks for a commit, a sha, or a
preflight result only on a committed ticket. It also states that a write charge names the
fence's fixture helpers and quotes any fixture-harness rule of the project. The charge section of the skill
states that each coverage row is observed red at least once and then green. The log names
the route of each red, and a TDD seam keeps the red sequence of `craft-tdd`. One anchor row
pins each new sentence, and one independent test expectation pins each new anchor row.

## Acceptance

- [ ] `bench test --check docs-currency-workflow` fails on the base guidance with the new anchor rows, and passes after the guidance edit.
- [ ] A probe that omits each new sentence turns `docs-currency-workflow` red, and the restore returns it to green.
- [ ] No other file states the conditional-return rule, the fixture-harness quote rule, or the red-route rule.
