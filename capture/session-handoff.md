# Session handoff

Repository: `77d1764072080892a27f7bfd09a71243-bf260a4fc120c3c5973ef54660eaede6` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/.bench/worktrees/bench-2826441890/881246de1c65ffcc8a4acb3ec49c6ff8-5f5b1987a92bcfa209d3dd9f26898586`
Branch: `bench/assign/881246de1c65ffcc8a4acb3ec49c6ff8/5f5b1987a92bcfa209d3dd9f26898586` — HEAD `d46b315c`, 5 dirty paths, 0 unpushed commits
Spec: `specs/module-size-split/spec.md` (Status: staged)
Gate: green (fresh, this worktree)

## State

This drain reconciled the roadmap against the tree; no shipped rows and
nothing new to retire. Both open capture entries merged into existing
decision-required rows instead of opening new ones. FT169 gained a
dirty-destination occurrence and a named sub-question — two candidate fixes —
inside its already-open landing-authority decision. FT8 gained the
escalation-ladder proposal and its evidence, with the review trigger brought
forward to 2026-08-30. The capture inbox and journal are both empty.

The flow report shows a positive net delta (opened 5, retired 1, target not
met). Two restructure candidates stand for a future `/bench-drain
--restructure`. Fold FT200, FT207, FT178, FT162, and FT173 into FT169's
worktree/landing-authority theme. Fold FT213, FT214, FT236, and FT237 into
the existing "Standards debt" craft-visit batch alongside FT117, FT179, and
FT94.

`## Recommended sequence` line 1 names `/bench-implement-spec` for FT238.
FT238's `Next:` is `ticket` (light path, no spec), and no `tickets/FT238*`
file exists yet. Flagged for the reviewer, not corrected this pass, since it
falls outside this drain's capture-and-journal scope.

`specs/module-size-split` stays staged and untouched by this pass.

## Next command

`/bench-implement-spec specs/module-size-split/spec.md`

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
