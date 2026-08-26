# Run the frontier in parallel by default

Blocked by: none
Writes: .agents/skills/bench-craft-tickets/SKILL.md, .agents/commands/bench-implement-spec.md

## What to build

`craft-tickets` owns the frontier rule. Today it permits parallel delegates
only where two tickets' `Writes:` notes are disjoint, so a coordinator
serializes by default. The reviewer wants the opposite default: every
unblocked frontier ticket runs in its own worktree at once, and the
coordinator serializes a pair only for a conflict it judges and names in the
charge. The named conflicts are an overlapping `Writes:` note, a contract one
ticket reads from the other's diff, and the machine's test budget. The
coordinator then folds each returned diff onto the one retained integration
source in `Blocked by:` order, one gate per commit.

`/bench-implement-spec` gains one pointer sentence in its Build section, so
the command and the skill state one rule. Both files sit at their prose
budgets, so the edit holds each file's line count.

## Acceptance

- [ ] `craft-tickets` states that frontier tickets run in parallel by default, one worktree each.
- [ ] `craft-tickets` names the three conflicts that serialize a pair and requires the charge to name the one that applies.
- [ ] `craft-tickets` states that the coordinator folds returned diffs in `Blocked by:` order, one gate per commit.
- [ ] `/bench-implement-spec` points at `craft-tickets` for the parallel default.
- [ ] The guidance prose budgets and the anchors registry stay green.
