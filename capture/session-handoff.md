# Session handoff

Repository: `d72c470fdf58bc0747305add9ff619f2-5880409cc8f500d39e38170e2485eef7` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/.bench/worktrees/bench-2826441890/d72c470fdf58bc0747305add9ff619f2-5880409cc8f500d39e38170e2485eef7`
Branch: `bench/assign/d72c470fdf58bc0747305add9ff619f2/5880409cc8f500d39e38170e2485eef7` — retained integration tip `03a5f736` (tickets 01–06 landed), plus the uncommitted ticket-07 evidence work
Spec: `specs/worktree-test-latency/spec.md` (Status: staged)
Frozen review base: `4a8aa16a`

## State

Tickets 01–06 are landed on the integration tip. The suite selects one Bench
executable and resolves effect inputs explicitly. Three pure packages own
landing, lifecycle, and reclaim policy, and one harness runs serial
proof-backed journeys.

Ticket 07 evidence is written and uncommitted:
`specs/worktree-test-latency/evidence/demand-reduction.md`, one `CHANGELOG.md`
entry, and this handoff. Measured at `03a5f736`: package spans 58.022, 56.779,
and 56.898 seconds (baseline median 125.790), whole suite 70.122 seconds. The
publication connection wait was not visible; the evidence claims nothing from
it.

These decisions stay closed: no `t.Parallel`, no scheduler, `-count=1` and the
one ordinary gate driver retained, real Git kept at representative journeys,
the publication wait and the WSL Go bootstrap owned by separate debug tracks.
Measured parallelism and the regression budget belong to the second spec.

## Next command

`/bench-review-implementation` over frozen base `4a8aa16a` and the source tip.

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
