# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `3c50349` plus the shaping commit, pushed through `4422d3a`
Spec: none staged; the map is spec-ready.
Gate: green at `edc15e9` — stale; diff since is docs/skills plus two conformance
anchor lines (no gate run this session — shaping only).

## State

- **The cost-follows-project-size map closed its FT91 arm.** All decisions are
  in `decisions/cost-follows-project-size.md` (#2 timing answer, #3 tier-split
  answer, rewritten Handoff): the gate splits into a dev tier and a ship tier
  (`bench prep-release`); fifteen-check fan-out is dead by measurement; the
  inner `go test` excludes release-only packages; per-check timing becomes
  permanent driver-owned gate output; pre-push stays fast-tier. Only ticket #6
  remains open — opportunistic cheap-tier retest, non-blocking.
- **Decisions that stay closed:** no build-approval prompt after final-check;
  ship tier never on pre-push; restaging is not check-weakening (authority
  binds at the release path's evidence refusal); `-count=1` and caching stay
  deferred, now with a measured price on the map.
- **Measured facts a spec session may want:** conformance phase ≈826 s, 99.8%
  in `checkPackageCoreAndGuards`; `checkReleasePreflight` ≈372 s;
  uncached `internal/preflight` is 676 s (slow, not hung — go-test cache
  normally hides it); the inner suite recurses through the conformance package
  on a cache miss. Probe file deleted; numbers live in the map's #2 answer.
- Known ambient facts, unchanged: 17 recovered worktrees from earlier sessions
  left untouched; structure budget violations stand.

## Next command

`/bench-write-spec` on a fresh mid-tier session, off
`decisions/cost-follows-project-size.md`'s Handoff (the FT91 tier split —
single spec, dependency order n/a).

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
