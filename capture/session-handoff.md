# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `ddd7b6f3` before this approved drain; assessment artifacts remain uncommitted outside its scope
Spec: `specs/single-build-serial-gate/spec.md` (Status: staged)
Gate: pending the approved drain's path-scoped `bench commit`

## State

The approved 2026-08-13 `/bench-what-next` drain removes the
canary-planted-reason-ownership retro. FT162 gains prospective-failure locality
and commit-timing receipts; FT106 gains a claim-path census; FT164 gains
synthetic-registry disclosure and aggregate-candidate repair rereads. The
producer-derived baseline reminder is already canonical and adds no new rule.

`ASSESSMENT.md` and `capture/FIXES.md` are pre-existing assessment artifacts
outside this drain and remain uncommitted. The staged serial-gate spec remains
the leading actionable work after this documentation commit.

Decisions staying closed: `internal/gittest` fence amendment; guards
real-stale fixture; coverage `why` rewording; FT175 shape unblocked by the
capstone landing per the 2026-08-02 reviewer ruling.

## Next command

`/bench-implement-spec` — `specs/single-build-serial-gate/spec.md`

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
