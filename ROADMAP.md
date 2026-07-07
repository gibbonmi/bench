# Roadmap

The working prioritization document: every row is open work, verified against
the tree; a row leaves when the work ships (spec-retire) or a
`/bench-what-next` reconcile removes it. Raw capture never lands here — it
goes to `IDEAS.md` and enters only through a reviewed drain. A row for spec'd
work names its spec path (`specs/<slug>.md`) — that path is what `bench status`
cross-checks against the tree, so a row that omits it is a visible choice to
stay outside the ambient check.

## Features, in priority order

**FT7 (implemented, awaiting post-hoc review) — `bench dashboard`**
(`specs/dashboard.md`). v1 is the data-faithful minimal HTML snapshot; the
rich visual treatment (ui_examples look, animated characters) is deferred as a
reviewer-taste follow-up in the spec's out-of-scope list. Remaining work:
`/bench-review-implementation`, then spec-retire on approval.

**FT22 (implemented, awaiting post-hoc review) — `bench spec history <slug>`**
(`specs/spec-history.md`). Retirement archaeology folded into the CLI (FT9
pattern). Remaining work: `/bench-review-implementation`, then spec-retire on
approval.

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

1. `/bench-review-implementation` — post-hoc three-axis review of the three
   implemented specs (`outline`, `dashboard`, `spec-history`).
2. `bench spec retire <slug>` per spec the review clears (reviewer call).
3. `/bench-what-next` — drain the three open journal entries from the batch.
