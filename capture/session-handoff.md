# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `5e5ba55`, 18 dirty paths, 1 unpushed commit
Spec: none staged.
Gate: green at `2108505` — stale, work tree `c23acc8`

## State

The 2026-08-22 drain is staged and uncommitted. It reconciles 71 rows (none
shipped), drains six ideas and seven journal entries, disposes both pending retros,
opens FT243 and FT244, and refreshes the sequence. The light-path item from the
journal (the `ste-prose.md` field-line note) landed on its own as the commit before
this one. The batch waits for the reviewer's approval, then one `bench commit` over
every touched path.

## Next command

`/bench-debug`

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
