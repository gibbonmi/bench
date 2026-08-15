# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `64c409b`, 5 dirty paths, 18 unpushed commits
Spec: none staged.
Gate: green at `e6e7c6b` — stale, work tree `6a4b4c1`

## State

Agent-performance capture is provider-specific under
`capture/agent-performance/`. Its shared contract bounds current aggregates and
evidence; `open-ai-models.md` carries the `gate-run-transaction` observations and
`claude-models.md` preserves unknowns until a Claude landing supplies evidence.
`/bench-final-check` now refreshes every participating provider during the
implementation retro, and `/bench-what-next` includes those durable scorecards in
the reviewed capture commit. Focused anchors, conformance, retros, and structure
checks are green; the whole-tree gate remains for the landing commit.

## Next command

`bench commit -m "capture: split agent performance by provider" .agents/commands/bench-final-check.md .agents/commands/bench-what-next.md .bench/BENCH.md capture/agent-performance.md capture/agent-performance/README.md capture/agent-performance/open-ai-models.md capture/agent-performance/claude-models.md capture/session-handoff.md`

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
