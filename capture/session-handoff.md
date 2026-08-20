# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `1766e4d`, clean tree, 4 unpushed commits
Spec: `specs/ft230-release-through-bench/spec.md` (Status: staged)
Gate: green at `269e698` — stale, work tree `ec20e95`

## State

FT230's `--full` run is in the review phase. The integration worktree (request
id `ft230-build`, label `ft230-release-through-bench`) holds three green
commits over frozen base `1766e4d1`: T1 `ca6e071f` (adapter selection), T2
`bc18b824` (workflow swap and conformance flip), and `2277209b` (fence
amendment — `tier_test.go` and the `mutable-workflow-action` canary anchor,
both gate-forced, flagged for reviewer veto). Review preflight is green at
base `1766e4d1` / tip `2277209b`; three sonnet axis delegates are reviewing.

Decisions that stay closed: `--adapter` defaults to `fixture` with no
environment twin; the reviewer's tag push is the attended act; promotion stays
out of CI; both step-name contract retirements are recorded in T2's commit.
`capture/learnings.md` holds one undrained entry.

## Next command

`git push` — the board's leading invocable signal (`git`).

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
