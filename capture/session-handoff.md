# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — pre-commit HEAD `4f2a969`
Spec: four staged specs (`exact-prospective-landing`, `go-build-cache-footprint`, `ft187-communication-surface-cut`, `pre-push-guard-visibility`)
Gate: full run owed on the drain commit; last cached verdict is stale

## State

**Phase reached: `/bench-what-next` drain over the ft156-anchor-registry close.**
The batch diff reconciles the goal tracks around FT156's shipped registry
(promoted `8e9f92d`, retired `e20d76e`), fixes FT107's stale lighter-path
pointer to FT144, drains three ideas (FT89 gains the CLI-inventory sync gate
check; the trace-bullet validation is dropped as already canonical; the
request-ID generator becomes FT196), and drains the ft156 retro (recs merged
into FT164, FT158, FT130, FT184; two dismissed as landed at `4f2a969`).
`capture/IDEAS.md` is empty and `capture/retros/` is gone.

Decisions that stay closed: FT156's registry/substring ruling lives in ADR 0012;
the staged-frontier-before-new-specs ordering stands.

## Next command

`/bench-implement-spec specs/exact-prospective-landing/spec.md` — FT188 first
to remove the writer lock, then FT195, FT187, FT135 per the refreshed
`## Recommended sequence`.

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
