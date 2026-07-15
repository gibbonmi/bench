---
description: Fixture preserving the map gate while omitting its reviewer-closed-forks fast path.
---

# /bench-write-spec fixture

This phase refuses to run without a complete decision map. This fixture keeps
the default map-required posture but intentionally omits the fast path for forks
the reviewer already closed in the current session.
