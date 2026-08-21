# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — clean once the drain commit lands.
Spec: none staged.
Gate: green at the drain commit.

## State

FT242 shipped at `84b7c4b0` and is retired: `bench worktree land` proves its own
executable before any repository proof, only where the repository declares Go
build inputs, with `--resume` exempt. That landing also carried reviewer-directed
housekeeping — the `regroup` example profile retired across all four registries,
`COMPLIANCE_ASSESSMENT.md` removed, `ui_example/` untracked but kept on disk.

The drain that follows it emptied every capture source. Both journal entries and
every retro recommendation have a disposition; the retro file is deleted and the
journal is back to its schema heading. Most of it merged as occurrence evidence
onto rows that already carried the same face: FT162 (`bench handoff` overwriting
the next command), FT224 (the `--discard-ignored` flag the release verb rejects),
FT233 (gate logs blocking release, and exit 1 after a completed publication),
FT236 (citation line drift), FT238 (`bench diff` rejecting `--source-tip`).
FT214 gained three map-discipline clauses and three occurrences.

One new row, FT243, and it is the sequence's top line: a dated `capture/learnings.md`
entry written as a bullet rather than a heading is skipped with no entry, no
malformed record, and no parse failure, while the source still reports as parsed
at full byte count. Both entries drained this pass were invisible that way.

Contestable calls left for veto: FT243 ranked above the two better-evidenced
refusal rows on severity rather than evidence; the retro's
housekeeping-convention recommendation dismissed rather than rowed; and FT233
and FT224 named as a restructure candidate but not merged, since a default drain
names rather than applies.

## Next command

`/bench-write-spec` — FT243: a capture entry the parser cannot see is reported
as zero, not as a failure.

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
