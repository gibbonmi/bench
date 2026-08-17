# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `73c97aa5` (frozen review base); integration worktree
`bench worktree path d58e5fed7d1634c79cdeacadfa519a8c` on
`bench/assign/…/d58e5fed7d1634c79cdeacadfa519a8c`, reviewed tip = the commit that
carries this file
Spec: `specs/spec-ticket-fence-reduction/spec.md` (Status: staged until `bench worktree land`)
Gate: green in the worktree at the tip

## State

`/bench-implement-spec --full` has finished build and review; the source is ready to
land. Ten tickets landed green (eight from the approved breakdown plus two repair
tickets the build and review added). Review round 1 (`opus`/high, three axes over
`73c97aa5..578811f2`) produced 8 raw findings → 4 repair targets, closed at `8be12dd5`;
round 2 over the repair delta was clean. A parallel session (`a78c0db0`) reviewed the
same source and stopped with five reviewer questions; all five are answered and
recorded in the spec: 73-line budget accepted, story 2 kept, `craft-spec | 150` row
and the four spec/handoff fences ratified, no rewrap. That session is abandoned; do
not resume it.

## Next command

From `~/workspace/bench` on `main`:
`bench worktree land --request 02c6f79b54a3505522af4c84014e0670 --base 73c97aa5 --source-tip <tip> --spec spec-ticket-fence-reduction -m "<msg>" <worktree path>`,
then `/bench-final-check`.

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
