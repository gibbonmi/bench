# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `2dd71f8`, 4 dirty paths, 6 unpushed commits
Spec: none staged.
Gate: green at `d79ad80` — stale, work tree `6456d77`

## State

The current `$bench-what-next` batch is uncommitted pending reviewer approval.
It removes FT175 at the reviewer's direction, reconciles every dependent
goal-track and sequence reference, dismisses the sole learning because it
repeats the phase's existing trusted-snapshot boundary, and leaves no ideas,
learnings, retros, occurrence discrepancies, or staged specs.

The refreshed sequence is FT171, FT203, then FT198. `ASSESSMENT.md` and
`capture/FIXES.md` remain pre-existing assessment artifacts outside this batch.
Both registered assignment worktrees remain foreign and untouched.

## Next command

`$bench-shape-idea`

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
