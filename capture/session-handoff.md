# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `6974e9f1` before this approved retirement; assessment artifacts remain uncommitted outside its scope
Spec: `specs/single-build-serial-gate/` retired after Fable confirmed its implementation already landed in ancestor `040ead11`
Gate: pending the retirement's path-scoped `bench commit`

## State

Fable/medium rejected the restored staged serial-gate plan because its reds and
migration tickets described behavior already present in the current tree. The
implemented-state correction landed as `6974e9f1`; this retirement removes the
stale spec and ticket DAG and removes their directives from `ROADMAP.md`.

FT171 remains open only for the gate-budget map's current work: reconcile landed
#23, mark lifecycle-removed #24 moot and remove it from #26's blockers, then
resume #25 and run #26 before pricing outer concurrency. Do not recreate or
implement the retired serial-gate plan.

`ASSESSMENT.md` and `capture/FIXES.md` remain pre-existing assessment artifacts
outside this retirement and uncommitted. Both registered assignment worktrees
remain foreign and untouched.

## Next command

`/bench-shape-idea` — reconcile FT171 #23–#26 in `decisions/gate-budget.md`

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
