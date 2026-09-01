# Record the checkout layout mutation red

Blocked by: single-source-prospective-checkout-name.md
Writes: specs/roadmap-light-path-fixes/tickets/single-source-prospective-checkout-name.md
Covers: none

## What to build

Change `CheckoutName` away from `"checkout"` and demonstrate that the
fixed-layout compatibility test becomes red. Restore production exactly, then
record the mutation site, focused command, and observed failure in the owning
ticket. The production and test expectation remain unchanged.

## Acceptance

- [ ] The named constant mutation makes the focused compatibility test red.
- [ ] The production source is restored byte-for-byte before publication.
- [ ] The owning ticket durably records the exact mutation, command, and red result.
