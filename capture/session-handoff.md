# Session handoff

Repository: `7e3c28702eb57ca1357eab2a37c5983b-8bcd0dd48a7788364b323e30dafa109b` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/.bench/worktrees/bench-2826441890/7e3c28702eb57ca1357eab2a37c5983b-8bcd0dd48a7788364b323e30dafa109b`
Branch: `bench/assign/7e3c28702eb57ca1357eab2a37c5983b/8bcd0dd48a7788364b323e30dafa109b` — HEAD `4fb6e2b`, 10 dirty paths, 20 unpushed commits
Spec: none staged.
Gate: red at `bb8b418` — current

## State

Destination reconciliation now preserves tracked edits that appear immediately
before its reset. The fix landed as `4fb6e2b9a38570e2cd43c2316746a38ef060c7bf`.

Bench now enforces the worktree-only phase rule at publication. `bench commit`
refuses the primary checkout and directs users to `bench worktree create`.
Status and commit share one primary-checkout classifier.

The latency decision map and its invocation census remain intact. The first of
two specs owns one selected test binary, explicit environment and directory
inputs, worktree owner seams, and before-and-after demand measurements.

The second spec follows measured first-spec results. It owns only necessary
pure-test parallelism and the slow-package regression budget. Descendant-spawning
journeys remain serial, and fresh runs keep `-count=1`.

## Next command

`$bench-write-spec decisions/worktree-test-latency.md --reviewer terra high`

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
