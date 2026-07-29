# Co-locate compiled decision maps

Blocked by: Contract live specs to folder only

## What to build

Keep pre-spec working maps at top-level `decisions/` through shaping. Make
`/bench-write-spec` move a closed source map and any map-owned assets into
`specs/<slug>/decisions/`, updating references in the same green change.
Treat the moved files as settled provenance outside `bench maps`, and rely on
whole-folder spec retirement to remove them with the tickets.

## Acceptance

- [x] Story 17 and both of its acceptance coverage rows are green or retain their recorded classification.
- [x] Decision #9 and the map Handoff record the approved lifecycle.
- [x] Write-spec, shape-idea, and final-check agree on the open, compiled, and retired lifecycle.
- [x] The FT154 source map and every in-scope affected reference moved together.
- [x] README, field guide, and changelog teach the user-visible lifecycle.
- [x] The project gate is green.
