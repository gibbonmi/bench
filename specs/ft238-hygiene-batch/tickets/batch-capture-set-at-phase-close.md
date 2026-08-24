# Batch the tracked capture set at phase close

Blocked by: none
Writes: AGENTS.md

## What to build

One line lands in the AGENTS.md phase-close paragraph. It makes the phase
close commit every tracked capture artifact of that close in one gate-priced
commit: the retro and the scorecard updates. The git-ignored files (the
handoff and the two inboxes) stay outside this rule, because they never
commit. The line removes the practice's dependence on session memory.

## Acceptance

- [ ] The phase-close paragraph carries the batch rule in one sentence set, in ASD-STE100.
- [ ] The rule names the tracked capture set and one commit; it does not name the git-ignored files as commit content.
