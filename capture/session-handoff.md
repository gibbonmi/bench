# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `0d196bc5`, uncommitted `$bench-what-next` batch awaiting approval
Spec: none staged; `spec-authoring-and-light-path` landed at `8c8e5325` and retired at `0d196bc5`.
Gate: cached green `cca928c` is stale; the uncommitted maintenance batch has not been gated.

## State

The proposed drain empties two ideas, two journal entries, and the
`spec-authoring-and-light-path` retro. It adds the live-tree ordinary-gate gap
to FT120, outcome-shaped guidance to FT100, landing-result clarity to FT169,
filesystem-state edge inventory to FT144, fixture-census ownership to FT174,
and new FT206 for destination metadata before exact-candidate review. The
journal rule is already shipped; its narrow A/B rerun requests no follow-up.

No roadmap row was retired or reworded: the trusted schema-4 context reported
61 live rows, no staged specs, no spec-history rows, and no occurrence
discrepancies. No restructure was requested. Reviewer approval is required
before the one scoped commit.

## Next command

`$bench-shape-idea` — FT189's reproduced upstream `git worktree list` hang needs a Bench-owned refusal or execution bound.

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
