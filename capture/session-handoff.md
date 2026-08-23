# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` once this drain lands; the source is the `drain-2026-08-23` worktree on base `a128cd3e`
Spec: none; `specs/` is empty
Gate: green on `a128cd3e` (2026-08-23)

## State

The 2026-08-23 drain is complete. The inbox, the journal, and the retro
folder are empty. The three parked ideas closed by implementation in
`parallel-landings`. The retro fed FT169, FT162, FT213, and FT238. The board
merge rule deferred by ADR 0014 now sits in FT169's authority decision.

The flow report shows a positive net delta (+2 over the window), so the next
drain owes reducing moves. The candidates are FT169, FT162, and FT238, which
each carry one face of the worktree request-token problem.

Closed decisions: merge composition stays the landing primitive, no rebase;
the journal union and the destination default; every phase lands through the
landing verb as guidance, not a hook.

## Next command

`/bench-write-spec` — FT113

## Shape

Rewritten in full at every phase close, pruned rather than accreted. A fresh
session pays for every line it reads cold; drop anything it would not act on.

Operational gotchas are placed by lifetime, not copied here. One that recurs across
phases belongs in `projects/benchkit.md`'s cold-session notes. One scoped to a build
belongs instead in that spec's coverage rows.

This file names at most when you'll hit one, never the command — a second copy
drifts from the source.

Keep the three sections above. **State** holds what is true now, including anything
uncommitted. **Next command** holds the exact harness-native invocation, not a
description of it. This section is the third.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
