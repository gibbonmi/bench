# Session handoff

Repository: `89d5d0eebf093c03c876ddd10b63473c-ee22771c3ed7bf922904bf893b610286` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/.bench/worktrees/bench-3325222104/89d5d0eebf093c03c876ddd10b63473c-ee22771c3ed7bf922904bf893b610286`
Branch: `bench/assign/89d5d0eebf093c03c876ddd10b63473c/ee22771c3ed7bf922904bf893b610286` — HEAD `db1d967`, clean tree, 49 unpushed commits
Spec: `specs/worktree-landed-retirement/spec.md` (Status: staged)
Gate: green at `4f15dfa` — stale, work tree `9df3007`

## State

Implementation is active on retained assignment
`ee22771c3ed7bf922904bf893b610286`. Ticket
`count-and-advertise-landed-assignments.md` landed green at `7f5b4798`; repair
`db1d9679` restored its staged ticket metadata after preflight correctly rejected that
out-of-fence path. The candidate remains based on `68ebb9ce`.

Frontier: `plan-the-landed-set-under-one-fingerprint.md`, then
`apply-the-landed-plan-and-settle-records.md`, then the independent final pair
`refuse-a-half-applied-landed-set.md` and `make-release-a-workflow-step.md`.

The full run uses fresh ticket write delegates and three independent Terra/high review
axes before `bench worktree land`; no review or landing has run yet.

## Next command

`$bench-implement-spec --full worktree-landed-retirement --reviewer terra high`

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
