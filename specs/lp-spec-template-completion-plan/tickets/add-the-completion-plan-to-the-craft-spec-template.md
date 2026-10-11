# Add the completion plan to the craft-spec template

Blocked by: none
Writes: .agents/skills/bench-craft-spec/SKILL.md, .agents/skills/bench-craft-spec/references/completion-plan.md
Covers: none

## What to build

Every build preflight mode refuses a spec that has no valid `bench-completion-plan` fence.
The spec template in `craft-spec` has no completion plan section, so a spec author writes no fence.
The template gets one completion plan line that points to `references/completion-plan.md`.
That reference gives the section heading, the fenced JSON skeleton, and one sentence for each key.
The keys and their meanings come from the plan parser in `internal/reviewrecord`.
The `craft-spec` skill stays at its prose budget of 155 lines.

## Acceptance

- [ ] `rg -n "bench-completion-plan" .agents/skills/bench-craft-spec/` finds the template line and the reference file.
- [ ] The reference skeleton holds each key that the plan parser requires for a version 1 plan.
- [ ] The `guidance-prose-budgets` check passes with `.agents/skills/bench-craft-spec/SKILL.md` at 155 lines.
