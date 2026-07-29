# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `6cd2657`, 1 dirty path, 1 unpushed commit
Spec: none staged.
Gate: pending at `23293e0`

## State

- **The 2026-07-29 drain is committed at `6cd2657`.** FT154 shipped and its
  row is gone; ideas, journal, and retros are all at zero. New rows FT162–166;
  clauses merged onto FT91, FT135, FT156, FT161; FT107 reworded. Decisions in
  that commit stay closed.

- **The gate cache reads `interrupted-pending` at `23293e0`.** A sanctioned
  defect repro (`bench gate --help`, now FT163) overwrote the green verdict;
  the tree is not suspect. FT163's own landing commit — or any gated commit —
  replaces it.

- **This repo's pre-push hook was refreshed by hand-copy from the current
  template** after `bench link` refused to converge the kit repo (symlink
  parent conflict, parked in `IDEAS.md`). `bench guards` now reports the full
  manifest.

- **`IDEAS.md` holds one parked line** (the `bench link` refusal), dirty in
  the main checkout, awaiting the next drain.

## Next command

`/bench-shape-idea`

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
