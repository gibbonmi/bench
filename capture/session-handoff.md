# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main`; the tail source is the `ft113-tail` worktree on base `148f3a68`
Spec: `specs/landing-authors-the-flip/` retired in this tail; it landed as `148f3a68` with `Status: implemented`
Gate: green on `148f3a68` (2026-08-23)

## State

FT113 landed as `148f3a68` from source `0e17d428..e673e2d5`. This tail retires
the spec folder, adds ADR 0015 (the landing verb is the one flip author), writes
the retro, and refreshes the Claude scorecard. The retire verb now names the
board remainder: the `ROADMAP.md` row `FT113` and `roadmap/FT113.md` stay for
the drain.

Pending for the drain: one parked idea about the landing verb's `next=`
sanitizer, one retro, and the FT113 board row and detail file.

## Next command

`/bench-drain`

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
