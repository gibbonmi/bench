# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `a189b87`, 24 dirty paths, 19 unpushed commits
Spec: `specs/axi-coherent-diff/spec.md` (Status: staged), `specs/axi-query-disclosure/spec.md` (Status: staged), `specs/single-build-serial-gate/spec.md` (Status: staged)
Gate: green at `d94dc30` — stale, work tree `244e150`

## State

The 2026-08-11 roadmap drain is complete. `bench-preflight` and
`remove-spec-build-lifecycle` are retired; their retros, five ideas, and seven
open learnings are fully dispositioned in `ROADMAP.md`. FT184 and FT196 left
with the deleted lifecycle. The remaining staged specs are
`axi-coherent-diff`, `axi-query-disclosure`, and `single-build-serial-gate`.

Next is the Pocock-alignment program's Spec C through FT107. Its reviewed source
is the decision map retained by
`bench spec history remove-spec-build-lifecycle`: adopt decisions #4, #5, #6,
#8, #9, and #10 as one kit-prose batch, re-reading both named upstream sources
for drift first. Gate authority and all three review axes stay closed. The
one-build Opus `/bench-debug` override is not a standing rule unless the
reviewer widens it during spec writing. After Spec C, rescope
`axi-coherent-diff` off the retired spec-build/`axi.Action` prerequisite.

## Next command

`$bench-write-spec`

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
