# Roadmap delivery commitment

Status: ready

## Destination

Decide how the shared Bench workflow protects a finite milestone and its ordered work list.
Keep useful findings without silently displacing committed work.

## Notes

Use craft-grill and craft-domain for reviewer decisions. Use craft-research for facts that the retained diagnosis does not settle.
The diagnosis worktree also owns this shaping pass.

This source is ready for specification. The reviewer approves commitment enforcement and closure as a bounded prerequisite before the quality milestone.
The reviewer confirms decisions #1 through #17 on 2026-10-04.
Milestone completion, intake, blocked work, concurrency, enforcement, closure, adoption, and delivery scope are decided.

## Decisions so far

- [Where does the commitment rule apply?](roadmap-delivery-commitment/tickets/1.md): the shared Bench workflow.
- [What work does the rule protect?](roadmap-delivery-commitment/tickets/2.md): a finite milestone and its ordered work list.
- [Who can change the commitment?](roadmap-delivery-commitment/tickets/3.md): the reviewer explicitly approves each displacement.

- [What proves that a milestone is complete?](roadmap-delivery-commitment/tickets/4.md): verify the stated outcome, not only row closure.
- [What happens to newly discovered work?](roadmap-delivery-commitment/tickets/5.md): retain assessed findings outside the commitment until explicit admission.
- [What happens when committed work is blocked?](roadmap-delivery-commitment/tickets/6.md): continue with the next independent committed outcome and preserve the blocked obligation.
- [How many committed outcomes can be active?](roadmap-delivery-commitment/tickets/7.md): one by default; the reviewer authorizes parallel outcomes.

- [How many milestones can be active?](roadmap-delivery-commitment/tickets/8.md): one active milestone per project, with explicit activation and switch approval.
- [How does Bench enforce the commitment?](roadmap-delivery-commitment/tickets/9.md): Bench commands refuse unauthorized work starts and commitment changes.
- [When do completed roadmap references close?](roadmap-delivery-commitment/tickets/10.md): verified delivery closes completed roadmap and sequence references.
- [How do existing projects adopt the commitment?](roadmap-delivery-commitment/tickets/11.md): the reviewer approves the initial commitment; existing authorized runs may finish.
- [Do commitment protection and closure ship together?](roadmap-delivery-commitment/tickets/12.md): commitment protection and reliable closure ship together.

- [Which work classes take priority?](roadmap-delivery-commitment/tickets/13.md): defects and refactoring generally precede new features; modernization and deepening work are the current priority.

- [What bounds the first quality milestone?](roadmap-delivery-commitment/tickets/14.md): FT376, FT373, and FT349, in order; survey outcomes stay uncommitted intake; verify delivery state first.
- [How do defects rank against refactoring?](roadmap-delivery-commitment/tickets/15.md): confirmed defects generally precede refactoring, which precedes new features.
- [When does commitment enforcement implementation run?](roadmap-delivery-commitment/tickets/16.md): commitment enforcement and closure run first as one bounded prerequisite; the quality milestone follows.

- [Does the reviewer confirm the complete shared understanding?](roadmap-delivery-commitment/tickets/17.md): the complete plan and the revised prerequisite order are approved.

## Not yet specified

## Spec-writer discretion

- No additional discretion is approved.

## Out of scope

- New features unrelated to the approved quality outcomes are outside the first milestone.
- The enforcement prerequisite does not expand the approved quality milestone.
- Automatic displacement through a routine drain is excluded.
- Automatic adoption of an existing sequence as an approved commitment is excluded.

## Sources

- Path: `docs/research/roadmap-focus-and-completion.md`
  Supports: the growth, priority continuity, intake, and closure diagnosis behind this map.
  Drift: a change to the inspected roadmap, drain rules, final-check rules, or flow implementation.
- Path: `.agents/commands/bench-drain.md`
  Supports: the current intake, sequence, approval, and reconciliation rules.
  Drift: a change to those rules.
- Path: `.agents/commands/bench-final-check.md`
  Supports: the current capture duties and deferred roadmap closure.
  Drift: a change to capture or closure duties.
- Path: `roadmap/FT373.md`
  Supports: the existing standard-library modernization scope.
  Drift: a change to its scope or delivery state.
- Path: `ROADMAP.md`
  Supports: the FT376, FT373, and FT349 owner rows and the current sequence as unapproved input.
  Drift: one of the three owner rows is changed, merged, or retired, or the sequence changes.
- Path: `specs/markdown-block-reader/spec.md`
  Supports: the staged FT358 spec stays uncommitted intake outside the first milestone.
  Drift: its scope or status changes.
