---
description: Fixture preserving the map gate while omitting its reviewer-closed-forks fast path.
---

# /bench-write-spec fixture

This phase refuses to run without a complete decision map. Default spec authoring
starts in a fresh mid-tier session. If any fork remains open, run
`/bench-shape-idea` and keep the normal map gate. This fixture intentionally omits
the same-session fast path for forks the reviewer already closed in the current
session.
