# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `d7a6e67`, 3 dirty paths, 7 unpushed commits
Spec: none staged.
Gate: green at `2180b3d` — stale, work tree `bbfa62c`

## State

FT171 closed at `2dd71f83`: gate-budget decision #27 retired #8 because the
serial baseline already meets the destination, and no spec follows. The roadmap
therefore carries no FT171 row or sequence item.

FT203 now owns decision #26's wider `internal/worktree` flake evidence: two reds
in six package runs on `a3b599ea`, including
`TestListCommandCheckedInCompletedAssignmentTerminalPair`, while three focused
repetitions passed. Debug the package family rather than overfitting to the
row's original named test. FT198 follows in the sequence.

There are no ideas, learnings, retros, occurrence discrepancies, or staged
specs. `ASSESSMENT.md` and `capture/FIXES.md` remain foreign assessment
artifacts; both registered assignment worktrees remain foreign and untouched.

## Next command

`$bench-debug`

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
