# Fence the family row for a per-member rule

Blocked by: none
Writes: .agents/skills/bench-craft-tickets/references/slicing-checks.md, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go
Covers: none

## What to build

An ft336 ticket added a rule for some worktree leaves. The rule states which leaves retire an
assignment. The leaf family table holds the one declaration of each leaf, but that table was
not in the ticket `Writes:` line. The author put the rule in the owner package instead. The
dispatcher test was also not in the fence, so a probe on the dispatcher handoff stayed silent.
The orchestrator widened the fence before review.

The slicing checks own the rules for a ticket `Writes:` line. Add one slicing rule in general
words. A rule for one member of an enumerated family goes on that member's row in the family
declaration table. The ticket `Writes:` line names that table and the test of the dispatcher
that reads it. The ticket-slicing anchor rows pin the sentence, so a removal turns the anchor
check red.

## Acceptance

- [ ] The slicing checks tell the slicer to put a per-member rule on the family declaration row and to fence that table and its dispatch test.
- [ ] No other guidance file states where a per-member rule of an enumerated family goes.
- [ ] A probe that removes the sentence turns the anchor check red.
