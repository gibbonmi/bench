# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — clean once the drain batch commit lands on `61f1db8c`.
Spec: none staged.
Gate: green.

## State

The 2026-08-20 `/bench-drain` pass is the batch this commit carries: the FT230
retro and both learnings entries are drained and removed, all capture sources
are empty, and no spec is staged. Dispositions: FT242 is new (decision
required — a spec amendment reaches the destination through one sanctioned
step: land adopts the source's spec bytes, or a `bench spec sync` verb);
merges landed on FT224 (spec-retire next-step omits the detail-file deletion),
FT162 (`bench handoff` overwrites the Next-command section), FT214 (two
craft-spec map clauses: name the test function you read, sweep for deleted
bytes), FT215 (capture-only fast lane joins the scoped-gate decision), and
FT238 (sixth piece: the phase-close capture-batch rule). The board otherwise
reconciled clean — nothing shipped since the last drain. The reviewer priced
FT242 to the front of the sequence.

## Next command

`/bench-write-spec` — FT242: a spec amendment reaches the destination through
one sanctioned step. The spec's first question is the shape: land adopts the
source's spec bytes, or a `bench spec sync` verb.

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
