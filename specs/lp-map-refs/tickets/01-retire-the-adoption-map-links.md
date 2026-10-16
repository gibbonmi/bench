# Retire the adoption map links in the decision-map index

Blocked by: none
Writes: capture/architecture-decision-maps-2026-10-06.md
Covers: none

## What to build

The decision-map index in `capture/architecture-decision-maps-2026-10-06.md` links the `adoption-lifecycle-consumption` spec and its C06 decision map.
The FT217 build landed that spec, and a later spec retirement removed the spec folder.
Thus both links point to paths that do not exist.

Replace the two links with plain text that names the retired spec.
Set the C06 row and the adoption row to the current state: the four tickets landed, and the spec is retired.
Correct the evidence row that says that the set holds one map for each candidate, because the C06 map is retired.

Do not change the FT343 rows.
In this index, `integrated` means that the plan is integrated into the planning branch.
Every staged plan in the index carries the same word, so the FT343 rows agree with the staged spec.
Do not change `specs/landing-test-efficiency/spec.md`.
Its `reviews/landing-test-efficiency.md` entry is the review pickup that its fence sentence names, and the build creates that file.

## Acceptance

- [ ] The index holds no link to a path under `specs/adoption-lifecycle-consumption/`.
- [ ] The C06 row and the adoption row state that the spec is retired.
- [ ] The evidence list states that the C06 map is retired with its spec.
- [ ] Each other row of the index is unchanged.
