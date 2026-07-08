# Roadmap

The working prioritization document: every row is open work, verified against
the tree; a row leaves when the work ships (spec-retire) or a
`/bench-what-next` reconcile removes it. Raw capture never lands here — it
goes to `IDEAS.md` and enters only through a reviewed drain. A row for spec'd
work names its spec path (`specs/<slug>.md`) — that path is what `bench status`
cross-checks against the tree, so a row that omits it is a visible choice to
stay outside the ambient check.

## Features, in priority order

**FT38 — dashboard visual identity pass.** `bench dashboard` v1 shipped
data-faithful and visually neutral; the original idea wanted an
ui_examples-inspired rich treatment with animated characters. Taste is a
reviewer call, so the work starts as a grill, not a build. Decision detail is
recoverable via `bench spec history dashboard`. Next: `/bench-shape-idea`.

**FT39 (parked pending repro) — concurrent-acquire contract's 60s wall-clock
window.** A batch session reported `bench_worktree_concurrent-acquire_contract`
red ("second acquire did not record within a minute") ~3 times in ~8 gate runs
under load; five fresh runs on 2026-07-07 (three idle, two gates racing) all
stayed green, so the red is unreproduced. The fixed one-minute spawn-to-record
deadline is the suspected weak point; candidate fix is an event-keyed or raised
deadline. Graduate only on a captured red with that exact message from a real
`bench gate` run.

**FT6 (LOW, parked pending evidence — leave parked):** `bench refs`, `bench
detect`, `bench doc`, `bench specs --retired`, doctor binary-presence row,
`conformanceFamilies`-vs-dispatch reconcile meta-check, and a per-anchor
bite-proof meta-test (canaries prove one needle per family today; graduate on
observed anchor rot). `bench symbols` is not carried; restore only if agents
demonstrably burn turns on symbol search.

**FT24 (parked pending upstream) — Codex agent-line guard parity.** Researched
2026-07-07: not implementable on current Codex — delegation never surfaces as a
matchable `tool_name` on a deny-capable hook event, and `SubagentStart` neither
carries the delegate's resolved model nor honors a deny (verdict recorded in
`.bench/BENCH-reference.md` Hook Layers). Graduate only when the Codex
changelog adds a spawn tool name or a deny-capable SubagentStart.

**FT8 (scheduled, not actionable) — Sonnet 5 mid-tier revisit.** Time-boxed to
2026-09-01 or the next frontier shift.

## Recommended sequence

1. `/bench-shape-idea` — FT38 dashboard visual identity: the open decisions
   are pure reviewer taste, so the grill comes before any build.
