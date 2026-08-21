# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — clean once the spec-staging commit lands.
Spec: `specs/land-executable-freshness/spec.md`, Status: staged, reviewer-approved with its one ticket.
Gate: green at the staging commit; `dist/bench` rebuilt fresh this session.

## State

FT242 was re-scoped this session: its original ask (a spec amendment reaches
the destination through one sanctioned step) already shipped as FT225
(`cb1462a6`/`33aa5258`, retired `bef52480`); the ft230 detour came from a stale
`dist/bench` enforcing the retired refusal. The staged spec carries the
residual: `bench worktree land` proves its own executable through
`freshness.Verify` before any repository proof, skips where
`scripts/go-build.inputs` is absent, and exempts `--resume`. One ticket
(`01-refuse-stale-landing-executable`, blocked by none) delivers the slice,
including the FT242 board rewrite in `roadmap/FT242.md` and `ROADMAP.md`.
Review: 2 iterations to accept; the ordering-row learning is in
`capture/learnings.md` awaiting the next drain. One commit remains unpushed
from before this session.

## Next command

`/bench-implement-spec` — specs/land-executable-freshness, fresh mid-tier
(opus) session, one retained integration source, one ticket.

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
