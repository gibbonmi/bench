# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `4016db5`, dirty handoff plus unrelated untracked gate assessment, 7 unpushed commits
Spec: `specs/ft91-artifact-hoist/spec.md` (Status: staged)
Gate: green for `4016db5` — final spec-closing gate pending

## State

- **FT91 implementation and semantic review are complete through three green
  ticket commits.**
  `8d074d3` adds the package-scoped lazy artifact set, belts, probes, helper,
  cleanup, and six consumer migrations. `187bc36` records the post-hoist
  measurements and the ≤60-second stop-rule miss. `4016db5` single-sources the
  shared-set count and makes the read-only belt capability-honest under UID 0.
- **Fresh semantic review ran on `gpt-5.6-sol`, high reasoning, fast yolo.**
  Standards and coverage findings above were repaired and re-gated; Spec had no
  findings. Its FIFO build-log-path observation remains a contestable judgment
  for reviewer veto, not an applied change.
- **Closed decisions stay closed.** Non-sharers retain private builds; production
  artifact scripts remain untouched; the separately scoped oracle-semantics work
  triggered by the 133–142-second result is not part of FT91.

## Next command

`$bench-final-check specs/ft91-artifact-hoist/spec.md`

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
