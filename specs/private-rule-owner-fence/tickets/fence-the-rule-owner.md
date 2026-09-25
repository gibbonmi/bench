# Fence the owner of a reused private rule

Blocked by: none
Writes: .agents/skills/bench-craft-tickets/references/slicing-checks.md, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go
Covers: none

## What to build

An ft336 ticket added a second caller of a private rule. The rule checks an assignment ID,
and the census package owns it. The ticket `Writes:` line did not name the census package or
the pool key package. The author copied the rule into a new file, and review found the
duplicate. The orchestrator widened the fence, so one exported predicate serves both callers.

The slicing checks own the rules for a ticket `Writes:` line. Add one slicing rule in general
words. A ticket that adds a second caller of a private rule in another package names that
package in `Writes:`. The ticket-slicing anchor rows pin the sentence, so a removal turns the
anchor check red.

## Acceptance

- [ ] The slicing checks tell the slicer to fence the owner package of a private rule that a ticket calls a second time.
- [ ] No other guidance file states that a second caller of a private rule fences the owner of that rule.
- [ ] A probe that removes the sentence turns the anchor check red.
