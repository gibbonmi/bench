# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`; integration worktree `bench worktree path spec-ste`
Branch: `bench/assign/e10f369f24649fcd21dedc65f5e36350/d5bd1e58c5c0e89e21ea6c277d0aab46` — review base `ba67efc4` (reviewed graph `4d8aa66e` merged with `main` `9618dde8`); tickets 01–03 landed through `55ed7a0c`
Spec: `specs/roadmap-flow/spec.md` — Status: staged, five tickets approved
Gate: green at every ticket commit on the integration worktree

## State

`/bench-implement-spec --full` is in its build phase. Tickets 01 (`bench
roadmap --flow`), 02 (`Next:` grammar, missing-line class behind the
`rowNextMissingEnforced` switch), and 03 (`retro-improvement-markers` check)
landed green. Ticket 04 builds on fable/high in its own worktree; ticket 05
then runs on sonnet/medium by reviewer decision (spec said opus). Reviewer for
the review phase: sonnet/high. Cap 3 fix loops per ticket. The delegate
worktrees `roadmap-flow-t01..t04` hold landed or in-flight diffs and are not
yet released.

## Next command

`/bench-implement-spec --full specs/roadmap-flow/spec.md --reviewer sonnet high`

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
