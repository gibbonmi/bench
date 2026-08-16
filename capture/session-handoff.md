# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `b6957e7`, 4 dirty paths, 0 unpushed commits
Spec: `specs/spec-ticket-fence-reduction/spec.md` (Status: staged)
Gate: green at `8516a2d` — stale, work tree `1261719`

## State

The `/bench-what-next` batch is prepared and uncommitted. It removes shipped FT210,
adds FT211 for the staged `spec-ticket-fence-reduction` build, folds the pending idea,
five learnings, and FT210 retro into existing roadmap owners, refreshes the OpenAI
scorecard, and drains ideas/learnings/retros to zero. Main at `b6957e7f` is pushed;
the five-path maintenance diff awaits reviewer approval.

## Next command

`$bench-implement-spec --full spec-ticket-fence-reduction --reviewer fable high`

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
