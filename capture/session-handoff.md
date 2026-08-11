# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `03b0e95`, clean tree, 10 unpushed commits
Spec: `specs/axi-coherent-diff/spec.md` (Status: staged), `specs/axi-query-disclosure/spec.md` (Status: staged), `specs/bench-preflight/spec.md` (Status: staged), `specs/single-build-serial-gate/spec.md` (Status: staged)
Gate: green at `a56d225` — stale, work tree `e139396`

## State

`/bench-implement-spec --full bench-preflight` phase reached: **all 7 tickets
landed green (fb827ebf..03b0e953), composed three-axis review delegate
(opus/high, fresh context) running**. All 25 coverage rows built; map
validates. Landed serially: 556bd98d routing-checker repair (PF23), 11e2c1d9
ResolveReviewBase export (PF25), c0fe8151 releasepreflight rename (PF24),
4dbe0902 preflight review command (PF1–PF8, PF15–PF17), a947821c bootstrap
diagnostics (PF12–PF14, PF22), 0284a3c9 build mode (PF9–PF11, PF21), 03b0e953
advertisement (PF18–PF20). Exit-report notes pending: C16's shell-route
mutation cannot bite (`bin/bench.sh` `*) route_binary` fallback dispatches
unknown tokens via the registry — pre-existing wrapper property); collapsed
per-story lines: stories 1/3 landed as one core ticket. Remaining: accepted
review findings as repair tickets (if any), final landing
`bench commit --spec bench-preflight`, inline final-check + retro. Parked
pre-reshape specs (`axi-coherent-diff`, `axi-query-disclosure`,
`single-build-serial-gate`) await re-rank, not active. Post-spec queue:
Spec C (doctrine adoption), then a `/bench-what-next` drain.

## Next command

`/bench-implement-spec --full bench-preflight`

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
