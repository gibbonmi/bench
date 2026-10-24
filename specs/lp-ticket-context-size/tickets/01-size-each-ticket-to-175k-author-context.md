# Size each spec-backed ticket to about 175k tokens of author context

Blocked by: none
Writes: .agents/skills/bench-craft-tickets/SKILL.md, .agents/commands/bench-write-spec.md
Covers: none

## What to build

The slicer sizes each spec-backed ticket to about 175k tokens of author context. `craft-tickets` states the rule once, because the slicer and the ticket review both read that skill. The slicer estimates a ticket from the byte total of its `Writes:` files plus its verification cycles. The rule replaces the earlier rewrite-ticket sizing sentence, which sized a ticket by its expected reads. The reviewer approval list shows the estimate of each ticket. `/bench-write-spec` points to that list and does not repeat its columns.

## Acceptance

- [ ] `craft-tickets` tells the slicer to split a spec-backed ticket whose estimate passes about 175k tokens.
- [ ] `craft-tickets` tells the slicer to estimate from the `Writes:` byte total plus the verification cycles.
- [ ] The reviewer approval list in `craft-tickets` shows a context estimate for each ticket.
- [ ] `/bench-write-spec` names the `craft-tickets` approval table and keeps its slicing session bound of about 300k.
