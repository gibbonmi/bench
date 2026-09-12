# 1. Ship the two Bench agent types

Blocked by: none
Writes: .claude/agents/bench-reviewer.md (new), .claude/agents/bench-writer.md (new), .claude/README.md, .bench/consumer-payload.json, package.json, internal/packagesurface/assets.go, internal/adopt/link_plan_test.go (new), internal/lines/lines_agentline_test.go, .agents/skills/bench-craft-delegate/SKILL.md, tests/canary/package-core-guard/kit-only-allowlist-emptied, tests/canary/package-core-guard/kit-only-asset-admitted, tests/canary/load-validity-metadata/invalid-json, tests/canary/workflow-guidance-anchors/delegate-cap-change-pinning-package, tests/canary/workflow-guidance-anchors/delegate-charge-effort-cap, tests/canary/workflow-guidance-anchors/delegate-coverage-row-charge, tests/canary/workflow-guidance-anchors/delegate-coverage-row-red-green, tests/canary/workflow-guidance-anchors/delegate-cross-harness-reviewer-pointer, tests/canary/workflow-guidance-anchors/delegate-exec-only-every-caller, tests/canary/workflow-guidance-anchors/delegate-model-id-escalation, tests/canary/workflow-guidance-anchors/delegate-own-family-native-surface, tests/canary/workflow-guidance-anchors/delegate-parallel-route-anchor, tests/canary/workflow-guidance-anchors/delegate-release-at-acceptance, tests/canary/workflow-guidance-anchors/delegate-resume-handoff-contents, tests/canary/workflow-guidance-anchors/delegate-self-probe-missing-row, tests/canary/workflow-guidance-anchors/delegate-stash-refusal-anchor, tests/canary/workflow-guidance-anchors/fix-pass-sentinel-anchor, tests/canary/workflow-guidance-anchors/prepared-review-delegate-handoff-route, tests/canary/workflow-guidance-anchors/prepared-triage-bounds, tests/canary/workflow-guidance-anchors/shared-worktree-path-pin, tests/canary/claude-agent-definitions/agent-unnamed-in-skill, tests/canary/claude-agent-definitions/model-declared, tests/canary/claude-agent-definitions/name-mismatch, tests/canary/claude-agent-definitions/shell-tool-absent, tests/canary/claude-agent-definitions/skill-names-missing-agent, tests/canary/claude-agent-definitions/spawning-tool, tests/canary/claude-agent-definitions/tools-absent
Covers: DB6, DB7, DB10, DB11, DB12, DB13, DB27, DB28

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
