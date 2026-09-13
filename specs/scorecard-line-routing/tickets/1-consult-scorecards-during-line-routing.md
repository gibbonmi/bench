# Consult scorecards during line routing

Blocked by: none
Writes: .agents/skills/bench-craft-line/SKILL.md, CHANGELOG.md
Covers: none

## What to build

When a repository has provider scorecards, `craft-line` uses them as advisory evidence before it selects a line.
The scorecards can suggest or validate a choice, but they do not indicate the task outcome.

## Acceptance

- [ ] `craft-line` reads an available provider scorecard before it selects a line.
- [ ] `craft-line` treats scorecards as suggestions or validation evidence.
- [ ] `craft-line` does not treat scorecards as task outcome indicators.
- [ ] User direction, the project `Lines`, and the current task signals remain authoritative.
- [ ] The changelog records the routing guidance change.
