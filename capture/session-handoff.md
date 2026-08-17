# Session handoff

Repository: `02c6f79b54a3505522af4c84014e0670-d58e5fed7d1634c79cdeacadfa519a8c` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/.bench/worktrees/bench-3325222104/02c6f79b54a3505522af4c84014e0670-d58e5fed7d1634c79cdeacadfa519a8c`
Branch: `bench/assign/02c6f79b54a3505522af4c84014e0670/d58e5fed7d1634c79cdeacadfa519a8c` — HEAD `5ced385`, clean tree, 2 unpushed commits
Spec: `specs/spec-ticket-fence-reduction/spec.md` (Status: staged)
Gate: green at `8f88a5a` — stale, work tree `18703da`

## State

`/bench-implement-spec --full` is in the build phase. Four of eight tickets are landed
green on the integration source: schema descriptor (`06f70f7e`), reduced header
(`2a333a81`), projection + fixture migration in one commit (`f7131404` — the
projection breaks callers only the migration ticket is fenced to repair), and the
template move (`5ced3850`). Remaining, in `Blocked by:` order:
`collapse-the-review-loop-and-narrow-reviewer`, `remake-craft-spec-and-craft-tickets-on-their-sources`,
`realign-the-consumers-glossary-and-docs`. Each is a `fable`/high write delegate in
the worktree; the coordinator probes and lands with `bench commit`. Review runs
`sonnet`/high (reviewer override of map #15's `fable`). Contestable, for veto: the
mid-build spec edits `9346dec6`/`22567a53` (needle inventory 47 not 20; craft-spec
takes a 150-line budget row; the 60-line row moved to the loop-collapse ticket).

## Next command

`/bench-implement-spec --full spec-ticket-fence-reduction --reviewer sonnet high`

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
