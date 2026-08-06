# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `31dd98a`, clean tree, 9 unpushed commits
Spec: `specs/ft187-communication-surface-cut/spec.md` (Status: staged), `specs/go-build-cache-footprint/spec.md` (Status: staged), `specs/pre-push-guard-visibility/spec.md` (Status: staged)
Gate: green at `2ca9812` — current

## State

**Phase reached: covers-traceability promoted, retired, and retro'd; next is the drain.**

The covers-traceability build is terminal: promotion squash `5ff1505`, retired
in `spec-retire: covers-traceability`, retro at
`capture/retros/covers-traceability.md`. Six review findings were left flagged
(not accepted) on the run's review receipt for reviewer veto — five Standards
judgment items and one low Coverage item (NBSP gap unanchors an annotation,
fail-closed); reopen any by reviewer request. The craft-tickets thin-by-default
slicing flip also landed (`1691017`), graduating its parked idea; its IDEAS.md
line and the new gate-invoke-to-Go idea await the drain. Reviewer decision made
in-session and closed: covers annotations are bracket-adjacent only.

Separate workflow, untouched: lifecycle run FT195 (`go-build-cache-footprint`)
still awaits its exact-candidate review, now against a moved main tip
(recomposition will discard nothing — no review is bound); its assignment
`f576124…` holds real uncommitted work in its worktree and must be
checkpointed or abandoned before FT195 promotes.

Unpushed commits await the reviewer's push decision.

## Next command

`/bench-what-next`

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
