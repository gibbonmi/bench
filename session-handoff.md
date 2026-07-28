# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — clean tree, ahead of origin
Spec: none active; `ft91-gate-phase-split` retired 2026-07-28 (stories 4/5
accepted as probed phases, story 9 dropped — map #4)
Gate: green, measured 2026-07-28 at 267 s wall (the instrumented run behind
map ticket #1)

## State

- **FT91's eighth arm is diagnosed and mostly ruled:
  `decisions/gate-critical-path.md`.** The gate is canary-bound (canary solo
  250 s of a 267 s gate): the 34 `behavior-owned` fixtures each re-run the
  full kit contract suite in their nested gates. Evidence asset:
  `decisions/assets/gate-critical-path-timeline.md`.
- **Reviewer rulings taken 2026-07-28, all recorded in the map:** stage 1 —
  behavior-owned fixtures scope to one contract package via subfamily
  directories (#6); stage 2 — canary nesting removed for that family, bites
  proven in-process at the owning contract test (#7; the gate-pipeline
  nesting clause transferred here and reopened); gate-pipeline.md stays
  closed and the phase-split spec retired (#4); **FT91's stop condition is a
  measured dev gate ≤60 s** — the oracle-semantics levers enter this map if
  the post-stage-2 re-measurement stays above it (#5).
- **Open tickets: #2 and #3 only.** #2 (Research, agent-alone) buckets the
  artifact suite build-vs-inspect; #3 (Grill, blocked by #2) rules which
  tests share one prepared artifact set. They gate only the artifact-hoist
  slice, not the two canary slices.
- **The roadmap's FT91 row still describes the phase-split spec as
  unretired** — stale as of this session; the next `/bench-what-next`
  reconcile owns the row rewrite. The tree wins meanwhile.
- **`decisions/cost-follows-project-size.md` stays open by design** (ticket
  #6 there, opportunistic cheap-tier retest).

## Next command

`/bench-write-spec` for stage 1 — the behavior-owned package-scoping slice —
on a fresh mid-tier session; map tickets #1/#6/#7 plus the timeline asset
carry the seams. The reviewer wants this moving now.

After it: `/bench-shape-idea` resume on `gate-critical-path` #2 (agent-alone
research), which unblocks #3 and the artifact-hoist slice; stage 2 specs the
moment stage 1 lands.

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
