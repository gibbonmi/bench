# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `102ae18`, clean tree, 17 unpushed commits
Spec: `specs/axi-aggregate-empty-migration/spec.md` (Status: staged), `specs/axi-bounded-projection-migration/spec.md` (Status: staged), `specs/axi-carriers-and-registry/spec.md` (Status: staged), `specs/axi-compatibility-oracle/spec.md` (Status: staged), `specs/axi-outcome-action-migration/spec.md` (Status: staged), `specs/single-build-serial-gate/spec.md` (Status: staged)
Gate: green at `b350bde` — stale, work tree `c470e97`

## State

FT173 implementation is underway on `specs/axi-compatibility-oracle/spec.md`,
first of the five staged foundation specs, followed by FT171's staged
`single-build-serial-gate` spec. Commit `102ae18e` landed the reviewed
pre-build repairs: pin advanced to `8ae1512f` (production bytes unchanged
since), release/fall-through/empty-form witnesses re-derived from production,
`pin-bound-and-byte-classes` split in two (12 tickets), breakdown re-reviewed
ready. `bench spec build start` is next; no assignment exists yet. Known
deferred repair: `specs/axi-bounded-projection-migration/tickets/contract-projection-routes.md`
still pins the retired subject `974020e4` — repair when that build starts.

## Next command

`/bench-implement-spec`

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
