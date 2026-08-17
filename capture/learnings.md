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

## 2026-08-17 — Inventory shared decision readers before ticket slicing [open]

The FT216 write-spec pass named the two obvious planners but missed the shared
`assignmentLanded` reader that preserved the formatted landedness protocol across
automatic cleanup, landed-set planning, resume, and list. The first ticket slice
also bundled automatic planning with the independently shippable landed-set
consumer, so the required review round needed a second acceptance iteration.

Before ownership fences or tickets are written, automatically census every
production reader of the decision fact being deepened, including helpers called
by named consumers. Put each shared reader under an exact fence, then ask whether
each consumer branch can land green without the others and split the DAG wherever
the answer is yes. Proposed rule: add the reader census and per-branch ship test to
`/bench-write-spec` immediately before the `craft-tickets` slice.
