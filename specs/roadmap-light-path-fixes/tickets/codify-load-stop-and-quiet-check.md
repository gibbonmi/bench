# Codify the load stop and quiet check

Blocked by: support-installed-lane-repair-commit.md
Writes: .agents/skills/bench-craft-delegate/SKILL.md, .agents/skills/bench-craft-line/SKILL.md, .agents/skills/bench-craft-delegate/references/delegation-discipline.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go
Covers: LF4

## What to build

Require a stop after two known-flaky refusals proven green in isolation.
Before aggregate grading, require returned delegates to have no live tests and
serialize the coordinator-owned resource.

## Acceptance

- [ ] The second proven flaky refusal stops and hands evidence to the reviewer.
- [ ] Aggregate grading waits until returned delegates own no live tests.
- [ ] Canonical anchors protect both rules without copying their full prose.

