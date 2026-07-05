# Prioritization assessment — 2026-07-05 (rev 2)

Current open work, verified against the tree. This file replaces the same-day
2026-07-05 assessment; R1 (benign stale gate status) and R2 (gate phase-level
concurrency) both shipped after that assessment was written — implemented,
spec-retired, and confirmed green — so they're removed here, not restated.
Rationale lives in git (`git log --grep=spec-retire`, decisions/gate-phase-concurrency.md).

## Ready to build

**S1 (LOW) — split `internal/gate/phases_test.go`.** 457 lines, over the
400-line structural cap (`bench structure` flags it now). Grew past the limit
during the just-shipped gate-phase-concurrency work. 14 test functions, cleanly
separable by responsibility (runner concurrency/output/cancel vs phase-table/
shellcheck/signal-handling) — split along those lines per the craft-seams
skill, don't fragment just to beat the line count.
Next action: direct fix-and-gate (small, mechanical, no spec needed).

**FT2 (MED) — adversarial gate pinning.** Hash-verify the gate outside the
writable tree in pre-push. Distinct threat model from the lazy-agent tripwire;
small (~6 edits) and closes the "determined agent weakens the gate" hole. Now
first in line since R1/R2 are done.
Next action: `/bench-write-spec`.

## Features, in priority order

**FT3 (MED-LOW) — `bench spec implemented` + `bench commit`.** Pair them: the
roadmap already notes commit could fold in the spec status flip. Replaces footgun
prose in the implement phase with two small wrappers over existing logic.

**FT4 (MED-LOW) — harness task list in `/bench-implement-spec`.** Per-harness
adapter (Claude hook + phase line; Codex native).

**FT5 (LOW) — `bench outline`.** Marginal for this repo, real as a kit
affordance for large/polyglot linked repos. Needs its grill (languages, on-demand
vs committed, prose anchors).

**FT6 (LOW, parked pending evidence — leave parked):** `bench refs`, `bench
detect`, `bench doc`, `bench specs --retired`, doctor binary-presence row,
`conformanceFamilies`-vs-dispatch reconcile meta-check. `bench symbols` is not
carried; restore only if agents demonstrably burn turns on symbol search.

**FT7 (LOW) — dashboard.** Low priority by declaration; unchanged.

**FT8 (scheduled, not actionable) — Sonnet 5 mid-tier revisit.** Time-boxed to
2026-09-01 or the next frontier shift; keep as is.

## Live repo signals

- Gate: confirmed green on a direct run this session. One contract test
  (`bench worktree concurrent-acquire`) failed once under full-gate load, then
  passed 3/3 in isolation and on a full-gate rerun — looks like a timing flake
  surfaced by phase concurrency (R2), not a real regression. Worth a learnings
  entry if it recurs; not actionable yet on a single occurrence.
- One open learnings entry (shared-tree contention → worktree discipline rule
  proposal) awaits `/bench-integrate-learnings`.
- Local `main` is 9 commits ahead of `origin/main`, unpushed.
- The persistent worktree under `~/.bench/worktrees/` is the warm pool by
  design (`internal/worktree/worktree.go`) — not a leak, no action needed.

## Recommended sequence

1. S1 — split `phases_test.go` (small, mechanical, clears the structure flag).
2. FT2 — adversarial gate pinning, shape then build.
3. FT3/FT4 by appetite.
