# 2. Pin the Bench agent definitions with a conformance check

Blocked by: 1-ship-bench-agent-types.md
Writes: .bench/structure.budgets, internal/conformance/claude_agent_definitions_test.go (new), internal/conformance/checks_test.go, internal/conformance/registry/registry.go, internal/conformance/registry_test.go, projects/benchkit.md, tests/canary/claude-agent-definitions (new), tests/canary/guidance-prose-budgets/over-budget-skill, tests/canary/line-routing/line-binding-prose-drift, tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading, tests/canary/workflow-guidance-anchors/benchkit-review-round-owner, tests/canary/workflow-guidance-anchors/benchkit-review-round-routing, tests/canary/workflow-guidance-anchors/benchkit-spec-ownership, tests/canary/workflow-guidance-anchors/benchkit-system-suite-route
Covers: DB1, DB2, DB3, DB4, DB5, DB8, DB9

## What to build

Add the `claude-agent-definitions` check over every adapter agent file whose basename starts with `bench-`.
It reports a name that differs from the basename, an absent or empty tool list, and a declared model.
It also reports a spawning, publishing, or blocking tool, and a list without the read tool or the shell tool.
It reads the delegate skill in both directions, so an unnamed agent and a named but missing agent each report.
Register the check in the inventory, the test binding map, the fixture map, and the profile check tables.
Add a canary family with one fixture per diagnostic.

## Acceptance

- [ ] The live tree passes the check through `bench test --check claude-agent-definitions`.
- [ ] Each canary fixture turns the check red with its expected diagnostic text.
- [ ] The conformance meta checks pass with the new family and its registry rows.
- [ ] A consumer-style agent file with another prefix in the fixture tree produces no diagnostic.
