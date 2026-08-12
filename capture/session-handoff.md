# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `e39d9f2`, clean tree, 52 unpushed commits
Spec: `specs/axi-query-disclosure/spec.md` (Status: staged), `specs/single-build-serial-gate/spec.md` (Status: staged)
Gate: green at `7fd5b28` — stale, work tree `20a8f05`

## State

The composed AXI query disclosure candidate includes the final-evidence closure
and is ready for full semantic review against the staged spec. No repair pickup
remains. The accepted AXI set, default-extraction disclosure scope,
compatibility partition, registry-owned membership, call-site action
provenance, QD1–QD6, and the signed QD5 ledger remain closed. Any stale paused
worktree diff is excluded from this candidate; foreign worktrees remain
preserved.

## Next command

`$bench-review-implementation specs/axi-query-disclosure/spec.md --full`

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
