# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `9deb062`, clean tree, 10 unpushed commits
Spec: `specs/exact-prospective-landing/spec.md` (Status: staged), `specs/ft187-communication-surface-cut/spec.md` (Status: staged), `specs/go-build-cache-footprint/spec.md` (Status: staged), `specs/pre-push-guard-visibility/spec.md` (Status: staged)
Gate: green at `254d2ef` — current

## State

**Phase reached: implementation is integrated and exact-candidate review is next.**

Run `869a2d33cab96f962882762a9ea56fd21c952976248e64fa917f6da6c48dbca3` has exact candidate `77bc0e6b5b1be8b753dbdac49e0b19105d17e67d`; all eight assignments are checkpointed, integrated, and released. The adapter's AC1-AC9 rows are green on the repaired composition. Focused R2, the six-row prospective fixture repro, `internal/commit`, full `internal/contract/runtime`, and `internal/landing` are green. Delegate and coordinator mutations independently red and were restored.

Critical fixture prerequisites were repaired through lifecycle assignments: Story 5 common-Git-dir coordination, pre-lock fault timing, prospective-safe runtime marker/diagnostic fixtures, and the R2 committed oracle base. No adjacent noncritical cleanup was taken. The candidate now requires the fresh exact-candidate Standards/Spec/Coverage review and receipt before promotion.

Closed decisions: preserve candidate identity during review; concrete review defects return through repair assignments; contestable design findings require reviewer veto. The large-diff cross-harness falsification offer remains unanswered, so promotion must not run until the reviewer answers it. If the next recomposed attempt fails, hand the preserved blocker and receipt to Fable rather than retrying locally.

## Next command

`$bench-implement-spec --full specs/exact-prospective-landing/spec.md`

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
