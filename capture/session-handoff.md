# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `48f93ed7`
Spec: none active
Gate: stale; the pending drain batch requires a fresh commit gate

## State

The landing-refusal-diagnostics spec is implemented and retired at `48f93ed7`.
The pending drain batch removes shipped FT233, folds its retro into FT164,
FT215, and FT224, folds the spec-review learning into FT214, and creates FT242
for the reproduced inherited-environment false green. No ideas, retros, or open
learnings remain after the batch.

## Next command

`$bench-write-spec` — FT242.

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
