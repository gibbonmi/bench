# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `72b87a1`, clean tree, 4 unpushed commits
Spec: `specs/ft91-artifact-suite/spec.md` (Status: staged)
Gate: green at `9e49084` — stale, work tree `0d84075`

## State

- **The `--full` run has reached final-check.** Build commits `6016be6` and
  `94b182e` landed the split and measurement. Top-tier semantic review
  (`ft91_semantic_review`) and the approved top-tier Codex CLI falsification
  found concrete topology, duplication, comment, and arithmetic defects.
- **The terminal repair pass is green at `72b87a1`.** The reviewer-requested
  Sol repair preserved the exact 33-test behavior while single-sourcing shared
  facts, structurally binding the four TestMain runners, rejecting extra
  packages and inline cache policy, and correcting the measurement arithmetic.
  Its first gate red exposed an over-broad shared-cache strip; the exact
  GOPROXY-off seam was repaired before the green landing.
- **Measured result remains:** focused suite 50.97 s with overlapping subject
  processes; fresh changed-tree gate 89.91 s versus 128 s. The dormant outer
  width cap is a separate reviewer decision; no scheduler policy changed.

## Next command

`$bench-implement-spec --full ft91-artifact-suite`

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
