# Prioritization assessment — 2026-07-05

Current open work, verified against the tree. This file replaces the 2026-07-04
assessment; everything that assessment carried and has since shipped — the no-grill
roadmap cleanup, the learnings-journal integration (review-findings persistence,
iteration-cap line declaration, the one-liner batch), the tests/ restructure
(test-suite-structure map + spec), and the binary auto-repair fallback — is removed
here, not restated. Rationale for removed items lives in git
(`git log --grep=spec-retire`, the 07-04 assessment's own history).

## Ready to build

**R1 (MED) — benign stale gate status.** Spec staged at
`specs/benign-stale-gate-status.md`: classify a stale gate verdict as capture-only
drift (`ROADMAP.md`, `.bench-notes.md` only) vs real drift, fail closed on any
untrusted comparison. The last open item from the 07-04 learnings integration.
Next action: `/bench-implement-spec`.

**R2 (MED) — gate phase-level concurrency.** Run the gate's
conformance/contract/canary phases in parallel; remeasure now that the test-layout
splits landed. Being shaped via `/bench-shape-idea` (this session).

## Features, in priority order (numbering continues the 07-04 list)

**FT2 (MED) — adversarial gate pinning.** Hash-verify the gate outside the
writable tree in pre-push. Distinct threat model from the lazy-agent tripwire;
small (~6 edits) and closes the "determined agent weakens the gate" hole.

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

*Live repo signals: one open learnings entry (shared-tree contention → worktree
discipline rule proposal) awaits `/bench-integrate-learnings`; unpushed commits
on main await a push.*

## Recommended sequence

1. R1 — implement the staged benign-stale-gate-status spec.
2. R2 — gate phase concurrency, once shaped.
3. FT2 (adversarial gate pinning), then FT3/FT4 by appetite.
