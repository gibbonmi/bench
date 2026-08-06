# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `5d66acf`, uncommitted roadmap-drain batch
Spec: `specs/ft187-communication-surface-cut/spec.md` (Status: staged), `specs/go-build-cache-footprint/spec.md` (Status: staged), `specs/pre-push-guard-visibility/spec.md` (Status: staged)
Gate: green at `5d66acf`

## State

**Phase reached: exact-prospective-landing retro drained into an uncommitted reviewer batch.**

The trusted schema-3 full roadmap snapshot found no ideas, open learnings, or occurrence discrepancies. FT188 was stale because its spec retired in `126c597`; the batch removes it, folds its retro's lifecycle, guidance, and immutable-snapshot coverage into FT162, FT164, and FT99, and removes the drained retro. The sequence now starts with FT195, followed by FT187 and FT135.

The batch awaits reviewer approval before its one green commit. Its source is the full snapshot captured in this run; do not reread or recreate its retro evidence.

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
