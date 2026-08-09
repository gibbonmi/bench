# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `14deeb42`; the `$bench-what-next` batch is uncommitted.
Spec: `byte-preserving-axi-foundation` remains staged.
Gate: green at `2000f0a` — stale against the pending batch.

## State

`branch-native-build-test-architecture` is retired in the pending roadmap-drain
batch. The two implementation retros and three open learnings are drained into
FT162, FT171, FT98, FT102, FT168, FT174, and FT184; the unrelated
`decisions/spec-build-review-gate-cadence.md` remains outside the batch.

The batch awaits reviewer approval. It must be committed once, on green, with a
subject ending `spec-retire: branch-native-build-test-architecture`.

## Next command

`$bench-shape-idea`

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
