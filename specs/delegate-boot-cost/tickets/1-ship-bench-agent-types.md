# 1. Ship the two Bench agent types

Blocked by: none
Writes: .claude/agents/bench-reviewer.md (new), .claude/agents/bench-writer.md (new), .claude/README.md, .bench/consumer-payload.json, package.json, internal/packagesurface/assets.go, internal/adopt/adopt_test.go, internal/lines/lines_agentline_test.go, .agents/skills/bench-craft-delegate/SKILL.md
Covers: DB6, DB7, DB10, DB11, DB12, DB13

## What to build

Add the reviewer and writer agent files to the Claude adapter, each with a name, a description, and a tool list, and no model.
Ship the directory through the consumer payload, the package files list, and the required-assets census.
State the routing rule in the delegate skill.
Reviews and diagnostics run as the reviewer type, and user-directed writes run as the writer type.
A fork runs only for work that needs the parent context.

Describe the agents directory in the Claude README.
Measure one "ok" boot per agent type from the usage line and record both counts for the retro.

## Acceptance

- [ ] The link plan for a fixture kit holds the reviewer agent file at its adapter path.
- [ ] The shipped-surface check passes with the agents tree row in the payload and the package files list.
- [ ] The required-assets census names both agent files.
- [ ] An envelope with the reviewer subagent type and a bound tier token returns a silent allow.
- [ ] The same envelope with no model returns the missing-model deny.
- [ ] The delegate skill names both agent basenames and states the fork rule within its line budget.
- [ ] A reviewer-type "ok" delegate reports under 15,000 tokens, and a writer-type one reports under 16,000.
