# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — clean, 22 unpushed commits ahead of `origin/main`
Spec: `specs/ft148-worktree-orphan-retirement.md` (Status: staged, approved
2026-07-27) — the next build
Gate: green at `08482c6`; every commit since is doc-only

## State

- **FT148's spec is staged and signed off, and the build has not started.**
  Twelve stories: a `created_at` field on the assignment record, an age-only
  orphan predicate, a labelling change in `PlanAutomatic`, orphan lines and a
  cap in the resume summary, ledger compaction for tree-gone rows, three kit-prose
  edits with conformance anchors, and one dead comment pointer. The map behind it
  is `decisions/worktree-orphan-retirement.md`; read its **Provenance** section
  before trusting any decision in it, because it was written in the same session
  as the spec.
- **Orphanhood is age alone, and that is the load-bearing fact.** There is no
  liveness signal: `bench worktree create` writes no lease, and a lease records a
  pid that dies the moment the create hook exits — a request-created worktree
  outliving its creating process is the design. Safety therefore rests on three
  things: `bounds.AssignmentStale` is 7 days, the sweep only ever reports, and the
  explicit cleanup recovers dirty work into a recovery ref before removing. A
  build that shortens the window or makes the sweep reap is a spec deviation, not
  a refinement.
- **The sweep must read the ledger, not a plan.** `PlanAutomatic` returns at its
  first retain verdict, before it reaches the assignment-state branch, so an
  orphan carrying ignored build output — the normal state of a worktree a shift
  ran in — never reaches the orphan label. Deriving the sweep's verdict from the
  plan's reason code reports nothing for exactly the population this row is about,
  while every coverage row stays green.
- **Decisions that stay closed.** Orphans route to `bench worktree clean`; a
  request-derivation override for `release` is rejected. An unstamped record
  counts as aged, which is what drains today's residue. The sweep reports and
  never removes. Three calls the roadmap row did not decide — the 7-day window,
  the summary cap, and ledger compaction — were put to the reviewer and approved
  2026-07-27.
- **A first draft was blocked by a mid-tier falsification pass and rewritten.**
  Its findings were verified against the tree before folding. Do not re-litigate
  them: the lease conjunct, the plan-derived sweep, a named conformance test that
  does not exist, and a compaction story carrying a sign-off it did not have are
  all already fixed in the staged spec.
- **The "preserved" wall at session start is expected.** The pool itself is
  drained, but 17 tree-missing assignment rows and their retained recovery refs
  survive. FT148 compacts one of them and bounds how loudly the rest print; the
  other 16 need FT98's landed proof for reshaped commits and are explicitly out of
  this spec's scope. One ref,
  `refs/bench/recovery/incident-20260712-ambient-probe`, matches no assignment and
  still wants a by-hand look.
- **`ft91-gate-phase-split` stays unretired on purpose** — retiring it destroys
  the veto surface on stories 4, 5, and 9.
- **Push needs `bench gate pin` first** — interactive TTY, so it is the
  reviewer's.

## Next command

`/bench-implement-spec` — in a **fresh mid-tier session**, not this one. Not
`bench shift`: stories 3, 7, and 8–10 are not cheap-line work, so the spec fails
`craft-line`'s venue-routing test.

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
