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

**FT39 — concurrent-acquire contract: replace the 60s wall-clock deadline.**
Graduated 2026-07-07: a real `bench gate` run went red with exactly "second
acquire did not record within a minute — the runs never overlapped"
(runtime_worktree_test.go:355, 60.06s) on a tree identical to a green-gated
one, meeting the row's captured-red criterion. The fixed one-minute
spawn-to-record deadline is the weak point; fix is an event-keyed or raised
deadline. Spec: `specs/concurrent-acquire-deadline.md`. Next:
`/bench-implement-spec`.

**FT41 — /bench-shape-idea resume mode carries the grill through unblocked
tickets.** Rule-shaped from the learnings journal: soften "resolve that one
ticket, then stop" so a running grill continues into newly-unblocked tickets
in the same sitting while the reviewer is present and answering. Kit edit
under `craft-synthesis`. Spec: `specs/shape-idea-grill-continuation.md`.
Next: `/bench-implement-spec`.

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

1. `/bench-implement-spec` — FT39 concurrent-acquire deadline: graduation evidence
   captured, flake still live in the gate.
2. `/bench-implement-spec` — FT41 shape-idea grill continuation: rule-shaped kit
   edit under the synthesis discipline.
