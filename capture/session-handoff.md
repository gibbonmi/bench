# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `d127289`, clean tree, 11 unpushed commits
Spec: `specs/pre-push-guard-visibility/spec.md` (Status: staged)
Gate: green at `02992a8` — current

## State

**Phase reached: FT135 repair tickets committed; recomposition blocked by unrelated working-tree dirt.**

Candidate `fcbe6cec459f6e4b4a3496d3b6faa6da68c123fd` has all prior assignments integrated and a fresh three-axis review receipt (`6efc8e082f5241ee79eec6301a44121281c9ee078988502fb260c486ba145f85`). It accepted P1: doctor repair overwrites an existing stale `0644` pre-push hook without restoring execute bits, and S1: the public `PrePushMarker` comment names removed `ClassifyPrePush`. It rejected the alleged unused `managedPrePushBody` fixture because `internal/contract/surface/doctor_rows_test.go` consumes it; Coverage found no miss across all 40 rows.

The two independently-green repair tickets are committed in `70aff37`: `restore-doctor-hook-execute-mode.md` (doctor repair plus its runtime contract) and `correct-hook-marker-comment.md` (current-only marker documentation). The attempted ticket commit's gate is green. `bench spec build promote pre-push-guard-visibility` must next recompose the candidate onto `70aff37`; it currently refuses because the checkout is dirty.

Preserve both dirty paths exactly: `.agents/skills/bench-craft-tickets/SKILL.md` is foreign concurrent work, and this handoff records the pause. Do not clear, stage, or commit either as part of FT135. Do not start a new lifecycle run. The final cross-harness falsification pass remains explicitly skipped; required lifecycle reviews and promotion are not skipped.

## Next command

`/bench-final-check`

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
