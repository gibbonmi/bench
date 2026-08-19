# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`; the build runs in the integration worktree with label
`ft226 implementation` (assignment `6a22e66f47a693f42e306c0ec030e1ba`, request
`ft226-build`) — address it by label, never by a pasted path
Branch: `main` — the source's frozen review base is `b2d3d625`, which a concurrent
session has since moved past
Spec: `specs/ft226-test-home-isolation/spec.md` — `Status: staged`
Phase reached: `/bench-implement-spec --full`, three tickets and the review repair
pass all landed green in the source; ready to reauthorize and land

## State

All three tickets are committed on green in the integration source, and the
three-axis review's repair pass has landed on top of them. The spec's Build
verification log carries every acceptance row's evidence — the sweep counts, the
three mutation probes, and SW3 — and is the source for all of it.

The review found twelve items across three axes, collapsing to eight repair
targets; `reviews/ft226-test-home-isolation.md` holds the pickup state and is
deleted by the commit that closes the findings. The reviewer approved the two
judgment calls: the spec's probe (b) command now takes its `GOTMPDIR` form, and
the test-local `withinDir` deliberately stays independent of the production
`insidePool`.

**The landing needs a reauthorize.** A concurrent session landed `e1b44e62` and
`a9e8e232` on `main` during this build, so the source's base is behind the
destination; run `bench worktree reauthorize` to the current base before
`bench worktree land`. Those two gate runs also leaked ten fresh `001-<digits>`
pool keys, because `main` does not yet carry ticket 01. A second sweep pass clears
them once this spec lands; the script was throwaway and never committed, so
regenerate it from the spec's SW1 predicate.

## Next command

`bench worktree reauthorize` to the current destination HEAD, then
`bench worktree land`, then `/bench-final-check`.

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
