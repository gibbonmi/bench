# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `2181369f`, drain committed in the next landing
Spec: five staged FT173 specs (`specs/axi-compatibility-oracle/spec.md` first,
then `axi-carriers-and-registry`, `axi-outcome-action-migration`,
`axi-bounded-projection-migration`, `axi-aggregate-empty-migration`) and
`specs/single-build-serial-gate/spec.md` (Status: staged, FT171)
Gate: green at `66ac6be` — the tree before this drain's uncommitted diff

## State

The completed `/bench-what-next` drain renumbers the duplicate second FT199 row
to FT200 (and its Dependencies entry), binds the staged
`single-build-serial-gate` spec to FT171 and makes it sequence item 2, and
drains the three parked ideas: FT201 owns the cancel-signal-literal conformance
check and the `Pdeathsig` decision, while FT120 owns the gitOutput/FIFO-orphan
claim and its required repro. The prior cancel-signal journal entry is
dismissed because `6fbf404` fixed its incident; its residual is in FT201. FT166
records the malformed bullet form that made it invisible to parser-backed
readers. The final journal entry folds into FT107: delegation that changes who
performs reviewer-requested work is surfaced before spawning. A stale FT135
historical spec-path reference is removed. FT173 remains top priority and its
five-foundation-spec order is unchanged.

## Next command

`$bench-implement-spec`

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
