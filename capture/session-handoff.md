# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `3f37bf73`
Spec: none (roadmap maintenance proposed)
Gate: unavailable

## State

This uncommitted drain batch removes shipped FT243 and its retired detail file,
clears the three journal entries and the `learnings-dated-line-visibility` retro,
and refreshes the sequence. The journal and retro evidence merge into existing
rows: FT238, FT214, FT233, FT224, FT113, and FT215.

The batch needs reviewer approval before its one green capture commit. FT214's
guidance change and FT215's scoped-gate design remain roadmap work; neither is
implemented by this maintenance pass. The Claude scorecard already carries the
landing evidence and remains persistent.

## Next command

`$bench-write-spec`

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
