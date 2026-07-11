# Roadmap

The working prioritization document: every row is open work, verified against
the tree; a row leaves when the work ships (spec-retire) or a
`/bench-what-next` reconcile removes it. Raw capture never lands here — it
goes to `IDEAS.md` and enters only through a reviewed drain. A row for spec'd
work names its spec path (`specs/<slug>.md`) — that path is what `bench status`
cross-checks against the tree, so a row that omits it is a visible choice to
stay outside the ambient check.

## Features, in priority order

**FT66 — bound the terminal repair pass in review guidance.** Recursive
review/fix rounds made the workflow unbounded. Add one clause to
`/bench-review-implementation`: integrate accepted findings, run focused
checks, run one final gate, then stop; open another review round only when that
gate fails or the reviewer requests one. Rule-shaped kit edit under
`craft-synthesis`. Next: `/bench-implement-spec` (lighter path).

**FT62 — structure debt: `internal/contract/helper.go` at 441/400.** Grown by
FT53's RunAtWithTimeout helper; one of two open `bench structure` issues. Split
along responsibility or propose a reviewer grant — check the dir's file-count
headroom before choosing (`internal/contract/runtime/` already holds a dir
grant, so a split into that dir deepens dir debt). Next:
`/bench-implement-spec` (lighter path).

**FT64 — structure debt: `internal/canary/canary_test.go` at 407/400.** The
targeted canary-phase work grew this past the file budget after FT62 was
recorded. Split along responsibility or propose a reviewer grant, checking
both file and directory budgets before choosing. Next: `/bench-implement-spec`
(lighter path).

**FT6 (LOW, parked pending evidence — leave parked):** `bench refs`, `bench
detect`, `bench doc`, `bench specs --retired`, doctor binary-presence row,
`conformanceFamilies`-vs-dispatch reconcile meta-check, and a per-anchor
bite-proof meta-test (canaries prove one needle per family today; graduate on
observed anchor rot). `bench symbols` is not carried; restore only if agents
demonstrably burn turns on symbol search.

**FT58 (parked pending evidence) — reclaim-lock protocol for lease Claim.**
Serialize dead-pid takeovers via a crash-safe lock file to close the residual
B-vs-C double-use window the FT45 no-clobber restore leaves open. Graduate on
an observed double-win (the two-reclaimer stress case red, or a real
double-use in the pool); until then the shipped identity-verify plus restore
is the accepted posture.

**FT24 (parked pending upstream) — Codex agent-line guard parity.** Researched
2026-07-07: not implementable on current Codex — delegation never surfaces as a
matchable `tool_name` on a deny-capable hook event, and `SubagentStart` neither
carries the delegate's resolved model nor honors a deny (verdict recorded in
`.bench/BENCH-reference.md` Hook Layers). Graduate only when the Codex
changelog adds a spawn tool name or a deny-capable SubagentStart.

**FT8 (scheduled, not actionable) — Sonnet 5 mid-tier revisit.** Time-boxed to
2026-09-01 or the next frontier shift.

**FT38 (tabled, revisit on or after 2026-08-09) — dashboard visual identity
pass.** `bench dashboard` v1 shipped data-faithful and visually neutral; the
original idea wanted a rich treatment with animated characters, reference
saved at `ui_example/` (Gather-style pixel office with activity feed).
Reviewer tabled it 2026-07-09 for at least a month. When it revives, the work
starts as a grill (`/bench-shape-idea`); decision detail recoverable via
`bench spec history dashboard`.

## Recommended sequence

1. `/bench-implement-spec` — FT66 add the terminal repair-pass bound under
   `craft-synthesis` (lighter path).
2. `/bench-implement-spec` — FT62 helper.go split-or-grant, checking both
   structure budgets before choosing.
3. `/bench-implement-spec` — FT64 canary_test.go split-or-grant, checking both
   structure budgets before choosing.
