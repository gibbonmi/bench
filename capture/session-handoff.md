# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `cd355f45`, tree carries the uncommitted spec below, this handoff, one `capture/learnings.md` entry from the review round, and one ambient parked idea in `capture/IDEAS.md`
Spec: `specs/ft226-test-home-isolation/spec.md` — `Status: staged`, awaiting reviewer sign-off on the spec-and-tickets pair
Gate: green at `cd355f45`

## State

**FT226 is specified, not approved.** The spec and its three serial tickets are
written, `bench coverage --check` is green (12 rows), `bench preflight build` is
green on every row-ownership check, and the mid-tier review round took two
iterations to accept (folded; see the spec's `Verification log` and the dated
`capture/learnings.md` entry). Nothing is committed: the spec
directory is untracked until you sign off.

What the build does: `reauthorizeFixture` binds a per-test `BENCH_HOME` like its
siblings (ticket 01); `internal/worktree` gains a `TestMain` that runs the
package under a process-private `BENCH_HOME` and exits red on residue, naming the
leaking test (ticket 02, proven by mutation); the operator's 1,690 orphaned
`001-<digits>` pool keys (63 MB) are swept once, plan-before-apply, under a
dangling-`gitdir:` predicate (ticket 03, destructive, outside the tree — your
spec sign-off is its approval).

Decisions left open for you, flagged in the spec: the oracle seam is in-package
`TestMain` rather than a gate-level private home for the test/race phases (priced
in Out of scope); the sweep is a throwaway script, not a `bench worktree` verb
(priced). Measured facts: the full ordinary suite under an empty sentinel
`BENCH_HOME` writes exactly ten `worktrees/001-*` entries, all from
`TestReauthorize*`, and nothing else.

`capture/IDEAS.md` holds one uncommitted parked idea from the drain session
(occurrence-ledger freeze check). `bench preflight build` reds on it as an
unauthorized dirty path until it is committed; commit it before the build starts.

## Next command

Sign off the spec and tickets (the approval table is in the write-spec session's
closing message; the spec file itself is the veto surface), commit the spec
directory, then start a fresh session with:

`/bench-implement-spec ft226-test-home-isolation`

Line for the build: `opus` / medium per the spec's story groups.

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
