# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `7bc6ee4e`; the roadmap drain is uncommitted pending reviewer approval
Spec: `specs/axi-coherent-diff/spec.md` (Status: staged), `specs/axi-query-disclosure/spec.md` (Status: staged), `specs/single-build-serial-gate/spec.md` (Status: staged)
Gate: cached green at pre-drain tree `51f4ad3`; this batch is not gated

## State

Pocock-guidance-doctrine is implemented at `a354cede`; this drain removes FT107
and retires its spec. Its two ideas route to FT89 (stale lifecycle references)
and FT100 (wrap-neutral prose weight). Its malformed journal bullets route to
FT164 (prove mutations landed) and FT169 (resolve the landing checkout). The
retro's cross-harness premise check routes to FT158 and its clean
conformance-probe recommendation routes to FT120; it carried no
repair-attribution table.

The minimum revisit dates for FT38 and FT170 have elapsed. This batch marks
both LOW and decision-required without putting either in the recommended
sequence; that activation and the FT120 routing remain reviewer-veto calls.
The refreshed sequence is FT173, FT171, then FT198. If approved, commit the
whole batch once on green with a subject ending
`spec-retire: pocock-guidance-doctrine`.

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
