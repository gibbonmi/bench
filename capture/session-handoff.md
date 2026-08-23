# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `3728719`, 3 dirty paths, 6 unpushed commits
Spec: none staged.
Gate: green at `e5ec6a3` — stale, work tree `02cfc81`

## State

The full `inherited-toolchain-environment` lifecycle is complete. Reviewed pair
`6c867eb5..d6918bad` landed green as `63dde6ae`; initial review found Standards
1, Spec 0, Coverage 4, and repair-scoped re-review found 0/0/0. Post-merge spec
and FT242 retirement landed green as `3728719d`. The retained source was released;
the landed-worktree sweep found only the unrelated dirty
`ste-prose-progressive-loading` assignment and retained it. The implementation
retro and OpenAI scorecard are pending capture for drain alongside the reviewer's
pre-existing `capture/IDEAS.md` entries; do not commit or revert them here.

## Next command

`$bench-drain`

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
