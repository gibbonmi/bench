# 5. Record the commitment scope in ADR 0028 and the delegate route in ADR 0023

Blocked by: 3-admit-light-path-landing.md
Writes: docs/adr/0028-the-commitment-gates-spec-implementations.md (new), docs/adr/0023-each-ticket-gets-a-fresh-author.md
Covers: LP47, LP48

## What to build

This ticket is the second ticket of review chunk LP-C2. It starts after the
LP-C1 checkpoint. It delivers the decision records that outlive this spec.

Add ADR 0028 in the shape of the current ADRs. It states that the commitment
gates spec implementations, and that a one-ticket light-path change needs none.
It states that a pinned row and the policy still change only through a plan. It
states that the drain dispatches each kept light-path fix to a fresh write
delegate. It holds no file path and no code snippet.

Change the last sentence of ADR 0023. It names the light-path delegate route
of the drain and of a phase that needs the fix.

## Acceptance

- [ ] ADR 0028 states that the commitment gates spec implementations and that a one-ticket light-path change needs none (LP47, review-owned).
- [ ] ADR 0023 names the drain's light-path delegate route (LP48, review-owned).
- [ ] `bench gate-prose . -- <each ADR>` passes for both ADRs.
