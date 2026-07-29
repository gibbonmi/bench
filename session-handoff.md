# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `89a6474`, clean tree, 5 unpushed commits
Spec: none staged.
Gate: green at `2aab503` — current

## State

- **The 2026-07-29 what-next pass is committed at `09b8fd0` (drain) plus two
  gated kit commits.** Ideas, journal, and retros are all at zero. The drain:
  FT135 gained the `bench link` symlink-refusal face; FT98 counts re-measured;
  FT97 merged into FT128; FT150 folded into FT141; false-green, load-red, and
  standards-debt theme sections added; sequence reordered fixes-first. The kit
  commits shipped FT137 (its row is gone): `/bench-what-next` now classifies
  every row fix/feature/decision-only on every run, and the board-restructuring
  pass is opt-in via `--restructure`. Decisions in those commits stay closed.

- **Flagged for post-hoc veto:** FT91's hoist grill dropped out of the
  recommended sequence's top three under the reviewer's fixes-first
  instruction; the row itself stays HIGH.

- **This repo's pre-push hook was refreshed by hand-copy from the current
  template** after `bench link` refused to converge the kit repo (symlink
  parent conflict, now FT135's third face). `bench guards` reports the full
  manifest.

## Next command

`/bench-implement-spec` — FT163 via the light path.

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
