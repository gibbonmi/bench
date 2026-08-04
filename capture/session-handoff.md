# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `ec3ae90`, clean tree, 23 unpushed commits
Spec: `specs/authoring-hardening/spec.md` (Status: staged), `specs/exact-prospective-landing/spec.md` (Status: staged), `specs/ft187-communication-surface-cut/spec.md` (Status: staged), `specs/pre-push-guard-visibility/spec.md` (Status: staged)
Gate: green at `b388a19` — stale, work tree `28691e8`

## State

**Phase reached: `/bench-implement-spec` opening on `specs/authoring-hardening/spec.md` (FT193).**
Two commits landed since the last handoff: `ec3ae90` recovered that spec and its eight
tickets from the stranded assignment branch `bench/assign/28d4fe19…/756572b7…` and added
FT193, and `eb8685f` is the concurrent session's recovery-fingerprint fix, unchanged here.

- The spec's own sequencing gate (that `recovery-discard` promote first, since two stories
  edit `internal/specbuild`) is satisfied, so the build is unblocked.
- FT174 no longer carries the `Assumptions:` question as open: the reviewer ruled on
  2026-08-04 to retire the field, and FT193 owns it along with the assign-time
  fence-versus-probe refusal.
- The build's ticket frontier is three ownership-disjoint tickets —
  `retire-assumptions-machinery` (`internal/specbuild`), `teach-size-split-signal`
  (`craft-tickets`), and `teach-spec-partition-signal` (`craft-spec`). The remaining five
  serialize behind them on shared fences.
- All capture sources are empty and stay that way; the seven reviewer-held refs under
  `refs/bench/recovery/` remain intact.

## Next command

`/bench-implement-spec specs/authoring-hardening/spec.md`

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
