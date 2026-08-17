# Learnings — usage journal

## 2026-08-17 — Refresh a ready deepening map before handoff [open]

A new `/bench-deepen` survey reverified a ready decision map whose active body
still centered already-landed candidates, while its remaining candidates were
spec-ready or light-path-ticket-ready and duplicate roadmap prose still
advertised open or conflicting decisions.

When a deepening report proves that an existing map's frontier is empty,
automatically rewrite that map to the current remaining work, replace duplicate
roadmap decision claims with pointers to the map, run `bench maps`, and rewrite
the session handoff. Do not ask the reviewer to re-grill settled predicates or
create a second map.

Proposed rule: add this current-state refresh to `/bench-deepen` after source
verification and before the exit handoff.
