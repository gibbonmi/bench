# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `df749df`, 3 dirty paths, 10 unpushed commits
Spec: `specs/ft128-agent-line-binding/spec.md` (Status: staged)
Gate: green at `2036d35` — current

## State

- **FT126 is closed.** Its roadmap row and pending retro are drained; the retro's
  four actionable recommendations now live on FT169, FT133, FT101, and FT120.
- **The recurrence migration contract follows current roadmap ownership.** Retired
  FT126 is no longer required by the one-time baseline inventory, while the eight
  remaining migrated rows keep their exact-count checks and mutation bite.
- **FT128 is the next actionable item.** Its staged spec owns the fork verdict,
  harness-native denial, and static command-token conformance sweep.
- The dev gate is green for the approved drain and conformance repair. Push and
  ship-tier verification have not run.

## Next command

`$bench-implement-spec ft128-agent-line-binding`

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
