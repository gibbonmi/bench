# Run build preflight on the final ticket graph

Blocked by: none
Writes: .agents/skills/bench-craft-tickets/references/slicing-checks.md, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go
Covers: none

## What to build

The slicing checks tell the slicer to run build preflight after each fence change and before
review. A change to a ticket or an acceptance row can leave the fence the same. A change after
review can also occur. Both changes can reach the build with a red preflight. The ft336 build
started on such a graph: two ticket paths had no `(new)` marker, and no ticket cited one row.

Widen the existing slicing rule in its owner file. The slicer runs `bench preflight build <slug>`
after each change to a ticket or an acceptance row. The spec goes to review, to sign-off, and
to its landing only on a green preflight of its current ticket graph. The ticket-slicing anchor
rows pin both sentences, so a removal turns the anchor check red.

## Acceptance

- [ ] The slicing checks state the rerun after each change to a ticket or an acceptance row.
- [ ] The slicing checks state that review, sign-off, and the landing wait for a green preflight.
- [ ] No other guidance file states when build preflight runs over the ticket graph.
- [ ] A probe that omits the green-preflight sentence turns the anchor check red.
