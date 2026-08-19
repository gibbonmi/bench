# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`; the build runs in the integration worktree with label
`ft226 implementation` (assignment `6a22e66f47a693f42e306c0ec030e1ba`, request
`ft226-build`) — address it by label, never by a pasted path
Branch: `main` — destination HEAD `a9e8e232`; the source's frozen review base is
`b2d3d625`, which a concurrent session has since moved past
Spec: `specs/ft226-test-home-isolation/spec.md` — `Status: staged`, verification log written
Phase reached: `/bench-implement-spec --full`, implementation complete, all three tickets landed green in the source; review not yet run

## State

All three tickets are committed on green in the integration source. Ticket 01
binds `reauthorizeFixture`'s `BENCH_HOME`. Ticket 02 adds
`internal/worktree/main_test.go` — a `TestMain` running the package under a
process-private `BENCH_HOME`, a residue predicate, and two tests. Ticket 03 swept
1,710 orphaned pool keys, 91 MB down to 51 MB, and wrote the spec's verification
log. Every acceptance row is covered; the evidence is in that log.

**A concurrent session landed `e1b44e62` and `a9e8e232` on `main` during this
build.** The source's review base `b2d3d625` is therefore behind the destination,
so the landing needs `bench worktree reauthorize` to the current base before
`bench worktree land`. Those two gate runs also leaked ten fresh `001-<digits>`
keys, because `main` does not yet carry ticket 01; a second sweep pass clears them
once this spec lands. The sweep script is throwaway and was never committed, so a
fresh session regenerates it from the spec's SW1 predicate.

Open for the reviewer, flagged in the verification log: the spec's probe (b)
command never reaches `TestMain`, because `TMPDIR=/nonexistent` fails in the go
driver first. The `GOTMPDIR=/tmp TMPDIR=/nonexistent` form does reach it and reds
correctly. DT4 holds; only the spec's probe text is wrong, and the log records the
correction as a finding for veto.

## Next command

`/bench-review-implementation` over the frozen base `b2d3d625` and the source tip,
then `bench worktree reauthorize` to the current destination HEAD, then
`bench worktree land`, then `/bench-final-check`.

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
