# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `0c02192`, 1 dirty path, 0 unpushed commits
Spec: none staged.
Gate: green at `469c6d0` — stale, work tree `3c590bb`

## State

- **The roadmap has 50 open rows.** The restructure folded eight duplicate rows
  into their owner surfaces; FT145 left after its ambient-state fix landed.
  FT113 now carries only the unresolved post-`bench commit --spec` cache decision.
- **Dependencies are explicit and single-sourced.** Literal means blocked;
  recommended means cheaper or less churn when specified after the prerequisite.
  FT107 on FT141 is the only literal edge.
- **Capture is empty.** There are no parked ideas, open learnings, pending
  retros, staged specs, or parse failures.
- **The next priced work is FT123 + FT124.** It follows the completed
  ambient-state batch and adds worktree label resolution plus structured test
  triage.

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
