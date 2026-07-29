# Standardize the final retro and drain

Blocked by: Surface pending implementation retros

## What to build

Make `/bench-final-check` rewrite `.bench/retros/<spec-slug>.md` after the
spec's landing gate and commit are green, using the required outcome, timings,
ticket/delegate comparison, coordinator-catch, and agent-experience sections.
Make `/bench-what-next` disposition every pending retro from roadmap context
and remove all drained retro files in its approved batch. Pin both owners with
paired conformance anchors and registered canary fixtures.

## Acceptance

- [x] Stories 21–22 and their prose acceptance-coverage rows are green.
- [x] Final-check owns after-green creation; what-next owns reviewed disposition and removal.
- [x] Every new anchor family has a classified fixture whose EXPECT matches its targeted diagnostic.
