# Roadmap

The working prioritization document: every row is open work, verified against
the tree; a row leaves when the work ships (spec-retire) or a
`/bench-what-next` reconcile removes it. Raw capture never lands here — it
goes to `IDEAS.md` and enters only through a reviewed drain.

## Ready to build

**WN (HIGH) — finish the what-next build.** `specs/what-next.md` stories 7–11:
the `/bench-what-next` phase, retirement of the old learnings-integration
phase, `/bench-shape-idea` retarget, gate/canary anchors, spec-retire wiring.
Stories 1–6 are shipped.
Next action: `/bench-implement-spec`.

**S1 (LOW) — split `internal/gate/phases_test.go`.** 457 lines, over the
400-line structural cap (`bench structure` flags it now). 14 test functions,
cleanly separable by responsibility (runner concurrency/output/cancel vs
phase-table/shellcheck/signal-handling) — split along those lines per the
craft-seams skill, don't fragment just to beat the line count.
Next action: direct fix-and-gate (small, mechanical, no spec needed).

**FT2 (MED) — adversarial gate pinning.** Hash-verify the gate outside the
writable tree in pre-push. Distinct threat model from the lazy-agent tripwire;
small (~6 edits) and closes the "determined agent weakens the gate" hole.
Next action: `/bench-write-spec`.

## Features, in priority order

**FT3 (MED-LOW) — `bench spec implemented` + `bench commit`.** Pair them:
commit could fold in the spec status flip. Replaces footgun prose in the
implement phase with two small wrappers over existing logic.

**FT4 (MED-LOW) — harness task list in `/bench-implement-spec`.** Per-harness
adapter (Claude hook + phase line; Codex native).

**FT5 (LOW) — `bench outline`.** Marginal for this repo, real as a kit
affordance for large/polyglot linked repos. Needs its grill (languages,
on-demand vs committed, prose anchors).

**FT6 (LOW, parked pending evidence — leave parked):** `bench refs`, `bench
detect`, `bench doc`, `bench specs --retired`, doctor binary-presence row,
`conformanceFamilies`-vs-dispatch reconcile meta-check. `bench symbols` is not
carried; restore only if agents demonstrably burn turns on symbol search.

**FT7 (LOW) — dashboard.** Low priority by declaration.

**FT8 (scheduled, not actionable) — Sonnet 5 mid-tier revisit.** Time-boxed to
2026-09-01 or the next frontier shift.

## Watch

- `bench worktree concurrent-acquire` contract test failed once under
  full-gate load, then passed 3/3 in isolation and on rerun — likely a timing
  flake surfaced by gate phase concurrency. Journal it if it recurs.

## Recommended sequence

1. WN what-next stories 7–11 — `/bench-implement-spec`
2. S1 split `internal/gate/phases_test.go` — `/bench-implement-spec`
3. FT2 adversarial gate pinning — `/bench-write-spec`
