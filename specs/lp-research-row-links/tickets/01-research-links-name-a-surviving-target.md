# Research links name a surviving target

Blocked by: none
Writes: docs/research/bench-efficiency-assessment-2026-10-06.md, docs/research/roadmap-review-2026-09-25.md
Covers: none

## What to build

Two research snapshots link to roadmap row files that the tree no longer holds.
A drain or a landing deleted each row, so each link is dead.
On 2026-10-27 the reviewer decided to repair every dead roadmap-row link in `docs/research`, not only the FT314 links.

Each dead link gets the surviving owner of its fact, in this order:

1. The ADR that records the shipped decision of the row.
2. The live roadmap row that absorbed the row, with merge evidence in the deleting landing.
3. No link: the link text stays as plain text, and the edit adds nothing else.

The research prose stays as it is, because each file is a dated evidence snapshot.
The edit changes only the link target, or it removes the link and keeps the link text.

The ADR and absorbing-row targets are:

| Dead row | New target | Evidence |
| --- | --- | --- |
| FT253 | ADR 0014 | Landing 35d5a8fa5 deleted the row and added the landing-window consequence to ADR 0014. |
| FT314 | ADR 0031 | Landing 8e1a88283 deleted the row and added ADR 0031. |
| FT336 | ADR 0025 | Close 38b0c0d84 deleted the row and added ADR 0025. |
| FT342 | ADR 0014 | Landing 35d5a8fa5 deleted the row and added the landing-window consequence to ADR 0014. |
| FT249 | FT265 | Landing 15ecad2d8 moved the inbox-reader occurrence into FT265. |
| FT292 | FT172 | Landing 15ecad2d8 moved the `bench roadmap decide` face and its occurrence into FT172. |
| FT296 | FT305 | Landing 15ecad2d8 moved the supervisor-lifetime face and its occurrence into FT305. |
| FT317 | FT318 | Drain fold c9d5c7123 folded the capability-blocked decision into FT318. |
| FT325 | FT326 | Landing 15ecad2d8 moved the wrap face and three occurrences into FT326. |
| FT339 | FT341 | Landing 15ecad2d8 moved the registry-attribution occurrence into FT341. |
| FT340 | FT89 | Landing 15ecad2d8 moved the drain-wording face and its occurrence into FT89. |
| FT388 | FT387 | Drain e77f4d3e5 folded FT388 into FT387, and FT387 cites the fold. |

These rows have no ADR and no live absorbing row, so their links become plain text:

- FT38, FT58, FT99, FT108, FT115, FT199, FT217, FT246, FT262, FT267, FT290, FT293, FT299, FT324, FT329, FT330, FT335, FT337, FT338, FT348, FT349, FT373, FT378, FT379, FT383.

The `FT339's directory query` link also becomes plain text.
Landing 15ecad2d8 moved that face into FT125, and drain 72e751b7a removed it from FT125.
FT330 folded into FT393, and FT393 shipped at 088cd6132 with no ADR.
FT378 shipped as `docs/known-issues.md`, which is not an ADR or a roadmap row.

## Acceptance

- [ ] A check finds no link in `docs/research` whose target `roadmap/FT<n>.md` does not exist.
- [ ] Each repointed link names a file that exists in the tree.
- [ ] Each changed line changes only a link target, or it removes a link and keeps its link text.
- [ ] The `bench commit` lane passes.
