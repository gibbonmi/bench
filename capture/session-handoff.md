# Session handoff

Repository: `b35df3e36b4767b8d0d21b99ad0c0e41-b1a61bd8a7c055204b9db0b9364a942d` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/.bench/worktrees/bench-2826441890/b35df3e36b4767b8d0d21b99ad0c0e41-b1a61bd8a7c055204b9db0b9364a942d`
Branch: worktree `drain-2026-08-24` — base `76208d5b` on `main`, the light-path commit `7704626a` plus the drain batch staged for one landing
Spec: none staged; the tickets-only folder `specs/spec-retire-primary-refusal/` closes at this landing
Gate: green at `7704626a` — the light-path commit's own run

## State

The 2026-08-24 drain is complete. The inbox was already empty. The three
learnings and the retro drained. The stale-premise pair merged into FT99, and
the request-token pair merged into FT238 as an occurrence.

The prose-bound
recommendation opened FT250 (`Next: kit-edit`). The retire-refusal
recommendation was built in this session as the light path. `bench spec
retire` now refuses the primary checkout, committed at `7704626a`.

The flow target does not hold (net +4). Two fold candidates await a reviewer
verdict: FT223 into FT141, and FT235 into FT238. Neither fold is applied.

## Next command

`/bench-implement-spec` — FT238, per the recommended sequence.

## Shape

Rewritten in full at every phase close, pruned rather than accreted. A fresh
session pays for every line it reads cold; drop anything it would not act on.

Operational gotchas are placed by lifetime, not copied here. One that recurs across
phases belongs in `projects/benchkit.md`'s cold-session notes. One scoped to a build
belongs instead in that spec's coverage rows.

This file names at most when you'll hit one, never the command — a second copy
drifts from the source.

Keep the three sections above. **State** holds what is true now, including anything
uncommitted. **Next command** holds the exact harness-native invocation, not a
description of it. This section is the third.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
