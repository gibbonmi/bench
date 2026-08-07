# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `/home/devuser/workspace/bench`
Branch: `main` — pre-commit HEAD `f5fbc457c678`, 27 unpushed commits
Specs: `specs/gate-decision-test-seam/spec.md` (Status: staged), `specs/pre-push-guard-visibility/spec.md` (Status: staged)
Gate: stale on the pre-commit tree; this handoff lands only through the roadmap batch's green `bench commit`

## State

**Phase reached: roadmap reconciled and capture drained.**

FT187 is promoted and retired. Its retro is drained into FT162's retained
promotion timing and memory evidence and FT174's declared-closure omission
review. FT171 now owns the staged gate-decision-test-seam slice; FT135 remains
staged behind it. The ideas inbox, learnings journal, occurrence ledger, and
retros directory are empty.

Closed decisions stay closed: the gate-decision slice preserves the existing
read-only decision seam and representative full-engine composition proofs;
pre-push guard visibility still follows its compiled decision map; FT187's
communication surface and atomic closure contract are shipped state rather than
work to reopen. The refreshed sequence drains both staged specs before new spec
authoring.

## Next command

`$bench-implement-spec`

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
