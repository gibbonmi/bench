# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `cfab082`, 1 dirty path, 11 unpushed commits
Spec: `specs/axi-aggregate-empty-migration/spec.md` (Status: staged), `specs/axi-bounded-projection-migration/spec.md` (Status: staged), `specs/axi-carriers-and-registry/spec.md` (Status: staged), `specs/axi-compatibility-oracle/spec.md` (Status: staged), `specs/axi-outcome-action-migration/spec.md` (Status: staged), `specs/single-build-serial-gate/spec.md` (Status: staged)
Gate: green at `685bfe0` — stale, work tree `7bd4c59`

## State

The `$bench-what-next` drain reconciled the just-retired
`spec-ticket-handoff-contract` spec and drained its retro. Its recommendations
merge into FT162 (already-owned promotion diagnostics), FT174 (ticket preflight
and compound-row mutation inventory), FT107 (exact prospective-source binary
guidance), and FT158 (one bounded post-review repair). The snapshot has no
ideas, learnings, or occurrence discrepancies.

The approved batch updates `ROADMAP.md`, removes the retro, and refreshes this
handoff. It needs one ordinary gated commit; no retirement suffix applies because
the retirement already landed as `cfab0821`.

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
