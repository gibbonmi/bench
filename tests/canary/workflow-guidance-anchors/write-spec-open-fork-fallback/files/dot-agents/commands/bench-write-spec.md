---
description: Fixture preserving the reviewer-closed fast path while omitting the open-fork fallback.
---

# /bench-write-spec fixture

This phase refuses to run without a complete map. Default spec authoring starts
in a fresh mid-tier session. The sole same-session exception applies when every
load-bearing fork has already been put to the reviewer and closed in the current
session: write those decisions directly into a new decision map with a complete
Handoff, continue from that file rather than unwritten grill memory, and compile
the spec without routing through `/bench-shape-idea`.
