# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `dae4387`, clean tree, 3 unpushed commits
Spec: `specs/bench-front-door/spec.md` (Status: staged)
Gate: green at `e3f0e62` — stale, work tree `19c0b52`

## State

**Reviewer override standing for this repo:** the 2026-08 capability audit's own priority
order (`docs/audits/2026-08-bench-capability/results-fable-high/proposed-roadmap.md`)
supersedes `ROADMAP.md`'s `## Recommended sequence` until A1–A11 are exhausted; recorded
in `ROADMAP.md`. A1 and A2 are landed; A3 (`bench-front-door`) is staged, unstarted.

**`/bench-implement-spec specs/bench-front-door/spec.md --full` was invoked and stopped
before the first ticket landed.** Sequence: created integration worktree at `dae4387`
(request `bench-front-door-impl`) → `bench preflight build bench-front-door` green →
charged ticket `01-extract-route-owner.md` (Line: opus/mid-tier, high effort) to a
write-delegate three times in a row. All three failed identically on API 529
(Overloaded) with nothing written to the worktree — same failure mode already recorded
in this spec's own verification log for the write-spec phase. This is a service
availability issue, not a diff-owned red; there is nothing to fix in the ticket or spec.
On asking, the reviewer chose to stop and hand off rather than keep retrying, wait, or
try a one-off top-tier delegate. The now-unused integration worktree was released
(`bench worktree release`); `main` is untouched and clean at `dae4387`.

**No decisions were closed or reopened this session.** The spec's five reviewer-visible
calls (Further notes) are all still open, unchanged from staging.

## Next command

`/bench-final-check` — the board's leading invocable signal (`git`).

## Shape

Rewritten in full at every phase close, pruned rather than accreted: a fresh
session pays for every line it reads cold, so drop anything it would not act on.
Operational gotchas are placed by lifetime, not copied here: one that recurs across
phases belongs in `projects/benchkit.md`'s cold-session notes, and one scoped to a
build belongs in that spec's coverage rows. This file names at most when you'll hit
one, never the command — a second copy drifts from the source.
Keep the three sections above — **State** (what is true now, including anything
uncommitted), **Next command** (the exact harness-native invocation, not a
description of it), and this one.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
