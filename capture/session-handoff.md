# Session handoff

Repository: `e10f369f24649fcd21dedc65f5e36350-d5bd1e58c5c0e89e21ea6c277d0aab46` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/.bench/worktrees/bench-2826441890/e10f369f24649fcd21dedc65f5e36350-d5bd1e58c5c0e89e21ea6c277d0aab46`
Branch: `bench/assign/e10f369f24649fcd21dedc65f5e36350/d5bd1e58c5c0e89e21ea6c277d0aab46` — HEAD `51d1da4`, 3 dirty paths, 2 unpushed commits
Spec: `specs/inherited-toolchain-environment/spec.md` (Status: staged), `specs/roadmap-flow/spec.md` (Status: staged)
Gate: green at `ba38063` — stale, work tree `4dac3bb`

## State

`/bench-implement-spec --full` is in its build phase. Line: opus/medium for
tickets 01, 02, 03, 05; fable/high for ticket 04; reviewer sonnet/high; cap 3
fix loops per ticket. Tickets 01, 02, 03 run as parallel write delegates in
their own worktrees and land serially on the integration source in ticket
order; 04 then 05 follow. The light-path kit edit requiring ASD-STE100 prose in
specs and tickets landed on `main` at `9618dde8`. `capture/learnings.md` holds
an open two-arm spec experiment entry for the next `/bench-drain`.

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
