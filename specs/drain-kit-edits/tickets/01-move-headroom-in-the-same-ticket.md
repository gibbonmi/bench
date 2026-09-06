# State that a ticket moves its headroom in the same ticket

Blocked by: none
Writes: .agents/skills/bench-craft-tickets/SKILL.md
Covers: none

## What to build

The lane grades growth per commit against the current tip, not against the
spec base. A ticket that adds a line to a file over its structure budget
therefore fails the lane unless the same ticket moves the headroom. The
bench-probe build paid two repair rounds for this rule. `craft-tickets` states
the rule once, under `## Draft the breakdown`, in one or two ASD-STE100
sentences. The rule names the same-ticket placement and the per-commit lane
base as its one-clause why.

## Acceptance

- [ ] `.agents/skills/bench-craft-tickets/SKILL.md` states that a ticket which adds a line to an over-budget file moves that file's headroom in the same ticket.
- [ ] The rule names the lane's per-commit growth base as its reason.
- [ ] `bench gate-prose . -- .agents/skills/bench-craft-tickets/SKILL.md` exits 0.
- [ ] `bench test --check guidance-prose-budgets` and `bench test --check kit-compliance` exit 0 over the worktree.
