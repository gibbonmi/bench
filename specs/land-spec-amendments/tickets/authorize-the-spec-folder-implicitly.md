# Authorize the spec folder implicitly

Blocked by: none
Writes: internal/preflight

## What to build

`paths-authorized` authorizes any changed path under the active spec's own
folder, `specs/<slug>/`, without a fence entry. The check derives the folder
from the spec path fact the gatherer already supplies — no second fact carries
it — and consults it beside the declared fence entries with the same
segment-boundary rule. The change lands once and serves all three
consumers — build preflight, review preflight, and the landing's
`AuthorizeReviewedSource` — because they share the gatherer-and-check pair.
Declared fence entries keep their exact semantics.

Covers: LS7, LS8, LS9, LS10, LS11.

## Acceptance

- [ ] Build preflight is green on a changed path under the active spec's folder with no self-fence entry (LS7).
- [ ] Review preflight is green on that same path the same way (LS8).
- [ ] A changed path under a different spec's folder stays red (LS9).
- [ ] A changed path under a sibling folder whose name extends the slug stays red (LS10).
- [ ] A path authorized only by a declared fence entry stays green, and one outside every fence and the spec folder stays red (LS11).
- [ ] The decision-domain tests stay repository-free: the implicit entry derives from the existing spec-path fact inside the `Decide` table harness.
