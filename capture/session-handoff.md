# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `b022697`, 3 dirty paths, 18 unpushed commits
Spec: `specs/axi-coherent-diff/spec.md` (Status: staged), `specs/axi-query-disclosure/spec.md` (Status: staged), `specs/bench-preflight/spec.md` (Status: staged), `specs/single-build-serial-gate/spec.md` (Status: staged)
Gate: green at `9bae8f2` — stale, work tree `08b996f`

## State

`/bench-implement-spec --full bench-preflight` phase reached: **complete**.
All 25 coverage rows landed (7 implementation + 4 review-repair tickets, every
landing gate-green); `bench preflight <build|review> <slug>` is live, advertised,
anchor-pinned, and the routing conformance checker reads `commandRegistry`
again. Composed review (opus/high, fresh context) returned 10 findings — 8
repaired and re-verified by a second fresh review, which added 6 minors left
for reviewer disposition (listed in `capture/retros/bench-preflight.md` with
the two held vetoes: test-comment ID density, mixed-tag map policy). Retro
written; drain owns it.

Parked pre-reshape specs (`axi-coherent-diff`, `axi-query-disclosure`,
`single-build-serial-gate`) await re-rank. Reviewer context note: the kit's
gate descends from no-mistakes, skills from the Pocock material, CLI shape
from AXI — hook-enforced preflight blocking (the no-mistakes posture) is
priced in this spec's Out of scope if ever wanted.

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
