# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `b6dd0a8`, clean tree, 14 unpushed commits
Spec: `specs/pre-push-guard-visibility/spec.md` (Status: staged)
Gate: green at `dca14fb` — stale, work tree `30aa496`

## State

**Phase reached: FT135 recomposition repair landed; promote is the next lifecycle operation.**

Candidate `fcbe6cec459f6e4b4a3496d3b6faa6da68c123fd` has all prior assignments integrated and a fresh three-axis review receipt (`6efc8e082f5241ee79eec6301a44121281c9ee078988502fb260c486ba145f85`). It accepted P1: doctor repair overwrites an existing stale `0644` pre-push hook without restoring execute bits, and S1: the public `PrePushMarker` comment names removed `ClassifyPrePush`. It rejected the alleged unused `managedPrePushBody` fixture because `internal/contract/surface/doctor_rows_test.go` consumes it; Coverage found no miss across all 40 rows.

The two independently-green repair tickets are committed in `70aff37`: `restore-doctor-hook-execute-mode.md` (doctor repair plus its runtime contract) and `correct-hook-marker-comment.md` (current-only marker documentation).

The recomposition marker-conflict defect from the debug receipt is fixed and landed green as HEAD `b6dd0a8`: `recomposePromotion` now bootstraps against `greenMarker(s.root, subject.branch)` instead of the run's recorded base, retains the recorded base until `finishRecomposition` completes, and carries the strong regression `TestPromoteRecompositionBootstrapsAgainstLiveMarkerAndRetainsRecordedBase` (red under the recorded-base mutation, verified). The fenced repair ticket is `specs/pre-push-guard-visibility/tickets/repair-recompose-bootstrap-marker.md`. Do not cherry-pick or absorb any dirty assignment worktree; the fix was authored fresh on main. Promotion now recomposes candidate `fcbe6cec` onto the moved tip; recomposition discards the bound review, so the round after promote is a fresh composed review, then promote again. The final cross-harness falsification pass remains explicitly skipped; required lifecycle reviews and promotion are not skipped.

## Next command

`bench spec build promote pre-push-guard-visibility`

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
