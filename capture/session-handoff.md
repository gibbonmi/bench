# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `0294a0a`, clean tree, 15 unpushed commits
Spec: none active (ft181-precondition-residuals promoted and retired)
Gate: green at `0294a0a` — current

## State

- `ft181-precondition-residuals` is done end-to-end: promoted at `e11df18` (candidate a25e2803, all 27 coverage rows dispositioned, 13 review findings resolved/endorsed), spec retired at `0294a0a`, retro written to `capture/retros/ft181-precondition-residuals.md`. The FT181 roadmap row and its run-window commit-sequencing rule retired with it.
- Kit guidance landed from the run's process findings (`2892501`, `60cb55d`): craft-tickets requires a `Contracts:` field per ticket; craft-spec gained the composition degenerate and the existing-control edge rule; craft-line makes a spec's per-story line a per-ticket ceiling; craft-delegate charges carry the fixture-and-seam inventory and focused checks only; implement-spec scales pre-build research to spec staleness.
- Capture holds 4 parked ideas (git FIFO-gitdir upstream hang; receipt-skeleton helper; fixture-inventory generator; injected-interface composition audit) and 4 open learnings (skipped derivation steps; finding triage; process self-audit; review-receipt disposition vocabulary) plus 1 pending retro — all awaiting the `/bench-what-next` drain.
- Decisions that stay closed: all rulings in the six active decision maps, including the 2026-08-03 gate-pipeline and gate-structure amendments; FT181's build decisions now live in the shipped code, its retained run record (`bench spec history ft181-precondition-residuals`), and the retro.
- 15 unpushed commits on `main`; the reviewer owns the push.

## Next command

`/bench-what-next`

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
