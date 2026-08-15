# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `cee1dfb`, clean, 1 unpushed commit
Spec: none staged.
Gate: green at `cee1dfb` — current exact tree

## State

FT189 remains retired. FT207 holds the unresolved decision whether malformed-admin
refusal fronts all worktree-mutating calls; FT185 remains the next ready build.
Provider scorecards remain under `capture/agent-performance/`: OpenAI retains the
gate-run-transaction aggregate and Claude remains intentionally unknown. Capture is
drained; the existing full gate is green on this exact tree.

## Next command

`/bench-shape-idea`

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
