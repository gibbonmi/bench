# Session handoff

Repository: `89d5d0eebf093c03c876ddd10b63473c-ee22771c3ed7bf922904bf893b610286` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/.bench/worktrees/bench-3325222104/89d5d0eebf093c03c876ddd10b63473c-ee22771c3ed7bf922904bf893b610286`
Branch: `bench/assign/89d5d0eebf093c03c876ddd10b63473c/ee22771c3ed7bf922904bf893b610286` — HEAD `7014bfa`, clean tree, 49 unpushed commits
Spec: `specs/worktree-landed-retirement/spec.md` (Status: staged)
Gate: green at `57162e3` — stale, work tree `9b8a805`

## State

Implementation is complete on retained assignment
`ee22771c3ed7bf922904bf893b610286`. All five tickets are committed serially and
green through the full gate. The frozen review base is `68ebb9ce`; the candidate
includes the approved ownership-fence repairs for the handoff and derived canary
inventory count.

Semantic review is next: run three independent Terra/high axes over the exact
base/source pair. No review, landing, retirement, retro, or push has run yet.

## Next command

`$bench-review-implementation worktree-landed-retirement --base 68ebb9cef9a39a3d35349b9dc4534dad2c044f33 --reviewer terra high`

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
