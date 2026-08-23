# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `8923d86`, 5 dirty paths, 1 unpushed commit
Spec: none staged.
Gate: red at `4bc4e72` — stale, work tree `7afa7d9`

## State

The drain batch reconciles 73 existing rows and retires none. It records the
2026-08-22 prose pass as a partial FT100 shipment. It merges the deferred STE
semicolon experiment into FT231. It also graduates the gate-process-boundary
learning into FT245 with `Next: kit-edit`. The idea inbox and learning journal
are empty. The recommended sequence remains FT169, FT113, then FT214.

Closed decisions: evaluate the semicolon rule only after FT231's harness lands;
do not treat partial gate output as a process boundary. The positive roadmap
flow delta requires a later reducing-moves proposal; this default drain does not
silently restructure the board.

## Next command

`/bench-debug`

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
